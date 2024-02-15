package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
)

func (n *Node) processExecControl(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the message.
	var req request.ExecControl
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack execute response: %w", err)
	}
	req.From = from

	log := n.log.With().Str("request", req.RequestID).Str("action", req.Action.String()).Logger()
	log.Info().Msg("received execution control message")

	// TODO: Remove limitation on allowed actions.
	if req.Action != request.ExecWait {
		return fmt.Errorf("unsupported action: %s", req.Action)
	}

	if n.isWorker() {
		exres, err := n.executor.ExecutionWait(req.RequestID)
		if err != nil {
			return fmt.Errorf("could not execute action: %w", err)
		}

		msg := response.ExecControl{
			Type:      blockless.MessageExecControlResponse,
			RequestID: req.RequestID,
			Action:    req.Action.String(),
			Results: execute.ResultMap{
				n.host.ID(): exres,
			},
		}

		err = n.send(ctx, req.From, msg)
		if err != nil {
			return fmt.Errorf("could not send exec control response: %w", err)
		}

		return nil
	}

	return n.headExecControl(req.RequestID, req.Action)
}

func (n *Node) headExecControl(requestID string, action request.ExecAction) error {

	n.detachedExecutionsLock.RLock()
	peers, ok := n.detachedExecutions[requestID]
	n.detachedExecutionsLock.RUnlock()

	if !ok || len(peers) == 0 {
		return errors.New("no known peers for detached execution")
	}

	req := request.ExecControl{
		Type:      blockless.MessageExecControl,
		RequestID: requestID,
		Action:    action,
	}

	err := n.sendToMany(context.Background(), peers, req)
	if err != nil {
		return fmt.Errorf("could not send execution control message to peers: %w", err)
	}

	return nil
}

func (n *Node) processExecControlResponse(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the message.
	var res response.ExecControl
	err := json.Unmarshal(payload, &res)
	if err != nil {
		return fmt.Errorf("could not unpack execute control response: %w", err)
	}
	res.From = from

	log := n.log.With().Interface("response", res).Str("request", res.RequestID).Str("action", res.Action).Logger()
	log.Info().Msg("received execution control response message")

	return nil
}
