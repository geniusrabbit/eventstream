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
	"sync/atomic"
	"time"

	"go.uber.org/zap"

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
	isWriting int32

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

	// Time ticker to pereodic data flush
	processTimer *time.Ticker

	// Logger object of log writing
	logger *zap.Logger
}

// NewStreamSQL creates streamer object for SQL based stream integration
func NewStreamSQL(id string, connector Connector, options ...Option) (eventstream.Streamer, error) {
	var opts Options
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, err
		}
	}
	if opts.QueryBuilder == nil {
		return nil, errInvalidQueryObject
	}
	return &StreamSQL{
		debug:         opts.Debug,
		id:            id,
		connector:     connector,
		blockSize:     opts.getBlockSize(),
		flushInterval: opts.getFlushInterval(),
		buffer:        make(chan eventstream.Message, opts.getBlockSize()*2),
		query:         opts.QueryBuilder,
		logger:        opts.getLogger(),
	}, nil
}

// ID returns unical stream identificator
func (s *StreamSQL) ID() string {
	return s.id
}

// Put message to stream
func (s *StreamSQL) Put(ctx context.Context, msg eventstream.Message) error {
	if s.debug {
		s.logger.Debug(`put-message`, zap.Any(`message`, msg))
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
			s.logger.Error(`write-buffer`, zap.Error(err))
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
func (s *StreamSQL) writeBuffer(flush bool) error {
	if !atomic.CompareAndSwapInt32(&s.isWriting, 0, 1) {
		return nil
	}

	var (
		err      error
		tx       *sql.Tx
		stmt     *sql.Stmt
		stop     = false
		conn     *sql.DB
		now      = time.Now()
		interval = now.Sub(s.writeLastTime)
	)

	defer func() {
		if rec := recover(); rec != nil {
			s.logger.Error(`write-buffer`, zap.Any(`error`, rec))
			if tx != nil {
				tx.Rollback()
			}
		}
		atomic.StoreInt32(&s.isWriting, 0)
	}()

	if !flush {
		if c := len(s.buffer); c < 1 || (s.blockSize > c && interval < s.flushInterval) {
			return err
		}
	}
	if conn, err = s.connector.Connection(); err != nil {
		return err
	}

	if s.debug {
		s.logger.Debug(`write-buffer`,
			zap.Bool(`hardflush`, flush),
			zap.Duration(`interval`, interval),
			zap.String(`query`, s.query.QueryString()),
		)
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
				s.logger.Debug(`write-message`, zap.Any(`message`, msg))
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
