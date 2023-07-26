package kcli

import (
	"fmt"
	"log"
	"os"
)

type (
	CmdName string

	Cmd interface {
		Name() string
		Usage() string
		Run([]string) error
		NArg() int
	}
)

const (
	BINARY_NAME = "kcli"
)

const (
	INIT CmdName = "init"
	ADD  CmdName = "add"
	HELP CmdName = "help"
)

var (
	CmdMap map[CmdName]Cmd
	logger *log.Logger
)

func init() {
	path, _ := os.Getwd()
	pkgName := getCurrDir(path)

	logger = log.New(os.Stdout, fmt.Sprintf("%v: ", pkgName), -1)
	logger.Println("path: ", path)
	logger.Println("current dir: ", pkgName)

	initCmd, err := NewInitCmd(logger, pkgName)
	if err != nil {
		panic(err)
	}
	CmdMap = map[CmdName]Cmd{
		HELP: &HelpCmd{},
		INIT: initCmd,
	}

}
