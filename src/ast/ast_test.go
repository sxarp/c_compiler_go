package ast

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/h"
)

func TestEval(t *testing.T) {
	addIns := func(nodes []*AST, code asm.Code) {
		code.Ins(asm.I().Ret())
	}

	nodes := []*AST{
		&AST{}, &AST{}, &AST{},
	}

	for _, node := range nodes {
		node.eval = addIns
	}

	a := AST{nodes: nodes}

	code := asm.New()
	a.Eval(code)

	numOfInst := 0

	code.ForEachInst(func(i asm.Fin) {
		*(&numOfInst)++
	})

	h.ExpectEq(t, numOfInst, 3)
}
