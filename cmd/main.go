package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, `这是一个「github.com/hollson/i18」多语言库的客户端工具.
Usage:
    i18n_cli <Command> [Option] <Param>...

Command:
    extract	从go源码提取「i18n.Message」,即预翻译的消息(不包含测试文件)
    merge	合并翻译文件

`)
}

type command interface {
	name() string
	parse(args []string) error
	execute() error
}

//go:generate go build -o $GOPATH/bin/i18n_cli
func main() {
	flags := flag.NewFlagSet("i18n_cli", flag.ContinueOnError)

	flags.Usage = usage
	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
	}

	if len(os.Args) == 1 {
		usage()
	}

	commands := []command{
		&mergeCommand{},
		&extractCommand{},
	}
	cmdName := flags.Arg(0)
	for _, cmd := range commands {
		if cmd.name() == cmdName {
			if err := cmd.parse(flags.Args()[1:]); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			if err := cmd.execute(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

