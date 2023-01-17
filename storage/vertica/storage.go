//
// @project geniusrabbit::eventstream 2017, 2019-2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019-2020
//

package vertica

import (
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/internal/utils"
	sqlstore "github.com/geniusrabbit/eventstream/storage/sql"
	"github.com/geniusrabbit/eventstream/stream"
)

const (
	queryPattern = `COPY {{target}} ({{fields}}) FROM STDIN DELIMITER '\t' NULL 'null'`
)

// Vertica storage object
type Vertica struct {
	debug   bool
	connect string
	conn    *sql.DB
}

// Open new vertica storage
func Open(connectURL string, options ...sqlstore.Option) (*Vertica, error) {
	var (
		conn        *sql.DB
		opts        sqlstore.Options
		urlObj, err = url.Parse(connectURL)
	)
	if err != nil {
		return nil, err
	}
	for _, o := range options {
		if err = o(&opts); err != nil {
			return nil, err
		}
	}
	if conn, err = verticaConnect(urlObj, opts.Debug); err != nil {
		return nil, err
	}
	return &Vertica{conn: conn, connect: connectURL, debug: opts.Debug}, nil
}

// Stream vertica processor
func (st *Vertica) Stream(options ...any) (strm eventstream.Streamer, err error) {
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
	if strm, err = sqlstore.New(st, queryPattern, &conf, storeOptions...); err != nil {
		return nil, err
	}
	return eventstream.NewStreamWrapper(strm, conf.Where, metricExec)
}

// Connection to clickhouse DB
func (st *Vertica) Connection() (_ *sql.DB, err error) {
	// Check current connection
	if st.conn != nil {
		if err = st.conn.Ping(); err != nil {
			st.conn.Close()
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

// verticaConnect source
// @param u URL vert://login:password@hostname:5433/name?sslmode=disable&idle=10&maxcon=30
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
