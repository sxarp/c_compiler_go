package psr

import (
	"github.com/sxarp/c_compiler_go/src/tok"
)

type ASTType struct{}

type AST struct {
	nodes []*AST
	token *tok.Token
	atype *ASTType
}

var TFail ASTType = ASTType{}
var Fail = AST{atype: &TFail}

func (a AST) Fail() bool {
	return a.atype == &TFail
}

type Parser struct {
	Call func([]tok.Token) (AST, []tok.Token)
}

func tokenTypeToPsr(tt *tok.TokenType) Parser {
	return Parser{Call: func(t []tok.Token) (AST, []tok.Token) {
		head, tail := tok.Ht(t)
		if head.Is(tt) {
			return AST{token: &head}, tail

		}

		return Fail, t
	},
	}
}

func (p Parser) decorate(f func(AST) AST) Parser {
	return p

}

var Plus Parser = tokenTypeToPsr(&tok.TPlus)
var Minus Parser = tokenTypeToPsr(&tok.TMinus)
var Int Parser = tokenTypeToPsr(&tok.TInt)
var LPar Parser = tokenTypeToPsr(&tok.TLPar)
var RPar Parser = tokenTypeToPsr(&tok.TRPar)
var EOF Parser = tokenTypeToPsr(&tok.TEOF)
