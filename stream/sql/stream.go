//
// @project geniusrabbit::eventstream 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

// TODO: Store all messages if not sended yet

package sql

import (
	"database/sql"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
)

// Connector to DB
type Connector interface {
	Connection() (*sql.DB, error)
}

// StreamSQL stream
type StreamSQL struct {
	sync.Mutex
	connector        Connector
	buffer           chan eventstream.Message
	blockSize        int
	writeMaxDuration time.Duration
	writeLastTime    time.Time
	query            stream.Query
	processTimer     *time.Ticker
	debug            bool
}

// NewStreamSQL streamer
func NewStreamSQL(connector Connector, blockSize int, duration time.Duration, query stream.Query, debug bool) (_ eventstream.Streamer, err error) {
	if blockSize < 1 {
		blockSize = 1000
	}

	if duration <= 0 {
		duration = time.Second * 1
	}

	return &StreamSQL{
		connector:        connector,
		buffer:           make(chan eventstream.Message, blockSize*2),
		blockSize:        blockSize,
		writeMaxDuration: duration,
		query:            query,
		debug:            debug,
	}, nil
}

// NewStreamSQLByRaw query
func NewStreamSQLByRaw(connector Connector, blockSize int, duration time.Duration, query string, fields interface{}, debug bool) (eventstream.Streamer, error) {
	q, err := stream.NewQueryByRaw(query, fields)
	if err != nil {
		return nil, err
	}
	return NewStreamSQL(connector, blockSize, duration, *q, debug)
}

// Put message to stream
func (s *StreamSQL) Put(msg eventstream.Message) error {
	s.buffer <- msg
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

// Check message value
func (s *StreamSQL) Check(msg eventstream.Message) bool {
	return true
}

// writeBuffer all data
func (s *StreamSQL) writeBuffer(flush bool) (err error) {
	s.Lock()
	defer s.Unlock()

	var conn *sql.DB
	if conn, err = s.connector.Connection(); err != nil {
		return
	}

	defer s.recError()

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

	if tx, err = conn.Begin(); err != nil {
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
				s.log(msg.JSON())
			}
			if _, err = stmt.Exec(s.query.ParamsBy(msg)...); err != nil {
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

///////////////////////////////////////////////////////////////////////////////
/// Logs
///////////////////////////////////////////////////////////////////////////////

func (s *StreamSQL) recError() {
	if rec := recover(); rec != nil {
		s.logError(rec)
		s.logError(string(debug.Stack()))
	}
}

func (s *StreamSQL) log(args ...interface{}) {
	if len(args) > 0 {
		log.Println("[clickhouse] ", fmt.Sprintln(args...))
	}
}

func (s *StreamSQL) logError(err interface{}) {
	if err != nil {
		log.Println("[clickhouse] ", err)
	}
}
