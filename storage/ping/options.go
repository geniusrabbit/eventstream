package ping

// Option configures ping storage open.
type Option func(*pinger)

// WithURL sets default ping URL (overrides connect when set).
func WithURL(url string) Option {
	return func(p *pinger) { p.URL = url }
}

// WithMethod sets default HTTP method.
func WithMethod(method string) Option {
	return func(p *pinger) { p.Method = method }
}

// StreamOption configures a ping sink stream.
type StreamOption func(*streamOptions)

type streamOptions struct {
	id          string
	url         string
	method      string
	contentType string
}

// WithStreamID sets the stream identifier.
func WithStreamID(id string) StreamOption {
	return func(o *streamOptions) { o.id = id }
}

// WithStreamURL sets per-stream URL template.
func WithStreamURL(url string) StreamOption {
	return func(o *streamOptions) { o.url = url }
}

// WithStreamMethod sets per-stream HTTP method.
func WithStreamMethod(method string) StreamOption {
	return func(o *streamOptions) { o.method = method }
}

// WithContentType sets request content type.
func WithContentType(ct string) StreamOption {
	return func(o *streamOptions) { o.contentType = ct }
}
