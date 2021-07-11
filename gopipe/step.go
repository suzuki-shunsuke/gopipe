package gopipe

import (
	"context"
	"errors"
)

type Step struct {
	Name   string
	Action func(ctx context.Context, args *Args) error
	If     StepIf
}

type StepIf func(ctx context.Context, args *Args) (bool, error)

func (step *Step) Validate() error {
	if step.Name == "" {
		return errors.New("Step.Name is required")
	}
	if step.Action == nil {
		return errors.New("Step.Action is required")
	}
	return nil
}

func (step *Step) Run(ctx context.Context, args *Args) error {
	if err := step.Validate(); err != nil {
		return err
	}
	return step.Action(ctx, args)
}
