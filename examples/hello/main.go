package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/suzuki-shunsuke/gopipe/gopipe"
)

func main() {
	if err := core(); err != nil {
		log.Fatal(err)
	}
}

func core() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	workflow := gopipe.Workflow{
		Tasks: []gopipe.Task{
			{
				Name: "init",
				Action: func(ctx context.Context) error {
					fmt.Println("init")
					return nil
				},
			},
			{
				Name: "test",
				Action: func(ctx context.Context) error {
					fmt.Println("test")
					return nil
				},
			},
		},
	}
	return workflow.Run(ctx)
}
