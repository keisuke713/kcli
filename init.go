package kcli

import (
	"log"
)

var _ Cmd = &InitCmd{}

type InitCmd struct {
	FT     *FileTemplate
	logger *log.Logger
}

func NewInitCmd(l *log.Logger, pkgName string) (*InitCmd, error) {
	ft, err := newFileTemplate(l, pkgName)
	if err != nil {
		return nil, err
	}

	return &InitCmd{
		ft,
		l,
	}, nil
}

func (i *InitCmd) Name() string {
	return string(INIT)
}

func (i *InitCmd) Usage() string {
	return "Init cli program"
}

func (i *InitCmd) Run(args []string) error {
	i.logger.Printf("%v: init dir and file\n", i.Name())
	i.FT.Build()
	return nil
}

func (i *InitCmd) NArg() int {
	return 0
}
