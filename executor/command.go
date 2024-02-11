package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Maelkum/overseer/job"

	"github.com/blocklessnetwork/b7s/models/execute"
)

// createJob will translate the execution request to a job specification.
func createJob(runtime string, req execute.Request) job.Job {

	// Setup stdin of the command.
	var stdin io.Reader
	if req.Config.Stdin != nil {
		stdin = strings.NewReader(*req.Config.Stdin)
	}

	job := job.Job{
		Exec: job.Command{
			WorkDir: req.Config.Runtime.Workdir,
			Path:    runtime,
			Args:    createArgs(req),
			Env:     createEnv(req),
		},
		Stdin: stdin,
	}

	return job
}

// createCmd will translate the execution request to a Cmd struct.
func createCmd(runtime string, req execute.Request) *exec.Cmd {

	// Setup stdin of the command.
	var stdin io.Reader
	if req.Config.Stdin != nil {
		stdin = strings.NewReader(*req.Config.Stdin)
	}

	args := createArgs(req)

	cmd := exec.Command(runtime, args...)
	cmd.Dir = req.Config.Runtime.Workdir
	cmd.Env = createEnv(req)
	cmd.Stdin = stdin

	return cmd
}

func createArgs(req execute.Request) []string {

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

	return args
}

func createEnv(req execute.Request) []string {

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

	return environ
}
