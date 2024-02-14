package executor

type Limiter interface {
	AssignToGroup(name string, pid uint64) error
}
