// +build nats allstorage all

// Package nats contains ints stream implementation
//
// @project geniusrabbit::eventstream 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019
//
package nats

import (
	"errors"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/notificationcenter"
)

// Errors set
var (
	ErrInvalidNATSScheme = errors.New(`[storage::nats] invalid NATS scheme`)
)

func init() {
	storage.RegisterConnector(connector, "nats")
}

func connector(conf *storage.Config) (_ eventstream.Storager, err error) {
	if !strings.HasPrefix(conf.Connect, "nats://") &&
		!strings.HasPrefix(conf.Connect, "natstream://") {
		return nil, ErrInvalidNATSScheme
	}
	var stream notificationcenter.Streamer
	if stream, err = connectNATS(conf.Connect); err != nil {
		return nil, err
	}
	return &NATS{stream: stream, debug: conf.Debug}, nil
}
