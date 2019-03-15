package gen

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/tp"
)

type Var struct {
	name string
	tp   tp.Type
}

type varProj struct {
	addr int
	tp   tp.Type
	seq  int
}

type SymTable struct {
	vars []Var
}

func newST() *SymTable {
	return &SymTable{}
}

const bpSize = 8

func (st *SymTable) find(name string) (*varProj, bool) {
	addr := bpSize
	for i, v := range st.vars {
		if v.name == name {
			return &varProj{seq: i, tp: v.tp, addr: addr}, true
		}
		addr += v.tp.Size()
	}

	return nil, false
}

func (st *SymTable) Count() int {
	return len(st.vars)
}

func (st *SymTable) Allocated() int {
	alloc := 0
	for _, val := range st.vars {
		alloc += val.tp.Size()
	}

	return alloc
}

// Get addr of symbol.
func (st *SymTable) AddrOf(name string) int {
	if ref, ok := st.find(name); ok {
		return ref.addr
	} else {
		panic(fmt.Sprintf("%s is not declared.", name))
	}
}

func (st *SymTable) RefOf(name string) int {
	if ref, ok := st.find(name); ok {
		return ref.seq
	} else {
		panic(fmt.Sprintf("%s is not declared.", name))
	}
}

// Declare symbol.
func (st *SymTable) DecOf(name string, t tp.Type) {
	if _, ok := st.find(name); ok {
		panic(fmt.Sprintf("%s is already declared.", name))
	} else {
		st.vars = append(st.vars, Var{tp: t, name: name})
	}
}
