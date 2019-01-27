package main

import (
	"bytes"
	"fmt"
)

func main() {
	r := ""
	asm := compile(r)
	fmt.Println(asm)
}

type strBuilder struct {
	b bytes.Buffer
}

func (b *strBuilder) write(s string) {
	b.b.WriteString(fmt.Sprintf("%s\n", s))
}

func (b *strBuilder) str() string {
	return b.b.String()
}

func (b *strBuilder) nr() {
	b.write("")
}

func compile(r string) string {

	var sb strBuilder = strBuilder{}

	sb.nr()
	sb.write(".intel_syntax noprefix")
	sb.write(".global main")
	sb.nr()
	sb.write("main:")
	sb.write("        mov rax, 42")
	sb.write("        ret")

	return sb.str()
}
