package gen

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/psr"
)

var orId = psr.OrId
var andId = psr.AndId

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

var adder = func(term *psr.Parser) psr.Parser {
	return andId().And(psr.Plus, false).And(term, true).
		SetEval(func(nodes []*ast.AST, code *asm.Code) {
			checkNodeCount(nodes, 1)
			nodes[0].Eval(code)
			code.
				Ins(asm.I().Pop().Rdi()).
				Ins(asm.I().Pop().Rax()).
				Ins(asm.I().Add().Rax().Rdi()).
				Ins(asm.I().Push().Rax())
		})

}

func GenParser() psr.Parser {

	numv := orId().Or(psr.Int)
	num := &numv

	term := orId()
	muls := andId()

	porm := orId().Or(psr.Plus).Or(psr.Minus)
	adder := andId().And(&porm, true).And(&muls, true)
	adds := andId().And(&muls, true).Rep(&adder).Trans(ast.PopSingle)

	mord := orId().Or(psr.Mul).Or(psr.Div)
	muler := andId().And(&mord, true).And(&term, true)
	muls = muls.And(&term, true).Rep(&muler).Trans(ast.PopSingle)

	parTerm := andId().And(psr.LPar, false).And(&adds, true).And(psr.RPar, false).Trans(ast.PopSingle)
	term = term.Or(&parTerm).Or(num)

	return andId().And(&adds, true).And(psr.EOF, false)

}
