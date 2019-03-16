package gen

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/tp"
)

type Var struct {
	name string
	tp   tp.Type
}

type VarProj struct {
	Addr int
	Type tp.Type
	Seq  int
}

type SymTable struct {
	vars []Var
}

func newST() *SymTable {
	return &SymTable{}
}

const bpSize = 8

func (st *SymTable) find(name string) (*VarProj, bool) {
	addr := bpSize
	for i, v := range st.vars {
		if v.name == name {
			return &VarProj{Seq: i, Type: v.tp, Addr: addr}, true
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

func (st *SymTable) RefOf(name string) *VarProj {
	if ref, ok := st.find(name); ok {
		return ref
	} else {
		panic(fmt.Sprintf("%s is not declared.", name))
	}
}

func (st *SymTable) Last() *VarProj {
	if v, ok := st.find(st.vars[len(st.vars)-1].name); ok {
		return v
	} else {
		panic("SymTable is empty.")
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
