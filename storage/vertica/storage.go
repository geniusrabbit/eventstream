// Package vertica provides a Vertica Storager backed by storage/sql streams.
package vertica

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/geniusrabbit/blaze-api/pkg/context/ctxlogger"
	"github.com/geniusrabbit/eventstream/internal/utils"
	"github.com/geniusrabbit/eventstream/storage"
	sqlstore "github.com/geniusrabbit/eventstream/storage/sql"
	"github.com/geniusrabbit/eventstream/stream"
)

const (
	queryPattern = `COPY {{target}} ({{fields}}) FROM STDIN DELIMITER '\t' NULL 'null'`
)

// Vertica is a Vertica connection that opens SQL COPY streams.
type Vertica struct {
	debug     bool
	connect   string
	initQuery []string
	conn      *sql.DB
}

// Open new vertica storage
func Open(ctx context.Context, connectURL string, options ...Option) (*Vertica, error) {
	var (
		conn        *sql.DB
		urlObj, err = url.Parse(connectURL)
		store       = Vertica{connect: connectURL}
	)
	if err != nil {
		return nil, err
	}
	for _, opt := range options {
		opt(&store)
	}
	if conn, err = verticaConnect(urlObj, store.debug); err != nil {
		return nil, err
	}

	store.conn = conn
	if len(store.initQuery) > 0 {
		ctxlogger.Get(ctx).Debug("init query",
			zap.String(`storage`, `vertica`),
			zap.String(`connect`, connectURL),
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

// Stream vertica processor
func (st *Vertica) Stream(options ...any) (stream.Streamer, error) {
	storeOptions := []sqlstore.Option{sqlstore.WithDebug(st.debug)}
	for _, opt := range options {
		switch o := opt.(type) {
		case sqlstore.Option:
			storeOptions = append(storeOptions, o)
		default:
			return nil, errors.Wrapf(storage.ErrInvalidOption, `vertica stream: %T`, opt)
		}
	}
	if len(storeOptions) == 1 {
		return nil, fmt.Errorf("vertica: stream requires sql options (query)")
	}
	return sqlstore.NewStreamSQL(st, storeOptions...)
}

// Connection to vertica DB
func (st *Vertica) Connection() (_ *sql.DB, err error) {
	// Check current connection
	if st.conn != nil {
		if err = st.conn.Ping(); err != nil {
			_ = st.conn.Close()
			err = nil
			st.conn = nil
		}
	}

	if st.conn == nil {
		urlObj, _ := url.Parse(st.connect)
		st.conn, err = verticaConnect(urlObj, st.debug)
	}
	return st.conn, err
}

// Close vertica connection
func (st *Vertica) Close() error {
	return st.conn.Close()
}

// verticaConnect opens a Vertica database handle.
// URL example: vert://login:password@hostname:5433/name?sslmode=disable&idle=10&maxcon=30
func verticaConnect(u *url.URL, debug bool) (*sql.DB, error) {
	var (
		query          = u.Query()
		idle           = utils.StringOrDefault(query.Get("idle"), "30")
		maxcon         = utils.StringOrDefault(query.Get("maxcon"), "0")
		lifetime       = utils.StringOrDefault(query.Get("lifetime"), "0")
		sslmode        = utils.StringOrDefault(query.Get("sslmode"), "disable")
		password, _    = u.User.Password()
		host, port, _  = net.SplitHostPort(u.Host)
		dataSourceName = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
			u.User.Username(), password, host, utils.StringOrDefault(port, "5432"), u.Path[1:], sslmode)
	)

	// Open connection
	conn, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err == nil {
		if count, _ := strconv.Atoi(idle); count >= 0 {
			conn.SetMaxIdleConns(count)
		}

		if count, _ := strconv.Atoi(maxcon); count >= 0 {
			conn.SetMaxOpenConns(count)
		}

		if lifetime, _ := strconv.Atoi(lifetime); lifetime >= 0 {
			conn.SetConnMaxLifetime(time.Duration(lifetime) * time.Second)
		}
	}
	return conn, err
}
