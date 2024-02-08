package executor

import (
	"github.com/Maelkum/overseer/job"
)

type Overseer interface {
	Run(job.Job) (job.State, error)
}
