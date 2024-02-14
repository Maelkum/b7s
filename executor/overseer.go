package executor

import (
	"github.com/Maelkum/overseer/job"
)

type Overseer interface {
	Run(job.Job) (job.State, error)
	Start(job.Job) (string, error)
	Wait(id string) (job.State, error)
}
