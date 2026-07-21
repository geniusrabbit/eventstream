package scenario

import (
	"fmt"
	"strings"
	"time"

	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/internal/utils"
)

func itemString(item Item, key string) string {
	v, _ := item[key].(string)
	return utils.PrepareValue(v)
}

func itemBool(item Item, key string) bool {
	switch v := item[key].(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(v, "true") || v == "1"
	default:
		return false
	}
}

func itemUint(item Item, key string) uint {
	switch v := item[key].(type) {
	case float64:
		if v > 0 {
			return uint(v)
		}
	case int:
		if v > 0 {
			return uint(v)
		}
	case int64:
		if v > 0 {
			return uint(v)
		}
	}
	return 0
}

func itemStringSlice(item Item, key string) []string {
	raw, ok := item[key]
	if !ok || raw == nil {
		return nil
	}
	switch v := raw.(type) {
	case []string:
		return v
	case []any:
		out := make([]string, 0, len(v))
		for _, x := range v {
			if s, ok := x.(string); ok {
				out = append(out, s)
			}
		}
		return out
	case string:
		return []string{v}
	default:
		return nil
	}
}

func itemDuration(item Item, key string) time.Duration {
	switch v := item[key].(type) {
	case string:
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	case float64:
		return time.Duration(v)
	case int:
		return time.Duration(v)
	case int64:
		return time.Duration(v)
	}
	return 0
}

func resolveDriver(item Item) string {
	driver := itemString(item, "driver")
	if driver != "" {
		return driver
	}
	connect := itemString(item, "connect")
	if i := strings.Index(connect, "://"); i > 0 {
		return connect[:i]
	}
	return ""
}

func itemMetrics(item Item) (metrics.MetricList, error) {
	raw, ok := item["metrics"]
	if !ok || raw == nil {
		return nil, nil
	}
	// Reuse JSON round-trip via Item.Decode into a tiny wrapper.
	var wrap struct {
		Metrics metrics.MetricList `json:"metrics"`
	}
	if err := (Item{"metrics": raw}).Decode(&wrap); err != nil {
		return nil, fmt.Errorf("metrics: %w", err)
	}
	return wrap.Metrics, nil
}
