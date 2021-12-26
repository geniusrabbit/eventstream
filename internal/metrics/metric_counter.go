package metrics

import (
	"context"

	"github.com/geniusrabbit/eventstream/internal/message"
	"github.com/geniusrabbit/eventstream/internal/patternkey"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metricCounter struct {
	counter *prometheus.CounterVec
	tags    *patternkey.PatterKeys
}

func counterFromMetrics(met *Metric) Metricer {
	return &metricCounter{
		counter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: met.Namespace,
			Subsystem: met.Subsystem,
			Name:      met.prometheusName(),
		}, met.Labels()),
		tags: met.tagValues(),
	}
}

func (met *metricCounter) Exec(ctx context.Context, msg message.Message) {
	met.counter.WithLabelValues(met.tags.Prepare(msg)...).Inc()
}
