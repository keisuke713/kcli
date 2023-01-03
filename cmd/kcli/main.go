package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"kcli"
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

	sub := args[1]
	f := flag.NewFlagSet(sub, flag.ContinueOnError)
	f.Usage = func() {
		if err := kcli.ShowUsage(stdout); err != nil {
			fmt.Fprintf(stderr, "failed to show usage: %v\n", err)
		}
	}
	var v bool
	f.BoolVar(&v, "v", false, "show debug print")
	if err := f.Parse(args[2:]); err == flag.ErrHelp {
		return err
	} else if err != nil {
		return fmt.Errorf("%q command with invalid args(%q): %w", sub, args[2:], err)
	}

	log.SetOutput(stderr)
	if !v {
		log.SetOutput(ioutil.Discard)
	}

	if cmd, ok := kcli.CmdMap[kcli.CmdName(sub)]; ok {
		args := f.Args()
		if len(args) != cmd.NArg() {
			return fmt.Errorf("%q command expects %d options, but actually %d options\n", cmd.Name(), cmd.NArg(), len(args))
		}
		if err := cmd.Run(args); err != nil {
			return fmt.Errorf("%q command failed: %w", sub, err)
		}
	} else {
		return fmt.Errorf("unknown command %q", sub)
	}
	return nil
}
