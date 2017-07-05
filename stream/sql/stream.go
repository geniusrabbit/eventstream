//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

// TODO: Store all messages if not sended yet

package eventstream

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/demdxx/gocast"
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
	when             *govaluate.EvaluableExpression
	query            stream.Query
	processTimer     *time.Ticker
}

// NewStreamSQL streamer
func NewStreamSQL(
	conn *sql.DB,
	blockSize int,
	duration time.Duration,
	whenCondition string,
	query stream.Query,
) (_ stream.Streamer, err error) {
	if blockSize < 1 {
		blockSize = 1000
	}

	if duration <= 0 {
		duration = time.Second * 1
	}

	var (
		when *govaluate.EvaluableExpression
	)

	if len(strings.TrimSpace(whenCondition)) > 0 {
		if when, err = govaluate.NewEvaluableExpression(whenCondition); nil != err {
			return
		}
	}

	return &StreamSQL{
		conn:             conn,
		buffer:           make(chan eventstream.Message, blockSize*2),
		blockSize:        blockSize,
		writeMaxDuration: duration,
		when:             when,
		query:            query,
	}, nil
}

// NewStreamSQLByRaw query
func NewStreamSQLByRaw(
	conn *sql.DB,
	blockSize int,
	duration time.Duration,
	when,
	query string,
	fields interface{},
) (stream.Streamer, error) {
	q, err := stream.NewQueryByRaw(query, fields)
	if nil != err {
		return nil, err
	}
	return NewStreamSQL(conn, blockSize, duration, when, *q)
}

// Check message value
func (s *StreamSQL) Check(msg eventstream.Message) bool {
	if nil != s.when {
		r, err := s.when.Evaluate(msg.Map())
		return err == nil && gocast.ToBool(r)
	}
	return true
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
			s.LogError(err)
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

///////////////////////////////////////////////////////////////////////////////
/// Logs
///////////////////////////////////////////////////////////////////////////////

// LogError message
func (s *StreamSQL) LogError(params ...interface{}) {
	if len(params) > 0 {
		log.Error(params...)
	}
}
