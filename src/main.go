package main

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/gen"
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
	if ast, rem := gen.Generator().Call(tokens); len(rem) == 0 {
		ast.Eval(acode.Main())

	} else {
		panic("Failed to parse!")
	}
	return acode.Str()
}
