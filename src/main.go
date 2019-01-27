package main

import (
	"fmt"
	"strconv"

	"github.com/sxarp/c_compiler_go/src/str"

	"github.com/sxarp/c_compiler_go/src/asm"
)

func main() {
	r := ""
	asm := compile(r)
	fmt.Println(asm)
}

func prec(sb *str.Builder) {
	sb.Nr()
	sb.Write(".intel_syntax noprefix")
	sb.Write(".global main")
	sb.Nr()

}
func compile(code string) string {

	var sb str.Builder = str.Builder{}

	prec(&sb)

	if i, err := strconv.Atoi(code); err != nil {
		panic("failed to parse code!")

	} else {
		sb.Write("main:")
		asm.I().Mov().Rax().Val(i).Write(&sb)
		asm.I().Ret().Write(&sb)

	}

	return sb.Str()
}
