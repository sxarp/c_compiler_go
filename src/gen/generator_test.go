package gen

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/h"
	"github.com/sxarp/c_compiler_go/src/psr"
	"github.com/sxarp/c_compiler_go/src/tok"
)

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

func TestAlphaToNum(t *testing.T) {
	h.Expecti(t, 0, alpaToNum('a'))
	h.Expecti(t, 25, alpaToNum('z'))
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

func TestAddsubs(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"1+1",
			[]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Push().Val(1),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Add().Rax().Rdi(),
				asm.I().Push().Rax(),
			},
			true,
		},
	} {
		compCode(t, addsubs(&numInt), c)
	}
}

func TestMuler(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"*2",
			[]asm.Fin{
				asm.I().Push().Val(2),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Mul().Rdi(),
				asm.I().Push().Rax(),
			},
			true,
		},
	} {
		compCode(t, muler(&numInt), c)
	}
}

func TestDiver(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"/2",
			[]asm.Fin{
				asm.I().Push().Val(2),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rdx().Val(0),
				asm.I().Div().Rdi(),
				asm.I().Push().Rax(),
			},
			true,
		},
	} {
		compCode(t, diver(&numInt), c)
	}
}

func TestMuldivs(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"1*2",
			[]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Push().Val(2),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Mul().Rdi(),
				asm.I().Push().Rax(),
			},
			true,
		},
	} {
		compCode(t, muldivs(&numInt), c)
	}
}

func TestPrologue(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"",
			[]asm.Fin{
				asm.I().Push().Rbp(),
				asm.I().Mov().Rbp().Rsp(),
				asm.I().Sub().Rsp().Val(16),
			},
			true,
		},
	} {
		compCode(t, prologue(2), c)
	}
}

func TestEpilogue(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"",
			[]asm.Fin{
				asm.I().Mov().Rsp().Rbp(),
				asm.I().Pop().Rbp(),
				asm.I().Ret(),
			},
			true,
		},
	} {
		compCode(t, epilogue, c)
	}
}

func TestPopRax(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"",
			[]asm.Fin{
				asm.I().Pop().Rax(),
			},
			true,
		},
	} {
		compCode(t, popRax, c)
	}
}

func TestFuncWrapper(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"1",
			[]asm.Fin{
				asm.I().Push().Rbp(),
				asm.I().Mov().Rbp().Rsp(),
				asm.I().Sub().Rsp().Val(8 * 26),
				asm.I().Push().Val(1),
				asm.I().Mov().Rsp().Rbp(),
				asm.I().Pop().Rbp(),
				asm.I().Ret(),
			},
			true,
		},
	} {
		compCode(t, funcWrapper(&numInt), c)
	}
}

func wrapInsts(insts []asm.Fin) []asm.Fin {
	return append(append(
		[]asm.Fin{
			asm.I().Push().Rbp(),
			asm.I().Mov().Rbp().Rsp(),
			asm.I().Sub().Rsp().Val(8 * 26),
		},
		insts...),
		[]asm.Fin{
			asm.I().Mov().Rsp().Rbp(),
			asm.I().Pop().Rbp(),
			asm.I().Ret(),
		}...)

}

func TestGenerator(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"1",
			wrapInsts([]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Pop().Rax(),
			}),
			true,
		},
		{

			"1+1",
			wrapInsts([]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Push().Val(1),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Add().Rax().Rdi(),
				asm.I().Push().Rax(),
				asm.I().Pop().Rax(),
			}),
			true,
		},
		{

			"(1+2)",
			wrapInsts([]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Push().Val(2),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Add().Rax().Rdi(),
				asm.I().Push().Rax(),
				asm.I().Pop().Rax(),
			}),
			true,
		},
		{

			"(1-(2))",
			wrapInsts([]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Push().Val(2),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Sub().Rax().Rdi(),
				asm.I().Push().Rax(),
				asm.I().Pop().Rax(),
			}),
			true,
		},
	} {
		compCode(t, Generator(), c)
	}
}
