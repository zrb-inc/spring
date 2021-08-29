package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"github.com/zrb-inc/spring/internals/parser"
)

// AstWriter write the parsed ast tree to file
type AstWriter interface {
	Write(n *parser.Node) error
}

type CompiledAstWriter struct{}

func (c *CompiledAstWriter) Write(n *parser.Node) error {
	if n.Type != parser.Entry {
		return nil
	}
	if n.Ast == nil {
		return errors.New("parser.Node.Ast empty")
	}
	b := bytes.NewBuffer([]byte{})

	if err := format.Node(b, token.NewFileSet(), n.Ast); err != nil {
		return err
	}

	pathname := n.GetFullPathName()
	pathnameWithoutShuffix := ""
	if strings.Contains(pathname, ".go") {
		pathnameWithoutShuffix = strings.ReplaceAll(pathname, ".go", "")
	} else {
		return nil
	}
	newPathname := fmt.Sprintf("%s_spring_generated.go", pathnameWithoutShuffix)

	ioutil.WriteFile(newPathname, b.Bytes(), os.FileMode(0644))

	signatureComment := `//+build spring

//this file rewrite by spring`

	oldCode := fmt.Sprintf("%s\n%s", signatureComment, b.String())
	ioutil.WriteFile(pathname, []byte(oldCode), os.FileMode(0644))

	return nil
}
