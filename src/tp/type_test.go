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
