package gopipe

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type Task struct {
	Name  string
	If    TaskIf
	Steps []*Step
	Args  GetArgs
}

type TaskIf func(ctx context.Context, args *Args) (bool, error)

func (task *Task) Validate() error {
	if task.Name == "" {
		return errors.New("Task.Name is required")
	}
	for i, step := range task.Steps {
		if err := step.Validate(); err != nil {
			return fmt.Errorf("step %s (%d) is invalid: %w", step.Name, i, err)
		}
	}
	return nil
}

func (task *Task) Run(ctx context.Context, args *Args) error {
	if err := task.Validate(); err != nil {
		return err
	}
	for _, step := range task.Steps {
		if step.If != nil {
			b, err := step.If(ctx, args)
			if err != nil {
				return err
			}
			if !b {
				log.Printf("===> step %s is skipped", step.Name)
				continue
			}
		}
		if err := step.Run(ctx, args); err != nil {
			return err
		}
	}
	return nil
}
