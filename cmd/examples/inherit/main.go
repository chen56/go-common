package main

import (
	"flag"
	"fmt"
	"github.com/chen56/go-common/cmd"
	"os"
)

func main() {
	type dockerCmd struct {
		cmd.Cmd
		logLevel *string
		debug    *bool
	}

	type imageCmd struct {
		cmd.Cmd
	}

	type buildCmd struct {
		cmd.Cmd
		tag   *string
		quiet *bool
	}

	var docker = &dockerCmd{}
	var build = &buildCmd{}
	var image = &imageCmd{}
	//docker
	//  image
	//     build
	docker.Add(image)
	image.Add(build)

	docker.Cmd = cmd.Cmd{
		Name: "docker",
		Init: func(flagSet *flag.FlagSet) error {
			docker.logLevel = flagSet.String("log-level", "debug", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)
			docker.debug = flagSet.Bool("debug", false, "Enable debug mode")
			return nil
		},
		Run: func(args []string) error {
			fmt.Println("docker...")
			return nil
		},
	}

	image.Cmd = cmd.Cmd{
		Name:  "image",
		Short: "Manage images",
	}

	build.Cmd = cmd.Cmd{
		Name:  "build",
		Short: "Build an image from a Dockerfile",
		Init: func(flagSet *flag.FlagSet) error {
			build.tag = flagSet.String("tag", "sss", "Name and optionally a tag in the 'name:tag' format")
			build.quiet = flagSet.Bool("quiet", false, "Suppress the build output and print image ID on success")
			return nil
		},
		Run: func(args []string) error {
			fmt.Printf("build tag=%s quiet=%t \n", *build.tag, *build.quiet)
			return nil
		},
	}

	cmd.Run(docker, os.Args...)

}
