package executor

import (
	"github.com/rs/zerolog"

	"github.com/blocklessnetwork/b7s/models/execute"
)

type simpleRunner struct {
	log        zerolog.Logger
	useLimiter bool
	limiter    Limiter
}

func (e *Executor) newSimpleRunner() *simpleRunner {

	sr := simpleRunner{
		log:        e.log,
		useLimiter: e.cfg.useLimiter,
		limiter:    e.cfg.Limiter,
	}
	return &sr
}

func (s *simpleRunner) Run(requestID string, runtime string, req execute.Request) (execute.RuntimeOutput, execute.Usage, error) {
	cmd := createCmd(runtime, req)
	return s.executeCommand(cmd)
}
