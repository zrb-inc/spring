package generator_test

import (
	"testing"

	"github.com/zrb-inc/spring/internals/generator"
)

func TestGenerate(t *testing.T) {
	generator.GenerateBean("Test", "j", "aa", "aaaaa", "alias")
}
