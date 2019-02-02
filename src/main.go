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

	acode := asm.New()

	if i, err := strconv.Atoi(pcode); err != nil {
		panic("failed to parse code!")

	} else {
		acode.Main().
			Ins(asm.I().Mov().Rax().Val(i)).
			Ins(asm.I().Ret())

	}

	return acode.Str()
}
