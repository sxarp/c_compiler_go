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

	termPlusMul := AndId().And(&term, true).And(Plus, false).And(&mul, true)
	termMinusMul := AndId().And(&term, true).And(Minus, false).And(&mul, true)
	add = add.Or(&termPlusMul).Or(&termMinusMul).Or(&term)

	termMulMul := AndId().And(&term, true).And(Mul, false).And(&mul, true)
	termDivMul := AndId().And(&term, true).And(Div, false).And(&mul, true)
	mul = mul.Or(&termMulMul).Or(&termDivMul).Or(&add)

	parTerm := AndId().And(LPar, true).And(&mul, true).And(RPar, true)
	term = term.Or(&parTerm).Or(num)

	return AndId().And(&mul, false).And(EOF, false)

}
