package worker

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/node/internal/pipeline"
	"github.com/blocklessnetwork/b7s/node/node"
)

// TODO: Perhaps create a map: message ID => handler

// processMessage will determine which message was received and how to process it.
func (w *Worker) processMessage(ctx context.Context, from peer.ID, payload []byte, pipeline pipeline.Pipeline) (procError error) {

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

	log := w.Log.With().Str("peer", from.String()).Str("type", msgType).Str("pipeline", pipeline.String()).Logger()

	if !node.CorrectPipeline(msgType, pipeline) {
		log.Warn().Msg("message not allowed on pipeline")
		return nil
	}

	log.Debug().Msg("received message from peer")

	switch msgType {
	case blockless.MessageHealthCheck:
		return node.HandleMessage(ctx, from, payload, w.processHealthCheck)
	case blockless.MessageInstallFunction:
		return node.HandleMessage(ctx, from, payload, w.processInstallFunction)
	case blockless.MessageRollCall:
		return node.HandleMessage(ctx, from, payload, w.processRollCall)
	case blockless.MessageExecute:
		return node.HandleMessage(ctx, from, payload, w.processExecute)
	case blockless.MessageFormCluster:
		return node.HandleMessage(ctx, from, payload, w.processFormCluster)
	case blockless.MessageDisbandCluster:
		return node.HandleMessage(ctx, from, payload, w.processDisbandCluster)

	default:
		return fmt.Errorf("unsupported message: %s", msgType)
	}
}
