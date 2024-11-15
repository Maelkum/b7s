package head

import (
	"github.com/armon/go-metrics/prometheus"
)

// Tracing span names.
const (
	// notifiee events
	spanPeerConnected    = "PeerConnected"
	spanPeerDisconnected = "PeerDisconnected"

	// TODO: Rename this to just execute.
	// execution events
	spanHeadExecute = "HeadExecute"
)

var (
	rollCallsPublishedMetric = []string{"node", "rollcalls", "published"}
	functionExecutionsMetric = []string{"node", "function", "executions"}
)

var Counters = []prometheus.CounterDefinition{
	{
		Name: rollCallsPublishedMetric,
		Help: "Number of roll calls this node issued.",
	},
	{
		Name: functionExecutionsMetric,
		Help: "Number of function executions.",
	},
}
