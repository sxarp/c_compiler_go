package gen

import (
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/psr"
)

var orId = psr.OrId
var andId = psr.AndId

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
