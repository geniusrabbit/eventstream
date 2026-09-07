// Package clickhouse provides a ClickHouse Storager backed by storage/sql streams.
package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/geniusrabbit/blaze-api/pkg/context/ctxlogger"
	"github.com/geniusrabbit/eventstream/storage"
	sqlstore "github.com/geniusrabbit/eventstream/storage/sql"
	"github.com/geniusrabbit/eventstream/stream"
)

const (
	queryPattern = `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`
)

// Clickhouse is a ClickHouse connection that opens SQL insert streams.
type Clickhouse struct {
	mu sync.Mutex

	debug   bool
	connect string
	conn    *sql.DB

	initQuery []string
}

// Open new clickhouse storage stream
func Open(ctx context.Context, connect string, options ...Option) (*Clickhouse, error) {
	var (
		err   error
		conn  *sql.DB
		store = Clickhouse{connect: connect}
	)
	for _, opt := range options {
		opt(&store)
	}
	if conn, err = sql.Open("clickhouse", connect); err != nil {
		return nil, err
	}
	store.conn = conn
	if len(store.initQuery) > 0 {
		ctxlogger.Get(ctx).Debug("init query",
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
func (c *Clickhouse) Stream(options ...any) (stream.Streamer, error) {
	storeOptions := []sqlstore.Option{sqlstore.WithDebug(c.debug)}
	for _, opt := range options {
		switch o := opt.(type) {
		case sqlstore.Option:
			storeOptions = append(storeOptions, o)
		default:
			return nil, errors.Wrapf(storage.ErrInvalidOption, `clickhouse stream: %T`, opt)
		}
	}
	if len(storeOptions) == 1 {
		return nil, fmt.Errorf("clickhouse: stream requires sql options (query)")
	}
	return sqlstore.NewStreamSQL(c, storeOptions...)
}

// Connection to clickhouse DB.
// Safe for concurrent use from multiple SQL streams sharing one storage handle.
func (c *Clickhouse) Connection() (_ *sql.DB, err error) {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()

	if conn != nil {
		if err = conn.Ping(); err == nil {
			return conn, nil
		}
		c.mu.Lock()
		if c.conn == conn {
			c.conn = nil
		}
		c.mu.Unlock()
		_ = conn.Close()
	}

	newConn, err := sql.Open("clickhouse", c.connect)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		if err = c.conn.Ping(); err == nil {
			_ = newConn.Close()
			return c.conn, nil
		}
		old := c.conn
		c.conn = nil
		_ = old.Close()
	}
	c.conn = newConn
	return c.conn, nil
}

// Close clickhouse connection
func (c *Clickhouse) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}
