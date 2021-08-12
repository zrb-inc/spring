package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/pkg/errors"
)

const GenerateBeanInject = `
package main

import spruntime "github.com/zrb-inc/spring/internals/runtime"

func Generate%s() {
    i := &spruntime.GeneralBeanDefinition{
        PackageName: "%s",
        Id:          "%s",
        DependsOn:   []string{},
        Primary:     false,
        Constrction: func() (interface{}, error) {
            return &%s{}, nil
        },
    }
    runtime.CompileApp.PushDefinition("%s", i)
}
`

func GenerateBean(
	funcName string,
	packageName string,
	idName string,
	structName string,
	aliasName string,
) (*ast.GenDecl, *ast.FuncDecl, error) {
	code := fmt.Sprintf(GenerateBeanInject, funcName, packageName, idName, structName, aliasName)
	a, err := parser.ParseFile(token.NewFileSet(), "", code, parser.ParseComments)
	if err != nil {
		return nil, nil, errors.Wrap(GenerateBeanErr, err.Error())
	}
	var fd *ast.FuncDecl
	var id *ast.GenDecl
	for _, v := range a.Decls {
		if vv, ok := v.(*ast.FuncDecl); ok {
			fd = vv
		}
		if vv, ok := v.(*ast.GenDecl); ok {
			id = vv
		}
	}
	return id, fd, nil
}
