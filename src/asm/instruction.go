package asm

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/str"
)

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

func OAdd() Ope {
	return Ope{i: "add"}
}

func OSub() Ope {
	return Ope{i: "sub"}
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

	sb := str.Builder{}

	sb.Put("        ")
	sb.Put(i.ope.str())

	if i.dest.nil() {
		return sb.Str()
	}

	sb.Put(" ").Put(i.dest.str()).Put(",")

	if !i.srcR.nil() {
		sb.Put(" ").Put(i.srcR.str())
		return sb.Str()
	}

	sb.Put(" ").Put(fmt.Sprintf("%v", i.srcI))

	return sb.Str()
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

func (i Fin) Write(sb *str.Builder) {
	sb.Write(i.str())
}

func I() Ini {
	return Ini{i: Ins{}}

}

func (i Ini) Ret() Fin {
	i.i.ope = ORet()
	return Fin{i: i.i}

}

// Initial Instractions
func (i Ini) Mov() Oped {
	i.i.ope = OMov()
	return Oped{i: i.i}
}

func (i Ini) Add() Oped {
	i.i.ope = OAdd()
	return Oped{i: i.i}
}

func (i Ini) Sub() Oped {
	i.i.ope = OSub()
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
