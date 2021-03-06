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
	h.ExpectEq(t, 1, st.RefOf("b").Seq)
	h.ExpectEq(t, 0, st.RefOf("a").Seq)
	h.ExpectEq(t, "a", st.LastRef())
	h.ExpectEq(t, 2, st.RefOf("c").Seq)
	h.ExpectEq(t, 24, st.RefOf("c").Addr)
	h.ExpectEq(t, 3, st.Count())
	h.ExpectEq(t, 24, st.Allocated())
	h.ExpectEq(t, 24, st.Last().Addr)

	al := 3
	st.OverWrite("b", tp.Array(tp.Int, al))
	h.ExpectEq(t, 32, st.RefOf("b").Addr)
	h.ExpectEq(t, 16+al*tp.Int.Size(), st.RefOf("c").Addr)
	h.ExpectEq(t, 32, st.RefOf("b").Addr)
}
