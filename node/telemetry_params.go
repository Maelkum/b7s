package node

import (
	"github.com/armon/go-metrics/prometheus"
)

// Tracing span names.
const (
	// notifiee events
	spanPeerConnected    = "PeerConnected"
	spanPeerDisconnected = "PeerDisconnected"
)

// Tracing span status messages.
const (
	spanStatusOK  = "message processed ok"
	spanStatusErr = "error processing message"
)

var (
	messagesProcessedMetric    = []string{"node", "messages", "processed"}
	messagesProcessedOkMetric  = []string{"node", "messages", "processed", "ok"}
	messagesProcessedErrMetric = []string{"node", "messages", "processed", "err"}
	subscriptionsMetric        = []string{"node", "topic", "subscriptions"}
)

var Counters = []prometheus.CounterDefinition{
	{
		Name: messagesProcessedMetric,
		Help: "Number of messages this node processed.",
	},
	{
		Name: messagesProcessedOkMetric,
		Help: "Number of messages successfully processed by the node.",
	},
	{
		Name: messagesProcessedErrMetric,
		Help: "Number of messages processed with an error.",
	},
	{
		Name: subscriptionsMetric,
		Help: "Number of topics this node subscribes to.",
	},
}
