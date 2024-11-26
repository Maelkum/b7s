package node

import (
	"github.com/armon/go-metrics/prometheus"
)

var (
	subscriptionsMetric = []string{"node", "topic", "subscriptions"}
)

var Counters = []prometheus.CounterDefinition{
	{
		Name: subscriptionsMetric,
		Help: "Number of topics this node subscribes to.",
	},
}
