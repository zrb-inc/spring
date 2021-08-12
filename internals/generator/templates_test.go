package generator_test

import (
	"go/printer"
	"go/token"
	"os"
	"testing"

	"github.com/zrb-inc/spring/internals/generator"
)

func TestGenerate(t *testing.T) {
	_, a, _ := generator.GenerateBean("Test", "j", "aa", "aaaaa", "alias")
	printer.Fprint(os.Stdout, token.NewFileSet(), a)
}
