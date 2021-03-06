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
	Row  int
	Col  int
}

func (t *Token) Val() string { return *(t.valp) }

func (t *Token) Vali() int {
	if t.tt.vali == nil {
		panic("Called Vali when vali is nil!")

	}

	return t.tt.vali(t.Val())
}

func (t *Token) Is(tt *TokenType) bool { return t.tt == tt }

func (t *Token) setRC(row, col int) { t.Row, t.Col = row, col }

var TFail = TokenType{literal: "FAIL"}
var Fail = Token{tt: &TFail, valp: &(TFail.literal)}

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

func (tt *TokenType) Str() string {
	if tt.regex != nil {
		return tt.regex.String()
	}
	return tt.literal
}

var TEOF = TokenType{literal: "EOF"}
var EOF = Token{tt: &TEOF, valp: &(TEOF.literal)}

func tokenizeLine(tokens []Token, line string, row int,
	tokenTypes []*TokenType, lineLen int) []Token {
	skipToken := TokenType{regex: regexp.MustCompile(`^[\s]+`)}

	for line != "" {
		t := Fail
		_, line = skipToken.match(line)

		for _, tt := range tokenTypes {

			col := lineLen - len(line)
			if t, line = tt.match(line); t != Fail {
				t.setRC(row, col)
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

	return tokens
}

func tokenize(tokenTypes []*TokenType, lines string) []Token {
	tokens := make([]Token, 0)
	scanner := bufio.NewScanner(strings.NewReader(lines))

	row := 0

	for scanner.Scan() {
		row++
		line := scanner.Text()
		lineLen := len(line)
		tokens = tokenizeLine(tokens, line, row, tokenTypes, lineLen)
	}

	return append(tokens, EOF)
}

var (
	TPlus  = TokenType{literal: "+"}
	TMinus = TokenType{literal: "-"}
	TInt   = TokenType{regex: regexp.MustCompile("^[0-9]+"),
		vali: func(s string) int {
			if i, err := strconv.Atoi(s); err != nil {
				panic(fmt.Sprintf("Failed to convert %s to int!", s))
			} else {
				return i
			}
		},
	}
	TVar   = TokenType{regex: regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*")}
	TLPar  = TokenType{literal: "("}
	TRPar  = TokenType{literal: ")"}
	TMul   = TokenType{literal: "*"}
	TDiv   = TokenType{literal: "/"}
	TSubs  = TokenType{literal: "="}
	TSemi  = TokenType{literal: ";"}
	TEq    = TokenType{literal: "=="}
	TNeq   = TokenType{literal: "!="}
	TLt    = TokenType{literal: "<"}
	TGt    = TokenType{literal: ">"}
	TCom   = TokenType{literal: ","}
	TRet   = TokenType{literal: "return"}
	TLBrc  = TokenType{literal: "{"}
	TRBrc  = TokenType{literal: "}"}
	TLSbr  = TokenType{literal: "["}
	TRSbr  = TokenType{literal: "]"}
	TIf    = TokenType{literal: "if"}
	TWhile = TokenType{literal: "while"}
	TFor   = TokenType{literal: "for"}
	TIntd  = TokenType{literal: "int"}
	TAmp   = TokenType{literal: "&"}
	TSys   = TokenType{literal: "syscall"}
)

var TokenTypes = []*TokenType{&TEq, &TNeq, &TLt, &TGt, &TSubs, &TPlus, &TMinus, &TInt, &TLPar, &TRPar,
	&TMul, &TRet, &TDiv, &TIf, &TWhile, &TFor, &TIntd, &TSys, &TVar, &TSemi, &TCom, &TLBrc, &TRBrc,
	&TLSbr, &TRSbr, &TAmp}

func Tokenize(s string) []Token { return tokenize(TokenTypes, s) }

func Ht(ts []Token) (Token, []Token) {
	if len(ts) == 0 {
		panic("Empty input tokens!")

	}

	return ts[0], ts[1:]
}
