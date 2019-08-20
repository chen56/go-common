package cmd

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

func init() {
}

func TestGetNextFlags(t *testing.T) {
	type cmd struct {
		name string
		args []string
	}
	spew.Dump()
	var testdatas = []struct {
		goal string
		args []string
		err  string
		cmds []cmd
	}{
		{
			"hello world",
			[]string{"docker", "-log-level=debug", "image", "build", "-tag=chen56/ubuntu:1.0", "./"},
			"",
			[]cmd{
				{"docker", []string{"-log-level=debug"}},
				{"image", nil},
				{"build", []string{"-tag=chen56/ubuntu:1.0", "./"}},
			},
		},
		{
			"父命令应print help",
			[]string{"docker"},
			flag.ErrHelp.Error(),
			[]cmd{
				{"docker", nil},
			},
		},
		{
			"父命令应print help, 即使有flag",
			[]string{"docker", "-log-level=debug"},
			flag.ErrHelp.Error(),
			[]cmd{
				{"docker", []string{"-log-level=debug"}},
			},
		},
		{
			"父命令应print help",
			[]string{"docker", "-log-level=debug", "image"},
			flag.ErrHelp.Error(),
			[]cmd{
				{"docker", []string{"-log-level=debug"}},
				{"image", nil},
			},
		},
		{
			"叶子命令",
			[]string{"docker", "-log-level=debug", "image", "build"},
			"",
			[]cmd{
				{"docker", []string{"-log-level=debug"}},
				{"image", nil},
				{"build", nil},
			},
		},
		{
			"叶子命令的选项",
			[]string{"docker", "-log-level=debug", "image", "build", "-tag", "chen56/ubuntu:1.0"},
			"",
			[]cmd{
				{"docker", []string{"-log-level=debug"}},
				{"image", nil},
				{"build", []string{"-tag", "chen56/ubuntu:1.0"}},
			},
		},
		{
			"叶子命令:help",
			[]string{"docker", "-log-level=debug", "image", "build", "-h"},
			flag.ErrHelp.Error(),
			[]cmd{
				{"docker", []string{"-log-level=debug"}},
				{"image", nil},
				{"build", []string{"-h"}},
			},
		},
		{
			"叶子命令:除了build起作用，ls也起作用",
			[]string{"docker", "-log-level=debug", "image", "ls"},
			"",
			[]cmd{
				{"docker", []string{"-log-level=debug"}},
				{"image", nil},
				{"ls", nil},
			},
		},
	}

	for i, testdata := range testdatas {
		var docker = newTestCmd()

		last, err := Run(docker, testdata.args...)
		fmt.Println("===========> ", last)

		for j, node := range last.nodes() {
			fmt.Println("node:", j, ".", node.path())

			msg := fmt.Sprintln(strconv.Itoa(i), ".", j, testdata.goal, "-", node.path())
			require.Equal(t, testdata.cmds[j].name, node.Name, msg)
			require.Equal(t, testdata.cmds[j].args, node.Args(), msg)
		}

		msg := strconv.Itoa(i) + ". " + testdata.goal + " - " + strings.Join(testdata.args, " ")
		require.Equal(t, len(testdata.cmds), len(last.nodes()), msg)
		if testdata.err == "" {
			require.NoError(t, err, msg)
		} else {
			fmt.Println("error====>", err)
			require.Contains(t, err.Error(), testdata.err, msg)
		}

	}
}

func newTestCmd() *Cmd {
	var logLevel *string
	var debug *bool
	docker := &Cmd{
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

	image := &Cmd{
		Name: "image",
	}

	var tag *string
	var quiet *bool
	build := &Cmd{
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

	ls := &Cmd{
		Name: "ls",
		Run: func(args []string) error {
			fmt.Printf("ls")
			return nil
		},
	}

	docker.Add(image)
	image.Add(build)
	image.Add(ls)
	return docker
}
