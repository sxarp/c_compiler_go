package asm

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/str"
)

type Reg struct {
	r string
}

func (r Reg) str() string { return r.r }
func (r Reg) nil() bool   { return r.r == "" }

// Registers
func Rax() Reg { return Reg{r: "rax"} }
func Rdi() Reg { return Reg{r: "rdi"} }
func Rdx() Reg { return Reg{r: "rdx"} }
func Rbp() Reg { return Reg{r: "rbp"} }
func Rsp() Reg { return Reg{r: "rsp"} }

type Ope struct {
	i string
}

func (i Ope) str() string { return i.i }

// Operators
func OMov() Ope  { return Ope{i: "mov"} }
func OAdd() Ope  { return Ope{i: "add"} }
func OSub() Ope  { return Ope{i: "sub"} }
func ORet() Ope  { return Ope{i: "ret"} }
func OPop() Ope  { return Ope{i: "pop"} }
func OPush() Ope { return Ope{i: "push"} }
func OMul() Ope  { return Ope{i: "mul"} }
func ODiv() Ope  { return Ope{i: "div"} }

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

	if i.ope == OPop() || i.ope == OPush() || i.ope == OMul() || i.ope == ODiv() {
		sb.Put(" ")
		if !i.srcR.nil() {
			sb.Put(i.srcR.str())
		} else {
			sb.Put(fmt.Sprintf("%d", i.srcI))

		}
		return sb.Str()
	}

	if i.dest.nil() {
		return sb.Str()
	}

	sb.Put(" ").Put(i.dest.str()).Put(",")

	if !i.srcR.nil() {
		sb.Put(" ").Put(i.srcR.str())
		return sb.Str()
	}

	sb.Put(" ").Put(fmt.Sprintf("%d", i.srcI))

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

func (i Fin) str() string { return i.i.str() }

func (lhs *Fin) Eq(rhs *Fin) bool { return lhs.str() == rhs.str() }

func (i Fin) Write(sb *str.Builder) { sb.Write(i.str()) }

func I() Ini { return Ini{i: Ins{}} }

func (i Ini) Ret() Fin {
	i.i.ope = ORet()
	return Fin{i: i.i}

}

func toOped(i Ins, o func() Ope) Oped {
	i.ope = o()
	return Oped{i: i}
}

func opeDested(i Ins, o func() Ope) Dested {
	i.ope = o()
	return Dested{i: i}
}

func regFin(i Ins, o func() Reg) Fin {
	i.srcR = o()
	return Fin{i: i}
}

func toDested(i Ins, o func() Reg) Dested {
	i.dest = o()
	return Dested{i: i}
}

func (i Ini) Mov() Oped { return toOped(i.i, OMov) }
func (i Ini) Add() Oped { return toOped(i.i, OAdd) }
func (i Ini) Sub() Oped { return toOped(i.i, OSub) }

func (i Ini) Pop() Dested  { return opeDested(i.i, OPop) }
func (i Ini) Push() Dested { return opeDested(i.i, OPush) }
func (i Ini) Mul() Dested  { return opeDested(i.i, OMul) }
func (i Ini) Div() Dested  { return opeDested(i.i, ODiv) }

func (i Oped) Rax() Dested { return toDested(i.i, Rax) }
func (i Oped) Rdx() Dested { return toDested(i.i, Rdx) }
func (i Oped) Rbp() Dested { return toDested(i.i, Rbp) }
func (i Oped) Rsp() Dested { return toDested(i.i, Rsp) }

func (i Dested) Val(s int) Fin {
	i.i.srcI = s
	return Fin{i: i.i}
}

func (i Dested) Rax() Fin { return regFin(i.i, Rax) }
func (i Dested) Rdi() Fin { return regFin(i.i, Rdi) }
func (i Dested) Rbp() Fin { return regFin(i.i, Rbp) }
func (i Dested) Rsp() Fin { return regFin(i.i, Rsp) }
