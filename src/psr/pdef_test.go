package psr

import (
	"fmt"
	"testing"

	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/tok"
)

func TestNum(t *testing.T) {
	var tokens []tok.Token
	var a ast.AST

	for _, c := range []struct {
		s string
		r bool
	}{
		{"1", true},
		{"1+1", true},
		{"1*1", true},
		{"1*(2+3)", true},
		{"1*2+3)", false},
		{"1*2+3*4*5", true},
		{"1*2++3*4*5", false},
		{"1*(2+3)*4+5", true},
		{"1*(2+3)/4+5", true},
		{"1*(2+3))*4+5", false},
		{"1/(2+3)*4+5/", false},
	} {
		tokens = tok.Tokenize(c.s)

		a, _ = GenParser().Call(tokens)
		ast.CheckAst(t, c.r, a)
		fmt.Println(a.Show())
	}

}

func TestNum2(t *testing.T) {
	var tokens []tok.Token
	var a ast.AST

	for _, c := range []struct {
		s string
		r bool
	}{
		{"1", true},
		{"1+2", true},
		{"1*2", true},
		{"1+2+3", true},
		{"1*2*3", true},
		{"1+2*3", true},
		{"1*2+3", true},
		{"1*(2+3)*4+5", true},
		{"1*2++3*4*5", false},
		{"1*(2+3)*4+5", true},
		{"1*(2+3)/4+5", true},
		{"1*2+3)", false},
	} {
		tokens = tok.Tokenize(c.s)

		a, _ = GenParser2().Call(tokens)
		ast.CheckAst(t, c.r, a)
		fmt.Println(a.Show())
	}

}
