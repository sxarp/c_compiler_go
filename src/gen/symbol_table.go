package gen

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/tp"
)

type Var struct {
	addr int
	seq  int
	tp   tp.Type
}

type SymTable struct {
	table map[string]Var
}

func newST() *SymTable {
	return &SymTable{make(map[string]Var)}
}

func (st *SymTable) Count() int {
	return len(st.table)
}

func (st *SymTable) Allocated() int {
	alloc := 0
	for _, val := range st.table {
		alloc += val.tp.Size()
	}

	return alloc
}

const bpSize = 8

// Get addr of symbol.
func (st *SymTable) AddrOf(s string) int {
	if ref, ok := st.table[s]; ok {
		return ref.addr + bpSize
	} else {
		panic(fmt.Sprintf("%s is not declared.", s))
	}
}

func (st *SymTable) RefOf(s string) int {
	if ref, ok := st.table[s]; ok {
		return ref.seq
	} else {
		panic(fmt.Sprintf("%s is not declared.", s))
	}
}

// Declare symbol.
func (st *SymTable) DecOf(s string, t tp.Type) {
	if _, ok := st.table[s]; ok {
		panic(fmt.Sprintf("%s is already declared.", s))
	} else {
		st.table[s] = Var{addr: st.Allocated(), tp: t, seq: st.Count()}
	}
}
