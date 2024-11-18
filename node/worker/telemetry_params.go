package worker

import (
	"github.com/armon/go-metrics/prometheus"
)

// Tracing span names.
const (
	spanWorkOrder = "WorkOrder"
)

// TODO: CHeck - duplicate span/metric definitions.

var (
	rollCallsSeenMetric      = []string{"node", "rollcalls", "seen"}
	rollCallsAppliedMetric   = []string{"node", "rollcalls", "applied"}
	functionExecutionsMetric = []string{"node", "function", "executions"}
)

var Counters = []prometheus.CounterDefinition{

	{
		Name: rollCallsSeenMetric,
		Help: "Number of roll calls seen by the node.",
	},
	{
		Name: rollCallsAppliedMetric,
		Help: "Number of roll calls this node applied to.",
	},
	{
		Name: functionExecutionsMetric,
		Help: "Number of function executions.",
	},
}
