package asm

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/h"
)

func TestIns(t *testing.T) {
	h.ExpectEq(t, "        ret", I().Ret().Str())
	h.ExpectEq(t, "        mov rax, 42", I().Mov().Rax().Val(42).Str())
	h.ExpectEq(t, "        mov rax, rax", I().Mov().Rax().Rax().Str())
	h.ExpectEq(t, "        add rax, 42", I().Add().Rax().Val(42).Str())
	h.ExpectEq(t, "        sub rax, 42", I().Sub().Rax().Val(42).Str())
	h.ExpectEq(t, "        pop 42", I().Pop().Val(42).Str())
	h.ExpectEq(t, "        push 42", I().Push().Val(42).Str())
	h.ExpectEq(t, "        push rax", I().Push().Rax().Str())
	h.ExpectEq(t, "        push rdi", I().Push().Rdi().Str())
	h.ExpectEq(t, "        pop rax", I().Pop().Rax().Str())
	h.ExpectEq(t, "        pop rdi", I().Pop().Rdi().Str())
	h.ExpectEq(t, "        mul rdi", I().Mul().Rdi().Str())
	h.ExpectEq(t, "        div rdi", I().Div().Rdi().Str())
	h.ExpectEq(t, "        mov rdx, 0", I().Mov().Rdx().Val(0).Str())
	h.ExpectEq(t, "        mov rbp, rsp", I().Mov().Rbp().Rsp().Str())
	h.ExpectEq(t, "        mov rsp, rbp", I().Mov().Rsp().Rbp().Str())
	h.ExpectEq(t, "        mov [rsp], rbp", I().Mov().Rsp().P().Rbp().Str())
	h.ExpectEq(t, "        mov rsp, [rbp]", I().Mov().Rsp().Rbp().P().Str())
	h.ExpectEq(t, "        mov [rdi], rax", I().Mov().Rdi().P().Rax().Str())
	h.ExpectEq(t, "        cmp rdi, rax", I().Cmp().Rdi().Rax().Str())
	h.ExpectEq(t, "        sete al", I().Sete().Al().Str())
	h.ExpectEq(t, "        setne al", I().Setne().Al().Str())
	h.ExpectEq(t, "        movzb rax, al", I().Movzb().Rax().Al().Str())
	h.ExpectEq(t, "        call foo", I().Call("foo").Str())
	h.ExpectEq(t, "        mov rsi, rcx", I().Mov().Rsi().Rcx().Str())
	h.ExpectEq(t, "        mov rcx, rsi", I().Mov().Rcx().Rsi().Str())
	h.ExpectEq(t, "        mov r8, r9", I().Mov().R8().R9().Str())
	h.ExpectEq(t, "        mov r9, r8", I().Mov().R9().R8().Str())
	h.ExpectEq(t, "        mov r10, r10", I().Mov().R10().R10().Str())
	h.ExpectEq(t, "main:", I().Label("main").Str())
	h.ExpectEq(t, "        je .Lend", I().Je(".Lend").Str())
	h.ExpectEq(t, "        jl .Lend", I().Jl(".Lend").Str())
	h.ExpectEq(t, "        jmp .Lend", I().Jmp(".Lend").Str())
	h.ExpectEq(t, "        syscall", I().Sys().Str())
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
		h.ExpectEq(t, comp.expect, comp.lhs.Eq((&comp.rhs)))
	}
}
