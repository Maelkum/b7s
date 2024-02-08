package executor

import (
	"fmt"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
)

// ExecuteFunction will run the Blockless function defined by the execution request.
func (e *Executor) ExecuteFunction(requestID string, req execute.Request) (execute.Result, error) {

	// Execute the function.
	out, usage, err := e.executeFunction(requestID, req)
	if err != nil {

		res := execute.Result{
			Code:      codes.Error,
			RequestID: requestID,
			Result:    out,
			Usage:     usage,
		}

		return res, fmt.Errorf("function execution failed: %w", err)
	}

	res := execute.Result{
		Code:      codes.OK,
		RequestID: requestID,
		Result:    out,
		Usage:     usage,
	}

	return res, nil
}

// executeFunction handles the actual execution of the Blockless function. It returns the
// execution information like standard output, standard error, exit code and resource usage.
func (e *Executor) executeFunction(requestID string, req execute.Request) (execute.RuntimeOutput, execute.Usage, error) {

	log := e.log.With().Str("request", requestID).Str("function", req.FunctionID).Logger()

	log.Info().Msg("processing execution request")

	// Generate paths for execution request.
	paths := e.generateRequestPaths(requestID, req.FunctionID, req.Method)

	err := e.cfg.FS.MkdirAll(paths.workdir, defaultPermissions)
	if err != nil {
		return execute.RuntimeOutput{}, execute.Usage{}, fmt.Errorf("could not setup working directory for execution (dir: %s): %w", paths.workdir, err)
	}
	// Remove all temporary files after we're done.
	defer func() {
		err := e.cfg.FS.RemoveAll(paths.workdir)
		if err != nil {
			log.Error().Err(err).Str("dir", paths.workdir).Msg("could not remove request working directory")
		}
	}()

	log.Debug().Str("dir", paths.workdir).Msg("working directory for the request")

	out, usage, err := e.executeWithOverseer(requestID, paths, req)
	if err != nil {
		return out, execute.Usage{}, fmt.Errorf("command execution failed: %w", err)
	}

	log.Info().Msg("command executed successfully")

	return out, usage, nil
}

func (e *Executor) executeWithOverseer(requestID string, paths requestPaths, req execute.Request) (execute.RuntimeOutput, execute.Usage, error) {

	job := e.createJob(paths, req)

	e.log.Debug().Interface("job", job).Str("request", requestID).Msg("job created")
	state, err := e.overseer.Run(job)
	if err != nil {
		// TODO: always return values here
		return execute.RuntimeOutput{}, execute.Usage{}, fmt.Errorf("job run failed: %w", err)
	}

	out := execute.RuntimeOutput{
		Stdout: state.Stdout,
		Stderr: state.Stderr,
	}

	// This should always be the case in the case where we `run` since we've waited for the process.
	if state.ExitCode != nil {
		out.ExitCode = *state.ExitCode
	} else {
		e.log.Warn().Str("request", requestID).Msg("exit code missing for executed process")
	}

	// TODO: (overseer) Collect usage info.
	usage := execute.Usage{}

	return out, usage, nil
}
