package ping

import (
	"context"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/message"
	"github.com/geniusrabbit/eventstream/internal/patternkey"
	"github.com/geniusrabbit/eventstream/internal/zlogger"
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
	var (
		err      error
		resp     *http.Response
		respData []byte
		nUrl     = p.url.Prepare(msg)
	)
	if strings.HasPrefix(nUrl, "//") {
		nUrl = "http:" + nUrl
	}
	if p.method == http.MethodGet {
		resp, err = p.httpClient.Get(nUrl)
	} else {
		resp, err = p.httpClient.Post(nUrl, p.contentType, strings.NewReader(msg.JSON()))
	}
	if resp != nil && resp.Body != nil && err == nil {
		respData, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
	}
	if err == nil {
		zlogger.FromContext(ctx).Info("ping",
			zap.String("url", nUrl),
			zap.String("method", p.method),
			zap.Int("status", resp.StatusCode),
			zap.String("response", string(respData)))
	} else {
		zlogger.FromContext(ctx).Error("ping",
			zap.String("url", nUrl),
			zap.String("method", p.method),
			zap.Error(err))
	}
	return err
}

// Run implements stream.Streamer.
func (p *pingStream) Run(ctx context.Context) error { return nil }

var _ eventstream.Streamer = (*pingStream)(nil)
