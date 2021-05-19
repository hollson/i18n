package main

import (
    "flag"
    "fmt"
    "os"

    "golang.org/x/text/language"
)

func mainUsage() {
    fmt.Fprintf(os.Stderr, `这是一个「github.com/hollson/i18」多语言库的客户端工具.
用法:
	i18n_cli <command> [arguments]

command列表:
	extract		从go源码提取「i18n.Message」,即预翻译的消息(不包含测试文件)
	merge		合并翻译文件
`)
}

type command interface {
    name() string
    parse(arguments []string) error
    execute() error
}

//go:generate go build -o $GOPATH/bin/i18n_cli
func main() {
    code := testableMain(os.Args[1:])
    // fmt.Println("Done !!!")
    os.Exit(code)
}

func testableMain(args []string) int {
    flags := flag.NewFlagSet("i18n_cli", flag.ContinueOnError)
    flags.Usage = mainUsage
    if err := flags.Parse(args); err != nil {
        if err == flag.ErrHelp {
            return 2
        }
        return 1
    }
    if flags.NArg() == 0 {
        mainUsage()
        return 2
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
                return 1
            }
            if err := cmd.execute(); err != nil {
                fmt.Fprintln(os.Stderr, err)
                return 1
            }
            return 0
        }
    }
    fmt.Fprintf(os.Stderr, "i18n_cli: unknown subcommand %s\n", cmdName)
    return 1
}

type languageTag language.Tag

func (lt languageTag) String() string {
    return lt.Tag().String()
}

func (lt *languageTag) Set(value string) error {
    t, err := language.Parse(value)
    if err != nil {
        return err
    }
    *lt = languageTag(t)
    return nil
}

func (lt languageTag) Tag() language.Tag {
    tag := language.Tag(lt)
    if tag.IsRoot() {
        return language.English
    }
    return tag
}
