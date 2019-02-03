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

func (a *AST) appendNode(an AST) {
	a.nodes = append(a.nodes, &an)
}

var TFail ASTType = ASTType{}
var Fail = AST{atype: &TFail}

func (a AST) Fail() bool {
	return a.atype == &TFail
}

type Parser struct {
	Call func([]tok.Token) (AST, []tok.Token)
}

func tokenTypeToPsr(tt *tok.TokenType) *Parser {
	return &Parser{Call: func(t []tok.Token) (AST, []tok.Token) {
		if len(t) == 0 {
			return Fail, t
		}

		head, tail := tok.Ht(t)
		if head.Is(tt) {
			return AST{token: &head}, tail

		}

		return Fail, t
	},
	}
}

func (p Parser) decorate(decorator func(AST) AST) Parser {
	call := func(t []tok.Token) (AST, []tok.Token) {
		ast, token := p.Call(t)

		ast = decorator(ast)

		return ast, token
	}

	return Parser{Call: call}

}

func AndId() Parser {
	return Parser{
		Call: func(t []tok.Token) (AST, []tok.Token) {
			return AST{}, t
		},
	}
}

func OrId() Parser {
	return Parser{
		Call: func(t []tok.Token) (AST, []tok.Token) {
			return Fail, t
		},
	}
}

func (lhsp Parser) And(rhsp *Parser, addNode bool) Parser {

	call := func(t []tok.Token) (AST, []tok.Token) {

		lhs, lhst := lhsp.Call(t)

		if lhs.Fail() {
			return Fail, t

		}

		rhs, rhst := rhsp.Call(lhst)

		if rhs.Fail() {
			return Fail, t

		}

		if addNode {
			lhs.appendNode(rhs)

		}

		return lhs, rhst

	}

	return Parser{Call: call}

}

func (lhsp Parser) Or(rhsp *Parser) Parser {

	call := func(t []tok.Token) (AST, []tok.Token) {

		if lhs, lhst := lhsp.Call(t); !lhs.Fail() {
			return lhs, lhst

		}

		if rhs, rhst := rhsp.Call(t); !rhs.Fail() {
			return rhs, rhst
		}

		return Fail, t

	}

	return Parser{Call: call}
}
