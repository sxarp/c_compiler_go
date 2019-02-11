package ast

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/tok"
)

type ASTType struct{}

type AST struct {
	nodes []*AST
	Token *tok.Token
	atype *ASTType
	eval  func([]*AST, *asm.Code)
}

func (a *AST) Node(i int) AST {
	return *(a.nodes[i])
}

func (a *AST) AppendNode(an AST) {
	a.nodes = append(a.nodes, &an)
}

var TFail ASTType = ASTType{}
var Fail = AST{atype: &TFail}

func (a AST) Fail() bool {
	return a.atype == &TFail
}

func (a AST) Show() string {
	label := "term"
	if a.Token != nil {
		label = a.Token.Val()

	}

	if len(a.nodes) == 0 {
		return label

	}

	rets := fmt.Sprintf("(%s", label)

	for _, n := range a.nodes {
		rets += " " + n.Show()

	}

	return rets + ")"
}

func PopSingle(a AST) AST {
	if len(a.nodes) == 1 {
		return *(a.nodes[0])
	} else {
		return a

	}
}

func (a AST) Eval(code *asm.Code) {

	if a.eval == nil {
		for _, node := range a.nodes {
			node.Eval(code)
		}

	} else {
		a.eval(a.nodes, code)

	}

}
