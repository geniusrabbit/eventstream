//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package storage

import (
	"fmt"
	"net/url"
)

type connector func(u *url.URL, debug bool) (interface{}, error)

var (
	connections = map[string]interface{}{}
	connectors  = map[string]connector{}
)

// RegisterConnector function
func RegisterConnector(c connector, scheme string, schemes ...string) {
	connectors[scheme] = c
	for _, sch := range schemes {
		connections[sch] = c
	}
}

// Register connection
func Register(name, connect string, debug bool) (err error) {
	var conn interface{}
	if conn, err = connection(connect, debug); nil == err {
		connections[name] = conn
	}
	return
}

// Get connection
func Get(name string) interface{} {
	conn, _ := connections[name]
	return conn
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func connection(connect string, debug bool) (interface{}, error) {
	var url, err = url.Parse(connect)
	if nil != err {
		return nil, err
	}

	if conn, _ := connectors[url.Scheme]; nil != conn {
		return conn(url, debug)
	}

	return nil, fmt.Errorf("Undefined stream scheme: %s", url.Scheme)
}
