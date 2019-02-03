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
