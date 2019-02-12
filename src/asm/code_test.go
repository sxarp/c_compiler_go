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
		lhs := New()
		rhs := New()

		for i, _ := range comp.lhs {
			lhs.Ins(comp.lhs[i])
			rhs.Ins(comp.rhs[i])
		}

		h.Expectt(t, comp.expect, lhs.Eq(&rhs))

	}

}
