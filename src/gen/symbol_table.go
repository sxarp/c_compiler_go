package gen

type SymTable struct {
	table map[string]int
}

func newST() *SymTable {
	return &SymTable{make(map[string]int)}
}

func (st *SymTable) Count() int {
	return len(st.table)
}

func (st *SymTable) RefOf(s string) int {
	if ref, ok := st.table[s]; ok {
		return ref
	} else {
		st.table[s] = st.Count()
		return st.RefOf(s)
	}
}
