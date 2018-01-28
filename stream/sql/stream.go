//
// @project geniusrabbit::eventstream 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

// TODO: Store all messages if not sended yet

package sql

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
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
	debug            bool
}

// NewStreamSQL streamer
func NewStreamSQL(conn *sql.DB, blockSize int, duration time.Duration, query stream.Query, debug bool) (_ eventstream.SimpleStreamer, err error) {
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
		debug:            debug,
	}, nil
}

// NewStreamSQLByRaw query
func NewStreamSQLByRaw(conn *sql.DB, blockSize int, duration time.Duration, query string, fields interface{}, debug bool) (eventstream.SimpleStreamer, error) {
	q, err := stream.NewQueryByRaw(query, fields)
	if err != nil {
		return nil, err
	}
	return NewStreamSQL(conn, blockSize, duration, *q, debug)
}

// Put message to stream
func (s *StreamSQL) Put(msg eventstream.Message) error {
	s.buffer <- msg
	return nil
}

// Close implementation
func (s *StreamSQL) Close() error {
	if s.processTimer != nil {
		s.processTimer.Stop()
		s.processTimer = nil
	}

	s.writeBuffer(true)
	close(s.buffer)
	return nil
}

// Run loop
func (s *StreamSQL) Run() error {
	if s.processTimer != nil {
		s.processTimer.Stop()
	}

	s.writeLastTime = time.Now()
	s.processTimer = time.NewTicker(time.Millisecond * 5)
	ch := s.processTimer.C

	for _, ok := <-ch; ok; {
		if err := s.writeBuffer(false); err != nil {
			return err
		}
	}
	return nil
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

	if tx, err = s.conn.Begin(); err != nil{
		return
	}

	if stmt, err = tx.Prepare(s.query.Q); err != nil {
		tx.Rollback()
		return
	}

	// Writing loop of prepared requests
	for !stop {
		select {
		case msg := <-s.buffer:
			if s.debug {
				log.Printf("[clickhouse] %s\n", msg.JSON())
			}
			if _, err = stmt.Exec(s.query.ParamsBy(msg)...);  err != nil {
				stop = true
			}
		default: 
			stop = true
		}
	}

	if err == nil {
		stmt.Exec()
		err = tx.Commit()
	} else {
		tx.Rollback()
	}

	s.writeLastTime = time.Now()
	return
}
