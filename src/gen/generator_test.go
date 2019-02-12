package gen

import (
	"fmt"
	"testing"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/h"
	"github.com/sxarp/c_compiler_go/src/tok"
)

func TestArithmetics(t *testing.T) {
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

		a, _ = GenParser().Call(tokens)
		ast.CheckAst(t, c.r, a)
		fmt.Println(a.Show())
	}

}

func TestNunInt(t *testing.T) {

	for _, c := range []struct {
		ins   []asm.Fin
		rcode string
		r     bool
	}{
		{
			[]asm.Fin{asm.I().Push().Val(42)},
			"42",
			true,
		},
		{
			[]asm.Fin{asm.I().Push().Val(42)},
			"43",
			false,
		},
	} {
		lhs := asm.New()
		for _, i := range c.ins {
			lhs.Ins(i)
		}

		rhs := asm.New()
		a, _ := numInt.Call(tok.Tokenize(c.rcode))
		a.Eval(&rhs)
		h.Expectt(t, c.r, lhs.Eq(&rhs))
	}
}
