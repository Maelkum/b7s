package api

import (
	"context"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/request"
)

type Node interface {
	ExecuteFunction(ctx context.Context, req execute.Request, subgroup string) (code codes.Code, requestID string, results execute.ResultMap, peers execute.Cluster, err error)
	ExecutionResult(id string) (execute.Result, bool)
	PublishFunctionInstall(ctx context.Context, uri string, cid string, subgroup string) error

	// TODO: This should return some information about the job - either stats or output.
	// Right now let's just return if we succeeded or not.
	ExecutionControl(id string, action request.ExecAction) error
}
