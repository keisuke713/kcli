package kcli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
)

var _ Cmd = &HelpCmd{}

type HelpCmd struct{}

func (h *HelpCmd) Name() string {
	return string(HELP)
}

func (h *HelpCmd) Usage() string {
	return "Show usage"
}

var desc = `Description: kcli is incredibly CLI tool that make it easy to create CLI outline
usage: kcli <subcommand> [<args>]

Subcommands:
`

func ShowUsage(w io.Writer) error {
	cms := make([]string, len(CmdMap))
	var i int
	for k := range CmdMap {
		cms[i] = string(k)
		i++
	}
	sort.Strings(cms)

	tw := tabwriter.NewWriter(w, 0, 4, 1, ' ', 0)
	fmt.Fprintf(tw, "%s", desc)
	for _, k := range cms {
		cn := CmdName(k)
		fmt.Fprintf(tw, "\t%s\t%s\n", CmdMap[cn].Name(), CmdMap[cn].Usage())
	}
	return tw.Flush()
}

func (h *HelpCmd) Run(args []string) error {
	return ShowUsage(os.Stdout)
}

func (h *HelpCmd) NArg() int {
	return 0
}
