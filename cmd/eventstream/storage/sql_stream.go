//go:build clickhouse || vertica || allsource || all
// +build clickhouse vertica allsource all

package storage

import (
	"fmt"
	"time"

	"github.com/geniusrabbit/eventstream/storage/sql"
)

type withQueryFn func(opts ...sql.QueryOption) sql.Option

func sqlStreamOpts(streamName string, p StreamParams, withQuery withQueryFn) ([]any, error) {
	opts := []any{
		sql.WithID(streamName),
		sql.WithDebug(p.Debug),
	}
	if p.BufferSize > 0 {
		opts = append(opts, sql.WithBlockSize(int(p.BufferSize)))
	}
	if p.WriteTimeout > 0 {
		opts = append(opts, sql.WithFlushIntervals(time.Duration(p.WriteTimeout)*time.Millisecond))
	}
	if p.SQLQuery != "" {
		qopts := []sql.QueryOption{}
		if p.IterateBy != "" {
			qopts = append(qopts, sql.QWithIterateBy(p.IterateBy))
		}
		if p.Fields != nil {
			qopts = append(qopts, sql.QWithMessageTmpl(p.Fields))
		}
		opts = append(opts, sql.WithQuery(p.SQLQuery, qopts...))
		return opts, nil
	}
	if p.Target == "" && p.Fields == nil {
		return nil, fmt.Errorf("stream %q: sql_query or target/fields required", streamName)
	}
	qopts := []sql.QueryOption{}
	if p.Target != "" {
		qopts = append(qopts, sql.QWithTarget(p.Target))
	}
	if p.IterateBy != "" {
		qopts = append(qopts, sql.QWithIterateBy(p.IterateBy))
	}
	if p.Fields != nil {
		qopts = append(qopts, sql.QWithMessageTmpl(p.Fields))
	}
	opts = append(opts, withQuery(qopts...))
	return opts, nil
}
