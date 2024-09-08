package metrics

import (
	"context"
	"errors"
	"strings"

	"github.com/geniusrabbit/eventstream/internal/message"
	"github.com/geniusrabbit/eventstream/internal/patternkey"
)

// ErrUnsupportedMetricType value
var ErrUnsupportedMetricType = errors.New(`unsupported metric type`)

// MetricType of the metric
type MetricType string

// Metric types list
const (
	MetricCounter MetricType = "counter"
)

// IsCounter type of the metric
func (mt MetricType) IsCounter() bool {
	return mt == MetricCounter
}

// Metricer executor
type Metricer interface {
	Exec(ctx context.Context, msg message.Message)
}

// MetricerList wrapper
type MetricerList []Metricer

// Execute all metrics frin the list
func (l MetricerList) Exec(ctx context.Context, msg message.Message) {
	for _, m := range l {
		m.Exec(ctx, msg)
	}
}

// Metric config type
type Metric struct {
	Namespace string              `json:"namespace,omitempty"`
	Subsystem string              `json:"subsystem,omitempty"`
	Name      string              `json:"name"`
	Type      MetricType          `json:"type"`
	Tags      []map[string]string `json:"tags,omitempty"`
}

// Labels list of tags
func (m *Metric) Labels() []string {
	if len(m.Tags) == 0 || len(m.Tags[0]) == 0 {
		return nil
	}
	labels := make([]string, 0, len(m.Tags[0]))
	for k := range m.Tags[0] {
		labels = append(labels, k)
	}
	return labels
}

func (m *Metric) tagValues() *patternkey.PatterKeys {
	if len(m.Tags) == 0 || len(m.Tags[0]) == 0 {
		return nil
	}
	vals := make([]string, 0, len(m.Tags[0]))
	for _, v := range m.Tags[0] {
		vals = append(vals, v)
	}
	return patternkey.PatternKeysFrom(vals...)
}

func (m *Metric) prometheusName() string {
	return strings.ReplaceAll(m.Name, ".", "_")
}

// Metric returns metric processor
func (m *Metric) Metric() (Metricer, error) {
	switch m.Type {
	case MetricCounter:
		return counterFromMetrics(m), nil
	default:
		return nil, ErrUnsupportedMetricType
	}
}

// MetricList extender
type MetricList []*Metric

// Metric executer of the
func (l MetricList) Metric() (Metricer, error) {
	if len(l) == 0 {
		return nil, nil
	}
	mList := make(MetricerList, 0, len(l))
	for _, m := range l {
		mt, err := m.Metric()
		if err != nil {
			return nil, err
		}
		mList = append(mList, mt)
	}
	return mList, nil
}
