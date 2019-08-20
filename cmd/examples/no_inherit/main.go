package main

import (
	"flag"
	"fmt"
	"github.com/chen56/go-common/cmd"
	"os"
)

func main() {
	var logLevel *string
	var debug *bool
	docker := &cmd.Cmd{
		Name: "docker",
		Init: func(flagSet *flag.FlagSet) error {
			logLevel = flagSet.String("log-level", "debug", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)
			debug = flagSet.Bool("debug", false, "Enable debug mode")
			return nil
		},
		Run: func(args []string) error {
			fmt.Println("docker: logLevel:", logLevel, "debug:", debug)
			return nil
		},
	}

	image := &cmd.Cmd{
		Name: "image",
	}

	var tag *string
	var quiet *bool
	build := &cmd.Cmd{
		Name: "build",
		Init: func(flagSet *flag.FlagSet) error {
			tag = flagSet.String("tag", "sss", "Name and optionally a tag in the 'name:tag' format")
			quiet = flagSet.Bool("quiet", false, "Suppress the build output and print image ID on success")
			return nil
		},
		Run: func(args []string) error {
			fmt.Printf("build tag=%s quiet=%t \n", *tag, *quiet)
			return nil
		},
	}
	ls := &cmd.Cmd{
		Name: "ls",
		Init: func(flagSet *flag.FlagSet) error {
			return nil
		},
		Run: func(args []string) error {
			fmt.Printf("ls")
			return nil
		},
	}

	docker.Add(image)
	image.Add(build)
	image.Add(ls)

	parsed, err := cmd.Run(docker, os.Args...)

	fmt.Println("=========")
	fmt.Println(parsed.Name)
	fmt.Println(err)
}
