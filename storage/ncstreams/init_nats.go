// +build nats allstorage all

// Package nats contains ints stream implementation
//
// @project geniusrabbit::eventstream 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019
//
package ncstreams

import (
	"github.com/geniusrabbit/eventstream/storage"
)

func init() {
	storage.RegisterConnector(connector, "nats")
}
