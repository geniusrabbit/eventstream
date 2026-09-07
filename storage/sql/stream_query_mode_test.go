package sql

import "testing"

func TestQueryUsesNativeInsertBatch(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		query string
		batch bool
	}{
		{
			name:  "insert values",
			query: "INSERT INTO t (a, b) VALUES (?, ?)",
			batch: true,
		},
		{
			name:  "insert select join",
			query: "INSERT INTO t (a, b, c) SELECT ?, x, y FROM o WHERE id = ?",
			batch: false,
		},
		{
			name:  "insert select no values keyword",
			query: "INSERT INTO dsp_stats.events_local (event_time) SELECT ? FROM dsp_stats.clicks_by_id",
			batch: false,
		},
		{
			name:  "plain select",
			query: "SELECT 1",
			batch: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := queryUsesNativeInsertBatch(tc.query); got != tc.batch {
				t.Fatalf("queryUsesNativeInsertBatch() = %v, want %v", got, tc.batch)
			}
		})
	}
}
