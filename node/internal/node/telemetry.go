package node

import (
	"context"
	"strings"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

func saveTraceContext(ctx context.Context, msg blockless.Message) {
	tmsg, ok := msg.(blockless.TraceableMessage)
	if !ok {
		return
	}

	t := tracing.GetTraceInfo(ctx)
	if !t.Empty() {
		tmsg.SaveTraceContext(t)
	}
}

type MessageSpanConfig struct {
	msgPipeline Pipeline
	receivers   []peer.ID
}

func (c *MessageSpanConfig) Pipeline(p Pipeline) *MessageSpanConfig {
	c.msgPipeline = p
	return c
}

func (c *MessageSpanConfig) Peer(id peer.ID) *MessageSpanConfig {
	if c.receivers == nil {
		c.receivers = make([]peer.ID, 0, 1)
	}

	c.receivers = append(c.receivers, id)
	return c
}

func (c *MessageSpanConfig) Peers(ids ...peer.ID) *MessageSpanConfig {
	if c.receivers == nil {
		c.receivers = make([]peer.ID, 0, len(ids))
	}

	c.receivers = append(c.receivers, ids...)
	return c
}

func (c *MessageSpanConfig) SpanOpts() []trace.SpanStartOption {

	attrs := []attribute.KeyValue{
		b7ssemconv.MessagePipeline.String(c.msgPipeline.ID.String()),
	}

	if c.msgPipeline.ID == PubSub {
		attrs = append(attrs, b7ssemconv.MessageTopic.String(c.msgPipeline.Topic))
	}

	if len(c.receivers) == 1 {
		attrs = append(attrs, b7ssemconv.MessagePeer.String(c.receivers[0].String()))
	} else if len(c.receivers) > 1 {
		attrs = append(attrs, b7ssemconv.MessagePeers.String(
			strings.Join(blockless.PeerIDsToStr(c.receivers), ","),
		))
	}

	return []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(attrs...),
	}
}
