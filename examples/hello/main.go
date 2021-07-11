package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/suzuki-shunsuke/gopipe/command"
	"github.com/suzuki-shunsuke/gopipe/gopipe"
)

func main() {
	if err := core(); err != nil {
		log.Fatal(err)
	}
}

func Hello(getArgs gopipe.GetArgs) *gopipe.Task {
	return &gopipe.Task{
		Name: "Hello",
		Steps: []*gopipe.Step{
			{
				Name: "Hello",
				Action: func(ctx context.Context, args *gopipe.Args) error {
					fmt.Println("Hello,", args.Get("name"))
					return nil
				},
			},
		},
		Args: getArgs,
	}
}

func core() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	args := &gopipe.Args{}
	args.Set("branch", "develop")
	workflow := gopipe.Workflow{
		Name: "main",
		Tasks: []*gopipe.Task{
			{
				Name: "init",
				Steps: []*gopipe.Step{
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
				Steps: []*gopipe.Step{
					{
						Name: "go vet",
						Action: command.Command(exec.Command("go", "vet"), nil,
							command.Dir("."),
							command.Envs(os.Environ()),
							command.Env("FOO", "FOO"),
						),
						If: func(ctx context.Context, args *gopipe.Args) (bool, error) {
							fmt.Println("name:", args.GetString("name"))
							return args.Get("branch") != "main", nil
						},
					},
				},
				Args: func(mArgs map[string]*gopipe.Args) (*gopipe.Args, error) {
					args := &gopipe.Args{}
					args.Set("branch", mArgs["_"].Get("branch"))
					args.Set("name", mArgs["init"].Get("name"))
					return args, nil
				},
			},
			Hello(func(mArgs map[string]*gopipe.Args) (*gopipe.Args, error) {
				args := &gopipe.Args{}
				args.Set("name", mArgs["init"].Get("name"))
				return args, nil
			}),
		},
	}
	return workflow.Run(ctx, args)
}
