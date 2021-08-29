package parser

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type NodeType int

func (t *NodeType) String() string {
	switch *t {
	case Entry:
		return "Entry"
	case Dictionary:
		return "Dictionary"
	}
	return ""
}

const (
	Entry NodeType = iota
	Dictionary
)

type Noder interface {
	GetPackageName() (string, error)
}

type Builder struct {
	Fs
}

func (b *Builder) TravelRoot(path string) (*Node, error) {
	files, err := b.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var modFile os.FileInfo
	for _, f := range files {
		if f.Name() == "go.mod" {
			modFile = f
		}
	}
	if modFile == nil {
		return nil, errors.New("no mod file found")
	}
	rootNode, err := b.Travel(path, nil)
	if err != nil {
		return nil, err
	}
	rootNode.ProjectManager = ProjectManager{
		ModFile: modFile.Name(),
	}
	return rootNode, nil
}

func (b *Builder) Travel(path string, node *Node) (*Node, error) {
	files, err := b.ReadDir(path)
	if err != nil {
		return nil, err
	}
	n := &Node{
		Parent:    node,
		Type:      Dictionary,
		Ast:       nil,
		Path:      path,
		Children:  []*Node{},
		Collector: nil,
	}
	for _, f := range files {
		if f.IsDir() {
			child, err := b.Travel(filepath.Join(path, f.Name()), n)
			if err != nil {
				return nil, err
			}
			n.PushChild(child)
		}
		if strings.Contains(f.Name(), SHUFFIX_GENERATOR_FILE) {
			continue
		}

		if strings.Contains(f.Name(), ".go") {
			//TODO ast implement
			ast, err := b.BuildAst(filepath.Join(path, f.Name()))
			if err != nil {
				return nil, err
			}
			fn := &Node{
				Parent:    n,
				Type:      Entry,
				Ast:       ast,
				Path:      filepath.Join(path, f.Name()),
				Children:  nil,
				Collector: NewCollector(ast),
			}
			n.PushChild(fn)
		}
	}
	return n, nil
}

func (b *Builder) BuildAst(filename string) (*ast.File, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	//read source code
	rawCode, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	return parser.ParseFile(token.NewFileSet(), "", rawCode, parser.ParseComments)

}

type Fs struct {
}

func (fs *Fs) ReadDir(dir string) ([]os.FileInfo, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Readdir(0)
}

// parser node
// Parent is current node's parent node
// Type Entry indicate that current node is an astable node
// Type Dictionary indicate that current node is an dictionary node
type Node struct {
	ProjectManager
	Parent *Node
	Type   NodeType
	//Ast node if Type == Entry
	Ast *ast.File
	//Dictionary path if Type == Dictionary
	Path string

	Children []*Node

	Collector *Collector
}

func (n *Node) String() string {
	return n.format(0)
}

func (n *Node) format(level int) string {
	buf := bytes.NewBuffer([]byte{})
	//indentFprintf(buf, "Type: %s\n", level, n.Type)
	indentFprintf(buf, "Path: %s\n", level, n.GetFullPathName())
	//indentFprintf(buf, "Ast: %+v\n", level, n.Ast)
	if n.Collector != nil {
		//indentFprintf(buf, "Annotations: %+v\n", level, n.Collector.GetAllAnnotationsString())
	}

	if n.Ast != nil {
		b := bytes.NewBufferString("")
		format.Node(b, token.NewFileSet(), n.Ast)
		s, _ := ioutil.ReadAll(b)
		indentFprintf(buf, "Source code ============\n%s", level, s)
	}

	if len(n.Children) > 0 {
		for idx, cn := range n.Children {
			nstr := cn.format(level + 1)
			indentFprintf(buf, "Children[%d]: %s\n", level, idx, nstr)
		}
	}
	return buf.String()
}

func (n *Node) GetFullPathName() string {
	if n.Parent == nil {
		return "."
	}
	path := fmt.Sprintf("%s/%s", n.Parent.GetFullPathName(), n.Path)
	return path
}

func indentFprintf(w io.Writer, format string, level int, args ...interface{}) {
	for i := 0; i < level; i++ {
		fmt.Fprint(w, "  ")
	}
	fmt.Fprintf(w, format, args...)
}

func (n *Node) Apply(fn func(*Node) error) error {
	if n == nil {
		return nil
	}
	fn(n)
	if len(n.Children) > 0 {
		for _, cn := range n.Children {
			fn(cn)
		}
	}
	return nil
}

func (n *Node) PushChild(node *Node) {
	n.Children = append(n.Children, node)
}

type ProjectManager struct {
	//go.mod file
	ModFile string
	Module  string
}

func (p *ProjectManager) SetModFile(filepath string) {
	p.ModFile = filepath
}

func (p *ProjectManager) GetPackageNameFromString(s string) (string, error) {
	//module github.com/zrb-inc/spring
	reg := regexp.MustCompile("^module *(.*)")
	match := reg.FindAllStringSubmatch(s, -1)

	if len(match) >= 1 && len(match[0]) > 1 {
		return match[0][1], nil
	}

	return "", errors.New("Package Not Found")
}

func (p *ProjectManager) GetPackageName() (string, error) {
	if p.Module != "" {
		return p.Module, nil
	}
	fd, err := os.OpenFile(p.ModFile, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	var b []byte
	if b, err = ioutil.ReadAll(fd); err != nil {
		return "", err
	}
	return p.GetPackageNameFromString(string(b))
}
