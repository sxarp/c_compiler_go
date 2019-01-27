package main

import (
	"fmt"
)

func main() {
	r := ""
	asm := compile(r)
	fmt.Println(asm)
}

func compile(r string) string {

	ret := `
.intel_syntax noprefix
.global main

main:
        mov rax, 42
        ret
`
	return ret
}
