//
// @project geniusrabbit::eventstream 2017 - 2023
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2023
//

package clickhouse

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/internal/zlogger"
	"github.com/geniusrabbit/eventstream/storage"
	sqlstore "github.com/geniusrabbit/eventstream/storage/sql"
	"github.com/geniusrabbit/eventstream/stream"
)

const (
	queryPattern = `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`
)

// Clickhouse storage object
type Clickhouse struct {
	debug   bool
	connect string
	conn    *sql.DB

	initQuery []string
}

// Open new clickhouse storage stream
func Open(ctx context.Context, connect string, options ...any) (*Clickhouse, error) {
	var (
		err    error
		conn   *sql.DB
		config storage.Config
		store  = Clickhouse{connect: connect}
	)
	for _, opt := range options {
		switch o := opt.(type) {
		case storage.Option:
			o(&config)
		case Option:
			o(&store)
		default:
			return nil, errors.Wrapf(storage.ErrInvalidOption, `%+v`, opt)
		}
	}
	if conn, err = sql.Open("clickhouse", connect); err != nil {
		return nil, err
	}
	store.debug = config.Debug
	store.conn = conn
	if len(store.initQuery) > 0 {
		zlogger.FromContext(ctx).Debug("init query",
			zap.String(`storage`, `clickhouse`),
			zap.String(`connect`, connect),
			zap.String(`init_query`, strings.Join(store.initQuery, "\n")),
		)
		tx, err := conn.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		for _, sqlQuery := range store.initQuery {
			if _, err = tx.ExecContext(ctx, sqlQuery); err != nil {
				_ = tx.Rollback()
				return nil, err
			}
		}
		if err = tx.Commit(); err != nil {
			return nil, err
		}
	}
	return &store, nil
}

// Stream clickhouse processor
func (c *Clickhouse) Stream(options ...any) (strm eventstream.Streamer, err error) {
	var (
		conf         stream.Config
		storeOptions []sqlstore.Option
		metricExec   metrics.Metricer
	)
	for _, opt := range options {
		switch o := opt.(type) {
		case stream.Option:
			o(&conf)
		case sqlstore.Option:
			storeOptions = append(storeOptions, o)
		case *stream.Config:
			conf = *o
		default:
			stream.WithObjectConfig(o)(&conf)
		}
	}
	if metricExec, err = conf.Metrics.Metric(); err != nil {
		return nil, err
	}
	if strm, err = sqlstore.New(c, queryPattern, &conf, storeOptions...); err != nil {
		return nil, err
	}
	cond, err := conf.Condition()
	if err != nil {
		return nil, err
	}
	return eventstream.NewStreamWrapper(strm, cond, metricExec), nil
}

// Connection to clickhouse DB
func (c *Clickhouse) Connection() (_ *sql.DB, err error) {
	// Check current connection
	if c.conn != nil {
		if err = c.conn.Ping(); err != nil {
			c.conn.Close()
			err = nil
			c.conn = nil
		}
	}
	if c.conn == nil {
		c.conn, err = sql.Open("clickhouse", c.connect)
	}
	return c.conn, err
}

// Close clickhouse connection
func (c *Clickhouse) Close() error {
	return c.conn.Close()
}
