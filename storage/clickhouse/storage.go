//
// @project geniusrabbit::eventstream 2017 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2020
//

package clickhouse

import (
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/internal/utils"
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
}

// Open new clickhouse storage stream
func Open(connect string, options ...interface{}) (*Clickhouse, error) {
	var (
		urlObj, err = url.Parse(connect)
		conn        *sql.DB
		config      storage.Config
	)
	if err != nil {
		return nil, err
	}
	for _, opt := range options {
		switch o := opt.(type) {
		case storage.Option:
			o(&config)
		default:
			return nil, errors.Wrapf(storage.ErrInvalidOption, `%+v`, opt)
		}
	}
	if conn, err = clickHouseConnect(urlObj, config.Debug); err != nil {
		return nil, err
	}
	return &Clickhouse{
		debug:   config.Debug,
		connect: connect,
		conn:    conn,
	}, nil
}

// Stream clickhouse processor
func (c *Clickhouse) Stream(options ...interface{}) (strm eventstream.Streamer, err error) {
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
	return eventstream.NewStreamWrapper(strm, conf.Where, metricExec)
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
		urlObj, _ := url.Parse(c.connect)
		c.conn, err = clickHouseConnect(urlObj, c.debug)
	}
	return c.conn, err
}

// Close clickhouse connection
func (c *Clickhouse) Close() error {
	return c.conn.Close()
}

// clickHouseConnect source
// URL example tcp://login:pass@hostname:port/name?sslmode=disable&idle=10&maxcon=30
func clickHouseConnect(u *url.URL, debug bool) (*sql.DB, error) {
	var (
		query         = u.Query()
		idle          = utils.StringOrDefault(query.Get("idle"), "30")
		maxcon        = utils.StringOrDefault(query.Get("maxcon"), "0")
		lifetime      = utils.StringOrDefault(query.Get("lifetime"), "0")
		host, port, _ = net.SplitHostPort(u.Host)
		dataSource    = fmt.Sprintf("tcp://%s:%s?database=%s", host, utils.StringOrDefault(port, "9000"), u.Path[1:])
	)

	// Open connection
	conn, err := sql.Open("clickhouse", dataSource)
	if err != nil {
		return nil, err
	}
	if count, _ := strconv.Atoi(idle); count >= 0 {
		conn.SetMaxIdleConns(count)
	}
	if count, _ := strconv.Atoi(maxcon); count >= 0 {
		conn.SetMaxOpenConns(count)
	}
	if lifetime, _ := strconv.Atoi(lifetime); lifetime >= 0 {
		conn.SetConnMaxLifetime(time.Duration(lifetime) * time.Second)
	}
	return conn, nil
}
