package gen

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/psr"
	"github.com/sxarp/c_compiler_go/src/tp"
)

var orIdt = psr.OrIdent
var andIdt = psr.AndIdent
var null = andIdt()

func checkNodeCount(nodes []*ast.AST, count int) {
	if l := len(nodes); l != count {
		panic(fmt.Sprintf("The number of nodes must be %d, got %d.", count, l))
	}
}

var numInt = andIdt().And(psr.Int, true).
	SetEval(func(nodes []*ast.AST, code asm.Code) {
		i := nodes[0].Token.Vali()
		code.Ins(asm.I().Push().Val(i))
	})

func binaryOperator(term *psr.Parser, operator *psr.Parser, insts []asm.Fin) psr.Parser {
	return andIdt().And(operator, false).And(term, true).
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
	addsub := orIdt().Or(&add).Or(&sub)
	return andIdt().And(term, true).Rep(&addsub).Trans(ast.PopSingle)
}

func muler(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Mul, []asm.Fin{asm.I().Mul().Rdi()})
}

func diver(term *psr.Parser) psr.Parser {
	return binaryOperator(term, psr.Div, []asm.Fin{asm.I().Mov().Rdx().Val(0), asm.I().Div().Rdi()})
}

func muldivs(term *psr.Parser) psr.Parser {
	mul, div := muler(term), diver(term)
	muldiv := orIdt().Or(&mul).Or(&div)
	return andIdt().And(term, true).Rep(&muldiv).Trans(ast.PopSingle)
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
	eqneq := orIdt().Or(&eq).Or(&neq)
	return andIdt().And(term, true).Rep(&eqneq).Trans(ast.PopSingle)
}

func syscaller(term *psr.Parser) psr.Parser {
	// Registers used to pass arguments to system call
	regs := []asm.Fin{
		asm.I().Pop().Rdi(),
		asm.I().Pop().Rsi(),
		asm.I().Pop().Rdx(),
		asm.I().Pop().R10(),
		asm.I().Pop().R8(),
		asm.I().Pop().R9(),
	}

	return andIdt().And(psr.Sys, false).And(&numInt, true).Rep(term).
		SetEval(func(nodes []*ast.AST, code asm.Code) {

			for i, node := range nodes[1:] {
				node.Eval(code)
				code.Ins(regs[i])
			}

			// Set syscall number
			nodes[0].Eval(code)
			code.Ins(asm.I().Pop().Rax())
			// Call syscall instruction
			code.Ins(asm.I().Sys())
			// Push returned value from system call
			code.Ins(asm.I().Push().Rax())
		})
}

func prologuer(st *SymTable) psr.Parser {
	return andIdt().SetEval(func(nodes []*ast.AST, code asm.Code) {
		code.
			Ins(asm.I().Push().Rbp()).
			Ins(asm.I().Mov().Rbp().Rsp()).
			Ins(asm.I().Sub().Rsp().Val(st.Allocated()))
	})
}

var epilogue psr.Parser = andIdt().SetEval(
	func(nodes []*ast.AST, code asm.Code) {
		code.
			Ins(asm.I().Mov().Rsp().Rbp()).
			Ins(asm.I().Pop().Rbp()).
			Ins(asm.I().Ret())
	})

var popRax psr.Parser = andIdt().SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Pop().Rax()) })

func loadValer(st *SymTable, sym *string) psr.Parser {
	return andIdt().SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			code.
				Ins(asm.I().Mov().Rax().Rbp()).
				Ins(asm.I().Sub().Rax().Val(st.RefOf(*sym).Addr)).
				Ins(asm.I().Push().Rax())
		})
}

func ptrAdder(st *SymTable, addv *psr.Parser) psr.Parser {
	return andIdt().And(psr.Mul, false).And(psr.LPar, false).
		And(psr.Var, true).And(psr.Plus, false).And(addv, true).
		And(psr.RPar, false).SetEval(func(nodes []*ast.AST, code asm.Code) {
		checkNodeCount(nodes, 2)
		val := st.RefOf(nodes[0].Token.Val())
		size := val.Type.Size()

		// load ptr
		code.
			Ins(asm.I().Mov().Rax().Rbp()).
			Ins(asm.I().Sub().Rax().Val(val.Addr)).
			Ins(asm.I().Push().Rax())

		// load value
		code.
			Ins(asm.I().Pop().Rax()).
			Ins(asm.I().Mov().Rax().Rax().P()).
			Ins(asm.I().Push().Rax())

		// eval add val
		nodes[1].Eval(code)

		// multiple add val by size
		code.
			Ins(asm.I().Pop().Rax()).
			Ins(asm.I().Mov().Rdi().Val(size)).
			Ins(asm.I().Mul().Rdi()).
			Ins(asm.I().Push().Rax())

		// add both values and push
		code.
			Ins(asm.I().Pop().Rdi()).
			Ins(asm.I().Pop().Rax()).
			Ins(asm.I().Add().Rax().Rdi()).
			Ins(asm.I().Push().Rax())
	})
}

func lvIdenter(st *SymTable) psr.Parser {
	var sym string
	loadVal := loadValer(st, &sym)

	return andIdt().And(psr.Var, true).And(&loadVal, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			sym = nodes[0].Token.Val()
			nodes[1].Eval(code)
		})
}

func rvAddrer(lvIdent *psr.Parser) psr.Parser {
	return andIdt().And(psr.Amp, false).And(lvIdent, true)
}

var deRefer psr.Parser = andIdt().SetEval(func(nodes []*ast.AST, code asm.Code) {
	code.
		Ins(asm.I().Pop().Rax()).
		Ins(asm.I().Mov().Rax().Rax().P()).
		Ins(asm.I().Push().Rax())
})

func ptrDeRefer(st *SymTable, lvIdent *psr.Parser) psr.Parser {
	astr := andIdt().And(psr.Mul, false).And(&deRefer, true)

	var deRefCount int
	astrs := andIdt().Rep(&astr).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			deRefCount = len(nodes)

			for _, node := range nodes {
				node.Eval(code)
			}
		})

	var sym string
	loadVal := loadValer(st, &sym)

	return andIdt().And(&astrs, true).And(psr.Var, true).And(&loadVal, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			sym = nodes[1].Token.Val()
			symType := st.RefOf(sym).Type

			nodes[2].Eval(code)
			nodes[0].Eval(code)

			for i := 0; i < deRefCount; i++ {
				if _, ok := symType.DeRef(); !ok {
					panic(fmt.Sprintf("Invalid pointer variable dereference of %s.", sym))
				}
			}
		})
}

func rvIdenter(ptrDeRef *psr.Parser) psr.Parser {
	return andIdt().And(ptrDeRef, true).And(&deRefer, true)
}

func assigner(lv *psr.Parser, rv *psr.Parser) psr.Parser {
	return andIdt().And(lv, true).And(psr.Subs, false).And(rv, true).SetEval(
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

func varDeclarer(st *SymTable, delm *psr.Parser) psr.Parser {
	var varType = tp.Int
	vp := &varType

	astrs := andIdt().Rep(psr.Mul).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			for range nodes {
				*vp = (*vp).Ptr()
			}
		})

	return andIdt().And(psr.Intd, false).And(&astrs, true).And(psr.Var, true).And(delm, false).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[0].Eval(nil)

			st.DecOf(nodes[1].Token.Val(), varType)
		})
}

func returner(term *psr.Parser) psr.Parser {
	retval := andIdt().And(term, true).And(&popRax, true)
	ret := orIdt().Or(&retval).Or(&null)
	return andIdt().And(psr.Ret, false).And(&ret, true).And(&epilogue, true)
}

var ifcount int

func ifer(condition *psr.Parser, body *psr.Parser) psr.Parser {
	return andIdt().And(psr.If, false).And(psr.LPar, false).And(condition, true).And(psr.RPar, false).
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
			ifcount++
		})
}

var whilecount int

func whiler(condition, body *psr.Parser) psr.Parser {
	return andIdt().And(psr.While, false).And(psr.LPar, false).And(condition, true).And(psr.RPar, false).
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

			whilecount++
		})
}

var forcount int

func forer(conditions, body *psr.Parser) psr.Parser {
	nullCond := orIdt().Or(conditions).Or(&null)
	semiCond := andIdt().And(&nullCond, true).And(psr.Semi, false)
	ini := andIdt().And(&semiCond, true).And(&popRax, true)
	incr := andIdt().And(&nullCond, true).And(&popRax, true)

	return andIdt().And(psr.For, false).And(psr.LPar, false).
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
	funcName := andIdt().And(psr.Var, true).
		SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Call(nodes[0].Token.Val())) })

	commed := andIdt().And(psr.Com, false).And(term, true).Trans(ast.PopSingle)

	argRegs := []asm.Dested{
		asm.I().Mov().Rdi(),
		asm.I().Mov().Rsi(),
		asm.I().Mov().Rdx(),
		asm.I().Mov().Rcx(),
		asm.I().Mov().R8(),
		asm.I().Mov().R9(),
	}

	argvs := andIdt().And(term, true).Rep(&commed).SetEval(func(nodes []*ast.AST, code asm.Code) {
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

	args := orIdt().Or(&argvs).Or(&null)

	return andIdt().And(&funcName, true).And(psr.LPar, false).And(&args, true).And(psr.RPar, false).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[1].Eval(code)
			nodes[0].Eval(code)
			code.Ins(asm.I().Push().Rax())
		})
}

func funcDefiner(bodyer func(*SymTable) psr.Parser) psr.Parser {
	var st = new(SymTable)

	funcName := andIdt().And(psr.Var, true).
		SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Label(nodes[0].Token.Val())) })

	argRegs := []asm.Fin{
		asm.I().Mov().Rax().P().Rdi(),
		asm.I().Mov().Rax().P().Rsi(),
		asm.I().Mov().Rax().P().Rdx(),
		asm.I().Mov().Rax().P().Rcx(),
		asm.I().Mov().Rax().P().R8(),
		asm.I().Mov().Rax().P().R9(),
	}

	varDeclare := varDeclarer(st, &null)
	argv := andIdt().And(&varDeclare, true).SetEval(func(nodes []*ast.AST, code asm.Code) {

		nodes[0].Eval(code)

		v := st.Last()

		if v.Seq >= 6 {
			panic("too many arguments")
		}

		code.
			Ins(asm.I().Mov().Rax().Rbp()).
			Ins(asm.I().Sub().Rax().Val(v.Addr)).
			Ins(argRegs[v.Seq])
	})

	commed := andIdt().And(psr.Com, false).And(&argv, true).Trans(ast.PopSingle)
	argvs := andIdt().And(&argv, true).Rep(&commed)
	args := orIdt().Or(&argvs).Or(&null)
	body := bodyer(st)
	prologue := prologuer(st)

	return andIdt().And(psr.Intd, false).And(&funcName, true).And(psr.LPar, false).And(&args, true).And(psr.RPar, false).
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
		ptrDeRef := ptrDeRefer(st, &lvIdent)

		rvAddr := rvAddrer(&lvIdent)
		rvIdent := rvIdenter(&ptrDeRef)
		rvVal := orIdt().Or(&rvAddr).Or(&rvIdent)

		ptrAddVal := orIdt().Or(&numInt).Or(&rvIdent)
		ptrAdd := ptrAdder(st, &ptrAddVal)
		rvPtrAdder := andIdt().And(&ptrAdd, true).And(&deRefer, true)

		var num, term, muls, adds, expr, eqs, call, ifex, while, forex, syscall psr.Parser

		num = orIdt().Or(&rvPtrAdder).Or(&numInt).Or(&syscall).Or(&call).Or(&rvVal)
		eqs, adds, muls = eqneqs(&adds), addsubs(&muls), muldivs(&term)

		parTerm := andIdt().And(psr.LPar, false).And(&adds, true).And(psr.RPar, false).Trans(ast.PopSingle)
		term = orIdt().Or(&num).Or(&parTerm)

		assign := assigner(&ptrDeRef, &expr)
		expr = orIdt().Or(&assign).Or(&eqs)
		call = funcCaller(&expr)
		syscall = syscaller(&expr)
		semi := andIdt().And(&expr, true).And(psr.Semi, false)
		varDeclare := varDeclarer(st, psr.Semi)

		line, ret := andIdt().And(&semi, true).And(&popRax, true), returner(&semi)

		body := orIdt().Or(&ifex).Or(&forex).Or(&while).Or(&ret).Or(&varDeclare).Or(&line)
		bodies := andIdt().Rep(&body)

		ifex, forex, while = ifer(&expr, &bodies), forer(&expr, &bodies), whiler(&expr, &bodies)

		return andIdt().And(&bodies, true)
	}

	function := funcDefiner(body)
	functions := andIdt().Rep(&function)

	return andIdt().And(&functions, true).And(psr.EOF, false)
}
