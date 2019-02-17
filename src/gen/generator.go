package gen

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/psr"
)

var orId = psr.OrId
var andId = psr.AndId

const wordSize = 8

func checkNodeCount(nodes []*ast.AST, count int) {
	if l := len(nodes); l != count {
		panic(fmt.Sprintf("The number of nodes must be %d, got %d.", count, l))
	}
}

var numInt = andId().And(psr.Int, true).
	SetEval(func(nodes []*ast.AST, code *asm.Code) {
		i := nodes[0].Token.Vali()
		code.Ins(asm.I().Push().Val(i))
	})

func binaryOperator(term *psr.Parser, operator *psr.Parser, insts []asm.Fin) psr.Parser {
	return andId().And(operator, false).And(term, true).
		SetEval(func(nodes []*ast.AST, code *asm.Code) {
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

func alpaToNum(alphabet rune) int {
	var alphabets = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'w', 'z'}

	for i, c := range alphabets {
		if c == alphabet {
			return i
		}
	}

	panic(fmt.Sprintf("Failed to convert %c into an integer.", alphabet))
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

func prologue(numberOfLocalVars int) psr.Parser {
	return andId().SetEval(func(nodes []*ast.AST, code *asm.Code) {
		code.
			Ins(asm.I().Push().Rbp()).
			Ins(asm.I().Mov().Rbp().Rsp()).
			Ins(asm.I().Sub().Rsp().Val(wordSize * numberOfLocalVars))
	})
}

var epilogue psr.Parser = andId().SetEval(
	func(nodes []*ast.AST, code *asm.Code) {
		code.
			Ins(asm.I().Mov().Rsp().Rbp()).
			Ins(asm.I().Pop().Rbp()).
			Ins(asm.I().Ret())
	})

var popRax psr.Parser = andId().SetEval(func(nodes []*ast.AST, code *asm.Code) { code.Ins(asm.I().Pop().Rax()) })

var lvIdent psr.Parser = andId().And(psr.SinVar, true).SetEval(
	func(nodes []*ast.AST, code *asm.Code) {
		checkNodeCount(nodes, 1)
		offSet := wordSize * (1 + alpaToNum([]rune(nodes[0].Token.Val())[0]))
		code.
			Ins(asm.I().Mov().Rax().Rbp()).
			Ins(asm.I().Sub().Rax().Val(offSet)).
			Ins(asm.I().Push().Rax())

	})

var rvIdent psr.Parser = andId().And(&lvIdent, true).SetEval(
	func(nodes []*ast.AST, code *asm.Code) {
		checkNodeCount(nodes, 1)
		nodes[0].Eval(code)
		code.
			Ins(asm.I().Pop().Rax()).
			Ins(asm.I().Mov().Rax().Rax().P()).
			Ins(asm.I().Push().Rax())
	})

func assigner(lv *psr.Parser, rv *psr.Parser) psr.Parser {
	return andId().And(lv, true).And(psr.Subs, false).And(rv, true).SetEval(
		func(nodes []*ast.AST, code *asm.Code) {
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

func funcWrapper(expr *psr.Parser) psr.Parser {
	pro := prologue(26)
	return andId().And(&pro, true).And(expr, true).And(&epilogue, true)
}

func Generator() psr.Parser {
	num := orId().Or(&numInt)

	term := orId()
	muls := andId()

	adds := addsubs(&muls)
	muls = muldivs(&term)

	parTerm := andId().And(psr.LPar, false).And(&adds, true).And(psr.RPar, false).Trans(ast.PopSingle)
	term = term.Or(&parTerm).Or(&num)

	expr := andId().And(&adds, true).And(&popRax, true)
	return funcWrapper(&expr).And(psr.EOF, false)
}
