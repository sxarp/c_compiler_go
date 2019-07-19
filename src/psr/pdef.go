package psr

import (
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/tok"
)

var (
	Plus  = TokenTypeToPsr(&tok.TPlus)
	Minus = TokenTypeToPsr(&tok.TMinus)
	Int   = TokenTypeToPsr(&tok.TInt)
	LPar  = TokenTypeToPsr(&tok.TLPar)
	RPar  = TokenTypeToPsr(&tok.TRPar)
	Mul   = TokenTypeToPsr(&tok.TMul)
	Div   = TokenTypeToPsr(&tok.TDiv)
	EOF   = TokenTypeToPsrWOE(&tok.TEOF)
)

// To show informative error messages.
func TokenTypeToPsrWOE(tt *tok.TokenType) *Parser {
	return &Parser{Call: func(t []tok.Token) (ast.AST, []tok.Token) {

		if len(t) == 0 {
			return ast.Fail, t
		}

		head, tail := tok.Ht(t)

		if head.Is(tt) {
			return ast.AST{Token: &head}, tail

		}

		return ast.Fail, t
	}}
}

func GenParser() Parser {
	numv := OrIdent().Or(Int)
	num := &numv

	mul := OrIdent()
	add := OrIdent()
	term := OrIdent()

	termPlusMul := AndIdent().And(&term, true).And(Plus, true).And(&mul, true)
	termMinusMul := AndIdent().And(&term, true).And(Minus, true).And(&mul, true)
	add = add.Or(&termPlusMul).Or(&termMinusMul).Or(&term)

	termMulMul := AndIdent().And(&term, true).And(Mul, true).And(&mul, true)
	termDivMul := AndIdent().And(&term, true).And(Div, true).And(&mul, true)
	mul = mul.Or(&termMulMul).Or(&termDivMul).Or(&add)

	pop := func(a ast.AST) ast.AST { return a.Node(0) }
	parTerm := AndIdent().And(LPar, false).And(&mul, true).And(RPar, false).Trans(pop)
	term = term.Or(&parTerm).Or(num)

	return AndIdent().And(&mul, true).And(EOF, false)

}

func GenParser2() Parser {
	numv := OrIdent().Or(Int)
	num := &numv

	term := OrIdent()
	muls := AndIdent()

	porm := OrIdent().Or(Plus).Or(Minus)
	adder := AndIdent().And(&porm, true).And(&muls, true)
	adds := AndIdent().And(&muls, true).Rep(&adder).Trans(ast.PopSingle)

	mord := OrIdent().Or(Mul).Or(Div)
	muler := AndIdent().And(&mord, true).And(&term, true)
	muls = muls.And(&term, true).Rep(&muler).Trans(ast.PopSingle)

	parTerm := AndIdent().And(LPar, false).And(&adds, true).And(RPar, false).Trans(ast.PopSingle)
	term = term.Or(&parTerm).Or(num)

	return AndIdent().And(&adds, true).And(EOF, false)

}
