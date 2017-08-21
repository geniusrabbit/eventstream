//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
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
	"github.com/geniusrabbit/eventstream/stream/vertica"
)

func init() {
	storage.RegisterConnector(connector, "vertica")
}

// Vertica storage object
type Vertica struct {
	conn *sql.DB
}

func connector(conf eventstream.ConfigItem, debug bool) (eventstream.Storager, error) {
	var (
		urlObj, err = url.Parse(conf.String("connection", ""))
		conn        *sql.DB
	)
	if err != nil {
		return nil, err
	}

	if conn, err = verticaConnect(urlObj, debug); err != nil {
		return nil, err
	}

	return &Vertica{conn: conn}, nil
}

// Stream vertica processor
func (st *Vertica) Stream(conf eventstream.ConfigItem) (eventstream.Streamer, error) {
	return vertica.New(st, st.conn, conf)
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
		idle           = defs(u.Query().Get("idle"), "30")
		maxcon         = defs(u.Query().Get("maxcon"), "0")
		lifetime       = defs(u.Query().Get("lifetime"), "0")
		password, _    = u.User.Password()
		sslmode        = defs(u.Query().Get("sslmode"), "disable")
		host, port, _  = net.SplitHostPort(u.Host)
		dataSourceName = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
			u.User.Username(), password, host, defs(port, "5432"), u.Path[1:], sslmode)
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

func defs(s, def string) string {
	if len(s) > 0 {
		return s
	}
	return def
}
