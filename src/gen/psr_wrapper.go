package gen

import (
	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/psr"
	"github.com/sxarp/c_compiler_go/src/tok"
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

func tokenTypeToComp(tt *tok.TokenType) *Compiler {
	return (*Compiler)(psr.TokenTypeToPsr(tt))
}

var (
	Int   = tokenTypeToComp(&tok.TInt)
	Intd  = tokenTypeToComp(&tok.TIntd)
	Plus  = tokenTypeToComp(&tok.TPlus)
	Minus = tokenTypeToComp(&tok.TMinus)
	Mul   = tokenTypeToComp(&tok.TMul)
	Div   = tokenTypeToComp(&tok.TDiv)
	Eq    = tokenTypeToComp(&tok.TEq)
	Neq   = tokenTypeToComp(&tok.TNeq)
	LPar  = tokenTypeToComp(&tok.TLPar)
	RPar  = tokenTypeToComp(&tok.TRPar)
	RBrc  = tokenTypeToComp(&tok.TRBrc)
	LBrc  = tokenTypeToComp(&tok.TLBrc)
	RSbr  = tokenTypeToComp(&tok.TRSbr)
	LSbr  = tokenTypeToComp(&tok.TLSbr)
	Amp   = tokenTypeToComp(&tok.TAmp)
	Subs  = tokenTypeToComp(&tok.TSubs)
	CVar  = tokenTypeToComp(&tok.TVar)
	Semi  = tokenTypeToComp(&tok.TSemi)
	Com   = tokenTypeToComp(&tok.TCom)
	Ret   = tokenTypeToComp(&tok.TRet)
	If    = tokenTypeToComp(&tok.TIf)
	While = tokenTypeToComp(&tok.TWhile)
	For   = tokenTypeToComp(&tok.TFor)
	Sys   = tokenTypeToComp(&tok.TSys)
	EOF   = (*Compiler)(psr.EOF)
)
