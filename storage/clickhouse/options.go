package clickhouse

import (
	sqlstore "github.com/geniusrabbit/eventstream/storage/sql"
)

// Option of the clickhouse storage
type Option func(store *Clickhouse)

// WithDebug enables storage debug mode
func WithDebug(debug bool) Option {
	return func(store *Clickhouse) {
		store.debug = debug
	}
}

// WithQuery SQL storage option
func WithQuery(opts ...sqlstore.QueryOption) sqlstore.Option {
	return sqlstore.WithQuery(queryPattern, opts...)
}

// WithInitQuery which will be executed after connection
func WithInitQuery(initQuery []string) Option {
	return func(store *Clickhouse) {
		store.initQuery = initQuery
	}
}
