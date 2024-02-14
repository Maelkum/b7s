package blockless

import (
	"github.com/blocklessnetwork/b7s/models/execute"
)

type Executor interface {
	ExecuteFunction(requestID string, request execute.Request) (execute.Result, error)
	SupportsLongRunningJobs() bool

	// Detached execution.
	ExecutionStart(requestID string, request execute.Request) error
	ExecutionWait(requestID string) (execute.Result, error)
	ExecutionStop(requestID string) (execute.Result, error)
	ExecutionStats(requestID string) (execute.Result, error)
}
