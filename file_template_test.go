package kcli

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func pointerOf[T any](e T) *T {
	return &e
}

var _ = Describe("newNodeName", func() {
	const (
		def = "kcli"
	)
	var (
		node *Node

		res string
	)
	JustBeforeEach(func() {
		res = newNodeName(node, def)
	})
	Context("ファイルの場合", func() {
		Context("Nameがnil", func() {
			BeforeEach(func() {
				node = &Node{
					Name: nil,
					Type: file,
				}
			})
			It("init.go", func() {
				Expect(res).To(Equal("kcli.go"))
			})
		})
		Context("Nameが指定されている場合", func() {
			BeforeEach(func() {
				node = &Node{
					Name: pointerOf("init"),
					Type: file,
				}
			})
			It("init.go", func() {
				Expect(res).To(Equal("init.go"))
			})
		})
	})
	Context("フォルダの場合", func() {
		BeforeEach(func() {
			node = &Node{
				Name: pointerOf("cmd"),
				Type: dir,
			}
		})
		It("cmd", func() {
			Expect(res).To(Equal("cmd"))
		})
	})
})

var _ = Describe("NewContent", func() {
	var (
		b           = &bytes.Buffer{}
		contentType ContentType
		param       contentTypeParam

		expect string

		ft, _ = newFileTemplate(nil, "kcli")
	)
	JustBeforeEach(func() {
		ft.newContent(b, contentType, param)
	})
	Context("main", func() {
		BeforeEach(func() {
			b.Reset()
			contentType = "main"
			param = contentTypeParam{
				PkgName: "kcli",
			}
			expect = `package main

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
`
		})
		It("main用のソースコード", func() {
			Expect(b.String()).To(Equal(expect))
		})
	})
	Context("core", func() {
		BeforeEach(func() {
			b.Reset()
			contentType = "core"
			param = contentTypeParam{
				PkgName: "kcli",
			}
			expect = `package kcli

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
    HELP      CmdName = "help"
)

var CmdMap = map[CmdName]Cmd{
    HELP:      &HelpCmd{},
}
`
		})
		It("core用のソースコード", func() {
			Expect(b.String()).To(Equal(expect))
		})
	})
	Context("help command", func() {
		BeforeEach(func() {
			b.Reset()
			contentType = "cmd"
			param = contentTypeParam{
				PkgName:      "kcli",
				CmdName:      "Help",
				ShortCmdName: "h",
			}
			expect = `package kcli

import (
    "log"
)

var _ Cmd = &HelpCmd{}

type HelpCmd struct {
    logger *log.Logger
}

func NewHelpCmd(l *log.Logger, pkgName string) (*HelpCmd, error) {
    return &HelpCmd{
        l,
    }, nil
}

func (h *HelpCmd) Name() string {
    return ""
}

func (h *HelpCmd) Usage() string {
    return ""
}

func (h *HelpCmd) Run(args []string) error {
    return nil
}

func (h *HelpCmd) NArg() int {
    return 0
}
`
		})
		It("コマンド用のソースコード", func() {
			Expect(b.String()).To(Equal(expect))
		})
	})
})
