//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package hdfs

import (
	"bufio"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/colinmarc/hdfs"
	"github.com/demdxx/gocast"
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/converter"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/stream"
	"github.com/labstack/gommon/log"
)

var (
	hdfsParamParser = regexp.MustCompile(`\$\{([^}:]+)(?:\|([^}]+))\}`)
)

// StreamHDFS object
type StreamHDFS struct {
	sync.Mutex
	client                  *hdfs.Client             // Connection
	fileNamePattern         string                   // example: path/{{date}}{{iterator}}.ext
	fileNameDatePattern     string                   // default: '2006-01-02_15'
	fileNameIteratorPattern string                   // default: '_%iter%' if it's zero then empty
	maxFileSize             int                      //
	blockSize               int                      //
	writeMaxDuration        time.Duration            //
	converter               converter.Converter      //
	recordSeparator         []byte                   //
	messageTimeParam        string                   // param name in Message object
	tmpDirectoryPath        string                   //
	record                  stream.Query             //
	buffer                  chan eventstream.Message //
	currentOutputName       string                   //
	currentOutput           io.WriteCloser           //
	currentOutputSize       int                      //
	currentOutputIterator   int64                    //
	currentOutputTime       time.Time                //
	processTimer            *time.Ticker
}

// New HDFS stream
func New(opt stream.Options) (stream.Streamer, error) {
	return NewStreamHDFS(
		storage.Get(opt.Connection).(*hdfs.Client),
		opt.Target,
		gocast.ToInt(opt.Get("maxsize")),
		gocast.ToInt(opt.Get("buffer")),
		time.Duration(gocast.ToInt(opt.Get("duration"))),
		converter.ByName(gocast.ToString(opt.Get("format"))),
		[]byte(gocast.ToString(opt.Get("separator"))),
		gocast.ToString(opt.Get("timefield")),
		gocast.ToString(opt.Get("tmpdir")),
		opt.Fields,
	)
}

// NewStreamHDFS streamer
func NewStreamHDFS(
	client *hdfs.Client,
	fileNamePattern string,
	maxFileSize int,
	blockSize int,
	duration time.Duration,
	converter converter.Converter,
	recordSeparator []byte,
	messageTimeParam string,
	tmpDirectoryPath string,
	fields interface{},
) (stream.Streamer, error) {
	var (
		fileNameDatePattern     string
		fileNameIteratorPattern string
	)

	// Prepare file pattern and extract format
	if all := hdfsParamParser.FindAllStringSubmatch(fileNamePattern, -1); len(all) > 0 {
		for _, item := range all {
			switch item[1] {
			case "date":
				fileNamePattern = strings.Replace(fileNamePattern, item[0], "{{date}}", -1)
				fileNameDatePattern = item[2]
			case "iterator":
				fileNamePattern = strings.Replace(fileNamePattern, item[0], "{{iterator}}", -1)
				fileNameIteratorPattern = item[2]
			}
		}
	}

	if blockSize < 1 {
		blockSize = 1000
	}

	if duration <= 0 {
		duration = time.Second * 1
	}

	if "" == fileNameDatePattern {
		fileNameDatePattern = "2006-01-02_15"
	}

	if "" == fileNameIteratorPattern {
		fileNameIteratorPattern = "_%iter%"
	}

	if len(recordSeparator) < 1 {
		recordSeparator = []byte("\n")
	}

	query, err := stream.NewQueryByRaw("", fields)
	if nil != err {
		return nil, err
	}

	return &StreamHDFS{
		client:                  client,
		fileNamePattern:         fileNamePattern,
		fileNameDatePattern:     fileNameDatePattern,
		fileNameIteratorPattern: fileNameIteratorPattern,
		maxFileSize:             maxFileSize,
		blockSize:               blockSize,
		writeMaxDuration:        duration,
		converter:               converter,
		recordSeparator:         recordSeparator,
		messageTimeParam:        messageTimeParam,
		tmpDirectoryPath:        tmpDirectoryPath,
		record:                  *query,
		buffer:                  make(chan eventstream.Message, blockSize*2),
	}, nil
}

// Put message to stream
func (s *StreamHDFS) Put(msg eventstream.Message) error {
	s.buffer <- msg
	return nil
}

// Close implementation
func (s *StreamHDFS) Close() error {
	if nil != s.processTimer {
		s.processTimer.Stop()
		s.processTimer = nil
	}

	s.writeBuffer(true)
	close(s.buffer)
	return nil
}

// Process loop
func (s *StreamHDFS) Process() {
	if nil != s.processTimer {
		s.processTimer.Stop()
	}

	s.processTimer = time.NewTicker(time.Millisecond * 5)
	ch := s.processTimer.C

	for _, ok := <-ch; ok; {
		if err := s.writeBuffer(false); nil != err {
			log.Error(err)
		}
	}
}

// writeBuffer all data
func (s *StreamHDFS) writeBuffer(flush bool) error {
	s.Lock()
	defer s.Unlock()

	if !flush {
		if c := len(s.buffer); c < 1 || (s.blockSize > c && time.Now().Sub(s.currentOutputTime) < s.writeMaxDuration) {
			return nil
		}
	}

	var (
		writer, err        = s.writer()
		count              int
		data               []byte
		currentOutputCache string
	)

	if nil != err {
		return err
	}

	if !s.currentOutputTime.IsZero() {
		currentOutputCache = s.currentOutputTime.Format(s.fileNameDatePattern)
	}

	for {
		msg, ok := <-s.buffer
		if !ok {
			break
		}

		var (
			lastTime, _ = s.msgTime(msg)
			lastTimeFmt = lastTime.Format(s.fileNameDatePattern)
		)

		// Init start params
		if "" == currentOutputCache {
			currentOutputCache = lastTimeFmt
			s.currentOutputTime = lastTime
		}

		// Writer file and create new storage for output
		if (s.maxFileSize > 0 && s.maxFileSize <= s.currentOutputSize) || lastTimeFmt != currentOutputCache {
			if err = s.writeIterationFile(); nil == err {
				currentOutputCache = lastTimeFmt
				s.currentOutputTime = lastTime
			} else {
				writer, err = s.writer()
			}
			if nil != err {
				return err
			}
		}

		// Prepare data for saving
		if data, err = s.converter.Marshal(s.record.Extract(msg)); nil != err {
			return err
		}

		// Write prepared data
		if count, err = writer.Write(data); nil != err {
			return err
		}

		// Counte count of data
		s.currentOutputSize += count

		// Write separator to store file
		if len(s.recordSeparator) > 0 {
			count, err = writer.Write(s.recordSeparator)
			s.currentOutputSize += count
		}

		if nil != err {
			break
		}
	}
	return err
}

// writeIterationFile
func (s *StreamHDFS) writeIterationFile() (err error) {
	s.currentOutputSize = 0
	if nil != s.currentOutput {
		err = s.currentOutput.Close()
		s.currentOutput = nil
	}

	if nil == err {
		var (
			targetName = s.targetName()
			dir        = filepath.Dir(targetName)
		)

		if len(dir) > 0 && "/" != dir {
			if err = s.client.MkdirAll(dir, 0755); nil != err {
				return
			}
		}

		if err = s.copyAndRemoveFile(s.currentOutputName, targetName); nil == err {
			s.currentOutputIterator++
		} else if os.IsExist(err) {
			var files []os.FileInfo

			for {
				s.currentOutputIterator++
				targetName = s.targetName()
				targetFileName := filepath.Base(targetName)
				dir2 := filepath.Dir(targetName)
				exists := false

				// Get dir list
				if dir2 != dir || len(files) < 1 {
					if dir2 != dir && len(dir2) > 0 && "/" != dir2 {
						if err = s.client.MkdirAll(dir2, 0755); nil != err {
							return
						}
						dir = dir2
					}

					if files, err = s.client.ReadDir(dir); nil != err {
						return
					}
				}

				for _, fl := range files {
					if targetFileName == fl.Name() {
						exists = true
						break
					}
				}

				if !exists {
					// Try to copy again
					if err = s.copyAndRemoveFile(s.currentOutputName, targetName); nil == err {
						s.currentOutputIterator++
					}
					fmt.Println("=====", err)
					break
				}
			} // end for
		} else {
			fmt.Println("=====", err)
		}
	}
	return
}

// Get writer object
func (s *StreamHDFS) writer() (io.Writer, error) {
	if nil != s.currentOutput {
		return s.currentOutput, nil
	}

	if "" == s.tmpDirectoryPath {
		s.tmpDirectoryPath = os.TempDir()
	}

	// Create TMP directory path
	if _, err := os.Stat(s.tmpDirectoryPath); os.IsNotExist(err) {
		if err = os.MkdirAll(s.tmpDirectoryPath, 0755); nil != err {
			return nil, err
		}
	}

	file, err := ioutil.TempFile(s.tmpDirectoryPath, "hdfs")
	if nil != err {
		return nil, err
	}

	var (
		sections                = strings.Split(s.fileNamePattern, ".")
		writer   io.WriteCloser = file
	)

	// Add compression
	switch sections[len(sections)-1] {
	case "gz", "gzip":
		writer = gzip.NewWriter(writer)
	case "z", "zip":
		writer = zlib.NewWriter(writer)
	case "lz", "lzw", "lzma":
		writer = lzw.NewWriter(writer, lzw.LSB, 8)
	default:
		writer = file
	}

	s.currentOutputSize = 0
	s.currentOutputName = file.Name()
	s.currentOutput = &buffWriter{
		buff:   bufio.NewWriter(writer),
		closer: buffCloser{writer, file},
	}
	return s.currentOutput, err
}

// get message time from message
func (s *StreamHDFS) msgTime(msg eventstream.Message) (time.Time, error) {
	if s.messageTimeParam != "" {
		if v := msg.ItemCast(s.messageTimeParam, eventstream.FieldTypeDate, ""); nil != v {
			return v.(time.Time), nil
		}
		return time.Time{}, eventstream.ErrInvalidMessageFieldType
	}
	return time.Now(), nil
}

// build name of target file for current iteration
func (s *StreamHDFS) targetName() string {
	var iterator string
	if s.currentOutputIterator > 0 {
		iterator = strings.Replace(
			s.fileNameIteratorPattern,
			"%iter%",
			strconv.FormatInt(s.currentOutputIterator, 10),
			-1,
		)
	}

	return strings.NewReplacer(
		"{{date}}", s.currentOutputTime.Format(s.fileNameDatePattern),
		"{{iterator}}", iterator,
	).Replace(s.fileNamePattern)
}

// copyAndRemoveFile local file to HDFS
func (s *StreamHDFS) copyAndRemoveFile(src, dst string) (err error) {
	file, _ := os.Open(src)
	defer file.Close()

	stat, _ := file.Stat()
	fmt.Println("=====", file.Name(), stat.Size())

	if err = s.client.CopyToRemote(src, dst); nil == err {
		os.Remove(src)
	} else if !os.IsExist(err) {
		s.client.Remove(dst)
	}
	return err
}
