package executor

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
)

func (e *Executor) ExecutionStart(requestID string, req execute.Request) error {

	if !e.SupportsLongRunningJobs() {
		return errors.New("executor doesn't support long running jobs")
	}

	log := e.log.With().Str("request", requestID).Str("function", req.FunctionID).Logger()

	log.Info().Msg("processing execution request")

	// Set paths for execution request.
	e.setRequestPaths(requestID, &req)

	workdir := req.Config.Runtime.Workdir
	err := e.cfg.FS.MkdirAll(workdir, defaultPermissions)
	if err != nil {
		return fmt.Errorf("could not setup working directory for execution (dir: %s): %w", workdir, err)
	}

	runtime := filepath.Join(e.cfg.RuntimeDir, e.cfg.ExecutableName)
	job := createJob(runtime, req)

	id, err := e.cfg.Overseer.Start(job)
	if err != nil {
		return fmt.Errorf("could not start execution: %w", err)
	}

	e.jobs.Lock()
	defer e.jobs.Unlock()

	e.jobs.jobs[requestID] = id

	return nil
}

func (e *Executor) ExecutionWait(requestID string) (execute.Result, error) {

	if !e.SupportsLongRunningJobs() {
		return execute.Result{}, errors.New("executor doesn't support long running jobs")
	}

	e.jobs.Lock()
	defer e.jobs.Unlock()

	jobID, ok := e.jobs.jobs[requestID]
	if !ok {
		return execute.Result{}, errors.New("no execution found")
	}

	state, err := e.cfg.Overseer.Wait(jobID)
	if err != nil {
		return execute.Result{}, fmt.Errorf("could not wait on job: %w", err)
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

	res := execute.Result{
		Code:      codes.OK,
		RequestID: requestID,
		Result:    out,
		Usage:     usage,
	}

	return res, nil
}

func (e *Executor) ExecutionStop(requestID string) (execute.Result, error) {
	return execute.Result{}, errors.New("TBD: not supported")
}

func (e *Executor) ExecutionStats(requestID string) (execute.Result, error) {
	return execute.Result{}, errors.New("TBD: not supported")
}
