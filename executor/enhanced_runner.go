package executor

import (
	"github.com/rs/zerolog"

	"github.com/blocklessnetwork/b7s/models/execute"
)

// enhancedRunenr uses an overseer to run the execution request.
type enhancedRunner struct {
	log      zerolog.Logger
	overseer Overseer
}

func (e *Executor) newEnhancedRunner() *enhancedRunner {

	if e.cfg.Overseer == nil {
		panic("cannot create enhanced runner without an overseer instance")
	}

	er := enhancedRunner{
		log:      e.log,
		overseer: e.cfg.Overseer,
	}

	return &er
}

func (e *enhancedRunner) Run(requestID string, runtime string, req execute.Request) (execute.RuntimeOutput, execute.Usage, error) {

	job := createJob(runtime, req)
	state, err := e.overseer.Run(job)
	if err != nil {
		e.log.Error().Err(err).Msg("job run failed")
		// NOTE: not returning here + preserving the execution error.
	}

	out := execute.RuntimeOutput{
		Stdout: state.Stdout,
		Stderr: state.Stderr,
	}

	// These are essentially the same types.
	ru := state.ResourceUsage
	usage := execute.Usage{
		WallClockTime: ru.WallClockTime,
		CPUUserTime:   ru.CPUUserTime,
		CPUSysTime:    ru.CPUSysTime,
		MemoryMaxKB:   ru.MemoryMaxKB,
	}

	// This should always be the case in the case where we `run` since we've waited for the process.
	if state.ExitCode != nil {
		out.ExitCode = *state.ExitCode
	} else {
		e.log.Warn().Str("request", requestID).Msg("exit code missing for executed process")
	}

	return out, usage, nil
}
