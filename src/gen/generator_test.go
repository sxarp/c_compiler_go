package gen

import (
	"fmt"
	"testing"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/h"
	"github.com/sxarp/c_compiler_go/src/psr"
	"github.com/sxarp/c_compiler_go/src/tok"
)

type psrTestCase struct {
	rcode    string
	ins      []asm.Fin
	tf       bool
	expected string
}

func compCode(t *testing.T, p psr.Parser, c psrTestCase) {
	t.Helper()
	lhs := asm.New()
	var finalInst asm.Fin
	var firstInst asm.Fin
	for i, ins := range c.ins {
		lhs.Ins(ins)
		finalInst = ins
		if i == 0 {
			firstInst = ins
		}
	}

	rhs := asm.New()
	a, rem := p.Call(tok.Tokenize(c.rcode))
	if len(rem) > 1 {
		t.Errorf("failed to parse")
	}

	a.Eval(rhs)

	if len(c.ins) > 0 {
		h.ExpectEq(t, c.tf, lhs.Eq(rhs))
		if c.tf && !lhs.Eq(rhs) {
			lhsasm := asm.NewBuilder(lhs)
			rhsasm := asm.NewBuilder(rhs)

			fmt.Println("Expected:----------------")
			fmt.Println(lhsasm.Str())
			fmt.Println("Got:---------------------")
			fmt.Println(rhsasm.Str())
		}
	}

	if c.expected != "" {
		ret := asm.I().Ret()
		if !finalInst.Eq(&ret) {
			rhs.Ins(asm.I().Pop().Rax()).Ins(asm.I().Ret())
		}

		ml := asm.I().Label("main")
		if !firstInst.Eq(&ml) {
			rrhs := asm.New()
			rrhs.Ins(asm.I().Label("main"))
			rrhs.Concat(rhs)
			rhs = rrhs
		}
		execInstComp(t, c.expected, rhs)
	}
}

func execInstComp(t *testing.T, expected string, insts *asm.Insts) {
	t.Helper()
	if gotValue := h.ExecCode(t, asm.NewBuilder(insts).Str(),
		"../../tmp", "insts"); gotValue != expected {
		t.Errorf("Expected %s, got %s.", expected, gotValue)
	}
}

func TestNunInt(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"42",
			[]asm.Fin{asm.I().Push().Val(42)},
			true,
			"42",
		},
		{

			"43",
			[]asm.Fin{asm.I().Push().Val(42)},
			false,
			"43",
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
			"",
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
			"",
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
			"2",
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
			"",
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
			"",
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
			"2",
		},
	} {
		compCode(t, muldivs(&numInt), c)
	}
}

func TestProloguer(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"",
			[]asm.Fin{
				asm.I().Push().Rbp(),
				asm.I().Mov().Rbp().Rsp(),
				asm.I().Sub().Rsp().Val(16),
			},
			true,
			"",
		},
	} {
		st := newST()
		st.RefOf("0")
		st.RefOf("1")
		compCode(t, prologuer(st), c)
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
			"",
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
			"",
		},
	} {
		compCode(t, popRax, c)
	}
}

func TestLvIdenter(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"a",
			[]asm.Fin{
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(wordSize * 1),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},

		{

			"z",
			[]asm.Fin{
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(wordSize * 1),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
	} {
		compCode(t, lvIdenter(newST()), c)
	}

}

func TestRvIdent(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"a",
			[]asm.Fin{
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(wordSize * 1),
				asm.I().Push().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rax().Rax().P(),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
	} {
		lvIdent := lvIdenter(newST())
		compCode(t, rvIdenter(&lvIdent), c)
	}

}

func TestAssigner(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"1=2",
			[]asm.Fin{
				asm.I().Push().Val(2),
				asm.I().Push().Val(1),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rdi().P().Rax(),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
	} {
		compCode(t, assigner(&numInt, &numInt), c)
	}

}

func TestEqer(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"2==1",
			[]asm.Fin{
				asm.I().Push().Val(2),
				asm.I().Push().Val(1),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Cmp().Rdi().Rax(),
				asm.I().Sete().Al(),
				asm.I().Movzb().Rax().Al(),
				asm.I().Push().Rax(),
			},
			true,
			"0",
		},
		{
			"2==2",
			[]asm.Fin{
				asm.I().Push().Val(2),
				asm.I().Push().Val(2),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Cmp().Rdi().Rax(),
				asm.I().Sete().Al(),
				asm.I().Movzb().Rax().Al(),
				asm.I().Push().Rax(),
			},
			true,
			"1",
		},
	} {
		eq := eqer(&numInt)
		psr := andId().And(&numInt, true).And(&eq, true)
		compCode(t, psr, c)
	}

}

func TestNeqer(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"2!=1",
			[]asm.Fin{
				asm.I().Push().Val(2),
				asm.I().Push().Val(1),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Cmp().Rdi().Rax(),
				asm.I().Setne().Al(),
				asm.I().Movzb().Rax().Al(),
				asm.I().Push().Rax(),
			},
			true,
			"1",
		},
		{
			"2!=2",
			[]asm.Fin{
				asm.I().Push().Val(2),
				asm.I().Push().Val(2),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Cmp().Rdi().Rax(),
				asm.I().Setne().Al(),
				asm.I().Movzb().Rax().Al(),
				asm.I().Push().Rax(),
			},
			true,
			"0",
		},
	} {
		neq := neqer(&numInt)
		psr := andId().And(&numInt, true).And(&neq, true)
		compCode(t, psr, c)
	}

}

func TestReturner(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"return",
			[]asm.Fin{
				asm.I().Mov().Rsp().Rbp(),
				asm.I().Pop().Rbp(),
				asm.I().Ret(),
			},
			true,
			"",
		},

		{
			"return 1",
			[]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rsp().Rbp(),
				asm.I().Pop().Rbp(),
				asm.I().Ret(),
			},
			true,
			"",
		},
	} {
		compCode(t, returner(&numInt), c)
	}
}

func TestIfer(t *testing.T) {
	for _, c := range []psrTestCase{
		{
			"if(0) { return 1} return 2",
			[]asm.Fin{},
			true,
			"2",
		},
		{
			"if(1) { return 1} return 2",
			[]asm.Fin{},
			true,
			"1",
		},
	} {
		prologue := prologuer(newST())
		ret := returner(&numInt)
		iF := ifer(&numInt, &ret)
		compCode(t, andId().And(&prologue, true).And(&iF, true).And(&ret, true), c)
	}

}

func TestFuncCaller(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"hoge()",
			[]asm.Fin{
				asm.I().Call("hoge"),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
		{
			"hoge(1)",
			[]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rdi().Rax(),
				asm.I().Call("hoge"),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
		{
			"hoge(1, 2)",
			[]asm.Fin{
				asm.I().Push().Val(2),
				asm.I().Push().Val(1),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rdi().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rsi().Rax(),
				asm.I().Call("hoge"),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
		{
			"hoge(1, 2, 3, 4, 5, 6)",
			[]asm.Fin{
				asm.I().Push().Val(6),
				asm.I().Push().Val(5),
				asm.I().Push().Val(4),
				asm.I().Push().Val(3),
				asm.I().Push().Val(2),
				asm.I().Push().Val(1),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rdi().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rsi().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rdx().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rcx().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().R8().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().R9().Rax(),
				asm.I().Call("hoge"),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
	} {
		compCode(t, funcCaller(&numInt), c)
	}

}

func TestFuncDefiner(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"hoge(){return}",
			[]asm.Fin{
				asm.I().Label("hoge"),
				asm.I().Push().Rbp(),
				asm.I().Mov().Rbp().Rsp(),
				asm.I().Sub().Rsp().Val(wordSize * 0),
				asm.I().Mov().Rsp().Rbp(),
				asm.I().Pop().Rbp(),
				asm.I().Ret(),
			},
			true,
			"",
		},
		{
			"main(a){ return 22}",
			[]asm.Fin{
				asm.I().Label("main"),
				asm.I().Push().Rbp(),
				asm.I().Mov().Rbp().Rsp(),
				asm.I().Sub().Rsp().Val(wordSize * 1),
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(wordSize * 1),
				asm.I().Mov().Rax().P().Rdi(),
				asm.I().Push().Val(22),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rsp().Rbp(),
				asm.I().Pop().Rbp(),
				asm.I().Ret(),
			},
			true,
			"22",
		},
	} {
		body := func(st *SymTable) psr.Parser { return returner(&numInt) }
		compCode(t, funcDefiner(body), c)
	}
}

func TestFuncDefineAndCall(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"main(){return id(11)}id(a){return a}",
			[]asm.Fin{
				asm.I().Label("main"),
				asm.I().Push().Rbp(),
				asm.I().Mov().Rbp().Rsp(),
				asm.I().Sub().Rsp().Val(wordSize * 0),
				asm.I().Push().Val(11),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rdi().Rax(),
				asm.I().Call("id"),
				asm.I().Push().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rsp().Rbp(),
				asm.I().Pop().Rbp(),
				asm.I().Ret(),
				asm.I().Label("id"),
				asm.I().Push().Rbp(),
				asm.I().Mov().Rbp().Rsp(),
				asm.I().Sub().Rsp().Val(wordSize * 1),
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(wordSize * 1),
				asm.I().Mov().Rax().P().Rdi(),
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(wordSize * 1),
				asm.I().Push().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rax().Rax().P(),
				asm.I().Push().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rsp().Rbp(),
				asm.I().Pop().Rbp(),
				asm.I().Ret(),
			},
			true,
			"11",
		},
		{
			"main(){return sub(11+1, 5)}sub(a, b){return a - b}",
			[]asm.Fin{},
			true,
			"7",
		},
		{
			`
	main(){return id(1,2,3,4,5,6)}
id(a, b, c, d, e, f){return a - b + c - d + e - f + 3}
`,
			[]asm.Fin{},
			true,
			"0",
		},
		{
			`
	main(){return id(1,2,3,4,5,6) - add(1, 2)}
id(a, b, c, d, e, f){return a - b + c - d + e - f + add(3, 4)}
add(a, b) { return a + b}
`,
			[]asm.Fin{},
			true,
			"1",
		},
	} {
		body := func(st *SymTable) psr.Parser {
			lvIdent := lvIdenter(st)
			rvIdent := rvIdenter(&lvIdent)
			var caller psr.Parser
			callerOrIdentOrNum := orId().Or(&caller).Or(&rvIdent).Or(&numInt)
			adds := addsubs(&callerOrIdentOrNum)
			caller = funcCaller(&adds)
			return returner(&adds)
		}
		fd := funcDefiner(body)
		compCode(t, andId().Rep(&fd), c)
	}
}

func wrapInsts(insts []asm.Fin) []asm.Fin {
	return append(append(
		[]asm.Fin{
			asm.I().Label("main"),
			asm.I().Push().Rbp(),
			asm.I().Mov().Rbp().Rsp(),
			asm.I().Sub().Rsp().Val(wordSize * 0),
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

			"main(){return 1;}",
			wrapInsts([]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Pop().Rax(),
			}),
			true,
			"1",
		},
		{

			"main(){return 1+1;}",
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
			"2",
		},
		{

			"main(){return (1+2);}",
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
			"3",
		},
		{

			"main(){return (1-(2));}",
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
			"255",
		},

		{
			"main() {return}",
			wrapInsts([]asm.Fin{}),
			true,
			"",
		},
	} {
		compCode(t, Generator(), c)
	}
}
