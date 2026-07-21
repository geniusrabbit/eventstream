// Package ping provides an HTTP “ping” Storager that POSTs messages to a URL.
package ping

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/demdxx/gocast/v2"
	"github.com/pkg/errors"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/patternkey"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/stream"
)

type pinger struct {
	URL    string
	Method string
}

// Open creates a ping storage.
func Open(_ context.Context, connect string, options ...Option) (eventstream.Storager, error) {
	p := &pinger{
		URL:    connect,
		Method: "GET",
	}
	for _, opt := range options {
		opt(p)
	}
	return p, nil
}

// Close implements eventstream.Storager.
func (p *pinger) Close() error {
	return nil
}

// Stream implements eventstream.Storager.
func (p *pinger) Stream(options ...any) (stream.Streamer, error) {
	var conf streamOptions
	for _, opt := range options {
		switch o := opt.(type) {
		case StreamOption:
			o(&conf)
		default:
			return nil, errors.Wrapf(storage.ErrInvalidOption, `ping stream: %T`, opt)
		}
	}
	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	if httpTransport.TLSClientConfig == nil {
		httpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		httpTransport.TLSClientConfig.InsecureSkipVerify = true
	}
	return &pingStream{
		id:          conf.id,
		url:         patternkey.PatternKeyFromTemplate(gocast.Or(conf.url, p.URL)),
		method:      gocast.Or(conf.method, p.Method),
		contentType: gocast.Or(conf.contentType, "application/json"),
		httpClient:  http.Client{Transport: httpTransport},
	}, nil
}

var _ eventstream.Storager = (*pinger)(nil)
