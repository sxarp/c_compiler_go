package asm

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/h"
)

func TestEq(t *testing.T) {
	for _, comp := range []struct {
		lhs    []Fin
		rhs    []Fin
		expect bool
	}{
		{
			[]Fin{I().Ret()}, []Fin{I().Ret()}, true,
		},
		{
			[]Fin{I().Mov().Rax().Val(23), I().Ret()},
			[]Fin{I().Mov().Rax().Val(23), I().Ret()}, true,
		},

		{
			[]Fin{I().Mov().Rax().Val(23), I().Ret()},
			[]Fin{I().Mov().Rax().Val(24), I().Ret()}, false,
		},

		{
			[]Fin{I().Mov().Rax().Val(23), I().Ret()},
			[]Fin{I().Mov().Rax().Val(23), I().Pop().Rax()}, false,
		},
	} {
		lhs := Insts{}
		rhs := Insts{}

		for i := range comp.lhs {
			lhs.Ins(comp.lhs[i])
			rhs.Ins(comp.rhs[i])
		}

		h.ExpectEq(t, comp.expect, lhs.Eq(&rhs))

	}

}

func TestFor(t *testing.T) {
	lhs := New()
	lhs.
		Ins(I().Ret()).
		Ins(I().Mov().Rax().Rdi()).
		Ins(I().Mov().Rax().Rax())

	rhs := New()
	rhs.Ins(I().Ret())
	rhs.Concat(New().
		Ins(I().Mov().Rax().Rdi()).
		Ins(I().Mov().Rax().Rax()))

	if !lhs.Eq(rhs) {
		t.Errorf("Expected to get the same Insts.")
	}
}
