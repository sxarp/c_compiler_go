package tok

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

	return t.tt.vali(t.Val())
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

func tokenize(tokenTypes []*TokenType, lines string) []Token {
	tokens := make([]Token, 0)
	scanner := bufio.NewScanner(strings.NewReader(lines))

	for scanner.Scan() {
		line := scanner.Text()

		for line != "" {
			t := Fail

			for _, tt := range tokenTypes {

				if t, line = tt.match(line); t != Fail {
					break

				}
			}

			if t == Fail {
				errsl := 10
				if len(line) < errsl {
					errsl = len(line)

				}
				panic(fmt.Sprintf("Failed to tokenize:[%s].", line[:errsl]))
			}

			tokens = append(tokens, t)
		}
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
var TVar TokenType = TokenType{regex: regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*")}
var TLPar TokenType = TokenType{literal: "("}
var TRPar TokenType = TokenType{literal: ")"}
var TMul TokenType = TokenType{literal: "*"}
var TDiv TokenType = TokenType{literal: "/"}
var TSubs TokenType = TokenType{literal: "="}
var TSemi TokenType = TokenType{literal: ";"}

var TokenTypes = []*TokenType{&TSubs, &TPlus, &TMinus, &TInt, &TLPar, &TRPar, &TMul, &TDiv, &TVar, &TSemi}

func Tokenize(s string) []Token {
	return tokenize(TokenTypes, s)
}

func Ht(ts []Token) (Token, []Token) {
	if len(ts) == 0 {
		panic("Empty input tokens!")

	}

	return ts[0], ts[1:]
}
