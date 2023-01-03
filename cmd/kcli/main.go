package main

import (
	"fmt"
	"io"
	"kcli"
	"os"

	_ "kcli"
)

const (
	ExitCodeOK              = 0
	ExitCodeParserFlagError = 1
)

func main() {
	if err := run(os.Stdout, os.Stderr, os.Args); err != nil {
		switch err := err.(type) {
		default:
			fmt.Fprintf(os.Stderr, "%v \n", err)
		}
		os.Exit(ExitCodeParserFlagError)
	}
	os.Exit(ExitCodeOK)
}

func run(stdout, stderr io.Writer, args []string) error {
	if len(args) < 2 {
		cmd := kcli.CmdMap[kcli.HELP]
		if err := cmd.Run(nil); err != nil {
			return fmt.Errorf("%s command failed: %w", cmd.Name(), err)
		}
		return nil
	}
	return nil
}
