package cmd

import (
	"flag"
	"fmt"
	"github.com/chen56/go-common/must"
	"github.com/pkg/errors"
	"os"
	"reflect"
	"strings"
)

var boolFlagType reflect.Type
var emptyOutput = &emptyOutputType{}

type Cmder interface {
	cmd() *Cmd
}

type emptyOutputType struct {
}

func (x *emptyOutputType) Write(p []byte) (n int, err error) {
	return len(p), err
}

func init() {
	s := flag.NewFlagSet("x", flag.ContinueOnError)
	s.Bool("bool", false, "")
	boolFlagType = reflect.TypeOf(s.Lookup("bool").Value)
}

// 一个命令结构体,命令可以是类似docker container build 这样的父子命令.
type Cmd struct {
	Init func(flagSet *flag.FlagSet) error

	Run func(args []string) error

	// 命令名
	Name string

	// 一行形式的Usage
	UsageLine string

	// 短描述.
	Short string

	// 长描述.
	Long string

	// 父子命令都可以有自己的flag.
	flagSet *flag.FlagSet

	// children cmd.
	children []*Cmd

	parent *Cmd

	// no flags args
	args []string

	cmder Cmder
}

// 运行命令
//   docker 必须是跟命令
//   args 应该传入os.Args, 例: cmd.Run(root,os.Args...)
// err 可能为 flag.ErrHelp , err已处理并打印相关信息，外层代码可以不处理
func Run(root Cmder, args ...string) (*Cmd, error) {

	rootCmd := root.cmd()

	must.True(rootCmd.isRoot(), "bug: `%s` should be root", rootCmd.Name)

	if len(args) == 0 {
		panic("bug: first args should be cmd")
	}

	parsed, err := parse(root, args...)
	if err != nil {
		if err == flag.ErrHelp {
			parsed.printHelp()
		} else {
			parsed.printError(err)
		}
		return parsed, err
	}

	// 试图执行命令
	if parsed.Run != nil {
		err = parsed.Run(parsed.args)
		if err != nil {
			// 命令运行失败，只打印错误本身，不需要printSeeHelp
			parsed.printf("%s\n", err.Error())
			//run失败退出应用
			os.Exit(1)
			//return parsed, err
		}
	}

	return parsed, nil
}

// Add
//  return cmder arg
func (x *Cmd) Add(cmder Cmder) Cmder {
	must.NotNil(cmder, "bug: cmder should not be nil")

	cmd := cmder.cmd()

	must.Nil(cmd.parent, "bug:`%s` already added , parent: `%s` ", cmd.parent, cmd.Name)

	cmd.cmder = cmder
	cmd.parent = x
	x.children = append(x.children, cmd)
	return cmder
}

func parse(rootCmder Cmder, args ...string) (*Cmd, error) {
	// init if need
	root := rootCmder.cmd()
	err := root.initTree(rootCmder)
	if err != nil {
		return root, err
	}

	// 先把参数，切割为父子命令
	parsed, err := root.split(args[1:])
	if err != nil {
		return parsed, err
	}

	// 通过setoutput(nil) 禁用FlagSet的出错print
	var quietParse = func(s *flag.FlagSet, arguments []string) error {
		var output = s.Output()
		defer s.SetOutput(output)

		s.SetOutput(emptyOutput)
		return s.Parse(arguments)
	}

	// 再把父子命令的flagSet逐个Parse
	for _, node := range parsed.nodes() {
		err := quietParse(node.flagSet, node.args)
		if err != nil {
			return node, err
		}
	}
	return parsed, nil
}

func (x *Cmd) initTree(cmder Cmder) error {
	x.cmder = cmder
	x.flagSet = flag.NewFlagSet(x.Name, flag.ContinueOnError)
	x.flagSet.Usage = func() {
		//x.printSeeHelp()
	}
	if x.Init != nil {
		err := x.Init(x.flagSet)
		if err != nil {
			return err
		}
	}
	for _, child := range x.children {
		if err := child.initTree(child.cmder); err != nil {
			return err
		}
	}
	return nil
}

func (x *Cmd) cmd() *Cmd {
	return x
}

func (x *Cmd) isRoot() bool {
	return x.parent == nil
}

func (x *Cmd) isLeaf() bool {
	return len(x.children) == 0
}
func (x *Cmd) Args() []string {
	return x.args
}

// build
func (x *Cmd) split(args []string) (last *Cmd, err error) {
	if x.isLeaf() {
		x.args = args

		//无参数统一为nil
		if len(x.args) == 0 {
			x.args = nil
		}
		return x, nil
	}

	//下面都是父命令的flags parse
	for {
		if len(args) == 0 { //父命令无后续参数，说明是要帮助
			return x, flag.ErrHelp
		}

		arg := args[0]
		args = args[1:]

		if len(arg) == 0 {
			return x, errors.New("invalid empty string '' args;")
		}

		if strings.Index(arg, "-") == 0 { //flag
			x.args = append(x.args, arg)

			if arg == "--" {
				return x, errors.New("only leaf cmd can contains '--' args")
			}

			if arg == "-h" || arg == "-help" || arg == "--help" {
				return x, flag.ErrHelp
			}

			flagName := arg[1:]

			var flagHasValue = false
			for i := 0; i < len(flagName); i++ { // equals cannot be first
				if flagName[i] == '=' {
					flagHasValue = true
					flagName = flagName[0:i]
					break
				}
			}
			if f := x.flagSet.Lookup(flagName); f == nil {
				return x, errors.Errorf("[%s]'s flag provided but not defined: -%s", strings.Join(x.path(), " "), flagName)
			} else {
				if boolFlagType == reflect.TypeOf(f.Value) { //bool value
					continue
				}
				if flagHasValue {
					continue
				}
				if len(args) == 0 {
					return x, errors.Errorf("[%s]'s flag needs an argument: -%s", strings.Join(x.path(), " "), flagName)
				}

				//token 向前走一步
				arg := args[0]
				args = args[1:]
				x.args = append(x.args, arg)
			}
		} else { //sub command
			// 用双引号可传入空字符串 `docker  ""`
			if len(arg) == 0 {
				return x, errors.Errorf("sub command '%s' not found;", arg)
			}
			child := x.getChild(arg)
			if child == nil {
				return x, errors.Errorf("'%s' sub command '%s' not found", x.Name, arg)
			}

			//继续parse下一个子命令
			return child.split(args)
		}
	}
}
func (x *Cmd) path() []string {
	if x.isRoot() {
		return []string{x.Name}
	} else {
		return append(x.parent.path(), x.Name)
	}
}
func (x *Cmd) nodes() []*Cmd {
	if x.isRoot() {
		return []*Cmd{x}
	} else {
		return append(x.parent.nodes(), x)
	}
}

func (x *Cmd) Root() *Cmd {
	if x.isRoot() {
		return x
	}
	return x.parent.Root()
}

func (x *Cmd) getChild(name string) *Cmd {
	for _, child := range x.children {
		if child.Name == name {
			return child
		}
	}
	return nil
}

func (x *Cmd) DumpNodes() string {
	s := ""
	for i, node := range x.nodes() {
		s += strings.Repeat("  ", i) + node.dumpCurrentNode() + "\n"
	}
	return s
}
func (x *Cmd) DumpTree() string {
	return x.dumpTree(0)
}
func (x *Cmd) dumpTree(level int) string {
	s := strings.Repeat("  ", level) + x.dumpCurrentNode() + "\n"
	for _, child := range x.children {
		s += child.dumpTree(level + 1)
	}
	return s
}
func (x *Cmd) dumpCurrentNode() string {
	s := strings.Join(x.path(), "/")
	flagSpecs := ""
	x.flagSet.VisitAll(func(i *flag.Flag) {
		flagSpecs += i.Name + ":" + i.DefValue + ","
	})
	s += fmt.Sprintf("%s Flags:{%s} args:%v", x.Name, flagSpecs, x.args)
	return s
}

func (x *Cmd) String() string {
	return strings.Join(x.path(), "/") + "\n"
}

func (x *Cmd) printSeeHelp() {
	x.printf("\n")
	x.printf("Run '%s -help' for more information on a command.\n",
		strings.Join(x.path(), " "))
}

func (x *Cmd) printf(format string, a ...interface{}) {
	fmt.Fprintf(x.flagSet.Output(), format, a...)
}

func (x *Cmd) printError(err error) {
	fmt.Fprintln(x.flagSet.Output(), err)
	x.printSeeHelp()
}

func (x *Cmd) printHelp() {
	f := x.flagSet
	{
		s := "Usage: "
		if x.UsageLine != "" {
			s += x.UsageLine
		} else { //default usage
			if x.isLeaf() { // leaf
				s += strings.Join(x.path(), " [FLAGS] ") + " [FLAGS] [ARGS]\n"
			} else { //parent
				s += strings.Join(x.path(), " [FLAGS] ") + " [FLAGS] SUB_COMMAND [FLAGS] [ARGS]\n"
			}
		}
		s += "\n"
		if x.Long != "" {
			s += "\n"
			s += x.Long
			s += "\n"
		} else {
			if x.Short != "" {
				s += "\n"
				s += x.Short
				s += "\n"
			}
		}
		x.printf(s)
	}

	x.printf("Flags:\n")
	f.PrintDefaults()

	if !x.isLeaf() {
		s := "\n"
		s += "Commands:\n"

		for _, child := range x.children {
			s += fmt.Sprintf("  %s %s\n", child.Name, child.Short)
		}
		s += "\n"
		x.printf(s)
	}
	x.printf("\n")
	x.printSeeHelp()
}
