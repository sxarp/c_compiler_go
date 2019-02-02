package tok

import (
	"testing"
)

func expectToken(t *testing.T, tt TokenType, inputStr, expectedVal, expectedStr string) {
	tok, outputStr := tt.match(inputStr)
	outputVal := tok.val()

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

}

func expectTokens(t *testing.T, tker func(string) []Token, s string, expectedTokens []string) {
	tokens := tker(s)

	tokenVals := make([]string, 0)

	for i, token := range tokens {
		tokenVals = append(tokenVals, token.val())

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
	GenTokenizer := func(tt []TokenType) func(string) []Token {
		return func(s string) []Token {
			return tokenizer(tt, s)
		}
	}

	expectTokens(t, GenTokenizer([]TokenType{TPlus}), "", []string{"EOF"})
	expectTokens(t, GenTokenizer([]TokenType{TPlus}), "+",
		[]string{"+", "EOF"})
	expectTokens(t, GenTokenizer([]TokenType{TPlus}), "++",
		[]string{"+", "+", "EOF"})
	expectTokens(t, GenTokenizer([]TokenType{TPlus, TMinus}), "+-",
		[]string{"+", "-", "EOF"})
	expectTokens(t, GenTokenizer([]TokenType{TPlus, TMinus, TInt}), "+23-11",
		[]string{"+", "23", "-", "11", "EOF"})

	expectTokens(t, Tokenizer, "+23-11", []string{"+", "23", "-", "11", "EOF"})
}
