package executor

import (
	"fmt"
	"path/filepath"

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

	// Set paths for execution request.
	e.setRequestPaths(requestID, &req)

	workdir := req.Config.Runtime.Workdir
	err := e.cfg.FS.MkdirAll(workdir, defaultPermissions)
	if err != nil {
		return execute.RuntimeOutput{}, execute.Usage{}, fmt.Errorf("could not setup working directory for execution (dir: %s): %w", workdir, err)
	}
	// Remove all temporary files after we're done.
	defer func() {
		err := e.cfg.FS.RemoveAll(workdir)
		if err != nil {
			log.Error().Err(err).Str("dir", workdir).Msg("could not remove request working directory")
		}
	}()

	log.Debug().Str("dir", workdir).Msg("working directory for the request")

	// TODO: Add support for different runtimes/binaries to execute.
	runtime := filepath.Join(e.cfg.RuntimeDir, e.cfg.ExecutableName)

	var runner Runner
	if e.useEnhancedRunner {
		runner = e.newEnhancedRunner()
	} else {
		runner = e.newSimpleRunner()
	}

	out, usage, err := runner.Run(requestID, runtime, req)
	if err != nil {
		return out, usage, fmt.Errorf("command execution failed: %w", err)
	}

	log.Info().Msg("command executed successfully")

	return out, usage, nil
}
