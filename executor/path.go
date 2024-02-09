package executor

import (
	"path/filepath"

	"github.com/blocklessnetwork/b7s/models/execute"
)

func (e *Executor) workdirPath(requestID string) string {
	return filepath.Join(e.cfg.WorkDir, "t", requestID)
}
func (e *Executor) setRequestPaths(requestID string, req *execute.Request) {

	workdir := e.workdirPath(requestID)
	req.Config.Runtime.FSRoot = filepath.Join(workdir, "fs")
	req.Config.Runtime.Input = filepath.Join(e.cfg.WorkDir, req.FunctionID, req.Method)
	req.Config.Runtime.DriversRootPath = e.cfg.DriversRootPath
}
