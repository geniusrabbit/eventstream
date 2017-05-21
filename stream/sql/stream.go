//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package eventstream

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
	"github.com/labstack/gommon/log"
)

// StreamSQL stream
type StreamSQL struct {
	sync.Mutex
	conn             *sql.DB
	buffer           chan eventstream.Message
	blockSize        int
	writeMaxDuration time.Duration
	writeLastTime    time.Time
	query            stream.Query
	processTimer     *time.Ticker
}

// NewStreamSQL streamer
func NewStreamSQL(conn *sql.DB, blockSize int, duration time.Duration, query stream.Query) stream.Streamer {
	if blockSize < 1 {
		blockSize = 1000
	}

	if duration <= 0 {
		duration = time.Second * 1
	}

	return &StreamSQL{
		conn:             conn,
		buffer:           make(chan eventstream.Message, blockSize*2),
		blockSize:        blockSize,
		writeMaxDuration: duration,
		query:            query,
	}
}

// NewStreamSQLByRaw query
func NewStreamSQLByRaw(conn *sql.DB, blockSize int, duration time.Duration, query string, fields interface{}) (stream.Streamer, error) {
	q, err := stream.NewQueryByRaw(query, fields)
	if nil != err {
		return nil, err
	}
	return NewStreamSQL(conn, blockSize, duration, *q), nil
}

// Put message to stream
func (s *StreamSQL) Put(msg eventstream.Message) error {
	s.buffer <- msg
	return nil
}

// Close implementation
func (s *StreamSQL) Close() error {
	if nil != s.processTimer {
		s.processTimer.Stop()
		s.processTimer = nil
	}

	s.writeBuffer(true)
	close(s.buffer)
	return nil
}

// Process loop
func (s *StreamSQL) Process() {
	if nil != s.processTimer {
		s.processTimer.Stop()
	}

	s.writeLastTime = time.Now()
	s.processTimer = time.NewTicker(time.Millisecond * 5)
	ch := s.processTimer.C

	for _, ok := <-ch; ok; {
		if err := s.writeBuffer(false); nil != err {
			log.Error(err)
		}
	}
}

// writeBuffer all data
func (s *StreamSQL) writeBuffer(flush bool) (err error) {
	s.Lock()
	defer s.Unlock()

	if !flush {
		if c := len(s.buffer); c < 1 || (s.blockSize > c && time.Now().Sub(s.writeLastTime) < s.writeMaxDuration) {
			return
		}
	}

	var (
		tx   *sql.Tx
		stmt *sql.Stmt
		stop = false
	)

	if tx, err = s.conn.Begin(); nil != err {
		return
	}

	if stmt, err = tx.Prepare(s.query.Q); nil != err {
		tx.Rollback()
		return
	}

	// Writing loop of prepared requests
	for !stop {
		select {
		case msg := <-s.buffer:
			fmt.Println("> =======", s.query.Q)
			fmt.Println("@ =======", msg)
			fmt.Println("$ =======", s.query.ParamsBy(msg))
			if _, err = stmt.Exec(s.query.ParamsBy(msg)...); nil != err {
				stop = true
			}
		default:
			stop = true
		}
	}

	if nil == err {
		stmt.Exec()
		err = tx.Commit()
	} else {
		tx.Rollback()
	}

	s.writeLastTime = time.Now()
	return
}
