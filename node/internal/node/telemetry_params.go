package node

import (
	"fmt"
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
)

// TODO: Descriptions for metrics
