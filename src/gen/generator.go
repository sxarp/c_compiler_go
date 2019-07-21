package gen

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/asm"
	"github.com/sxarp/c_compiler_go/src/ast"
	"github.com/sxarp/c_compiler_go/src/tp"
)

func checkNodeCount(nodes []*ast.AST, count int) {
	if l := len(nodes); l != count {
		panic(fmt.Sprintf("The number of nodes must be %d, got %d.", count, l))
	}
}

var numInt = andIdt().And(Int, true).
	SetEval(func(nodes []*ast.AST, code asm.Code) {
		i := nodes[0].Token.Vali()
		code.Ins(asm.I().Push().Val(i))
	})

func binaryOperator(term *Compiler, operator *Compiler, insts []asm.Fin) Compiler {
	return andIdt().And(operator, false).And(term, true).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 1)
			nodes[0].Eval(code)

			code.Ins(asm.I().Pop().Rdi(), asm.I().Pop().Rax())
			code.Ins(insts...)
			code.Ins(asm.I().Push().Rax())
		})

}

func adder(term *Compiler) Compiler {
	return binaryOperator(term, Plus, []asm.Fin{asm.I().Add().Rax().Rdi()})
}

func subber(term *Compiler) Compiler {
	return binaryOperator(term, Minus, []asm.Fin{asm.I().Sub().Rax().Rdi()})
}

func addsubs(term *Compiler) Compiler {
	add, sub := adder(term), subber(term)
	addsub := orIdt().Or(&add).Or(&sub)
	return andIdt().And(term, true).Rep(&addsub).Trans(ast.PopSingle)
}

func muler(term *Compiler) Compiler {
	return binaryOperator(term, Mul, []asm.Fin{asm.I().Mul().Rdi()})
}

func diver(term *Compiler) Compiler {
	return binaryOperator(term, Div, []asm.Fin{asm.I().Mov().Rdx().Val(0), asm.I().Div().Rdi()})
}

func muldivs(term *Compiler) Compiler {
	mul, div := muler(term), diver(term)
	muldiv := orIdt().Or(&mul).Or(&div)
	return andIdt().And(term, true).Rep(&muldiv).Trans(ast.PopSingle)
}

func eqer(term *Compiler) Compiler {
	return binaryOperator(term, Eq, []asm.Fin{
		asm.I().Cmp().Rdi().Rax(),
		asm.I().Sete().Al(),
		asm.I().Movzb().Rax().Al()})
}

func neqer(term *Compiler) Compiler {
	return binaryOperator(term, Neq, []asm.Fin{
		asm.I().Cmp().Rdi().Rax(),
		asm.I().Setne().Al(),
		asm.I().Movzb().Rax().Al()})
}

func eqneqs(term *Compiler) Compiler {
	eq, neq := eqer(term), neqer(term)
	eqneq := orIdt().Or(&eq).Or(&neq)
	return andIdt().And(term, true).Rep(&eqneq).Trans(ast.PopSingle)
}

func syscaller(term *Compiler) Compiler {
	// Registers used to pass arguments to system call
	regs := []asm.Fin{
		asm.I().Pop().Rdi(),
		asm.I().Pop().Rsi(),
		asm.I().Pop().Rdx(),
		asm.I().Pop().R10(),
		asm.I().Pop().R8(),
		asm.I().Pop().R9(),
	}

	return andIdt().And(Sys, false).And(&numInt, true).Rep(term).
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

func prologuer(st *SymTable) Compiler {
	return andIdt().SetEval(func(nodes []*ast.AST, code asm.Code) {
		code.Ins(
			asm.I().Push().Rbp(),
			asm.I().Mov().Rbp().Rsp(),
			asm.I().Sub().Rsp().Val(st.Allocated()))
	})
}

var epilogue = andIdt().SetEval(
	func(nodes []*ast.AST, code asm.Code) {
		code.Ins(
			asm.I().Mov().Rsp().Rbp(),
			asm.I().Pop().Rbp(),
			asm.I().Ret())
	})

var popRax = andIdt().SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Pop().Rax()) })

func loadValer(st *SymTable, sym *string) Compiler {
	return andIdt().SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			code.Ins(
				asm.I().Mov().Rax().Rbp(),
				asm.I().Sub().Rax().Val(st.RefOf(*sym).Addr),
				asm.I().Push().Rax())
		})
}

func ptrAdder(st *SymTable, ptr *Compiler, addv *Compiler) Compiler {
	var size int

	fetchSize := func(nodes []*ast.AST, code asm.Code) {
		checkNodeCount(nodes, 2)

		// eval pointer value
		nodes[0].Eval(code)

		// then get the size of last referenced variable
		size = st.RefOf(st.LastRef()).Type.Size()

		// eval add val
		nodes[1].Eval(code)
	}

	ptrAdded := andIdt().And(Mul, false).And(LPar, false).And(ptr, true).
		And(Plus, false).And(addv, true).And(RPar, false).SetEval(fetchSize)

	array := andIdt().And(ptr, true).And(LSbr, false).And(addv, true).And(RSbr, false).SetEval(fetchSize)

	ptrArray := orIdt().Or(&ptrAdded).Or(&array)

	return andIdt().And(&ptrArray, true).SetEval(func(nodes []*ast.AST, code asm.Code) {
		checkNodeCount(nodes, 1)
		nodes[0].Eval(code)

		// multiple add val by size
		code.Ins(
			asm.I().Pop().Rax(),
			asm.I().Mov().Rdi().Val(size),
			asm.I().Mul().Rdi(),
			asm.I().Push().Rax())

		// add both values and push
		code.Ins(
			asm.I().Pop().Rdi(),
			asm.I().Pop().Rax(),
			asm.I().Add().Rax().Rdi(),
			asm.I().Push().Rax())
	})
}

func lvIdenter(st *SymTable) Compiler {
	var sym string
	loadVal := loadValer(st, &sym)

	return andIdt().And(CVar, true).And(&loadVal, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			sym = nodes[0].Token.Val()
			nodes[1].Eval(code)
		})
}

func rvAddrer(lvIdent *Compiler) Compiler {
	return andIdt().And(Amp, false).And(lvIdent, true)
}

var deRefer = andIdt().SetEval(func(nodes []*ast.AST, code asm.Code) {
	code.Ins(
		asm.I().Pop().Rax(),
		asm.I().Mov().Rax().Rax().P(),
		asm.I().Push().Rax())
})

func ptrDeRefer(st *SymTable, lvIdent *Compiler) Compiler {
	astr := andIdt().And(Mul, false).And(&deRefer, true)

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

	return andIdt().And(&astrs, true).And(CVar, true).And(&loadVal, true).SetEval(
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

func rvIdenter(st *SymTable, ptrDeRef *Compiler) Compiler {
	return andIdt().And(ptrDeRef, true).And(&deRefer, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[0].Eval(code)

			// skip dereference for array variables
			// so that they behaves like pointer variables
			if !st.RefOf(st.LastRef()).Type.IsArray() {
				nodes[1].Eval(code)
			}
		})
}

func assigner(lv *Compiler, rv *Compiler) Compiler {
	return andIdt().And(lv, true).And(Subs, false).And(rv, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)

			nodes[1].Eval(code) // Evaluate right value
			nodes[0].Eval(code) // Evaluate left value
			code.Ins(
				asm.I().Pop().Rdi(),           // load lv to rdi
				asm.I().Pop().Rax(),           // load rv to rax
				asm.I().Mov().Rdi().P().Rax(), // mv rax to [lv]
				asm.I().Push().Rax())

		})
}

func varDeclarer(st *SymTable) Compiler {
	var varType = tp.Int
	vp := &varType

	astrs := andIdt().Rep(Mul).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			for range nodes {
				*vp = (*vp).Ptr()
			}
		})

	return andIdt().And(Intd, false).And(&astrs, true).And(CVar, true).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[0].Eval(nil)

			st.DecOf(nodes[1].Token.Val(), varType)
		})
}

func arrayDeclarer(varDeclare *Compiler, st *SymTable) Compiler {
	array := andIdt().And(LSbr, false).And(Int, true).And(RSbr, false).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 1)
			v := st.Last()
			st.OverWrite(v.Name, tp.Array(v.Type, nodes[0].Token.Vali()))
		})

	return andIdt().And(varDeclare, true).Rep(&array)
}

func returner(term *Compiler) Compiler {
	retval := andIdt().And(term, true).And(&popRax, true)
	ret := orIdt().Or(&retval).Or(&null)
	return andIdt().And(Ret, false).And(&ret, true).And(&epilogue, true)
}

var ifcount int

func ifer(condition *Compiler, body *Compiler) Compiler {
	return andIdt().And(If, false).And(LPar, false).And(condition, true).And(RPar, false).
		And(LBrc, false).And(body, true).And(RBrc, false).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[0].Eval(code)

			label := fmt.Sprintf("iflabel_%d", ifcount)
			code.Ins(
				asm.I().Pop().Rax(),
				asm.I().Cmp().Rax().Val(0),
				asm.I().Je(label))

			nodes[1].Eval(code)
			code.Ins(asm.I().Label(label))
			ifcount++
		})
}

var whilecount int

func whiler(condition, body *Compiler) Compiler {
	return andIdt().And(While, false).And(LPar, false).And(condition, true).And(RPar, false).
		And(LBrc, false).And(body, true).And(RBrc, false).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			begin, end := fmt.Sprintf("while_begin_%d", whilecount),
				fmt.Sprintf("while_end_%d", whilecount)

			code.Ins(asm.I().Label(begin))

			// Evaluate the condition part.
			nodes[0].Eval(code)

			// If the condition part is evaluated as zero, then go to end.
			code.Ins(
				asm.I().Pop().Rax(),
				asm.I().Cmp().Rax().Val(0),
				asm.I().Je(end))

			// Evaluate the body part.
			nodes[1].Eval(code)

			// Unconditional jump to begin.
			code.Ins(asm.I().Jmp(begin))

			code.Ins(asm.I().Label(end))

			whilecount++
		})
}

var forcount int

func forer(conditions, body *Compiler) Compiler {
	nullCond := orIdt().Or(conditions).Or(&null)
	semiCond := andIdt().And(&nullCond, true).And(Semi, false)
	ini := andIdt().And(&semiCond, true).And(&popRax, true)
	incr := andIdt().And(&nullCond, true).And(&popRax, true)

	return andIdt().And(For, false).And(LPar, false).
		And(&ini, true).And(&semiCond, true).And(&incr, true).And(RPar, false).
		And(LBrc, false).And(body, true).And(RBrc, false).SetEval(
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
			code.Ins(
				asm.I().Pop().Rax(),
				asm.I().Cmp().Rax().Val(0),
				asm.I().Je(end))

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

func funcCaller(term *Compiler) Compiler {
	funcName := andIdt().And(CVar, true).
		SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Call(nodes[0].Token.Val())) })

	commed := andIdt().And(Com, false).And(term, true).Trans(ast.PopSingle)

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

	return andIdt().And(&funcName, true).And(LPar, false).And(&args, true).And(RPar, false).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[1].Eval(code)
			nodes[0].Eval(code)
			code.Ins(asm.I().Push().Rax())
		})
}

func funcDefiner(bodyer func(*SymTable) Compiler) Compiler {
	var st = new(SymTable)

	funcName := andIdt().And(CVar, true).
		SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(asm.I().Label(nodes[0].Token.Val())) })

	argRegs := []asm.Fin{
		asm.I().Mov().Rax().P().Rdi(),
		asm.I().Mov().Rax().P().Rsi(),
		asm.I().Mov().Rax().P().Rdx(),
		asm.I().Mov().Rax().P().Rcx(),
		asm.I().Mov().Rax().P().R8(),
		asm.I().Mov().Rax().P().R9(),
	}

	varDeclare := varDeclarer(st)

	argv := andIdt().And(&varDeclare, true).SetEval(func(nodes []*ast.AST, code asm.Code) {

		nodes[0].Eval(code)

		v := st.Last()

		if v.Seq >= 6 {
			panic("too many arguments")
		}

		code.Ins(
			asm.I().Mov().Rax().Rbp(),
			asm.I().Sub().Rax().Val(v.Addr),
			argRegs[v.Seq])
	})

	commed := andIdt().And(Com, false).And(&argv, true).Trans(ast.PopSingle)
	argvs := andIdt().And(&argv, true).Rep(&commed)
	args := orIdt().Or(&argvs).Or(&null)
	body := bodyer(st)
	prologue := prologuer(st)

	return andIdt().And(Intd, false).And(&funcName, true).And(LPar, false).And(&args, true).And(RPar, false).
		And(&prologue, true).And(LBrc, false).And(&body, true).And(RBrc, false).
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

func Generator() Compiler {
	body := func(st *SymTable) Compiler {
		lvIdent := lvIdenter(st)
		ptrDeRef := ptrDeRefer(st, &lvIdent)

		rvAddr := rvAddrer(&lvIdent)
		rvIdent := rvIdenter(st, &ptrDeRef)
		rvVal := orIdt().Or(&rvAddr).Or(&rvIdent)

		ptrAddVal := orIdt().Or(&numInt).Or(&rvIdent)
		ptrAdd := ptrAdder(st, &rvVal, &ptrAddVal)
		rvPtrAdder := andIdt().And(&ptrAdd, true).And(&deRefer, true)

		leftVal := orIdt().Or(&ptrAdd).Or(&ptrDeRef)

		var num, term, muls, adds, expr, eqs, call, ifex, while, forex, syscall Compiler

		num = orIdt().Or(&rvPtrAdder).Or(&numInt).Or(&syscall).Or(&call).Or(&rvVal)
		eqs, adds, muls = eqneqs(&adds), addsubs(&muls), muldivs(&term)

		parTerm := andIdt().And(LPar, false).And(&adds, true).And(RPar, false).Trans(ast.PopSingle)
		term = orIdt().Or(&num).Or(&parTerm)

		assign := assigner(&leftVal, &expr)
		expr = orIdt().Or(&assign).Or(&eqs)
		call = funcCaller(&expr)
		syscall = syscaller(&expr)
		semi := andIdt().And(&expr, true).And(Semi, false)

		varDeclare := varDeclarer(st)
		arrayDeclare := arrayDeclarer(&varDeclare, st)
		semiVarDeclare := andIdt().And(&arrayDeclare, true).And(Semi, false)

		line, ret := andIdt().And(&semi, true).And(&popRax, true), returner(&semi)

		body := orIdt().Or(&ifex).Or(&forex).Or(&while).Or(&ret).Or(&semiVarDeclare).Or(&line)
		bodies := andIdt().Rep(&body)

		ifex, forex, while = ifer(&expr, &bodies), forer(&expr, &bodies), whiler(&expr, &bodies)

		return andIdt().And(&bodies, true)
	}

	function := funcDefiner(body)
	functions := andIdt().Rep(&function)

	return andIdt().And(&functions, true).And(EOF, false)
}
