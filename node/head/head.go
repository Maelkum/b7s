package head

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/response"
	"github.com/blocklessnetwork/b7s/node/internal/node"
	"github.com/blocklessnetwork/b7s/node/internal/waitmap"
)

type HeadNode struct {
	node.Core

	rollCall           *rollCallQueue
	consensusResponses *waitmap.WaitMap[string, response.FormCluster]
	executeResponses   *waitmap.WaitMap[string, execute.NodeResult]

	cfg Config
}

func New(node node.Core, options ...Option) (*HeadNode, error) {

	// Initialize config.
	cfg := DefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	// TODO: Make sure default topic is included in main (also worker).

	err := cfg.Valid()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// TODO: Tracing.

	head := &HeadNode{
		Core: node,
		cfg:  cfg,

		rollCall:           newQueue(rollCallQueueBufferSize),
		consensusResponses: waitmap.New[string, response.FormCluster](0),
		executeResponses:   waitmap.New[string, execute.NodeResult](executionResultCacheSize),
	}

	return head, nil
}

func newRequestID() string {
	return uuid.New().String()
}
