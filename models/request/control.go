package request

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

type ExecControl struct {
	Type      string     `json:"type,omitempty"`
	From      peer.ID    `json:"from,omitempty"`
	RequestID string     `json:"request_id,omitempty"`
	Action    ExecAction `json:"action,omitempty"`
}

type ExecAction uint

const (
	ExecStat ExecAction = iota + 1
	ExecWait
	ExecKill
)
