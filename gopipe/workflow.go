package gopipe

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Workflow struct {
	Name  string
	Tasks []*Task
}

func (workflow *Workflow) Validate() error {
	if workflow.Name == "" {
		return errors.New("Workflow.Name is required")
	}
	names := make(map[string]struct{}, len(workflow.Tasks))
	for _, task := range workflow.Tasks {
		if _, ok := names[task.Name]; ok {
			return fmt.Errorf("task name %s is duplicated", task.Name)
		}
		names[task.Name] = struct{}{}
		if err := task.Validate(); err != nil {
			return fmt.Errorf("task %s is invalid: %w", task.Name, err)
		}
	}
	return nil
}

func (workflow *Workflow) Run(ctx context.Context, args *Args) error {
	log.Print("===> workflow " + workflow.Name + " starts")
	startTime := time.Now()
	defer func() {
		log.Printf("===> %v: workflow %s ended", time.Since(startTime), workflow.Name)
	}()
	if err := workflow.Validate(); err != nil {
		return err
	}
	argsMap := map[string]*Args{
		"_": args,
	}
	for _, task := range workflow.Tasks {
		t := time.Now()
		log.Print("===> task " + task.Name + " starts")
		args := &Args{}
		if task.Args != nil {
			a, err := task.Args(argsMap)
			if err != nil {
				return err
			}
			args = a
		}
		if task.If != nil {
			b, err := task.If(ctx, args)
			if err != nil {
				return err
			}
			if !b {
				log.Printf("===> task %s is skipped", task.Name)
				continue
			}
		}
		err := task.Run(ctx, args)
		duration := time.Since(t)
		argsMap[task.Name] = args
		if err != nil {
			return fmt.Errorf("%v: task %s is failed: %w", duration, task.Name, err)
		}
		log.Printf("===> %v: task %s ended", duration, task.Name)
	}
	return nil
}
