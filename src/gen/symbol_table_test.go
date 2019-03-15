package gen

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/h"
	"github.com/sxarp/c_compiler_go/src/tp"
)

func TestDecOf(t *testing.T) {
	st := newST()
	st.DecOf("a", tp.Int)
	st.DecOf("b", tp.Int)
	st.DecOf("c", tp.Int)
	h.ExpectEq(t, 1, st.RefOf("b"))
	h.ExpectEq(t, 0, st.RefOf("a"))
	h.ExpectEq(t, 2, st.RefOf("c"))
	h.ExpectEq(t, 24, st.AddrOf("c"))
	h.ExpectEq(t, 3, st.Count())
	h.ExpectEq(t, 24, st.Allocated())
}
