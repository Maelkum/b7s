package worker

import (
	"context"

	"github.com/blocklessnetwork/b7s/models/response"
	"github.com/libp2p/go-libp2p/core/peer"
)

func (w *Worker) processHealthCheck(ctx context.Context, from peer.ID, _ response.Health) error {
	w.Log().Trace().Stringer("from", from).Msg("peer health check received")
	return nil
}
