package kcli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type (
	FileTemplate struct {
		Tree     map[NodeName]Node      `yaml:"Tree"`
		Template map[ContentType]string `yaml:"Template"`
		logger   *log.Logger
		pkgName  string
	}

	Node struct {
		Name        *string
		Type        NodeType
		Type2       ContentType
		ContentType *string
		// ContentType ContentType
		Child *Node
	}
	contentType上手にマッピングできない問題

	contentTypeParam struct {
		PkgName      string
		CmdName      string
		ShortCmdName string
	}

	NodeName    string
	NodeType    string
	ContentType string
	extention   string
)

const (
	templatePath = "static/file_template.yaml"
)

// const (
// 	cmd  NodeName = "cmd"
// 	core NodeName = "core"
// 	help NodeName = "help"
// )

const (
	main ContentType = "main"
	core ContentType = "core"
	cmd  ContentType = "cmd"
)

const (
	file NodeType = "file"
	dir  NodeType = "dir"
)

const (
	goo extention = ".go"
)

func newFileTemplate(l *log.Logger, pkgName string) (*FileTemplate, error) {
	var ft FileTemplate
	bs, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}
	yaml.Unmarshal(bs, &ft)
	fmt.Println("jfdklsjflkasjdfl")
	fmt.Println(ft.Tree["core"])
	fmt.Println(ft.Tree["help"])
	ft.logger = logger
	ft.pkgName = pkgName
	return &ft, nil
}

func (ft *FileTemplate) Build() {
	for _, node := range ft.Tree {
		ft.build("", &node)
	}
}

func (ft *FileTemplate) build(path string, n *Node) {
	if n == nil {
		return
	}

	nn := newNodeName(n, ft.pkgName)
	newPath := filepath.Join(path, nn)
	if n.IsFile() {
		f, err := os.Create(newPath)
		if err != nil {
			ft.logger.Fatalf("failed to create file which is named %v \n err: %v \n", nn, err)
		}
		ft.logger.Printf("create file which is named %v \n", nn)
		// エラーが出なかったらファイルに書き込む
		// 各ノードにmain,core,cmdのどれかを判別するためのフィールドを作るe
		// それを元にmapを用意
		if n.ContentType == nil {
			ft.logger.Fatalf("%v doesn't have content \n", nn)
		}
		bf := &bytes.Buffer{}
		for i, c := range nn {
			if i == 0 {
				bf.WriteRune(c + 20)
			}
			bf.WriteRune(c)
		}
		ctp := contentTypeParam{
			PkgName:      ft.pkgName,
			CmdName:      bf.String(),
			ShortCmdName: nn[0:1],
		}
		ft.newContent(f, ContentType(*n.ContentType), ctp)
		return
	}
	ft.logger.Printf("dir path: %v \n", newPath)
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		ft.logger.Fatalf("failed to create dir which is named %v  \n err: %v \n", nn, err)
	}
	ft.logger.Printf("create dir which is named %v \n", nn)
	ft.build(newPath, n.Child)
}

func (ft *FileTemplate) newContent(w io.Writer, contentType ContentType, param contentTypeParam) {
	var (
		content string
		ok      bool
	)
	if content, ok = ft.Template[contentType]; !ok {
		return
	}

	tmpl, err := template.New("contentTemplate").Parse(content)
	if err != nil {
		return
	}
	tmpl.Execute(w, param)
}

func (n *Node) IsFile() bool {
	return n.Type == file
}

func newNodeName(n *Node, def string) string {
	builder := &strings.Builder{}
	builder.WriteString(orDefault(n.Name, def))

	if n.IsFile() {
		builder.WriteString(string(goo))
	}
	return builder.String()
}

func orDefault[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}
