package asm

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/h"
)

func TestIns(t *testing.T) {
	h.Expects(t, "        ret", I().Ret().str())
	h.Expects(t, "        mov rax, 42", I().Mov().Rax().Val(42).str())
	h.Expects(t, "        mov rax, rax", I().Mov().Rax().Rax().str())
	h.Expects(t, "        add rax, 42", I().Add().Rax().Val(42).str())
	h.Expects(t, "        sub rax, 42", I().Sub().Rax().Val(42).str())
	h.Expects(t, "        pop 42", I().Pop().Val(42).str())
	h.Expects(t, "        push 42", I().Push().Val(42).str())
	h.Expects(t, "        push rax", I().Push().Rax().str())
	h.Expects(t, "        push rdi", I().Push().Rdi().str())
	h.Expects(t, "        pop rax", I().Pop().Rax().str())
	h.Expects(t, "        pop rdi", I().Pop().Rdi().str())
}

func TestFinEq(t *testing.T) {
	for _, comp := range []struct {
		lhs    Fin
		rhs    Fin
		expect bool
	}{
		{I().Sub().Rax().Val(42), I().Sub().Rax().Val(42), true},
		{I().Sub().Rax().Val(42), I().Sub().Rax().Val(43), false},
		{I().Push().Rax(), I().Push().Rax(), true},
		{I().Push().Rax(), I().Push().Rdi(), false},
	} {
		h.Expectt(t, comp.expect, comp.lhs.Eq((&comp.rhs)))
	}
}
