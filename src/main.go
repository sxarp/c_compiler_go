package main

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/em"
	"github.com/sxarp/c_compiler_go/src/gen"
	"github.com/sxarp/c_compiler_go/src/tok"
)

func main() {
	r := ""
	asm := compile(r)
	fmt.Println(asm)
}

func compile(tcode string) string {
	tokens := tok.Tokenize(tcode)

	insts := asm.New()
	if ast, rem := gen.Generator().Call(tokens); len(rem) == 0 {
		ast.Eval(insts)

	} else {
		panic(em.EM.Message())
	}

	return asm.NewBuilder(insts).Main().Str()
}
