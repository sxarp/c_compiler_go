package tok

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/h"
)

func expectToken(t *testing.T, tt TokenType, inputStr, expectedVal, expectedStr string) {
	tok, outputStr := tt.match(inputStr)
	outputVal := tok.Val()

	if outputVal != expectedVal || outputStr != expectedStr {
		t.Errorf("Expected %s, %s for input %s, got %s, %s.", expectedVal, expectedStr, inputStr, outputVal, outputStr)

	}

}

func TestTokens(t *testing.T) {
	expectToken(t, TPlus, "+", "+", "")
	expectToken(t, TPlus, "+1234", "+", "1234")
	expectToken(t, TPlus, "1+1234", "FAIL", "1+1234")
	expectToken(t, TPlus, "", "FAIL", "")

	expectToken(t, TMinus, "-", "-", "")
	expectToken(t, TMinus, "-123", "-", "123")
	expectToken(t, TMinus, "1-1234", "FAIL", "1-1234")
	expectToken(t, TMinus, "", "FAIL", "")

	expectToken(t, TInt, "1", "1", "")
	expectToken(t, TInt, "123", "123", "")
	expectToken(t, TInt, "123x", "123", "x")
	expectToken(t, TInt, "x123", "FAIL", "x123")
	expectToken(t, TInt, "", "FAIL", "")

	tk, _ := TInt.match("10")
	h.Expecti(t, 10, tk.Vali())

	expectToken(t, TSinVar, "a", "a", "")
	expectToken(t, TSinVar, "za", "z", "a")
	expectToken(t, TSinVar, "A", "FAIL", "A")
}

func expectTokens(t *testing.T, tker func(string) []Token, s string, expectedTokens []string) {
	tokens := tker(s)

	tokenVals := make([]string, 0)

	for i, token := range tokens {
		tokenVals = append(tokenVals, token.Val())

		expectedVal := "nil"
		if len(expectedTokens) > i {
			expectedVal = expectedTokens[i]

		}

		if tokenVals[len(tokenVals)-1] != expectedVal {
			t.Errorf("Expected %v, got %v.", expectedTokens[:i+1], tokenVals)

		}
	}
}

func TestTokenizer(t *testing.T) {
	GenTokenizer := func(tt []*TokenType) func(string) []Token {
		return func(s string) []Token {
			return tokenize(tt, s)
		}
	}

	expectTokens(t, GenTokenizer([]*TokenType{&TPlus}), "", []string{"EOF"})
	expectTokens(t, GenTokenizer([]*TokenType{&TPlus}), "+",
		[]string{"+", "EOF"})
	expectTokens(t, GenTokenizer([]*TokenType{&TPlus}), "++",
		[]string{"+", "+", "EOF"})
	expectTokens(t, GenTokenizer([]*TokenType{&TPlus, &TMinus}), "+-",
		[]string{"+", "-", "EOF"})
	expectTokens(t, GenTokenizer([]*TokenType{&TPlus, &TMinus, &TInt}), "+23-11",
		[]string{"+", "23", "-", "11", "EOF"})
	expectTokens(t, Tokenize, "+23-11", []string{"+", "23", "-", "11", "EOF"})
	expectTokens(t, Tokenize, "a+b", []string{"a", "+", "b", "EOF"})
}

func TestHt(t *testing.T) {
	tokens := Tokenize("+-")

	head, tail := Ht(tokens)
	h.Expects(t, "+", head.Val())

	head, tail = Ht(tail)
	h.Expects(t, "-", head.Val())

	head, tail = Ht(tail)
	h.Expects(t, "EOF", head.Val())

	if len(tail) != 0 {
		t.Errorf("Expected empty slice, got %v.", t)

	}

}

func TestIs(t *testing.T) {
	tokens := Tokenize("+-1")

	h.Expectt(t, true, tokens[0].Is(&TPlus))
	h.Expectt(t, false, tokens[0].Is(&TMinus))
	h.Expectt(t, true, tokens[1].Is(&TMinus))
	h.Expectt(t, true, tokens[2].Is(&TInt))
	h.Expectt(t, true, tokens[3].Is(&TEOF))

}
