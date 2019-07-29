//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
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
	"github.com/geniusrabbit/eventstream/storage"
	sqlstore "github.com/geniusrabbit/eventstream/storage/sql"
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

func connector(conf *storage.Config) (eventstream.Storager, error) {
	var (
		urlObj, err = url.Parse(conf.Connect)
		conn        *sql.DB
	)

	if err != nil {
		return nil, err
	}

	if conn, err = verticaConnect(urlObj, conf.Debug); err != nil {
		return nil, err
	}

	return &Vertica{conn: conn, connect: conf.Connect, debug: conf.Debug}, nil
}

// Stream vertica processor
func (st *Vertica) Stream(conf interface{}) (strm eventstream.Streamer, err error) {
	var confObj = conf.(*storage.StreamConfig)
	if strm, err = sqlstore.New(st, confObj, queryPattern); err != nil {
		return nil, err
	}
	return eventstream.NewStreamWrapper(strm, confObj.Where)
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

///////////////////////////////////////////////////////////////////////////////
/// Connection helpers
///////////////////////////////////////////////////////////////////////////////

// verticaConnect source
// @param u URL vert://login:password@hostname:5433/name?sslmode=disable&idle=10&maxcon=30
func verticaConnect(u *url.URL, debug bool) (*sql.DB, error) {
	var (
		query          = u.Query()
		idle           = defString(query.Get("idle"), "30")
		maxcon         = defString(query.Get("maxcon"), "0")
		lifetime       = defString(query.Get("lifetime"), "0")
		sslmode        = defString(query.Get("sslmode"), "disable")
		password, _    = u.User.Password()
		host, port, _  = net.SplitHostPort(u.Host)
		dataSourceName = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
			u.User.Username(), password, host, defString(port, "5432"), u.Path[1:], sslmode)
	)

	// Open connection
	conn, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	conn.Ping()
	{
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
	return conn, nil
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func defString(s, def string) string {
	if len(s) > 0 {
		return s
	}
	return def
}
