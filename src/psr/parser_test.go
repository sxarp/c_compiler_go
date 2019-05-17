package psr

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/tok"
)

func TestBaseParser(t *testing.T) {

	tokens := tok.Tokenize("1+(11-5)")

	var a ast.AST

	for _, c := range []struct {
		p        *Parser
		expected bool
	}{
		{Int, true}, {Int, false}, {Plus, true}, {LPar, true}, {Int, true}, {Minus, true},
		{Int, true}, {RPar, true}, {EOF, true},
	} {

		a, tokens = c.p.Call(tokens)
		ast.CheckAst(t, c.expected, a)

	}

}

func TestAnd(t *testing.T) {
	tokens := tok.Tokenize("1+3")
	var a ast.AST

	p := AndId().And(Int, true).And(Plus, false).And(Int, true).And(EOF, false)
	a, tokens = p.Call(tokens)
	ast.CheckAst(t, true, a)

	if len(tokens) != 0 {
		t.Errorf("Tokens must be consumed, got %v", tokens)

	}

	if i := a.Node(0).Token.Vali(); i != 1 {
		t.Errorf("Expected 1, got %d", i)

	}

	if i := a.Node(1).Token.Vali(); i != 3 {
		t.Errorf("Expected 3, got %d", i)

	}

	tokens = tok.Tokenize("1+1(")
	a, _ = p.Call(tokens)
	ast.CheckAst(t, false, a)

}

func TestOr(t *testing.T) {
	tokens := tok.Tokenize("1+3")
	var a ast.AST

	porm := OrId().Or(Plus).Or(Minus)

	p := AndId().And(Int, true).
		And(&porm, false).
		And(Int, true).And(EOF, false)

	a, tokens = p.Call(tokens)
	ast.CheckAst(t, true, a)

	if len(tokens) != 0 {
		t.Errorf("Tokens must be consumed, got %v", tokens)

	}

	if i := a.Node(0).Token.Vali(); i != 1 {
		t.Errorf("Expected 1, got %d", i)

	}

	if i := a.Node(1).Token.Vali(); i != 3 {
		t.Errorf("Expected 1, got %d", i)

	}

	tokens = tok.Tokenize("1-3")
	plus := OrId().Or(Plus)
	p = AndId().And(Int, true).
		And(&plus, false).
		And(Int, true).And(EOF, false)

	a, _ = p.Call(tokens)
	ast.CheckAst(t, false, a)

}

func TestRecc(t *testing.T) {
	tokens := tok.Tokenize("(((((+)))))")

	parser := OrId()
	par := AndId().And(LPar, true).And(&parser, true).And(RPar, true)
	parser = parser.Or(Plus).Or(&par)
	final := AndId().And(&parser, false).And(EOF, false)

	a, _ := final.Call(tokens)
	ast.CheckAst(t, true, a)

	tokens = tok.Tokenize("(((((+))))")
	a, _ = final.Call(tokens)
	ast.CheckAst(t, false, a)

}

func TestRep(t *testing.T) {
	tokens := tok.Tokenize("1-1+2-4+5")
	var a ast.AST

	porm := OrId().Or(Plus).Or(Minus)
	add := AndId().And(&porm, true).And(Int, true)

	p := AndId().And(Int, true).
		Rep(&add).And(EOF, false)

	a, tokens = p.Call(tokens)
	ast.CheckAst(t, true, a)

	if len(tokens) != 0 {
		t.Errorf("Tokens must be consumed, got %v", tokens)

	}

}
