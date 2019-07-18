package gen

import (
	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/psr"
)

type Compiler psr.Parser

var orIdt = func() Compiler { return Compiler(psr.OrIdent()) }
var andIdt = func() Compiler { return Compiler(psr.AndIdent()) }
var null = andIdt()

func p(c Compiler) psr.Parser    { return psr.Parser(c) }
func pp(c *Compiler) *psr.Parser { return (*psr.Parser)(c) }

func (c Compiler) And(cp *Compiler, addNode bool) Compiler {
	return Compiler(p(c).And(pp(cp), addNode))
}

func (c Compiler) Or(cp *Compiler) Compiler {
	return Compiler(p(c).Or(pp(cp)))
}

func (c Compiler) Rep(cp *Compiler) Compiler {
	return Compiler(p(c).Rep(pp(cp)))
}

func (c Compiler) Trans(f func(ast.AST) ast.AST) Compiler {
	return Compiler(p(c).Trans(f))
}

func (c Compiler) SetEval(f func(nodes []*ast.AST, code asm.Code)) Compiler {
	return Compiler(p(c).SetEval(f))
}

var (
	Int   = (*Compiler)(psr.Int)
	Intd  = (*Compiler)(psr.Intd)
	Plus  = (*Compiler)(psr.Plus)
	Minus = (*Compiler)(psr.Minus)
	Mul   = (*Compiler)(psr.Mul)
	Div   = (*Compiler)(psr.Div)
	Eq    = (*Compiler)(psr.Eq)
	Neq   = (*Compiler)(psr.Neq)
	LPar  = (*Compiler)(psr.LPar)
	RPar  = (*Compiler)(psr.RPar)
	RBrc  = (*Compiler)(psr.RBrc)
	LBrc  = (*Compiler)(psr.LBrc)
	Amp   = (*Compiler)(psr.Amp)
	Subs  = (*Compiler)(psr.Subs)
	CVar  = (*Compiler)(psr.Var)
	Semi  = (*Compiler)(psr.Semi)
	Com   = (*Compiler)(psr.Com)
	Ret   = (*Compiler)(psr.Ret)
	If    = (*Compiler)(psr.If)
	While = (*Compiler)(psr.While)
	For   = (*Compiler)(psr.For)
	Sys   = (*Compiler)(psr.Sys)
	EOF   = (*Compiler)(psr.EOF)
)
