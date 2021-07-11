package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
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
		Name: "main",
		Tasks: []gopipe.Task{
			{
				Name: "init",
				Steps: []gopipe.Step{
					{
						Name: "init",
						Action: func(ctx context.Context, args *gopipe.Args) error {
							fmt.Println("init")
							args.Set("name", "yoo")
							return nil
						},
					},
				},
			},
			{
				Name: "go vet",
				Steps: []gopipe.Step{
					{
						Name: "go vet",
						Action: gopipe.Command(exec.Command("go", "vet"), nil,
							gopipe.Dir("."),
							gopipe.Env("FOO", "FOO")),
						If: func(ctx context.Context, args *gopipe.Args) (bool, error) {
							fmt.Println("name:", args.GetString("name"))
							return args.Get("branch") != "main", nil
						},
					},
				},
			},
		},
	}
	args := &gopipe.Args{}
	args.Set("branch", "main")
	return workflow.Run(ctx, args)
}
