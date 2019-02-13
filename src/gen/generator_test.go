package gen

import (
	"fmt"
	"testing"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/h"
	"github.com/sxarp/c_compiler_go/src/psr"
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

type psrTestCase struct {
	rcode string
	ins   []asm.Fin
	tf    bool
}

func compCode(t *testing.T, p psr.Parser, c psrTestCase) {
	lhs := asm.New()
	for _, i := range c.ins {
		lhs.Ins(i)
	}

	rhs := asm.New()
	a, _ := p.Call(tok.Tokenize(c.rcode))
	a.Eval(&rhs)
	h.Expectt(t, c.tf, lhs.Eq(&rhs))

}

func TestNunInt(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"42",
			[]asm.Fin{asm.I().Push().Val(42)},
			true,
		},
		{

			"43",
			[]asm.Fin{asm.I().Push().Val(42)},
			false,
		},
	} {
		compCode(t, numInt, c)
	}
}

func TestAdder(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"+2",
			[]asm.Fin{
				asm.I().Push().Val(2),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Add().Rax().Rdi(),
				asm.I().Push().Rax(),
			},
			true,
		},
	} {
		compCode(t, adder(&numInt), c)
	}
}

func TestSubber(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"-2",
			[]asm.Fin{
				asm.I().Push().Val(2),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Sub().Rax().Rdi(),
				asm.I().Push().Rax(),
			},
			true,
		},
	} {
		compCode(t, subber(&numInt), c)
	}
}

func Testaddsubs(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"1+1",
			[]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Push().Val(1),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Sub().Rax().Rdi(),
				asm.I().Push().Rax(),
			},
			true,
		},
	} {
		compCode(t, addsubs(&numInt), c)
	}
}
