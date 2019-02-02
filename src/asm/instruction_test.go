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
}
