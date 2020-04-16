//
// @project geniusrabbit::eventstream 2017 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2020
//

// TODO: Store all messages if not sended yet

package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/geniusrabbit/eventstream"
)

var (
	errInvalidQueryObject = errors.New("[sql] invalid query object")
)

// Connector to DB
type Connector interface {
	Connection() (*sql.DB, error)
}

// StreamSQL stream
type StreamSQL struct {
	// Debug mode of the stream
	debug bool

	// ID of the stream
	id string

	// Connector interface
	connector Connector

	buffer        chan eventstream.Message
	blockSize     int // size of the block suitable to save into DB
	flushInterval time.Duration
	writeLastTime time.Time

	// Query prepared data formater object
	query *Query

	// Time ticker
	processTimer *time.Ticker

	isWriting int32
}

// NewStreamSQL creates streamer object for SQL based stream integration
func NewStreamSQL(id string, connector Connector, options ...Option) (eventstream.Streamer, error) {
	stream := &StreamSQL{
		id:        id,
		connector: connector,
	}
	for _, opt := range options {
		if err := opt(stream); err != nil {
			return nil, err
		}
	}
	if stream.query == nil {
		return nil, errInvalidQueryObject
	}
	if stream.blockSize < 1 {
		stream.blockSize = 1000
	}
	if stream.flushInterval <= 0 {
		stream.flushInterval = time.Second * 1
	}
	stream.buffer = make(chan eventstream.Message, stream.blockSize*2)
	return stream, nil
}

// ID returns unical stream identificator
func (s *StreamSQL) ID() string {
	return s.id
}

// Put message to stream
func (s *StreamSQL) Put(ctx context.Context, msg eventstream.Message) error {
	if s.debug {
		log.Println("[stream] put message", msg)
	}
	s.buffer <- msg
	return nil
}

// Run SQL writer daemon
func (s *StreamSQL) Run(ctx context.Context) error {
	if s.processTimer != nil {
		s.processTimer.Stop()
	}

	s.writeLastTime = time.Now()
	s.processTimer = time.NewTicker(time.Millisecond * 50)
	ch := s.processTimer.C

	for _, ok := <-ch; ok; {
		if err := s.writeBuffer(false); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 50)
	}
	return nil
}

// Check message value
func (s *StreamSQL) Check(ctx context.Context, msg eventstream.Message) bool {
	return true
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

// writeBuffer all data
func (s *StreamSQL) writeBuffer(flush bool) (err error) {
	if !atomic.CompareAndSwapInt32(&s.isWriting, 0, 1) {
		return err
	}

	var (
		tx   *sql.Tx
		stmt *sql.Stmt
		stop = false
		conn *sql.DB
		now  = time.Now()
	)

	defer func() {
		if rec := recover(); rec != nil {
			s.logError(rec)
			s.logError(string(debug.Stack()))
			if tx != nil {
				tx.Rollback()
			}
		}
		atomic.StoreInt32(&s.isWriting, 0)
	}()

	if !flush {
		if c := len(s.buffer); c < 1 || (s.blockSize > c && now.Sub(s.writeLastTime) < s.flushInterval) {
			return err
		}
	}

	if conn, err = s.connector.Connection(); err != nil {
		return err
	}

	if s.debug {
		log.Println("[stream] write buffer", flush, now.Sub(s.writeLastTime))
	}

	if tx, err = conn.Begin(); err != nil {
		return err
	}

	if stmt, err = tx.Prepare(s.query.QueryString()); err != nil {
		tx.Rollback()
		return err
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
	return err
}

///////////////////////////////////////////////////////////////////////////////
/// Logs
///////////////////////////////////////////////////////////////////////////////

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
