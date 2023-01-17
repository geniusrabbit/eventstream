package clickhouse

import (
	sqlstore "github.com/geniusrabbit/eventstream/storage/sql"
)

// Option of the clickhouse storage
type Option func(store *Clickhouse)

// WithQueryByTarget SQL storage option
func WithQueryByTarget(target string, fields any) sqlstore.Option {
	return sqlstore.WithQueryByPattern(queryPattern, target, fields)
}

// WithInitQuery which will be executed after connection
func WithInitQuery(initQuery []string) Option {
	return func(store *Clickhouse) {
		store.initQuery = initQuery
	}
}
