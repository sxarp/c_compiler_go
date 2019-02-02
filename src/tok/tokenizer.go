package tok

import (
	"fmt"
	"regexp"
)

type TokenType struct {
	regex   *regexp.Regexp
	literal string
	vali    func(string) int
}

type Token struct {
	tt   *TokenType
	valp *string
}

func (t *Token) val() string {
	return *(t.valp)

}

var TFail TokenType = TokenType{literal: "FAIL"}
var Fail Token = Token{tt: &TFail, valp: &(TFail.literal)}

func (tt *TokenType) match(s string) (Token, string) {
	if tt.literal != "" {
		tll := len(tt.literal)
		if len(s) >= tll && s[:tll] == tt.literal {

			return Token{tt: tt, valp: &tt.literal}, s[tll:]

		}
	}

	if tt.regex != nil {
		foundStr := tt.regex.FindString(s)
		if foundStr != "" {
			fsl := len(foundStr)
			return Token{tt: tt, valp: &foundStr}, s[fsl:]
		}
	}

	return Fail, s
}

var TEOF TokenType = TokenType{literal: "EOF"}
var EOF Token = Token{tt: &TEOF, valp: &(TEOF.literal)}

func tokenizer(tokenTypes []TokenType, s string) []Token {
	tokens := make([]Token, 0)

	for s != "" {
		t := Fail

		for _, tt := range tokenTypes {

			if t, s = tt.match(s); t != Fail {
				break

			}
		}

		if t == Fail {
			errsl := 10
			if len(s) < errsl {
				errsl = len(s)

			}
			panic(fmt.Sprintf("Failed to tokenize:[%s].", s[:errsl]))
		}

		tokens = append(tokens, t)
	}

	return append(tokens, EOF)
}

var TPlus TokenType = TokenType{literal: "+"}
var TMinus TokenType = TokenType{literal: "-"}
var TInt TokenType = TokenType{regex: regexp.MustCompile("^[0-9]+")}

var TokenTypes = []TokenType{TPlus, TMinus, TInt}

func Tokenizer(s string) []Token {
	return tokenizer(TokenTypes, s)

}
