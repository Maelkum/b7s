package node

import (
	"fmt"

	"github.com/armon/go-metrics/prometheus"
)

// Tracing span names.
const (
	// message events
	spanMessageSend    = "MessageSend"
	spanMessagePublish = "MessagePublish"
	spanMessageProcess = "MessageProcess"

	// notifiee events
	spanPeerConnected    = "PeerConnected"
	spanPeerDisconnected = "PeerDisconnected"
)

func msgProcessSpanName(msgType string) string {
	return fmt.Sprintf("%s %s", spanMessageProcess, msgType)
}

func msgSendSpanName(prefix string, msgType string) string {
	return fmt.Sprintf("%s %s", prefix, msgType)
}

var (
	messagesSentMetric      = []string{"node", "messages", "sent"}
	messagesPublishedMetric = []string{"node", "messages", "published"}
	subscriptionsMetric     = []string{"node", "topic", "subscriptions"}
	directMessagesMetric    = []string{"node", "direct", "messages"}
	topicMessagesMetric     = []string{"node", "topic", "messages"}

	nodeInfoMetric = []string{"node", "info"}
)

// TODO: Descriptions for metrics

var Counters = []prometheus.CounterDefinition{
	{
		Name: directMessagesMetric,
		Help: "Number of direct messages this node received.",
	},
	{
		Name: topicMessagesMetric,
		Help: "Number of topic messages this node received.",
	},
	{
		Name: messagesSentMetric,
		Help: "Number of messages sent.",
	},
	{
		Name: messagesPublishedMetric,
		Help: "Number of messages published.",
	},
}

var Gauges = []prometheus.GaugeDefinition{
	{
		Name: nodeInfoMetric,
		Help: "Information about the b7s node.",
	},
}
