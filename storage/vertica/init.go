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

	"github.com/geniusrabbit/eventstream/storage"
)

func init() {
	storage.RegisterConnector(verticaConnect, "vertica", "vert", "vt")
}

// verticaConnect source
// @param u URL vert://login:password@hostname:5433/name?sslmode=disable&idle=10&maxcon=30
func verticaConnect(u *url.URL, debug bool) (interface{}, error) {
	var (
		idle          = defs(u.Query().Get("idle"), "30")
		maxcon        = defs(u.Query().Get("maxcon"), "0")
		lifetime      = defs(u.Query().Get("lifetime"), "0")
		password, _   = u.User.Password()
		sslmode       = defs(u.Query().Get("sslmode"), "disable")
		host, port, _ = net.SplitHostPort(u.Host)
	)

	if len(port) < 1 {
		port = "5433"
	}

	// Compile connection string
	dataSourceName := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		u.User.Username(), password, host, port, u.Path[1:], sslmode)

	// Open connection
	conn, err := sql.Open("postgres", dataSourceName)
	if nil == err {
		conn.Ping()

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
	return nil, err
}

func defs(s, def string) string {
	if len(s) > 0 {
		return s
	}
	return def
}
