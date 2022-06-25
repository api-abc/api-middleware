package worker

import "context"

type ExecutionFn func(s string) (interface{}, error)

type JobDescriptor struct {
	ID   int
	Type string
}

type Job struct {
	Descriptor JobDescriptor
	ExecFn     ExecutionFn
	Args       string
}

type Result struct {
	Value      interface{}
	Err        error
	Descriptor JobDescriptor
}

func (j Job) execute(ctx context.Context) Result {
	value, err := j.ExecFn(j.Args)
	if err != nil {
		return Result{
			Err:        err,
			Descriptor: j.Descriptor,
		}
	}

	return Result{
		Value:      value,
		Descriptor: j.Descriptor,
	}
}
