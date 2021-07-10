package gopipe

import (
	"context"
	"errors"
)

type Task struct {
	Name   string
	Action func(ctx context.Context) error
}

func (task *Task) Validate() error {
	if task.Action == nil {
		return errors.New("Task.Action is required")
	}
	return nil
}

func (task *Task) Run(ctx context.Context) error {
	if err := task.Validate(); err != nil {
		return err
	}
	return task.Action(ctx)
}
