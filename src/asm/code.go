package asm

import (
	"github.com/sxarp/c_compiler_go/src/str"
)

type code struct {
	code *str.Builder
}

func New() code {
	var sb str.Builder = str.Builder{}
	sb.Nr()
	sb.Write(".intel_syntax noprefix")
	sb.Write(".global main")
	sb.Nr()

	return code{code: &sb}
}

func (c *code) Main() *code {
	c.code.Write("main:")

	return c
}

func (c *code) Ins(i Fin) *code {
	c.code.Write(i.str())

	return c
}

func (c *code) Str() string {
	return c.code.Str()

}
