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
	Name string
}

type SymTable struct {
	vars []Var
}

func newST() *SymTable {
	return &SymTable{}
}

func (st *SymTable) find(name string) (*VarProj, bool) {
	addr := 0
	for i, v := range st.vars {
		addr += v.tp.Size()
		if v.name == name {
			return &VarProj{Seq: i, Type: v.tp, Addr: addr, Name: name}, true
		}
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
	ref, ok := st.find(name)
	if ok {
		return ref
	}

	panic(fmt.Sprintf("%s is not declared.", name))
}

func (st *SymTable) Last() *VarProj {
	v, ok := st.find(st.vars[len(st.vars)-1].name)
	if ok {
		return v
	}

	panic("SymTable is empty.")
}

// Declare symbol.
func (st *SymTable) DecOf(name string, t tp.Type) {
	if _, ok := st.find(name); ok {
		panic(fmt.Sprintf("%s is already declared.", name))
	} else {
		st.vars = append(st.vars, Var{tp: t, name: name})
	}
}

func (st *SymTable) OverWrite(name string, nt tp.Type) {
	if vp, ok := st.find(name); ok {
		st.vars[vp.Seq] = Var{tp: nt, name: name}
		return
	}

	panic("found no symbol to overwite")
}
