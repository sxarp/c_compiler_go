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
func Al() Reg  { return Reg{r: "al"} }
func Rsi() Reg { return Reg{r: "rsi"} }
func Rcx() Reg { return Reg{r: "rcx"} }
func R8() Reg  { return Reg{r: "r8"} }
func R9() Reg  { return Reg{r: "r9"} }
func R10() Reg { return Reg{r: "r10"} }

type Ope struct {
	i string
}

func (i Ope) str() string { return i.i }

// Operators
func OMov() Ope   { return Ope{i: "mov"} }
func OAdd() Ope   { return Ope{i: "add"} }
func OSub() Ope   { return Ope{i: "sub"} }
func ORet() Ope   { return Ope{i: "ret"} }
func OPop() Ope   { return Ope{i: "pop"} }
func OPush() Ope  { return Ope{i: "push"} }
func OMul() Ope   { return Ope{i: "mul"} }
func ODiv() Ope   { return Ope{i: "div"} }
func OCmp() Ope   { return Ope{i: "cmp"} }
func OSete() Ope  { return Ope{i: "sete"} }
func OSetne() Ope { return Ope{i: "setne"} }
func OMovzb() Ope { return Ope{i: "movzb"} }
func OSys() Ope   { return Ope{i: "syscall"} }

type Ins struct {
	ope  Ope
	dest Reg
	srcR Reg
	srcI int

	destP bool
	srcP  bool

	toS   func() string
	label string
}

func (i Ins) str() string {
	sb := str.Builder{}

	if i.label != "" {
		sb.Put(i.label + ":")
		return sb.Str()
	}

	sb.Put("        ")
	sb.Put(i.ope.str())

	if i.toS != nil {
		sb.Put(" ")
		sb.Put(i.toS())
		return sb.Str()
	}

	includeOpe := func(o Ope) bool {
		retv := false
		for _, op := range []func() Ope{OPop, OPush, OMul, ODiv, OSete, OSetne} {
			if o == op() {
				retv = true
			}
		}
		return retv
	}

	if includeOpe(i.ope) {
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

	dest := i.dest.str()
	if i.destP {
		dest = fmt.Sprintf("[%s]", dest)
	}
	sb.Put(" ").Put(dest).Put(",")

	if !i.srcR.nil() {
		src := i.srcR.str()
		if i.srcP {
			src = fmt.Sprintf("[%s]", src)
		}
		sb.Put(" ").Put(src)
		return sb.Str()
	}

	sb.Put(" ").Put(fmt.Sprintf("%d", i.srcI))

	return sb.Str()
}

type Ini Ins

type Oped Ins

type Dested Ins

type Fin Ins

func (i Fin) Str() string { return Ins(i).str() }

func (i *Fin) Eq(rhs *Fin) bool { return i.Str() == rhs.Str() }

func (i Fin) Write(sb *str.Builder) { sb.Write(i.Str()) }

func I() Ini { return Ini{} }

func (i Ini) Ret() Fin {
	i.ope = ORet()
	return Fin(i)
}

func (i Ini) Sys() Fin {
	i.ope = OSys()
	return Fin(i)
}

func (i Ini) Call(name string) Fin {
	i.ope = Ope{i: "call"}
	i.toS = func() string { return name }
	return Fin(i)
}

func (i Ini) Je(name string) Fin {
	i.ope = Ope{i: "je"}
	i.toS = func() string { return name }
	return Fin(i)
}

func (i Ini) Jl(name string) Fin {
	i.ope = Ope{i: "jl"}
	i.toS = func() string { return name }
	return Fin(i)
}

func (i Ini) Jmp(name string) Fin {
	i.ope = Ope{i: "jmp"}
	i.toS = func() string { return name }
	return Fin(i)
}

func (i Ini) Label(name string) Fin {
	i.label = name
	return Fin(i)
}

func toOped(i Ini, o func() Ope) Oped {
	i.ope = o()
	return Oped(i)
}

func opeDested(i Ini, o func() Ope) Dested {
	i.ope = o()
	return Dested(i)
}

func regFin(i Dested, o func() Reg) Fin {
	i.srcR = o()
	return Fin(i)
}

func toDested(i Oped, o func() Reg) Dested {
	i.dest = o()
	return Dested(i)
}

func (i Ini) Mov() Oped   { return toOped(i, OMov) }
func (i Ini) Add() Oped   { return toOped(i, OAdd) }
func (i Ini) Sub() Oped   { return toOped(i, OSub) }
func (i Ini) Cmp() Oped   { return toOped(i, OCmp) }
func (i Ini) Movzb() Oped { return toOped(i, OMovzb) }

func (i Ini) Pop() Dested   { return opeDested(i, OPop) }
func (i Ini) Push() Dested  { return opeDested(i, OPush) }
func (i Ini) Mul() Dested   { return opeDested(i, OMul) }
func (i Ini) Div() Dested   { return opeDested(i, ODiv) }
func (i Ini) Sete() Dested  { return opeDested(i, OSete) }
func (i Ini) Setne() Dested { return opeDested(i, OSetne) }

func (i Dested) P() Dested {
	i.destP = true
	return i
}

func (i Oped) Rax() Dested { return toDested(i, Rax) }
func (i Oped) Rdx() Dested { return toDested(i, Rdx) }
func (i Oped) Rbp() Dested { return toDested(i, Rbp) }
func (i Oped) Rsp() Dested { return toDested(i, Rsp) }
func (i Oped) Rdi() Dested { return toDested(i, Rdi) }
func (i Oped) Al() Dested  { return toDested(i, Al) }
func (i Oped) Rsi() Dested { return toDested(i, Rsi) }
func (i Oped) Rcx() Dested { return toDested(i, Rcx) }
func (i Oped) R8() Dested  { return toDested(i, R8) }
func (i Oped) R9() Dested  { return toDested(i, R9) }
func (i Oped) R10() Dested { return toDested(i, R10) }

func (i Dested) Val(s int) Fin {
	i.srcI = s
	return Fin(i)
}

func (i Dested) Rax() Fin { return regFin(i, Rax) }
func (i Dested) Rdi() Fin { return regFin(i, Rdi) }
func (i Dested) Rdx() Fin { return regFin(i, Rdx) }
func (i Dested) Rbp() Fin { return regFin(i, Rbp) }
func (i Dested) Rsp() Fin { return regFin(i, Rsp) }
func (i Dested) Al() Fin  { return regFin(i, Al) }
func (i Dested) Rsi() Fin { return regFin(i, Rsi) }
func (i Dested) Rcx() Fin { return regFin(i, Rcx) }
func (i Dested) R8() Fin  { return regFin(i, R8) }
func (i Dested) R9() Fin  { return regFin(i, R9) }
func (i Dested) R10() Fin { return regFin(i, R10) }

func (i Fin) P() Fin {
	i.srcP = true
	return i
}
