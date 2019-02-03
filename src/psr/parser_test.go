package psr

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/tok"
)

func checkAst(t *testing.T, success bool, a AST) {
	if success && a.Fail() {
		t.Errorf("Expected to succeed in parsing, got failed.")

	}

	if !success && !a.Fail() {
		t.Errorf("Expected to fail at parsing, got %v.", a.token.Val())

	}
}

func TestBaseParser(t *testing.T) {

	tokens := tok.Tokenize("1+(11-5)")

	var ast AST

	for _, c := range []struct {
		p        Parser
		expected bool
	}{
		{Int, true}, {Int, false}, {Plus, true}, {LPar, true}, {Int, true}, {Minus, true},
		{Int, true}, {RPar, true}, {EOF, true},
	} {

		ast, tokens = c.p.Call(tokens)
		checkAst(t, c.expected, ast)

	}

}

func TestAnd(t *testing.T) {
	tokens := tok.Tokenize("1+3")
	var ast AST

	p := AndId.And(Int, true).And(Plus, false).And(Int, true).And(EOF, false)
	ast, tokens = p.Call(tokens)
	checkAst(t, true, ast)

	if len(tokens) != 0 {
		t.Errorf("Tokens must be consumed, got %v", tokens)

	}

	if i := ast.nodes[0].token.Vali(); i != 1 {
		t.Errorf("Expected 1, got %d", i)

	}

	if i := ast.nodes[1].token.Vali(); i != 3 {
		t.Errorf("Expected 1, got %d", i)

	}

	tokens = tok.Tokenize("1+1(")
	ast, tokens = p.Call(tokens)
	checkAst(t, false, ast)

}

func TestOr(t *testing.T) {
	tokens := tok.Tokenize("1+3")
	var ast AST

	p := AndId.And(Int, true).
		And(OrId.Or(Plus).Or(Minus), false).
		And(Int, true).And(EOF, false)

	ast, tokens = p.Call(tokens)
	checkAst(t, true, ast)

	if len(tokens) != 0 {
		t.Errorf("Tokens must be consumed, got %v", tokens)

	}

	if i := ast.nodes[0].token.Vali(); i != 1 {
		t.Errorf("Expected 1, got %d", i)

	}

	if i := ast.nodes[1].token.Vali(); i != 3 {
		t.Errorf("Expected 1, got %d", i)

	}

	tokens = tok.Tokenize("1-3")
	p = AndId.And(Int, true).
		And(OrId.Or(Plus), false).
		And(Int, true).And(EOF, false)

	ast, tokens = p.Call(tokens)
	checkAst(t, false, ast)

}
