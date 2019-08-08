package psr

import (
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/em"
	"github.com/sxarp/c_compiler_go/src/tok"
)

type Parser struct {
	Call func([]tok.Token) (ast.AST, []tok.Token)
}

func TokenTypeToPsr(tt *tok.TokenType) *Parser {
	return &Parser{Call: func(t []tok.Token) (ast.AST, []tok.Token) {

		if len(t) == 0 {
			return ast.Fail, t
		}

		head, tail := tok.Ht(t)

		if head.Is(tt) {
			return ast.AST{Token: &head}, tail

		}

		// For displaying error messages.
		em.EM.Set(tt, &head)

		return ast.Fail, t
	}}
}

func AndIdent() Parser {
	return Parser{
		Call: func(t []tok.Token) (ast.AST, []tok.Token) {
			return ast.AST{}, t
		},
	}
}

func OrIdent() Parser {
	return Parser{
		Call: func(t []tok.Token) (ast.AST, []tok.Token) {
			return ast.Fail, t
		},
	}
}

func (lhsp Parser) And(rhsp *Parser, addNode bool) Parser {

	call := func(t []tok.Token) (ast.AST, []tok.Token) {

		lhs, lhst := lhsp.Call(t)

		if lhs.Fail() {
			return ast.Fail, t

		}

		rhs, rhst := rhsp.Call(lhst)

		if rhs.Fail() {
			return ast.Fail, t

		}

		if addNode {
			lhs.AppendNode(rhs)

		}

		return lhs, rhst

	}

	return Parser{Call: call}

}

func (lhsp Parser) Seq(rhsp *Parser) Parser {
	return lhsp.And(rhsp, false)
}

func (lhsp Parser) Or(rhsp *Parser) Parser {

	call := func(t []tok.Token) (ast.AST, []tok.Token) {

		if lhs, lhst := lhsp.Call(t); !lhs.Fail() {
			return lhs, lhst

		}

		if rhs, rhst := rhsp.Call(t); !rhs.Fail() {
			return rhs, rhst
		}

		return ast.Fail, t

	}

	return Parser{Call: call}
}

func (lhsp Parser) Rep(rhsp *Parser) Parser {

	call := func(lest []tok.Token) (ast.AST, []tok.Token) {

		var node ast.AST

		lhs, lhst := lhsp.Call(lest)
		if lhs.Fail() {
			return lhs, lest
		}

		node, lest = lhs, lhst

		for {

			rhs, rhst := rhsp.Call(lest)

			if rhs.Fail() {
				return node, lest
			}
			lest = rhst
			node.AppendNode(rhs)
		}

	}

	return Parser{Call: call}
}

func (lhsp Parser) Trans(f func(ast.AST) ast.AST) Parser {

	call := func(t []tok.Token) (ast.AST, []tok.Token) {

		lhs, lhst := lhsp.Call(t)

		if lhs.Fail() {
			return lhs, t
		}

		return f(lhs), lhst
	}

	return Parser{Call: call}
}

func (lhsp Parser) SetEval(f ast.Eval) Parser {
	return lhsp.Trans(func(a ast.AST) ast.AST {
		a.SetEval(f)
		return a
	})
}
