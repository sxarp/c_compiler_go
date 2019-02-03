package tok

import (
	"fmt"
	"regexp"
	"strconv"
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

func (t *Token) Val() string {
	return *(t.valp)

}

func (t *Token) Vali() int {
	if t.tt.vali == nil {
		panic("Called Vali when vali is nil!")

	}

	return t.tt.vali(*(t.valp))

}

func (t *Token) Is(tt *TokenType) bool {
	return t.tt == tt

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

func tokenize(tokenTypes []*TokenType, s string) []Token {
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
var TInt TokenType = TokenType{regex: regexp.MustCompile("^[0-9]+"),
	vali: func(s string) int {
		if i, err := strconv.Atoi(s); err != nil {
			panic(fmt.Sprintf("Failed to convert %s to int!", s))
		} else {
			return i
		}
	},
}
var TLPar TokenType = TokenType{literal: "("}
var TRPar TokenType = TokenType{literal: ")"}

var TokenTypes = []*TokenType{&TPlus, &TMinus, &TInt, &TLPar, &TRPar}

func Tokenize(s string) []Token {
	return tokenize(TokenTypes, s)

}

func Ht(ts []Token) (Token, []Token) {
	if len(ts) == 0 {
		panic("Empty input tokens!")

	}

	return ts[0], ts[1:]

}
