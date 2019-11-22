package clickhouse

import (
	sqlstore "github.com/geniusrabbit/eventstream/storage/sql"
)

// WithQueryByTarget SQL storage option
func WithQueryByTarget(target string, fields interface{}) sqlstore.Option {
	return sqlstore.WithQueryByPattern(queryPattern, target, fields)
}
