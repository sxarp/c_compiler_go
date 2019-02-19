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
		if c == ' ' || c == '\n' {
			continue
		}

		rets = rets + string(c)

	}
	return rets

}

func compile(tcode string) string {
	pcode := preprocess(tcode)
	tokens := tok.Tokenize(pcode)

	insts := asm.New()
	if ast, rem := gen.Generator().Call(tokens); len(rem) == 0 {
		ast.Eval(insts)

	} else {
		panic("Failed to parse!")
	}

	return asm.NewBuilder(insts).Main().Str()
}
