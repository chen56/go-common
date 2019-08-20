# cmd

父子命令形式的flag封装.

go内置的flag包封装的足够简单，我们并不希望因为父子命令的需求而引入比较复杂的第三方包，
简单封装一下flag.FlagSet,就足以达到很不错的效果: 

父子命令帮助：

```
$ cd examples/no_inherit
$ ./docker -h
Usage: docker [FLAGS] SUB_COMMAND [FLAGS] [ARGS]

Flags:
  -debug
        Enable debug mode
  -log-level string
        Set the logging level ("debug"|"info"|"warn"|"error"|"fatal") (default "debug")

Commands:
  build Build an image from a Dockerfile
  ps List containers



Run 'docker -help' for more information on a command.
```

子命令执行

```bash
$ ./docker ps mysql_1
ps: [mysql_1] 


```

以上范例只有一层父子命令，实际可以多层嵌套，比如：
```bash
$ ./docker image inspect mysql_1
```
