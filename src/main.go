package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/em"
	"github.com/sxarp/c_compiler_go/src/gen"
	"github.com/sxarp/c_compiler_go/src/tok"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	b, err := ioutil.ReadAll(r)
	fatal(err)

	src := string(b)
	asm := compile(src)

	fmt.Println(asm)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func compile(tcode string) string {
	tokens := tok.Tokenize(tcode)

	insts := asm.New()
	if ast, rem := gen.Generator().Call(tokens); len(rem) == 0 {
		ast.Eval(insts)

	} else {
		panic(em.EM.Message())
	}

	return asm.NewBuilder(insts).Str()
}
