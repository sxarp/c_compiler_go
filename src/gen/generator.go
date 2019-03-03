package gen

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/psr"
)

var orId = psr.OrId
var andId = psr.AndId
var null = andId()

const wordSize = 8

func checkNodeCount(nodes []*ast.AST, count int) {
	if l := len(nodes); l != count {
		panic(fmt.Sprintf("The number of nodes must be %d, got %d.", count, l))
	}
}

var numInt = andId().And(psr.Int, true).
	SetEval(func(nodes []*ast.AST, code asm.Code) {
		i := nodes[0].Token.Vali()
		code.Ins(asm.I().Push().Val(i))
	})

func binaryOperator(term *psr.Parser, operator *psr.Parser, insts []asm.Fin) psr.Parser {
	return andId().And(operator, false).And(term, true).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 1)
			nodes[0].Eval(code)

			code.
				Ins(asm.I().Pop().Rdi()).
				Ins(asm.I().Pop().Rax())

			for _, i := range insts {
				code.Ins(i)
			}

			code.Ins(asm.I().Push().Rax())
		})

}

func adder(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Plus, []asm.Fin{asm.I().Add().Rax().Rdi()})
}

func subber(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Minus, []asm.Fin{asm.I().Sub().Rax().Rdi()})
}

func addsubs(term *psr.Parser) psr.Parser {
	add, sub := adder(term), subber(term)
	addsub := orId().Or(&add).Or(&sub)
	return andId().And(term, true).Rep(&addsub).Trans(ast.PopSingle)
}

func muler(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Mul, []asm.Fin{asm.I().Mul().Rdi()})
}

func diver(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Div, []asm.Fin{asm.I().Mov().Rdx().Val(0), asm.I().Div().Rdi()})
}

func muldivs(term *psr.Parser) psr.Parser {
	mul, div := muler(term), diver(term)
	muldiv := orId().Or(&mul).Or(&div)
	return andId().And(term, true).Rep(&muldiv).Trans(ast.PopSingle)
}

func eqer(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Eq, []asm.Fin{
		asm.I().Cmp().Rdi().Rax(),
		asm.I().Sete().Al(),
		asm.I().Movzb().Rax().Al()})
}

func neqer(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Neq, []asm.Fin{
		asm.I().Cmp().Rdi().Rax(),
		asm.I().Setne().Al(),
		asm.I().Movzb().Rax().Al()})
}

func eqneqs(term *psr.Parser) psr.Parser {
	eq, neq := eqer(term), neqer(term)
	eqneq := orId().Or(&eq).Or(&neq)
	return andId().And(term, true).Rep(&eqneq).Trans(ast.PopSingle)
}

func prologuer(st *SymTable) psr.Parser {
	return andId().SetEval(func(nodes []*ast.AST, code asm.Code) {
		code.
			Ins(asm.I().Push().Rbp()).
			Ins(asm.I().Mov().Rbp().Rsp()).
			Ins(asm.I().Sub().Rsp().Val(wordSize * st.Count()))
	})
}

var epilogue psr.Parser = andId().SetEval(
	func(nodes []*ast.AST, code asm.Code) {
		code.
			Ins(asm.I().Mov().Rsp().Rbp()).
			Ins(asm.I().Pop().Rbp()).
			Ins(asm.I().Ret())
	})

var popRax psr.Parser = andId().SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Pop().Rax()) })

func lvIdenter(st *SymTable) psr.Parser {
	return andId().And(psr.Var, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 1)
			offSet := wordSize * (1 + st.RefOf(nodes[0].Token.Val()))
			code.
				Ins(asm.I().Mov().Rax().Rbp()).
				Ins(asm.I().Sub().Rax().Val(offSet)).
				Ins(asm.I().Push().Rax())

		})
}

func rvIdenter(lvIdent *psr.Parser) psr.Parser {
	return andId().And(lvIdent, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 1)
			nodes[0].Eval(code)
			code.
				Ins(asm.I().Pop().Rax()).
				Ins(asm.I().Mov().Rax().Rax().P()).
				Ins(asm.I().Push().Rax())
		})

}

func assigner(lv *psr.Parser, rv *psr.Parser) psr.Parser {
	return andId().And(lv, true).And(psr.Subs, false).And(rv, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)

			nodes[1].Eval(code) // Evaluate right value
			nodes[0].Eval(code) // Evaluate left value
			code.
				Ins(asm.I().Pop().Rdi()).           // load lv to rdi
				Ins(asm.I().Pop().Rax()).           // load rv to rax
				Ins(asm.I().Mov().Rdi().P().Rax()). // mv rax to [lv]
				Ins(asm.I().Push().Rax())

		})
}

func returner(term *psr.Parser) psr.Parser {
	retval := andId().And(term, true).And(&popRax, true)
	ret := orId().Or(&retval).Or(&null)
	return andId().And(psr.Ret, false).And(&ret, true).And(&epilogue, true)
}

var ifcount int

func ifer(condition *psr.Parser, body *psr.Parser) psr.Parser {
	return andId().And(psr.If, false).And(psr.LPar, false).And(condition, true).And(psr.RPar, false).
		And(psr.LBrc, false).And(body, true).And(psr.RBrc, false).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[0].Eval(code)

			label := fmt.Sprintf("iflabel_%d", ifcount)
			code.
				Ins(asm.I().Pop().Rax()).
				Ins(asm.I().Cmp().Rax().Val(0)).
				Ins(asm.I().Je(label))

			nodes[1].Eval(code)
			code.Ins(asm.I().Label(label))
			ifcount += 1
		})
}

var whilecount int

func whiler(condition, body *psr.Parser) psr.Parser {
	return andId().And(psr.While, false).And(psr.LPar, false).And(condition, true).And(psr.RPar, false).
		And(psr.LBrc, false).And(body, true).And(psr.RBrc, false).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			begin, end := fmt.Sprintf("while_begin_%d", whilecount),
				fmt.Sprintf("while_end_%d", whilecount)

			code.Ins(asm.I().Label(begin))

			// Evaluate the condition part.
			nodes[0].Eval(code)

			// If the condition part is evaluated as zero, then go to end.
			code.
				Ins(asm.I().Pop().Rax()).
				Ins(asm.I().Cmp().Rax().Val(0)).
				Ins(asm.I().Je(end))

			// Evaluate the body part.
			nodes[1].Eval(code)

			// Unconditional jump to begin.
			code.Ins(asm.I().Jmp(begin))

			code.Ins(asm.I().Label(end))

			whilecount += 1
		})
}

var forcount int

func forer(conditions, body *psr.Parser) psr.Parser {
	nullCond := orId().Or(conditions).Or(&null)
	semiCond := andId().And(&nullCond, true).And(psr.Semi, false)
	ini := andId().And(&semiCond, true).And(&popRax, true)
	incr := andId().And(&nullCond, true).And(&popRax, true)

	return andId().And(psr.For, false).And(psr.LPar, false).
		And(&ini, true).And(&semiCond, true).And(&incr, true).And(psr.RPar, false).
		And(psr.LBrc, false).And(body, true).And(psr.RBrc, false).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 4)

			begin, end := fmt.Sprintf("for_begin_%d", forcount),
				fmt.Sprintf("for_end_%d", forcount)

			// Evaluate the initialization part.
			nodes[0].Eval(code)

			code.Ins(asm.I().Label(begin))

			// Evaluate the condition part.
			nodes[1].Eval(code)

			// If condition part is evaluated as zero, then go to end.
			code.
				Ins(asm.I().Pop().Rax()).
				Ins(asm.I().Cmp().Rax().Val(0)).
				Ins(asm.I().Je(end))

			// Evaluate the increment part.
			nodes[2].Eval(code)

			// Evaluate the body part.
			nodes[3].Eval(code)

			// Unconditional jump to begin.
			code.Ins(asm.I().Jmp(begin))

			code.Ins(asm.I().Label(end))

			forcount++
		})
}

func funcCaller(term *psr.Parser) psr.Parser {
	funcName := andId().And(psr.Var, true).
		SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Call(nodes[0].Token.Val())) })

	commed := andId().And(psr.Com, false).And(term, true).Trans(ast.PopSingle)

	argRegs := []asm.Dested{
		asm.I().Mov().Rdi(),
		asm.I().Mov().Rsi(),
		asm.I().Mov().Rdx(),
		asm.I().Mov().Rcx(),
		asm.I().Mov().R8(),
		asm.I().Mov().R9(),
	}

	argvs := andId().And(term, true).Rep(&commed).SetEval(func(nodes []*ast.AST, code asm.Code) {
		if len(nodes) > 6 {
			panic("too many arguments")
		}

		// Evaluate args from right to left and push into the stack.
		for i := range nodes {
			nodes[len(nodes)-i-1].Eval(code)
		}

		for i := range nodes {
			code.Ins(asm.I().Pop().Rax()).Ins(argRegs[i].Rax())
		}
	})

	args := orId().Or(&argvs).Or(&null)

	return andId().And(&funcName, true).And(psr.LPar, false).And(&args, true).And(psr.RPar, false).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[1].Eval(code)
			nodes[0].Eval(code)
			code.Ins(asm.I().Push().Rax())
		})
}

func funcDefiner(bodyer func(*SymTable) psr.Parser) psr.Parser {
	var st = new(SymTable)

	funcName := andId().And(psr.Var, true).
		SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Label(nodes[0].Token.Val())) })

	argRegs := []asm.Fin{
		asm.I().Mov().Rax().P().Rdi(),
		asm.I().Mov().Rax().P().Rsi(),
		asm.I().Mov().Rax().P().Rdx(),
		asm.I().Mov().Rax().P().Rcx(),
		asm.I().Mov().Rax().P().R8(),
		asm.I().Mov().Rax().P().R9(),
	}

	argv := andId().And(psr.Var, true).SetEval(func(nodes []*ast.AST, code asm.Code) {
		seqNum := st.RefOf(nodes[0].Token.Val())

		if seqNum >= 6 {
			panic("too many arguments")
		}

		offSet := wordSize * (1 + seqNum)
		code.
			Ins(asm.I().Mov().Rax().Rbp()).
			Ins(asm.I().Sub().Rax().Val(offSet)).
			Ins(argRegs[seqNum])
	})

	commed := andId().And(psr.Com, false).And(&argv, true).Trans(ast.PopSingle)
	argvs := andId().And(&argv, true).Rep(&commed)
	args := orId().Or(&argvs).Or(&null)
	body := bodyer(st)
	prologue := prologuer(st)

	return andId().And(psr.Intd, false).And(&funcName, true).And(psr.LPar, false).And(&args, true).And(psr.RPar, false).
		And(&prologue, true).And(psr.LBrc, false).And(&body, true).And(psr.RBrc, false).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 4)

			*st = *newST()      // Initialize SymTable.
			nodes[0].Eval(code) // Start defining function.

			bottom := asm.New()
			nodes[1].Eval(bottom) // Evaluete argvs.
			nodes[3].Eval(bottom) // Evaluate body.

			// Evaluate the prologue AST afterwards so that the symbol table can emit
			// the correct number of variables declared, which is used by the prologue code
			// to determine the size of the stack to allocate for the variables.
			nodes[2].Eval(code)

			insts := code.(*asm.Insts)
			insts.Concat(bottom)
		})
}

func Generator() psr.Parser {
	body := func(st *SymTable) psr.Parser {
		lvIdent := lvIdenter(st)
		rvIdent := rvIdenter(&lvIdent)

		var term, muls, adds, expr, eqs, call, ifex, while, forex psr.Parser
		num := orId().Or(&numInt).Or(&call).Or(&rvIdent)

		eqs = eqneqs(&adds)
		adds = addsubs(&muls)
		muls = muldivs(&term)

		parTerm := andId().And(psr.LPar, false).And(&adds, true).And(psr.RPar, false).Trans(ast.PopSingle)
		term = orId().Or(&num).Or(&parTerm)

		assign := assigner(&lvIdent, &expr)
		expr = orId().Or(&assign).Or(&eqs)
		call = funcCaller(&expr)
		semi := andId().And(&expr, true).And(psr.Semi, false)

		line := andId().And(&semi, true).And(&popRax, true)
		ret := returner(&semi)

		body := orId().Or(&ifex).Or(&forex).Or(&while).Or(&ret).Or(&line)
		bodies := andId().Rep(&body)

		ifex = ifer(&expr, &bodies)
		forex = forer(&expr, &bodies)
		while = whiler(&expr, &bodies)

		return andId().And(&bodies, true)
	}

	function := funcDefiner(body)
	functions := andId().Rep(&function)

	return andId().And(&functions, true).And(psr.EOF, false)
}
