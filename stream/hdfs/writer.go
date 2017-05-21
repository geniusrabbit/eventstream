//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package hdfs

import (
	"bufio"
	"io"
)

// buffWriter and closer
type buffWriter struct {
	buff   *bufio.Writer
	closer io.Closer
}

func (b buffWriter) Write(p []byte) (int, error) {
	return b.buff.Write(p)
}

func (b buffWriter) Close() (err error) {
	if nil != b.buff && nil != b.closer {
		_, err = b.buff.Flush(), b.closer.Close()
		b.buff, b.closer = nil, nil
	}
	return
}

type buffCloser []io.Closer

func (b buffCloser) Close() (err error) {
	for _, c := range b {
		if nil != c {
			err = c.Close()
		}
	}
	return
}
