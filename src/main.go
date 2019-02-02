package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/sxarp/c_compiler_go/src/asm"
)

func main() {
	r := ""
	asm := compile(r)
	fmt.Println(asm)
}

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

var TFail TokenType = TokenType{
	literal: "FAIL",
}

var Fail Token = Token{
	tt:   &TFail,
	valp: &(TFail.literal),
}

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

var TPlus TokenType = TokenType{literal: "+"}

var TMinus TokenType = TokenType{literal: "-"}

var TInt TokenType = TokenType{regex: regexp.MustCompile("^[0-9]+")}

func compile(tcode string) string {

	acode := asm.New()

	if i, err := strconv.Atoi(tcode); err != nil {
		panic("failed to parse code!")

	} else {
		acode.Main().
			Ins(asm.I().Mov().Rax().Val(i)).
			Ins(asm.I().Ret())

	}

	return acode.Str()
}
