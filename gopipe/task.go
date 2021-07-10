package gopipe

import (
	"context"
	"errors"
)

type Task struct {
	Name   string
	Action func(ctx context.Context) error
}

func (task *Task) Run(ctx context.Context) error {
	if task.Action == nil {
		return errors.New("Task.Action is required")
	}
	return task.Action(ctx)
}
