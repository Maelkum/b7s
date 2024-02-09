package executor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Maelkum/overseer/job"

	"github.com/blocklessnetwork/b7s/models/execute"
)

// createJob will translate the execution request to a job specification.
func (e *Executor) createJob(workdir string, req execute.Request) job.Job {

	// Prepare CLI arguments.
	// Append the input argument first.
	var args []string
	args = append(args, req.Config.Runtime.Input)

	// Append the arguments for the runtime.
	runtimeFlags := runtimeFlags(req.Config.Runtime, req.Config.Permissions)
	args = append(args, runtimeFlags...)

	// Separate runtime arguments from the function arguments.
	args = append(args, "--")

	// Function arguments.
	for _, param := range req.Parameters {
		if param.Value != "" {
			args = append(args, param.Value)
		}
	}

	// Setup stdin of the command.
	var stdin io.Reader
	if req.Config.Stdin != nil {
		stdin = strings.NewReader(*req.Config.Stdin)
	}

	// Setup environment.
	// First, pass through our environment variables.
	environ := os.Environ()

	// Second, set the variables set in the execution request.
	names := make([]string, 0, len(req.Config.Environment))
	for _, env := range req.Config.Environment {
		e := fmt.Sprintf("%s=%s", env.Name, env.Value)
		environ = append(environ, e)

		names = append(names, env.Name)
	}

	// Third and final - set the `BLS_LIST_VARS` variable with
	// the list of names of the variables from the execution request.
	blsList := strings.Join(names, ";")
	blsEnv := fmt.Sprintf("%s=%s", blsListEnvName, blsList)
	environ = append(environ, blsEnv)

	job := job.Job{
		Exec: job.Command{
			WorkDir: workdir,
			// TODO: Add support for different runtimes/binaries to execute.
			Path: filepath.Join(e.cfg.RuntimeDir, e.cfg.ExecutableName),
			Args: args,
			Env:  environ,
		},
		Stdin: stdin,
	}

	return job
}
