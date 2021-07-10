package gopipe

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Workflow struct {
	Name  string
	Tasks []Task
}

func (workflow *Workflow) Run(ctx context.Context) error {
	log.Print("===> workflow " + workflow.Name + " starts")
	for _, task := range workflow.Tasks {
		t := time.Now()
		log.Print("===> task " + task.Name + " starts")
		err := task.Run(ctx)
		duration := time.Since(t)
		if err != nil {
			return fmt.Errorf("%v: task %s is failed: %w", duration, task.Name, err)
		}
		log.Printf("===> %v: task %s ended", duration, task.Name)
	}
	return nil
}
