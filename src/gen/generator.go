package gen

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/psr"
)

var orId = psr.OrId
var andId = psr.AndId
var null = andId()

const wordSize = 8

func checkNodeCount(nodes []*ast.AST, count int) {
	if l := len(nodes); l != count {
		panic(fmt.Sprintf("The number of nodes must be %d, got %d.", count, l))
	}
}

var numInt = andId().And(psr.Int, true).
	SetEval(func(nodes []*ast.AST, code asm.Code) {
		i := nodes[0].Token.Vali()
		code.Ins(asm.I().Push().Val(i))
	})

func binaryOperator(term *psr.Parser, operator *psr.Parser, insts []asm.Fin) psr.Parser {
	return andId().And(operator, false).And(term, true).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 1)
			nodes[0].Eval(code)

			code.
				Ins(asm.I().Pop().Rdi()).
				Ins(asm.I().Pop().Rax())

			for _, i := range insts {
				code.Ins(i)
			}

			code.Ins(asm.I().Push().Rax())
		})

}

func adder(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Plus, []asm.Fin{asm.I().Add().Rax().Rdi()})
}

func subber(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Minus, []asm.Fin{asm.I().Sub().Rax().Rdi()})
}

func addsubs(term *psr.Parser) psr.Parser {
	add, sub := adder(term), subber(term)
	addsub := orId().Or(&add).Or(&sub)
	return andId().And(term, true).Rep(&addsub).Trans(ast.PopSingle)
}

func muler(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Mul, []asm.Fin{asm.I().Mul().Rdi()})
}

func diver(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Div, []asm.Fin{asm.I().Mov().Rdx().Val(0), asm.I().Div().Rdi()})
}

func muldivs(term *psr.Parser) psr.Parser {
	mul, div := muler(term), diver(term)
	muldiv := orId().Or(&mul).Or(&div)
	return andId().And(term, true).Rep(&muldiv).Trans(ast.PopSingle)
}

func eqer(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Eq, []asm.Fin{
		asm.I().Cmp().Rdi().Rax(),
		asm.I().Sete().Al(),
		asm.I().Movzb().Rax().Al()})
}

func neqer(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Neq, []asm.Fin{
		asm.I().Cmp().Rdi().Rax(),
		asm.I().Setne().Al(),
		asm.I().Movzb().Rax().Al()})
}

func eqneqs(term *psr.Parser) psr.Parser {
	eq, neq := eqer(term), neqer(term)
	eqneq := orId().Or(&eq).Or(&neq)
	return andId().And(term, true).Rep(&eqneq).Trans(ast.PopSingle)
}

func prologuer(st *SymTable) psr.Parser {
	return andId().SetEval(func(nodes []*ast.AST, code asm.Code) {
		code.
			Ins(asm.I().Push().Rbp()).
			Ins(asm.I().Mov().Rbp().Rsp()).
			Ins(asm.I().Sub().Rsp().Val(wordSize * st.Count()))
	})
}

var epilogue psr.Parser = andId().SetEval(
	func(nodes []*ast.AST, code asm.Code) {
		code.
			Ins(asm.I().Mov().Rsp().Rbp()).
			Ins(asm.I().Pop().Rbp()).
			Ins(asm.I().Ret())
	})

var popRax psr.Parser = andId().SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Pop().Rax()) })

func lvIdenter(st *SymTable) psr.Parser {
	return andId().And(psr.Var, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 1)
			offSet := wordSize * (1 + st.RefOf(nodes[0].Token.Val()))
			code.
				Ins(asm.I().Mov().Rax().Rbp()).
				Ins(asm.I().Sub().Rax().Val(offSet)).
				Ins(asm.I().Push().Rax())

		})
}

func rvIdenter(lvIdent *psr.Parser) psr.Parser {
	return andId().And(lvIdent, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 1)
			nodes[0].Eval(code)
			code.
				Ins(asm.I().Pop().Rax()).
				Ins(asm.I().Mov().Rax().Rax().P()).
				Ins(asm.I().Push().Rax())
		})

}

func assigner(lv *psr.Parser, rv *psr.Parser) psr.Parser {
	return andId().And(lv, true).And(psr.Subs, false).And(rv, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)

			nodes[1].Eval(code) // Evaluate right value
			nodes[0].Eval(code) // Evaluate left value
			code.
				Ins(asm.I().Pop().Rdi()).           // load lv to rdi
				Ins(asm.I().Pop().Rax()).           // load rv to rax
				Ins(asm.I().Mov().Rdi().P().Rax()). // mv rax to [lv]
				Ins(asm.I().Push().Rax())

		})
}

func returner(term *psr.Parser) psr.Parser {
	retval := andId().And(term, true).And(&popRax, true)
	ret := orId().Or(&retval).Or(&null)
	return andId().And(psr.Ret, false).And(&ret, true).And(&epilogue, true)
}

var argRegs = []asm.Fin{
	asm.I().Mov().Rdi().Rax(),
	asm.I().Mov().Rsi().Rax(),
	asm.I().Mov().Rdx().Rax(),
	asm.I().Mov().Rcx().Rax(),
	asm.I().Mov().R8().Rax(),
	asm.I().Mov().R9().Rax(),
}

func funcCaller(term *psr.Parser) psr.Parser {
	funcName := andId().And(psr.Var, true).
		SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Call(nodes[0].Token.Val())) })

	commed := andId().And(psr.Com, false).And(term, true).Trans(ast.PopSingle)
	argvs := andId().And(term, true).Rep(&commed).SetEval(func(nodes []*ast.AST, code asm.Code) {
		if len(nodes) > 6 {
			panic("too many arguments")
		}

		// Evaluate args from right to left and push into the stack.
		for i := range nodes {
			nodes[len(nodes)-i-1].Eval(code)
		}

		for i := range nodes {
			code.Ins(asm.I().Pop().Rax()).Ins(argRegs[i])
		}
	})

	args := orId().Or(&argvs).Or(&null)

	return andId().And(&funcName, true).And(psr.LPar, false).And(&args, true).And(psr.RPar, false).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[1].Eval(code)
			nodes[0].Eval(code)
			code.Ins(asm.I().Push().Rax())
		})
}

func funcDefiner() psr.Parser {
	funcName := andId().And(psr.Var, true).
		SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Label(nodes[0].Token.Val())) })

	commed := andId().And(psr.Com, false).And(psr.Var, true).Trans(ast.PopSingle)
	argvs := andId().And(psr.Var, true).Rep(&commed).SetEval(func(nodes []*ast.AST, code asm.Code) {
		if len(nodes) > 6 {
			panic("too many arguments")
		}
	})

	args := orId().Or(&argvs).Or(&null)

	return andId().And(&funcName, true).And(psr.LPar, false).And(&args, true).And(psr.RPar, false)
}

func funcWrapper(expr *psr.Parser, st *SymTable) psr.Parser {
	prologue := prologuer(st)
	return andId().And(&prologue, true).And(expr, true).And(&epilogue, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 3)

			insts := code.(*asm.Insts)
			insts.Ins(asm.I().Label("main"))

			bottom := asm.New()
			nodes[1].Eval(bottom)
			nodes[2].Eval(bottom)

			// Evaluate the prologue AST afterwards so that the symbol table can emit
			// the correct number of variables declared, which is used by the prologue code
			// to determine the size of the stack to allocate for the variables.
			nodes[0].Eval(insts)

			insts.Concat(bottom)
		})
}

func Generator() psr.Parser {
	st := newST()
	lvIdent := lvIdenter(st)
	rvIdent := rvIdenter(&lvIdent)

	num := orId().Or(&numInt).Or(&rvIdent)

	var term, muls, adds, expr, eqs psr.Parser

	eqs = eqneqs(&adds)
	adds = addsubs(&muls)
	muls = muldivs(&term)

	parTerm := andId().And(psr.LPar, false).And(&adds, true).And(psr.RPar, false).Trans(ast.PopSingle)
	term = orId().Or(&num).Or(&parTerm)

	assign := assigner(&lvIdent, &expr)
	expr = orId().Or(&assign).Or(&eqs)
	semi := andId().And(&expr, true).And(psr.Semi, false)

	line := andId().And(&semi, true).And(&popRax, true)
	ret := returner(&semi)

	retLine := orId().Or(&ret).Or(&line)
	retLines := andId().Rep(&retLine)

	return funcWrapper(&retLines, st).And(psr.EOF, false)
}
