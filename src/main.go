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

func (b *strBuilder) put(s string) *strBuilder {
	b.b.WriteString(s)

	return b
}

func (b *strBuilder) write(s string) {
	b.put(fmt.Sprintf("%s\n", s))
}

func (b *strBuilder) str() string {
	return b.b.String()
}

func (b *strBuilder) nr() {
	b.write("")
}

type Reg struct {
	r string
}

func (r Reg) str() string {
	return r.r

}

func (r Reg) nil() bool {
	return r.r == ""

}

func Rax() Reg {
	return Reg{r: "rax"}
}

type Ope struct {
	i string
}

func (i Ope) str() string {
	return i.i

}

func OMov() Ope {
	return Ope{i: "mov"}
}

func ORet() Ope {
	return Ope{i: "ret"}
}

type Ins struct {
	ope  Ope
	dest Reg
	srcR Reg
	srcI int
}

func (i Ins) str() string {

	var sb strBuilder = strBuilder{}
	sb.put("        ")
	sb.put(i.ope.str())

	if i.dest.nil() {
		return sb.str()
	}

	sb.put(" ").put(i.dest.str()).put(",")

	if !i.srcR.nil() {
		sb.put(" ").put(i.srcR.str())
		return sb.str()
	}

	sb.put(" ").put(fmt.Sprintf("%v", i.srcI))

	return sb.str()
}

type Ini struct {
	i Ins
}

type Oped struct {
	i Ins
}

type Dested struct {
	i Ins
}

type Fin struct {
	i Ins
}

func (i Fin) str() string {
	return i.i.str()
}

func (i Fin) write(sb *strBuilder) {
	sb.write(i.str())
}

func I() Ini {
	return Ini{i: Ins{}}

}

func (i Ini) Ret() Fin {
	i.i.ope = ORet()
	return Fin{i: i.i}

}

func (i Ini) Mov() Oped {
	i.i.ope = OMov()
	return Oped{i: i.i}
}

func (i Oped) Rax() Dested {
	i.i.dest = Rax()
	return Dested{i: i.i}
}

func (i Dested) Val(s int) Fin {
	i.i.srcI = s
	return Fin{i: i.i}
}

func (i Dested) Rax() Fin {
	i.i.srcR = Rax()
	return Fin{i: i.i}
}

func compile(r string) string {

	var sb strBuilder = strBuilder{}

	sb.nr()
	sb.write(".intel_syntax noprefix")
	sb.write(".global main")
	sb.nr()
	sb.write("main:")

	I().Mov().Rax().Val(42).write(&sb)
	I().Ret().write(&sb)

	return sb.str()
}
