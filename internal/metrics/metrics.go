package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsCost = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "demo-go",
		Name:      "requests-total",
		Help:      "all requests cost",
		Buckets:   prometheus.ExponentialBuckets(0.001, 5, 7),
	}, []string{"method", "path"})
)

func init() {
	prometheus.MustRegister(RequestsCost)
}
