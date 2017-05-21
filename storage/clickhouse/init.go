//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package clickhouse

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
	storage.RegisterConnector(clickHouseConnect, "clickhouse", "ch")
}

// clickHouseConnect source
// URL example tcp://login:pass@hostname:port/name?sslmode=disable&idle=10&maxcon=30
func clickHouseConnect(u *url.URL, debug bool) (interface{}, error) {
	var (
		idle          = defs(u.Query().Get("idle"), "30")
		maxcon        = defs(u.Query().Get("maxcon"), "0")
		lifetime      = defs(u.Query().Get("lifetime"), "0")
		host, port, _ = net.SplitHostPort(u.Host)
	)

	if len(port) < 1 {
		port = "9000"
	}

	// Compile connection string
	dataSourceName := fmt.Sprintf("tcp://%s:%s?database=%s", host, port, u.Path[1:])

	// Open connection
	conn, err := sql.Open("clickhouse", dataSourceName)
	if nil == err {
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
