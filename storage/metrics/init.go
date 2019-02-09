// +build metrics allsource all
//
// @project geniusrabbit::eventstream 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019
//

package metrics

import (
	"errors"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/notificationcenter"
)

// Errors set
var (
	ErrUndefinedMetricsEngine = errors.New(`[storage::metrics] undefined metrics engine or wrong "connect"`)
)

func init() {
	storage.RegisterConnector(connector, "metrics")
}

func connector(conf *storage.Config) (_ eventstream.Storager, err error) {
	var (
		logger notificationcenter.Logger
	)
	switch {
	case strings.HasPrefix(conf.Connect, "nats://"):
		logger, err = connectNATS(conf.Connect)
	case strings.HasPrefix(conf.Connect, "statsd://"):
		logger, err = connectStatsD(conf.Connect)
	default:
		return nil, ErrUndefinedMetricsEngine
	}

	if err != nil {
		return nil, err
	}

	return &Metrics{metrica: logger, debug: conf.Debug}, nil
}
