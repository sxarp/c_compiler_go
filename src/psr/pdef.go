package psr

import (
	"github.com/sxarp/c_compiler_go/src/tok"
)

var Plus *Parser = tokenTypeToPsr(&tok.TPlus)
var Minus *Parser = tokenTypeToPsr(&tok.TMinus)
var Int *Parser = tokenTypeToPsr(&tok.TInt)
var LPar *Parser = tokenTypeToPsr(&tok.TLPar)
var RPar *Parser = tokenTypeToPsr(&tok.TRPar)
var Mul *Parser = tokenTypeToPsr(&tok.TMul)
var Div *Parser = tokenTypeToPsr(&tok.TDiv)
var EOF *Parser = tokenTypeToPsr(&tok.TEOF)

func GenParser() Parser {
	numv := OrId().Or(Int)
	num := &numv

	mul := OrId()
	add := OrId()
	term := OrId()

	termPlusMul := AndId().And(&term, true).And(Plus, true).And(&mul, true)
	termMinusMul := AndId().And(&term, true).And(Minus, true).And(&mul, true)
	add = add.Or(&termPlusMul).Or(&termMinusMul).Or(&term)

	termMulMul := AndId().And(&term, true).And(Mul, true).And(&mul, true)
	termDivMul := AndId().And(&term, true).And(Div, true).And(&mul, true)
	mul = mul.Or(&termMulMul).Or(&termDivMul).Or(&add)

	parTerm := AndId().And(LPar, false).And(&mul, true).And(RPar, false).Trans(func(a AST) AST { return *(a.nodes[0]) })
	term = term.Or(&parTerm).Or(num)

	return AndId().And(&mul, true).And(EOF, false)

}
