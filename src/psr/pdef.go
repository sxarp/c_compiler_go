package psr

import (
	"github.com/sxarp/c_compiler_go/src/ast"
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
var Semi *Parser = tokenTypeToPsr(&tok.TSemi)
var Subs *Parser = tokenTypeToPsr(&tok.TSubs)
var Var *Parser = tokenTypeToPsr(&tok.TVar)

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

	pop := func(a ast.AST) ast.AST { return a.Node(0) }
	parTerm := AndId().And(LPar, false).And(&mul, true).And(RPar, false).Trans(pop)
	term = term.Or(&parTerm).Or(num)

	return AndId().And(&mul, true).And(EOF, false)

}

func GenParser2() Parser {
	numv := OrId().Or(Int)
	num := &numv

	term := OrId()
	muls := AndId()

	porm := OrId().Or(Plus).Or(Minus)
	adder := AndId().And(&porm, true).And(&muls, true)
	adds := AndId().And(&muls, true).Rep(&adder).Trans(ast.PopSingle)

	mord := OrId().Or(Mul).Or(Div)
	muler := AndId().And(&mord, true).And(&term, true)
	muls = muls.And(&term, true).Rep(&muler).Trans(ast.PopSingle)

	parTerm := AndId().And(LPar, false).And(&adds, true).And(RPar, false).Trans(ast.PopSingle)
	term = term.Or(&parTerm).Or(num)

	return AndId().And(&adds, true).And(EOF, false)

}
