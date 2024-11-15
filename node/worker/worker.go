package worker

import (
	"fmt"
	"slices"
	"sync"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/node/internal/waitmap"
	"github.com/blocklessnetwork/b7s/node/node"
)

type Worker struct {
	node.Core

	cfg Config

	executor blockless.Executor
	fstore   FStore

	sema chan struct{}
	wg   *sync.WaitGroup
	//	subgroups  workSubgroups
	attributes *attributes.Attestation

	// TODO: Update cluster map, don't use this.
	// clusters maps request ID to the cluster the node belongs to.
	clusters map[string]consensusExecutor

	// TODO: This no longer needs to be a waitmap for the worker.
	executeResponses *waitmap.WaitMap[string, execute.ResultMap]

	// clusterLock is used to synchronize access to the `clusters` map.
	clusterLock sync.RWMutex
}

func New(node node.Core, fstore FStore, executor blockless.Executor, options ...Option) (*Worker, error) {

	// Initialize config.
	cfg := DefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	err := cfg.Valid()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// TODO: Do this in main.
	// Ensure default topic is included in the topic list.
	defaultSubscription := slices.Contains(cfg.Topics, blockless.DefaultTopic)
	if !defaultSubscription {
		cfg.Topics = append(cfg.Topics, blockless.DefaultTopic)
	}

	// TODO: Tracing

	worker := &Worker{
		Core:     node,
		fstore:   fstore,
		executor: executor,

		clusters: make(map[string]consensusExecutor),
	}

	return worker, nil
}
