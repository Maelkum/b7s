package head

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/node/internal/node"
)

// processMessage will determine which message was received and how to process it.
func (h *HeadNode) processMessage(ctx context.Context, from peer.ID, payload []byte, pipeline node.Pipeline) (procError error) {

	// Determine message type.
	msgType, err := node.GetMessageType(payload)
	if err != nil {
		return fmt.Errorf("could not unpack message: %w", err)
	}

	// n.metrics.IncrCounterWithLabels(messagesProcessedMetric, 1, []metrics.Label{{Name: "type", Value: msgType}})
	// defer func() {
	// 	switch procError {
	// 	case nil:
	// 		n.metrics.IncrCounterWithLabels(messagesProcessedOkMetric, 1, []metrics.Label{{Name: "type", Value: msgType}})
	// 	default:
	// 		n.metrics.IncrCounterWithLabels(messagesProcessedErrMetric, 1, []metrics.Label{{Name: "type", Value: msgType}})
	// 	}
	// }()

	// ctx, err = tracing.TraceContextFromMessage(ctx, payload)
	// if err != nil {
	// 	n.log.Error().Err(err).Msg("could not get trace context from message")
	// }

	// ctx, span := n.tracer.Start(ctx, msgProcessSpanName(msgType), msgProcessSpanOpts(from, msgType, pipeline)...)
	// defer span.End()
	// // NOTE: This function checks the named return error value in order to set the span status accordingly.
	// defer func() {
	// 	if procError == nil {
	// 		span.SetStatus(otelcodes.Ok, spanStatusOK)
	// 		return
	// 	}

	// 	if allowErrorLeakToTelemetry {
	// 		span.SetStatus(otelcodes.Error, procError.Error())
	// 		return
	// 	}

	// 	span.SetStatus(otelcodes.Error, spanStatusErr)
	// }()

	log := h.Log().With().Stringer("peer", from).Str("type", msgType).Stringer("pipeline", pipeline).Logger()

	if !node.CorrectPipeline(msgType, pipeline) {
		log.Warn().Msg("message not allowed on pipeline")
		return nil
	}

	log.Debug().Msg("received message from peer")

	switch msgType {

	// TODO: Consider function install.
	// case blockless.MessageInstallFunction:
	// 	return handleMessage(ctx, from, payload, n.processInstallFunction)

	case blockless.MessageInstallFunctionResponse:
		return node.HandleMessage(ctx, from, payload, h.processInstallFunctionResponse)
	case blockless.MessageRollCallResponse:
		return node.HandleMessage(ctx, from, payload, h.processRollCallResponse)
	case blockless.MessageExecute:
		return node.HandleMessage(ctx, from, payload, h.processExecute)
	case blockless.MessageExecuteResponse:
		return node.HandleMessage(ctx, from, payload, h.processExecuteResponse)
	case blockless.MessageFormClusterResponse:
		return node.HandleMessage(ctx, from, payload, h.processFormClusterResponse)

	default:
		return fmt.Errorf("unsupported message: %s", msgType)
	}
}
