package tp

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/h"
)

func TestEq(t *testing.T) {
	h.ExpectEq(t, Int.Eq(Int), true)
	h.ExpectEq(t, Int.Ptr().Eq(Int), false)
	h.ExpectEq(t, Int.Eq(Int.Ptr()), false)
	h.ExpectEq(t, Int.Ptr().Eq(Int.Ptr()), true)
	h.ExpectEq(t, Int.Ptr().Ptr().Eq(Int.Ptr()), false)
	h.ExpectEq(t, Int.Ptr().Eq(Int.Ptr().Ptr()), false)
	h.ExpectEq(t, Int.Ptr().Ptr().Eq(Int.Ptr().Ptr()), true)
}

func TestAddUnit(t *testing.T) {
	h.ExpectEq(t, Int.AddUnit(), 1)
	h.ExpectEq(t, Int.Ptr().AddUnit(), 8)
}

func TestDeRef(t *testing.T) {
	tipe, ok := Int.Ptr().Ptr().DeRef()
	h.ExpectEq(t, true, ok)

	tipe, ok = tipe.DeRef()
	h.ExpectEq(t, true, ok)

	tipe, ok = tipe.DeRef()
	h.ExpectEq(t, false, ok)
}
