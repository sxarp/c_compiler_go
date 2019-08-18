package gen

import (
	"fmt"
	"testing"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/h"
	"github.com/sxarp/c_compiler_go/src/tok"
	"github.com/sxarp/c_compiler_go/src/tp"
)

type psrTestCase struct {
	rcode    string
	ins      []asm.Fin
	tf       bool
	expected string
}

func compCode(t *testing.T, p Compiler, c psrTestCase) {
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
		ret := i().Ret()
		if !finalInst.Eq(&ret) {
			rhs.Ins(i().Pop().Rax()).Ins(i().Ret())
		}

		ml := i().Label("main")
		if !firstInst.Eq(&ml) {
			rrhs := asm.New()
			rrhs.Ins(i().Label("main"))
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
			[]asm.Fin{i().Push().Val(42)},
			true,
			"42",
		},
		{

			"43",
			[]asm.Fin{i().Push().Val(42)},
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
				i().Push().Val(2),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Add().Rax().Rdi(),
				i().Push().Rax(),
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
				i().Push().Val(2),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Sub().Rax().Rdi(),
				i().Push().Rax(),
			},
			true,
			"",
		},
	} {
		compCode(t, subber(&numInt), c)
	}
}

func TestLter(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"1<2",
			[]asm.Fin{},
			true,
			"1",
		},
		{
			"2<1",
			[]asm.Fin{},
			true,
			"0",
		},
	} {
		lt := lter(&numInt)
		compCode(t, ai().And(&numInt).And(&lt), c)
	}
}

func TestAddsubs(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"1+1",
			[]asm.Fin{
				i().Push().Val(1),
				i().Push().Val(1),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Add().Rax().Rdi(),
				i().Push().Rax(),
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
				i().Push().Val(2),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Mul().Rdi(),
				i().Push().Rax(),
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
				i().Push().Val(2),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Mov().Rdx().Val(0),
				i().Div().Rdi(),
				i().Push().Rax(),
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
				i().Push().Val(1),
				i().Push().Val(2),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Mul().Rdi(),
				i().Push().Rax(),
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
				i().Push().Val(1),
				i().Pop().Rdi(),
				i().Push().Val(256),
				i().Pop().Rsi(),
				i().Push().Val(2),
				i().Pop().Rdx(),
				i().Push().Val(1),
				i().Pop().Rax(),
				i().Sys(),
				i().Push().Rax(),
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
				i().Push().Rbp(),
				i().Mov().Rbp().Rsp(),
				i().Sub().Rsp().Val(16),
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
				i().Mov().Rsp().Rbp(),
				i().Pop().Rbp(),
				i().Ret(),
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
				i().Pop().Rax(),
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
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(tp.Int.Size() * 1),
				i().Push().Rax(),
				i().Push().Val(1),
				i().Pop().Rax(),
				i().Mov().Rdi().Val(8),
				i().Mul().Rdi(),
				i().Push().Rax(),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Add().Rax().Rdi(),
				i().Push().Rax(),
			},
			true,
			"",
		},
		{
			"a[1]",
			[]asm.Fin{
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(tp.Int.Size() * 1),
				i().Push().Rax(),
				i().Push().Val(1),
				i().Pop().Rax(),
				i().Mov().Rdi().Val(8),
				i().Mul().Rdi(),
				i().Push().Rax(),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Add().Rax().Rdi(),
				i().Push().Rax(),
			},
			true,
			"",
		},
		{
			"b[2]",
			[]asm.Fin{
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(tp.Int.Size() * 4),
				i().Push().Rax(),
				i().Push().Val(2),
				i().Pop().Rax(),
				i().Mov().Rdi().Val(8),
				i().Mul().Rdi(),
				i().Push().Rax(),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Add().Rax().Rdi(),
				i().Push().Rax(),
			},
			true,
			"",
		},
	} {
		st := newST()
		st.DecOf("a", tp.Int.Ptr())
		st.DecOf("b", tp.Array(tp.Int, 3))
		lv := lvIdenter(st)
		compCode(t, ptrAdder(st, &lv, &numInt), c)
	}

}

func TestLvIdenter(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"a",
			[]asm.Fin{
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(tp.Int.Size() * 1),
				i().Push().Rax(),
			},
			true,
			"",
		},
		{

			"b",
			[]asm.Fin{
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(tp.Int.Size() * 2),
				i().Push().Rax(),
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
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(tp.Int.Size() * 1),
				i().Push().Rax(),
				i().Pop().Rax(),
				i().Mov().Rax().Rax().P(),
				i().Push().Rax(),
			},
			true,
			"",
		},
	} {
		st := newST()
		st.DecOf("a", tp.Int)
		lvIdent := lvIdenter(st)
		compCode(t, rvIdenter(st, &lvIdent), c)
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
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(st.RefOf("a").Addr),
				i().Push().Rax(),
			},
			true,
			"",
		},
		{
			"*ap",
			[]asm.Fin{
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(st.RefOf("ap").Addr),
				i().Push().Rax(),
				i().Pop().Rax(),
				i().Mov().Rax().Rax().P(),
				i().Push().Rax(),
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
				i().Push().Val(2),
				i().Push().Val(1),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Mov().Rdi().P().Rax(),
				i().Push().Rax(),
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
		compCode(t, varDeclarer(st), c)
		h.ExpectEq(t, true, st.RefOf("a").Type.Eq(varType))
		varType = varType.Ptr()
	}
}

func TestArrayDeclarer(t *testing.T) {
	expectedTypes := []tp.Type{
		tp.Int,
		tp.Array(tp.Int, 3),
		tp.Array(tp.Array(tp.Int.Ptr().Ptr(), 3), 4),
	}

	for i, c := range []psrTestCase{
		{
			"int a",
			[]asm.Fin{},
			true,
			"",
		},
		{
			"int a[3]",
			[]asm.Fin{},
			true,
			"",
		},
		{
			"int **a[3][4]",
			[]asm.Fin{},
			true,
			"",
		},
	} {
		st := newST()
		varDeclare := varDeclarer(st)
		compCode(t, arrayDeclarer(&varDeclare, st), c)

		expectedType := expectedTypes[i]
		h.ExpectEq(t, true, expectedType.Eq(st.Last().Type))
		h.ExpectEq(t, expectedType.Size(), st.Last().Type.Size())
	}
}

func TestEqer(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"2==1",
			[]asm.Fin{
				i().Push().Val(2),
				i().Push().Val(1),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Cmp().Rdi().Rax(),
				i().Sete().Al(),
				i().Movzb().Rax().Al(),
				i().Push().Rax(),
			},
			true,
			"0",
		},
		{
			"2==2",
			[]asm.Fin{
				i().Push().Val(2),
				i().Push().Val(2),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Cmp().Rdi().Rax(),
				i().Sete().Al(),
				i().Movzb().Rax().Al(),
				i().Push().Rax(),
			},
			true,
			"1",
		},
	} {
		eq := eqer(&numInt)
		psr := ai().And(&numInt).And(&eq)
		compCode(t, psr, c)
	}

}

func TestNeqer(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"2!=1",
			[]asm.Fin{
				i().Push().Val(2),
				i().Push().Val(1),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Cmp().Rdi().Rax(),
				i().Setne().Al(),
				i().Movzb().Rax().Al(),
				i().Push().Rax(),
			},
			true,
			"1",
		},
		{
			"2!=2",
			[]asm.Fin{
				i().Push().Val(2),
				i().Push().Val(2),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Cmp().Rdi().Rax(),
				i().Setne().Al(),
				i().Movzb().Rax().Al(),
				i().Push().Rax(),
			},
			true,
			"0",
		},
	} {
		neq := neqer(&numInt)
		psr := ai().And(&numInt).And(&neq)
		compCode(t, psr, c)
	}

}

func TestReturner(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"return",
			[]asm.Fin{
				i().Mov().Rsp().Rbp(),
				i().Pop().Rbp(),
				i().Ret(),
			},
			true,
			"",
		},

		{
			"return 1",
			[]asm.Fin{
				i().Push().Val(1),
				i().Pop().Rax(),
				i().Mov().Rsp().Rbp(),
				i().Pop().Rbp(),
				i().Ret(),
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
		prologue, ret := prologuer(newST()), returner(&numInt)
		compCode(t, ai().And(&prologue, ifer(&numInt, &ret).P(), &ret), c)
	}
}

func TestWhiler(t *testing.T) {
	for _, c := range []psrTestCase{
		{
			"while(0) { 1 }",
			[]asm.Fin{
				i().Label("while_begin_0"),
				i().Push().Val(0),
				i().Pop().Rax(),
				i().Cmp().Rax().Val(0),
				i().Je("while_end_0"),
				i().Push().Val(1),
				i().Jmp("while_begin_0"),
				i().Label("while_end_0"),
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
				i().Push().Val(0),
				i().Pop().Rax(),
				i().Label("for_begin_0"),
				i().Push().Val(1),
				i().Pop().Rax(),
				i().Cmp().Rax().Val(0),
				i().Je("for_end_0"),
				i().Push().Val(3),
				i().Push().Val(2),
				i().Pop().Rax(),
				i().Jmp("for_begin_0"),
				i().Label("for_end_0"),
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
				i().Call("hoge"),
				i().Push().Rax(),
			},
			true,
			"",
		},
		{
			"hoge(1)",
			[]asm.Fin{
				i().Push().Val(1),
				i().Pop().Rax(),
				i().Mov().Rdi().Rax(),
				i().Call("hoge"),
				i().Push().Rax(),
			},
			true,
			"",
		},
		{
			"hoge(1, 2)",
			[]asm.Fin{
				i().Push().Val(2),
				i().Push().Val(1),
				i().Pop().Rax(),
				i().Mov().Rdi().Rax(),
				i().Pop().Rax(),
				i().Mov().Rsi().Rax(),
				i().Call("hoge"),
				i().Push().Rax(),
			},
			true,
			"",
		},
		{
			"hoge(1, 2, 3, 4, 5, 6)",
			[]asm.Fin{
				i().Push().Val(6),
				i().Push().Val(5),
				i().Push().Val(4),
				i().Push().Val(3),
				i().Push().Val(2),
				i().Push().Val(1),
				i().Pop().Rax(),
				i().Mov().Rdi().Rax(),
				i().Pop().Rax(),
				i().Mov().Rsi().Rax(),
				i().Pop().Rax(),
				i().Mov().Rdx().Rax(),
				i().Pop().Rax(),
				i().Mov().Rcx().Rax(),
				i().Pop().Rax(),
				i().Mov().R8().Rax(),
				i().Pop().Rax(),
				i().Mov().R9().Rax(),
				i().Call("hoge"),
				i().Push().Rax(),
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
				i().Label("hoge"),
				i().Push().Rbp(),
				i().Mov().Rbp().Rsp(),
				i().Sub().Rsp().Val(tp.Int.Size() * 0),
				i().Mov().Rsp().Rbp(),
				i().Pop().Rbp(),
				i().Ret(),
			},
			true,
			"",
		},
		{
			"int main(int a){ return 22}",
			[]asm.Fin{
				i().Label("main"),
				i().Push().Rbp(),
				i().Mov().Rbp().Rsp(),
				i().Sub().Rsp().Val(tp.Int.Size() * 1),
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(tp.Int.Size() * 1),
				i().Mov().Rax().P().Rdi(),
				i().Push().Val(22),
				i().Pop().Rax(),
				i().Mov().Rsp().Rbp(),
				i().Pop().Rbp(),
				i().Ret(),
			},
			true,
			"22",
		},
	} {
		body := func(st *SymTable) Compiler { return returner(&numInt) }
		compCode(t, funcDefiner(body), c)
	}
}

func TestFuncDefineAndCall(t *testing.T) {

	for _, c := range []psrTestCase{
		{
			"int main(){return id(11)}int id(int a){return a}",
			[]asm.Fin{
				i().Label("main"),
				i().Push().Rbp(),
				i().Mov().Rbp().Rsp(),
				i().Sub().Rsp().Val(tp.Int.Size() * 0),
				i().Push().Val(11),
				i().Pop().Rax(),
				i().Mov().Rdi().Rax(),
				i().Call("id"),
				i().Push().Rax(),
				i().Pop().Rax(),
				i().Mov().Rsp().Rbp(),
				i().Pop().Rbp(),
				i().Ret(),
				i().Label("id"),
				i().Push().Rbp(),
				i().Mov().Rbp().Rsp(),
				i().Sub().Rsp().Val(tp.Int.Size() * 1),
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(tp.Int.Size() * 1),
				i().Mov().Rax().P().Rdi(),
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(tp.Int.Size() * 1),
				i().Push().Rax(),
				i().Pop().Rax(),
				i().Mov().Rax().Rax().P(),
				i().Push().Rax(),
				i().Pop().Rax(),
				i().Mov().Rsp().Rbp(),
				i().Pop().Rbp(),
				i().Ret(),
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
		body := func(st *SymTable) Compiler {
			lvIdent := lvIdenter(st)
			ptrDeRef := ptrDeRefer(st, &lvIdent)

			rvAddr := rvAddrer(&lvIdent)
			rvIdent := rvIdenter(st, &ptrDeRef)
			rvVal := oi().Or(&rvAddr).Or(&rvIdent)

			var caller Compiler
			callerOrIdentOrNum := oi().Or(&caller).Or(&rvVal).Or(&numInt)
			adds := addsubs(&callerOrIdentOrNum)
			caller = funcCaller(&adds)
			return returner(&adds)
		}
		fd := funcDefiner(body)
		compCode(t, ai().Rep(&fd), c)
	}
}

func wrapInsts(insts []asm.Fin) []asm.Fin {
	return append(append(
		[]asm.Fin{
			i().Label("main"),
			i().Push().Rbp(),
			i().Mov().Rbp().Rsp(),
			i().Sub().Rsp().Val(tp.Int.Size() * 0),
		},
		insts...),
		[]asm.Fin{
			i().Mov().Rsp().Rbp(),
			i().Pop().Rbp(),
			i().Ret(),
		}...)

}

func TestGenerator(t *testing.T) {

	for _, c := range []psrTestCase{
		{

			"int main(){return 1;}",
			wrapInsts([]asm.Fin{
				i().Push().Val(1),
				i().Pop().Rax(),
			}),
			true,
			"1",
		},
		{

			"int main(){return 1+1;}",
			wrapInsts([]asm.Fin{
				i().Push().Val(1),
				i().Push().Val(1),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Add().Rax().Rdi(),
				i().Push().Rax(),
				i().Pop().Rax(),
			}),
			true,
			"2",
		},
		{

			"int main(){return (1+2);}",
			wrapInsts([]asm.Fin{
				i().Push().Val(1),
				i().Push().Val(2),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Add().Rax().Rdi(),
				i().Push().Rax(),
				i().Pop().Rax(),
			}),
			true,
			"3",
		},
		{

			"int main(){return (1-(2));}",
			wrapInsts([]asm.Fin{
				i().Push().Val(1),
				i().Push().Val(2),
				i().Pop().Rdi(),
				i().Pop().Rax(),
				i().Sub().Rax().Rdi(),
				i().Push().Rax(),
				i().Pop().Rax(),
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
