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

var numInt = ai().And(Int).
	SetEval(func(nodes []*ast.AST, code asm.Code) {
		n := nodes[0].Token.Vali()
		code.Ins(i().Push().Val(n))
	})

func binaryOperator(term *Compiler, operator *Compiler, insts []asm.Fin) Compiler {
	return ai().Seq(operator).And(term).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 1)
			nodes[0].Eval(code)

			code.Ins(i().Pop().Rdi(), i().Pop().Rax())
			code.Ins(insts...)
			code.Ins(i().Push().Rax())
		})
}

func adder(term *Compiler) Compiler {
	return binaryOperator(term, Plus, []asm.Fin{i().Add().Rax().Rdi()})
}

func subber(term *Compiler) Compiler {
	return binaryOperator(term, Minus, []asm.Fin{i().Sub().Rax().Rdi()})
}

func lter(term *Compiler) Compiler {
	return binaryOperator(term, Lt, []asm.Fin{
		i().Cmp().Rax().Rdi(),
		i().Jl("0f"),
		i().Mov().Rax().Val(0),
		i().Jmp("1f"),
		i().Label("0"),
		i().Mov().Rax().Val(1),
		i().Label("1"),
	})
}

func addsubs(term *Compiler) Compiler {
	return ai().And(term).
		Rep(oi().Or(adder(term).P(), subber(term).P()).P()).
		Trans(ast.PopSingle)
}

func muler(term *Compiler) Compiler {
	return binaryOperator(term, Mul, []asm.Fin{i().Mul().Rdi()})
}

func diver(term *Compiler) Compiler {
	return binaryOperator(term, Div, []asm.Fin{i().Mov().Rdx().Val(0), i().Div().Rdi()})
}

func muldivs(term *Compiler) Compiler {
	return ai().And(term).
		Rep(oi().Or(muler(term).P(), diver(term).P()).P()).
		Trans(ast.PopSingle)
}

func eqer(term *Compiler) Compiler {
	return binaryOperator(term, Eq, []asm.Fin{
		i().Cmp().Rdi().Rax(),
		i().Sete().Al(),
		i().Movzb().Rax().Al()})
}

func neqer(term *Compiler) Compiler {
	return binaryOperator(term, Neq, []asm.Fin{
		i().Cmp().Rdi().Rax(),
		i().Setne().Al(),
		i().Movzb().Rax().Al()})
}

func eqneqs(term *Compiler) Compiler {
	return ai().And(term).
		Rep(oi().Or(eqer(term).P(), neqer(term).P(), lter(term).P()).P()).
		Trans(ast.PopSingle)
}

func syscaller(term *Compiler) Compiler {
	// Registers used to pass arguments to system call
	regs := []asm.Fin{
		i().Pop().Rdi(),
		i().Pop().Rsi(),
		i().Pop().Rdx(),
		i().Pop().R10(),
		i().Pop().R8(),
		i().Pop().R9(),
	}

	return ai().Seq(Sys).And(&numInt).Rep(term).
		SetEval(func(nodes []*ast.AST, code asm.Code) {

			for j, node := range nodes[1:] {
				node.Eval(code)
				code.Ins(regs[j])
			}

			// Set syscall number
			nodes[0].Eval(code)
			code.Ins(i().Pop().Rax())
			// Call syscall instruction
			code.Ins(i().Sys())
			// Push returned value from system call
			code.Ins(i().Push().Rax())
		})
}

func prologuer(st *SymTable) Compiler {
	return ai().SetEval(func(nodes []*ast.AST, code asm.Code) {
		code.Ins(
			i().Push().Rbp(),
			i().Mov().Rbp().Rsp(),
			i().Sub().Rsp().Val(st.Allocated()))
	})
}

var epilogue = ai().SetEval(
	func(nodes []*ast.AST, code asm.Code) {
		code.Ins(
			i().Mov().Rsp().Rbp(),
			i().Pop().Rbp(),
			i().Ret())
	})

var popRax = ai().SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(i().Pop().Rax()) })

func loadValer(st *SymTable, sym *string) Compiler {
	return ai().SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			code.Ins(
				i().Mov().Rax().Rbp(),
				i().Sub().Rax().Val(st.RefOf(*sym).Addr),
				i().Push().Rax())
		})
}

func ptrAdder(st *SymTable, ptr *Compiler, addv *Compiler) Compiler {
	var size int

	fetchSize := func(nodes []*ast.AST, code asm.Code) {
		checkNodeCount(nodes, 2)

		// eval pointer value
		nodes[0].Eval(code)

		// then get the size of last referenced pointer variable
		elm, _ := st.RefOf(st.LastRef()).Type.DeRef()
		size = elm.Size()

		// eval add val
		nodes[1].Eval(code)
	}

	ptrAdded := ai().Seq(Mul, LPar).And(ptr).
		Seq(Plus).And(addv).Seq(RPar).SetEval(fetchSize)
	array := ai().And(ptr).Seq(LSbr).And(addv).Seq(RSbr).SetEval(fetchSize)
	ptrArray := oi().Or(&ptrAdded, &array)

	return ai().And(&ptrArray).SetEval(func(nodes []*ast.AST, code asm.Code) {
		checkNodeCount(nodes, 1)
		nodes[0].Eval(code)

		// multiple add val by size
		code.Ins(
			i().Pop().Rax(),
			i().Mov().Rdi().Val(size),
			i().Mul().Rdi(),
			i().Push().Rax())

		// add both values and push
		code.Ins(
			i().Pop().Rdi(),
			i().Pop().Rax(),
			i().Add().Rax().Rdi(),
			i().Push().Rax())
	})
}

func lvIdenter(st *SymTable) Compiler {
	var sym string
	return ai().And(CVar).And(loadValer(st, &sym).P()).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			sym = nodes[0].Token.Val()
			nodes[1].Eval(code)
		})
}

func rvAddrer(lvIdent *Compiler) Compiler {
	return ai().Seq(Amp).And(lvIdent)
}

var deRefer = ai().SetEval(func(nodes []*ast.AST, code asm.Code) {
	code.Ins(
		i().Pop().Rax(),
		i().Mov().Rax().Rax().P(),
		i().Push().Rax())
})

func ptrDeRefer(st *SymTable, lvIdent *Compiler) Compiler {
	var (
		deRefCount int
		sym        string
	)

	astrs := ai().Rep(ai().Seq(Mul).And(&deRefer).P()).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			deRefCount = len(nodes)

			for _, node := range nodes {
				node.Eval(code)
			}
		})

	return ai().And(&astrs, CVar, loadValer(st, &sym).P()).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			sym = nodes[1].Token.Val()
			symType := st.RefOf(sym).Type

			nodes[2].Eval(code)
			nodes[0].Eval(code)

			for j := 0; j < deRefCount; j++ {
				if _, ok := symType.DeRef(); !ok {
					panic(fmt.Sprintf("Invalid pointer variable dereference of %s.", sym))
				}
			}
		})
}

func rvIdenter(st *SymTable, ptrDeRef *Compiler) Compiler {
	return ai().And(ptrDeRef, &deRefer).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[0].Eval(code)

			// skip dereference for array variables
			// so that they behave like pointer variables
			if !st.RefOf(st.LastRef()).Type.IsArray() {
				nodes[1].Eval(code)
			}
		})
}

func assigner(lv *Compiler, rv *Compiler) Compiler {
	return ai().And(lv).Seq(Subs).And(rv).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)

			nodes[1].Eval(code) // Evaluate right value
			nodes[0].Eval(code) // Evaluate left value
			code.Ins(
				i().Pop().Rdi(),           // load lv to rdi
				i().Pop().Rax(),           // load rv to rax
				i().Mov().Rdi().P().Rax(), // mv rax to [lv]
				i().Push().Rax())
		})
}

func varDeclarer(st *SymTable) Compiler {
	var varType = tp.Int
	vp := &varType

	astrs := ai().Rep(Mul).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			for range nodes {
				*vp = (*vp).Ptr()
			}
		})

	return ai().Seq(Intd).And(&astrs, CVar).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[0].Eval(nil)
			st.DecOf(nodes[1].Token.Val(), varType)
		})
}

func arrayDeclarer(varDeclare *Compiler, st *SymTable) Compiler {
	array := ai().Seq(LSbr).And(Int).Seq(RSbr).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 1)
			v := st.Last()
			st.OverWrite(v.Name, tp.Array(v.Type, nodes[0].Token.Vali()))
		})

	return ai().And(varDeclare).Rep(&array)
}

func returner(term *Compiler) Compiler {
	return ai().Seq(Ret).
		And(oi().Or(ai().And(term, &popRax).P(), &null).P(), &epilogue)
}

var ifcount int

func ifer(condition *Compiler, body *Compiler) Compiler {
	return ai().Seq(If, LPar).And(condition).Seq(RPar, LBrc).
		And(body).Seq(RBrc).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[0].Eval(code)

			label := fmt.Sprintf("iflabel_%d", ifcount)
			code.Ins(
				i().Pop().Rax(),
				i().Cmp().Rax().Val(0),
				i().Je(label))

			nodes[1].Eval(code)
			code.Ins(i().Label(label))
			ifcount++
		})
}

var whilecount int

func whiler(condition, body *Compiler) Compiler {
	return ai().Seq(While, LPar).And(condition).Seq(RPar, LBrc).
		And(body).Seq(RBrc).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			begin, end := fmt.Sprintf("while_begin_%d", whilecount),
				fmt.Sprintf("while_end_%d", whilecount)

			code.Ins(i().Label(begin))

			// Evaluate the condition part.
			nodes[0].Eval(code)

			// If the condition part is evaluated as zero, then go to end.
			code.Ins(
				i().Pop().Rax(),
				i().Cmp().Rax().Val(0),
				i().Je(end))

			// Evaluate the body part.
			nodes[1].Eval(code)

			// Unconditional jump to begin.
			code.Ins(i().Jmp(begin))

			code.Ins(i().Label(end))

			whilecount++
		})
}

var forcount int

func forer(conditions, body *Compiler) Compiler {
	nullCond := oi().Or(conditions, &null)
	semiCond := ai().And(&nullCond).Seq(Semi)
	ini := ai().And(&semiCond, &popRax)
	incr := ai().And(&nullCond, &popRax)

	return ai().Seq(For).Seq(LPar).
		And(&ini, &semiCond).And(&incr).Seq(RPar, LBrc).And(body).
		Seq(RBrc).SetEval(func(nodes []*ast.AST, code asm.Code) {
		checkNodeCount(nodes, 4)

		begin, end := fmt.Sprintf("for_begin_%d", forcount),
			fmt.Sprintf("for_end_%d", forcount)

		// Evaluate the initialization part.
		nodes[0].Eval(code)

		code.Ins(i().Label(begin))

		// Evaluate the condition part.
		nodes[1].Eval(code)

		// If condition part is evaluated as zero, then go to end.
		code.Ins(
			i().Pop().Rax(),
			i().Cmp().Rax().Val(0),
			i().Je(end))

		// Evaluate the body part.
		nodes[3].Eval(code)

		// Evaluate the increment part.
		nodes[2].Eval(code)

		// Unconditional jump to begin.
		code.Ins(i().Jmp(begin))

		code.Ins(i().Label(end))

		forcount++
	})
}

func funcCaller(term *Compiler) Compiler {
	funcName := ai().And(CVar).
		SetEval(func(nodes []*ast.AST, code asm.Code) { code.Ins(i().Call(nodes[0].Token.Val())) })
	commed := ai().Seq(Com).And(term).Trans(ast.PopSingle)

	argRegs := []asm.Dested{
		i().Mov().Rdi(),
		i().Mov().Rsi(),
		i().Mov().Rdx(),
		i().Mov().Rcx(),
		i().Mov().R8(),
		i().Mov().R9(),
	}

	args := oi().Or(ai().And(term).Rep(&commed).SetEval(
		func(nodes []*ast.AST, code asm.Code) {
			if len(nodes) > 6 {
				panic("too many arguments")
			}

			// Evaluate args from right to left and push into the stack.
			for j := range nodes {
				nodes[len(nodes)-j-1].Eval(code)
			}

			for j := range nodes {
				code.Ins(i().Pop().Rax()).Ins(argRegs[j].Rax())
			}
		}).P(), &null)

	return ai().And(&funcName).Seq(LPar).And(&args).Seq(RPar).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			checkNodeCount(nodes, 2)
			nodes[1].Eval(code)
			nodes[0].Eval(code)
			code.Ins(i().Push().Rax())
		})
}

func funcDefiner(bodyer func(*SymTable) Compiler) Compiler {
	var st = new(SymTable)

	argRegs := []asm.Fin{
		i().Mov().Rax().P().Rdi(),
		i().Mov().Rax().P().Rsi(),
		i().Mov().Rax().P().Rdx(),
		i().Mov().Rax().P().Rcx(),
		i().Mov().Rax().P().R8(),
		i().Mov().Rax().P().R9(),
	}

	funcLabel := ai().And(CVar).
		SetEval(func(nodes []*ast.AST, code asm.Code) {
			code.Ins(i().Label(nodes[0].Token.Val()))
		})
	argv := ai().And(varDeclarer(st).P()).SetEval(func(nodes []*ast.AST, code asm.Code) {

		nodes[0].Eval(code)

		v := st.Last()

		if v.Seq >= 6 {
			panic("too many arguments")
		}

		code.Ins(
			i().Mov().Rax().Rbp(),
			i().Sub().Rax().Val(v.Addr),
			argRegs[v.Seq])
	})
	args := oi().Or(ai().And(&argv).
		Rep(ai().Seq(Com).And(&argv).Trans(ast.PopSingle).P()).P(), &null)

	return ai().Seq(Intd).And(&funcLabel).
		Seq(LPar).And(&args).Seq(RPar).
		And(prologuer(st).P()).Seq(LBrc).And(bodyer(st).P()).Seq(RBrc).
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
		var term, muls, adds, expr, bodies Compiler

		ptrRef := ptrDeRefer(st, lvIdenter(st).P())
		val := oi().Or(rvAddrer(lvIdenter(st).P()).P(), rvIdenter(st, &ptrRef).P())
		adds, muls, ptrAdd := addsubs(&muls), muldivs(&term), ptrAdder(st, &val, &expr)

		term = oi().Or(
			ai().And(&ptrAdd, &deRefer).P(),
			&numInt, syscaller(&expr).P(), funcCaller(&expr).P(), &val,
			ai().Seq(LPar).And(&adds).Seq(RPar).Trans(ast.PopSingle).P())
		expr = oi().Or(assigner(oi().Or(&ptrAdd, &ptrRef).P(), &expr).P(), eqneqs(&adds).P())
		semi := ai().And(&expr).Seq(Semi)

		bodies = ai().Rep(oi().Or(
			ifer(&expr, &bodies).P(), forer(&expr, &bodies).P(),
			whiler(&expr, &bodies).P(), returner(&semi).P(),
			ai().And(arrayDeclarer(varDeclarer(st).P(), st).P()).Seq(Semi).P(),
			ai().And(&semi, &popRax).P()).P())

		return ai().And(&bodies)
	}

	return ai().And(ai().Rep(funcDefiner(body).P()).P()).Seq(EOF)
}
