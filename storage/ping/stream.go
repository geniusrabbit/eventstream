package ping

import (
	"context"
	"net/http"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/message"
	"github.com/geniusrabbit/eventstream/internal/patternkey"
	"github.com/geniusrabbit/eventstream/internal/zlogger"
	"go.uber.org/zap"
)

type pingStreamConfig struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	ContentType string `json:"content_type"`
}

type pingStream struct {
	id string

	url         *patternkey.PatternKey
	method      string
	contentType string

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
	if strings.HasPrefix(nUrl, "//") {
		nUrl = "http:" + nUrl
	}
	if p.method == http.MethodGet {
		_, err := p.httpClient.Get(nUrl)
		return err
	}
	resp, err := p.httpClient.Post(nUrl, p.contentType, strings.NewReader(msg.JSON()))
	if err == nil {
		_ = resp.Body.Close()
		zlogger.FromContext(ctx).Info("ping",
			zap.String("method", p.method),
			zap.String("url", nUrl),
			zap.Int("status", resp.StatusCode))
	} else {
		zlogger.FromContext(ctx).Error("ping",
			zap.String("method", p.method),
			zap.String("url", nUrl),
			zap.Error(err))
	}
	return err
}

// Run implements stream.Streamer.
func (p *pingStream) Run(ctx context.Context) error { return nil }

var _ eventstream.Streamer = (*pingStream)(nil)
