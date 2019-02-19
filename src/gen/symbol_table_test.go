package gen

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/h"
)

func TestRefOf(t *testing.T) {
	st := newST()
	h.Expecti(t, 0, st.RefOf("a"))
	h.Expecti(t, 1, st.RefOf("b"))
	h.Expecti(t, 2, st.RefOf("c"))
	h.Expecti(t, 1, st.RefOf("b"))
	h.Expecti(t, 2, st.RefOf("c"))
	h.Expecti(t, 0, st.RefOf("a"))
	h.Expecti(t, 3, st.Count())
}
