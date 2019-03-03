package gen

import "fmt"

type SymTable struct {
	table map[string]int
}

func newST() *SymTable {
	return &SymTable{make(map[string]int)}
}

func (st *SymTable) Count() int {
	return len(st.table)
}

// Get reference of symbol.
func (st *SymTable) RefOf(s string) int {
	if ref, ok := st.table[s]; ok {
		return ref
	} else {
		panic(fmt.Sprintf("%s is not declared.", s))
	}
}

// Declare symbol.
func (st *SymTable) DecOf(s string) {
	if _, ok := st.table[s]; ok {
		panic(fmt.Sprintf("%s is already declared.", s))
	} else {
		st.table[s] = st.Count()
	}
}
