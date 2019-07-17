package gen

import (
	"fmt"
	"testing"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/h"
	"github.com/sxarp/c_compiler_go/src/psr"
	"github.com/sxarp/c_compiler_go/src/tok"
	"github.com/sxarp/c_compiler_go/src/tp"
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

func TestSyscaller(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"syscall 1 1 256 2",
			[]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Pop().Rdi(),
				asm.I().Push().Val(256),
				asm.I().Pop().Rsi(),
				asm.I().Push().Val(2),
				asm.I().Pop().Rdx(),
				asm.I().Push().Val(1),
				asm.I().Pop().Rax(),
				asm.I().Sys(),
				asm.I().Push().Rax(),
			},
			true,
			"2",
		},
	} {
		compCode(t, syscaller(&numInt), c)
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
		st.DecOf("0", tp.Int)
		st.DecOf("1", tp.Int)
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

func TestPtrAdder(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"*(a+1)",
			[]asm.Fin{
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(tp.Int.Size() * 1),
				asm.I().Push().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rax().Rax().P(),
				asm.I().Push().Rax(),
				asm.I().Push().Val(1),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rdi().Val(8),
				asm.I().Mul().Rdi(),
				asm.I().Push().Rax(),
				asm.I().Pop().Rdi(),
				asm.I().Pop().Rax(),
				asm.I().Add().Rax().Rdi(),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
	} {
		st := newST()
		st.DecOf("a", tp.Int)
		compCode(t, ptrAdder(st, &numInt), c)
	}

}

func TestLvIdenter(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"a",
			[]asm.Fin{
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(tp.Int.Size() * 1),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
		{

			"b",
			[]asm.Fin{
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(tp.Int.Size() * 2),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
	} {
		st := newST()
		st.DecOf("a", tp.Int)
		st.DecOf("b", tp.Int)
		compCode(t, lvIdenter(st), c)
	}

}

func TestRvIdent(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"a",
			[]asm.Fin{
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(tp.Int.Size() * 1),
				asm.I().Push().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rax().Rax().P(),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
	} {
		st := newST()
		st.DecOf("a", tp.Int)
		lvIdent := lvIdenter(st)
		compCode(t, rvIdenter(&lvIdent), c)
	}

}

func TestPtrDeRefer(t *testing.T) {
	st := newST()
	st.DecOf("a", tp.Int)
	st.DecOf("ap", tp.Int.Ptr())

	for _, c := range []psrTestCase{
		{
			"a",
			[]asm.Fin{
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(st.RefOf("a").Addr),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
		{
			"*ap",
			[]asm.Fin{
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(st.RefOf("ap").Addr),
				asm.I().Push().Rax(),
				asm.I().Pop().Rax(),
				asm.I().Mov().Rax().Rax().P(),
				asm.I().Push().Rax(),
			},
			true,
			"",
		},
	} {
		lvIdent := lvIdenter(st)
		compCode(t, ptrDeRefer(st, &lvIdent), c)
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

func TestVarDeclarer(t *testing.T) {
	varType := tp.Int

	for _, c := range []psrTestCase{
		{
			"int a",
			[]asm.Fin{},
			true,
			"",
		},
		{
			"int *a",
			[]asm.Fin{},
			true,
			"",
		},
		{
			"int **a",
			[]asm.Fin{},
			true,
			"",
		},
	} {
		st := newST()
		compCode(t, varDeclarer(st, &null), c)
		h.ExpectEq(t, true, st.RefOf("a").Type.Eq(varType))
		varType = varType.Ptr()
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
		psr := andIdt().And(&numInt, true).And(&eq, true)
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
		psr := andIdt().And(&numInt, true).And(&neq, true)
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

		{
			"if(0) { return 1} if(1) { return 3} return 4",
			[]asm.Fin{},
			true,
			"3",
		},

		{
			"if(0) { return 1} if(0) { return 3} return 4",
			[]asm.Fin{},
			true,
			"4",
		},
	} {
		prologue := prologuer(newST())
		ret := returner(&numInt)
		iF := ifer(&numInt, &ret)
		ifRet := orIdt().Or(&iF).Or(&ret)
		compCode(t, andIdt().And(&prologue, true).Rep(&ifRet), c)
	}

}

func TestWhiler(t *testing.T) {
	for _, c := range []psrTestCase{
		{
			"while(0) { 1 }",
			[]asm.Fin{
				asm.I().Label("while_begin_0"),
				asm.I().Push().Val(0),
				asm.I().Pop().Rax(),
				asm.I().Cmp().Rax().Val(0),
				asm.I().Je("while_end_0"),
				asm.I().Push().Val(1),
				asm.I().Jmp("while_begin_0"),
				asm.I().Label("while_end_0"),
			},
			true,
			"",
		},
	} {
		compCode(t, whiler(&numInt, &numInt), c)
	}

}

func TestForer(t *testing.T) {
	for _, c := range []psrTestCase{
		{
			"for(0;1;2) { 3 }",
			[]asm.Fin{
				asm.I().Push().Val(0),
				asm.I().Pop().Rax(),
				asm.I().Label("for_begin_0"),
				asm.I().Push().Val(1),
				asm.I().Pop().Rax(),
				asm.I().Cmp().Rax().Val(0),
				asm.I().Je("for_end_0"),
				asm.I().Push().Val(2),
				asm.I().Pop().Rax(),
				asm.I().Push().Val(3),
				asm.I().Jmp("for_begin_0"),
				asm.I().Label("for_end_0"),
			},
			true,
			"",
		},
	} {
		compCode(t, forer(&numInt, &numInt), c)
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
			"int hoge(){return}",
			[]asm.Fin{
				asm.I().Label("hoge"),
				asm.I().Push().Rbp(),
				asm.I().Mov().Rbp().Rsp(),
				asm.I().Sub().Rsp().Val(tp.Int.Size() * 0),
				asm.I().Mov().Rsp().Rbp(),
				asm.I().Pop().Rbp(),
				asm.I().Ret(),
			},
			true,
			"",
		},
		{
			"int main(int a){ return 22}",
			[]asm.Fin{
				asm.I().Label("main"),
				asm.I().Push().Rbp(),
				asm.I().Mov().Rbp().Rsp(),
				asm.I().Sub().Rsp().Val(tp.Int.Size() * 1),
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(tp.Int.Size() * 1),
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
			"int main(){return id(11)}int id(int a){return a}",
			[]asm.Fin{
				asm.I().Label("main"),
				asm.I().Push().Rbp(),
				asm.I().Mov().Rbp().Rsp(),
				asm.I().Sub().Rsp().Val(tp.Int.Size() * 0),
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
				asm.I().Sub().Rsp().Val(tp.Int.Size() * 1),
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(tp.Int.Size() * 1),
				asm.I().Mov().Rax().P().Rdi(),
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(tp.Int.Size() * 1),
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
			"int main(){return sub(11+1, 5)} int sub(int a, int b){return a - b}",
			[]asm.Fin{},
			true,
			"7",
		},
		{
			`
	int main(){return id(1,2,3,4,5,6)}
int id(int a, int b, int c, int d, int e, int f){return a - b + c - d + e - f + 3}
`,
			[]asm.Fin{},
			true,
			"0",
		},
		{
			`
	int main(){return id(1,2,3,4,5,6) - add(1, 2)}
int id(int a, int b, int c, int d, int e, int f){return a - b + c - d + e - f + add(3, 4)}
int add(int a, int b) { return a + b}
`,
			[]asm.Fin{},
			true,
			"1",
		},
		{
			`
	int main(){return id(1, 2, 3)}
int id(int a, int b, int c){return idp(&a, &b, c)}
int idp(int *a, int *b, int c) { return *a + *b + c}
`,
			[]asm.Fin{},
			true,
			"6",
		},
	} {
		body := func(st *SymTable) psr.Parser {
			lvIdent := lvIdenter(st)
			ptrDeRef := ptrDeRefer(st, &lvIdent)

			rvAddr := rvAddrer(&lvIdent)
			rvIdent := rvIdenter(&ptrDeRef)
			rvVal := orIdt().Or(&rvAddr).Or(&rvIdent)

			var caller psr.Parser
			callerOrIdentOrNum := orIdt().Or(&caller).Or(&rvVal).Or(&numInt)
			adds := addsubs(&callerOrIdentOrNum)
			caller = funcCaller(&adds)
			return returner(&adds)
		}
		fd := funcDefiner(body)
		compCode(t, andIdt().Rep(&fd), c)
	}
}

func wrapInsts(insts []asm.Fin) []asm.Fin {
	return append(append(
		[]asm.Fin{
			asm.I().Label("main"),
			asm.I().Push().Rbp(),
			asm.I().Mov().Rbp().Rsp(),
			asm.I().Sub().Rsp().Val(tp.Int.Size() * 0),
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

			"int main(){return 1;}",
			wrapInsts([]asm.Fin{
				asm.I().Push().Val(1),
				asm.I().Pop().Rax(),
			}),
			true,
			"1",
		},
		{

			"int main(){return 1+1;}",
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

			"int main(){return (1+2);}",
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

			"int main(){return (1-(2));}",
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
			"int main() {return}",
			wrapInsts([]asm.Fin{}),
			true,
			"",
		},
	} {
		compCode(t, Generator(), c)
	}
}
