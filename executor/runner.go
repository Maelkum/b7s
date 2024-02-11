package executor

import (
	"github.com/blocklessnetwork/b7s/models/execute"
)

// Runner is be used to actually run/execute the request.
type Runner interface {
	Run(id string, runtime string, req execute.Request) (execute.RuntimeOutput, execute.Usage, error)
}
