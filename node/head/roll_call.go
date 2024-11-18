package head

import (
	"context"
	"fmt"
	"time"

	"github.com/armon/go-metrics"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/consensus"
	"github.com/blocklessnetwork/b7s/consensus/pbft"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
)

func (h *HeadNode) executeRollCall(
	ctx context.Context,
	requestID string,
	functionID string,
	nodeCount int,
	consensusAlgo consensus.Type,
	topic string,
	attributes *execute.Attributes,
	timeout int,
) ([]peer.ID, error) {

	// Create a logger with relevant context.
	log := h.Log().With().
		Str("request", requestID).
		Str("function", functionID).
		Int("node_count", nodeCount).
		Str("topic", topic).
		Logger()

	log.Info().Msg("performing roll call for request")

	h.rollCall.create(requestID)
	defer h.rollCall.remove(requestID)

	err := h.publishRollCall(ctx, requestID, functionID, consensusAlgo, topic, attributes)
	if err != nil {
		return nil, fmt.Errorf("could not publish roll call: %w", err)
	}

	log.Info().Msg("roll call published")

	// Limit for how long we wait for responses.
	t := h.cfg.RollCallTimeout
	if timeout > 0 {
		t = time.Duration(timeout) * time.Second
	}

	tctx, exCancel := context.WithTimeout(ctx, t)
	defer exCancel()

	// Peers that have reported on roll call.
	var reportingPeers []peer.ID
rollCallResponseLoop:
	for {
		// Wait for responses from nodes who want to work on the request.
		select {
		// Request timed out.
		case <-tctx.Done():

			// -1 means we'll take any peers reporting
			if len(reportingPeers) >= 1 && nodeCount == -1 {
				log.Info().Msg("enough peers reported for roll call")
				break rollCallResponseLoop
			}

			log.Warn().Msg("roll call timed out")
			return nil, blockless.ErrRollCallTimeout

		case reply := <-h.rollCall.responses(requestID):

			// Check if this is the reply we want - shouldn't really happen.
			if reply.FunctionID != functionID {
				log.Info().Str("peer", reply.From.String()).Str("function_got", reply.FunctionID).Msg("skipping inadequate roll call response - wrong function")
				continue
			}

			// Check if we are connected to this peer.
			// Since we receive responses to roll call via direct messages - should not happen.
			if !h.Connected(reply.From) {
				h.Log().Info().Str("peer", reply.From.String()).Msg("skipping roll call response from unconnected peer")
				continue
			}

			log.Info().Str("peer", reply.From.String()).Msg("roll called peer chosen for execution")

			reportingPeers = append(reportingPeers, reply.From)

			// -1 means we'll take any peers reporting
			if len(reportingPeers) >= nodeCount && nodeCount != -1 {
				log.Info().Msg("enough peers reported for roll call")
				break rollCallResponseLoop
			}
		}
	}

	if consensusAlgo == consensus.PBFT && len(reportingPeers) < pbft.MinimumReplicaCount {
		return nil, fmt.Errorf("not enough peers reported for PBFT consensus (have: %v, need: %v)", len(reportingPeers), pbft.MinimumReplicaCount)
	}

	return reportingPeers, nil
}

// publishRollCall will create a roll call request for executing the given function.
// On successful issuance of the roll call request, we return the ID of the issued request.
func (h *HeadNode) publishRollCall(ctx context.Context, requestID string, functionID string, consensus consensus.Type, topic string, attributes *execute.Attributes) error {

	h.Metrics().IncrCounterWithLabels(rollCallsPublishedMetric, 1, []metrics.Label{{Name: "function", Value: functionID}})

	// Create a roll call request.
	rollCall := request.RollCall{
		Origin:     h.Host().ID(),
		FunctionID: functionID,
		RequestID:  requestID,
		Consensus:  consensus,
		Attributes: attributes,
	}

	if topic == "" {
		topic = blockless.DefaultTopic
	}

	// Publish the mssage.
	err := h.PublishToTopic(ctx, topic, &rollCall)
	if err != nil {
		return fmt.Errorf("could not publish to topic: %w", err)
	}

	return nil
}

func (h *HeadNode) processRollCallResponse(ctx context.Context, from peer.ID, res response.RollCall) error {

	log := h.Log().With().
		Str("request", res.RequestID).
		Stringer("peer", from).Logger()

	log.Debug().Msg("processing peer's roll call response")

	// Check if the response is adequate.
	if res.Code != codes.Accepted {
		log.Info().Stringer("code", res.Code).Msg("skipping inadequate roll call response - unwanted code")
		return nil
	}

	// Check if there's an active roll call already.
	exists := h.rollCall.exists(res.RequestID)
	if !exists {
		log.Info().Msg("no pending roll call for the given request, dropping response")
		return nil
	}

	log.Info().Msg("recording roll call response")

	rres := rollCallResponse{
		From:     from,
		RollCall: res,
	}

	// Record the response.
	h.rollCall.add(res.RequestID, rres)

	return nil
}
