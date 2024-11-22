package node

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	pp "github.com/blocklessnetwork/b7s/node/internal/pipeline"
	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
)

const (
	tracerName = "b7s.Node"
)

func msgProcessSpanOpts(from peer.ID, msgType string, pipeline pp.Pipeline) []trace.SpanStartOption {

	attrs := []attribute.KeyValue{
		b7ssemconv.MessagePeer.String(from.String()),
		b7ssemconv.MessageType.String(msgType),
		b7ssemconv.MessagePipeline.String(pipeline.ID.String()),
	}

	if pipeline.ID == pp.PubSub {
		attrs = append(attrs, b7ssemconv.MessageTopic.String(pipeline.Topic))
	}

	return []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(attrs...),
	}
}
