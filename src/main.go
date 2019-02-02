package main

import (
	"fmt"
	"strconv"

	"github.com/sxarp/c_compiler_go/src/asm"
)

func main() {
	r := ""
	asm := compile(r)
	fmt.Println(asm)
}

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
