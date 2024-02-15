package response

import (
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/execute"
)

type ExecControl struct {
	Type      string            `json:"type,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	From      peer.ID           `json:"from,omitempty"`
	Action    string            `json:"action,omitempty"`
	Results   execute.ResultMap `json:"results,omitempty"`

	Message string `json:"message,omitempty"` // Used to communicate the reason for failure to the user.
}
