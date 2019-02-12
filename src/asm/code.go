package asm

import (
	"github.com/sxarp/c_compiler_go/src/str"
)

type Code struct {
	code  *str.Builder
	insts []Fin
}

func New() Code {
	var sb str.Builder = str.Builder{}
	sb.Nr()
	sb.Write(".intel_syntax noprefix")
	sb.Write(".global main")
	sb.Nr()

	return Code{code: &sb}
}

func (c *Code) Main() *Code {
	c.code.Write("main:")

	return c
}

func (c *Code) Ins(i Fin) *Code {
	c.insts = append(c.insts, i)

	return c
}

func (c *Code) ForEachInst(f func(Fin)) {

	for _, i := range c.insts {
		f(i)
	}

}

func (c *Code) Str() string {
	c.ForEachInst(func(i Fin) {
		c.code.Write(i.str())
	})

	return c.code.Str()

}

func (lhs *Code) Eq(rhs *Code) bool {
	if len(lhs.insts) != len(rhs.insts) {
		return false
	}

	eq := true
	for i, li := range lhs.insts {
		eq = eq && li.Eq(&(rhs.insts[i]))

	}

	return eq

}
