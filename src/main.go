package main

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/tok"
)

func main() {
	r := ""
	asm := compile(r)
	fmt.Println(asm)
}

func preprocess(s string) string {
	rets := ""

	for _, c := range s {
		if c == ' ' {
			continue
		}

		rets = rets + string(c)

	}
	return rets

}

func compile(tcode string) string {
	pcode := preprocess(tcode)
	tokens := tok.Tokenize(pcode)
	acode := asm.New()

	var token tok.Token

	acode.Main()

Loop:
	for {
		token, tokens = tok.Ht(tokens)

		switch {
		case token.Is(&tok.TInt):
			acode.Ins(asm.I().Mov().Rax().Val(token.Vali()))
		case token.Is(&tok.TPlus):
			token, tokens = tok.Ht(tokens)
			if token.Is(&tok.TInt) {
				acode.Ins(asm.I().Add().Rax().Val(token.Vali()))

			} else {
				panic("Expected Int token!")

			}
		case token.Is(&tok.TMinus):
			token, tokens = tok.Ht(tokens)
			if token.Is(&tok.TInt) {
				acode.Ins(asm.I().Sub().Rax().Val(token.Vali()))

			} else {
				panic("Expected Int token!")

			}
		case token.Is(&tok.TEOF):
			acode.Ins(asm.I().Ret())
			break Loop
		default:
			panic("Invalid token!")

		}

	}
	return acode.Str()
}
