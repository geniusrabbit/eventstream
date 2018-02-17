//
// @project geniusrabbit::eventstream 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

package clickhouse

import (
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/errors"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/stream/clickhouse"
)

func init() {
	storage.RegisterConnector(connector, "clickhouse")
}

// Clickhouse storage object
type Clickhouse struct {
	debug   bool
	connect string
	conn    *sql.DB
}

func connector(conf eventstream.ConfigItem, debug bool) (eventstream.Storager, error) {
	var (
		connect     = conf.String("connect", "")
		urlObj, err = url.Parse(connect)
		conn        *sql.DB
	)

	if connect == "" {
		return nil, errors.ErrConnectionIsNotDefined
	}

	if err != nil {
		return nil, err
	}

	if conn, err = clickHouseConnect(urlObj, debug); err != nil {
		return nil, err
	}

	return &Clickhouse{conn: conn, connect: connect, debug: debug}, nil
}

// Close vertica connection
func (c *Clickhouse) Close() error {
	return c.conn.Close()
}

// Stream clickhouse processor
func (c *Clickhouse) Stream(conf eventstream.ConfigItem) (_ eventstream.Streamer, err error) {
	var (
		simple eventstream.SimpleStreamer
	)

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

	if err == nil {
		simple, err = clickhouse.New(c, c.conn, conf, c.debug)
	}

	if err != nil {
		return nil, err
	}
	return eventstream.NewStreamWrapper(simple, conf.String("where", ""))
}

///////////////////////////////////////////////////////////////////////////////
/// Connection helpers
///////////////////////////////////////////////////////////////////////////////

// clickHouseConnect source
// URL example tcp://login:pass@hostname:port/name?sslmode=disable&idle=10&maxcon=30
func clickHouseConnect(u *url.URL, debug bool) (*sql.DB, error) {
	var (
		idle           = defs(u.Query().Get("idle"), "30")
		maxcon         = defs(u.Query().Get("maxcon"), "0")
		lifetime       = defs(u.Query().Get("lifetime"), "0")
		host, port, _  = net.SplitHostPort(u.Host)
		dataSourceName = fmt.Sprintf("tcp://%s:%s?database=%s", host, defs(port, "9000"), u.Path[1:])
	)

	// Open connection
	conn, err := sql.Open("clickhouse", dataSourceName)
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

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func defs(s, def string) string {
	if len(s) > 0 {
		return s
	}
	return def
}
