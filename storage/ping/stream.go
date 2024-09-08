package ping

import (
	"context"
	"net/http"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/message"
	"github.com/geniusrabbit/eventstream/internal/patternkey"
)

type pingStreamConfig struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

type pingStream struct {
	id string

	url    *patternkey.PatternKey
	method string

	httpClient http.Client
}

// Check implements stream.Streamer.
func (p *pingStream) Check(ctx context.Context, msg message.Message) bool { return true }

// Close implements stream.Streamer.
func (p *pingStream) Close() error { return nil }

// ID returns unical stream identificator
func (p *pingStream) ID() string { return p.id }

// Put implements stream.Streamer.
func (p *pingStream) Put(ctx context.Context, msg message.Message) error {
	nUrl := p.url.Prepare(msg)
	if p.method == http.MethodGet {
		_, err := p.httpClient.Get(nUrl)
		return err
	}
	_, err := p.httpClient.Post(nUrl, "application/json", strings.NewReader(msg.JSON()))
	return err
}

// Run implements stream.Streamer.
func (p *pingStream) Run(ctx context.Context) error { return nil }

var _ eventstream.Streamer = (*pingStream)(nil)
