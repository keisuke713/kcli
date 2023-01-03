package kcli

type (
	CmdName string

	Cmd interface {
		Name() string
		Usage() string
		Run([]string) error
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

var CmdMap = map[CmdName]Cmd{
	HELP: &HelpCmd{},
}
