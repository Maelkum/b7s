package bls

import (
	"context"

	"github.com/Maelkum/b7s/models/execute"
)

type Executor interface {
	ExecuteFunction(ctx context.Context, requestID string, request execute.Request) (execute.Result, error)
}
