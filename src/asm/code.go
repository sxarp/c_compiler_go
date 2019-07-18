package asm

import (
	"github.com/sxarp/c_compiler_go/src/str"
)

type Code interface {
	Ins(...Fin) Code
	For(func(int, Fin))
}

type Insts struct {
	insts []Fin
}

func New() *Insts {
	return &Insts{}
}

type Builder struct {
	code  *str.Builder
	insts *Insts
}

func NewBuilder(is *Insts) *Builder {
	var sb str.Builder = str.Builder{}
	sb.Nr()
	sb.Write(".intel_syntax noprefix")
	sb.Write(".global main")
	sb.Nr()

	return &Builder{code: &sb, insts: is}
}

func (is *Insts) Ins(fs ...Fin) Code {
	is.insts = append(is.insts, fs...)
	return is
}

func (is *Insts) For(f func(int, Fin)) {
	for c, i := range is.insts {
		f(c, i)
	}
}

func (is *Insts) Concat(c Code) { c.For(func(i int, f Fin) { is.Ins(f) }) }

func (is *Insts) ForEachInst(f func(Fin)) {

	for _, is := range is.insts {
		f(is)
	}
}

func (b *Builder) Str() string {
	b.insts.ForEachInst(func(i Fin) {
		b.code.Write(i.str())
	})

	return b.code.Str()
}

func (is *Insts) Eq(rhs *Insts) bool {
	if len(is.insts) != len(rhs.insts) {
		return false
	}

	eq := true
	for i, li := range is.insts {
		eq = eq && li.Eq(&(rhs.insts[i]))

	}

	return eq
}
