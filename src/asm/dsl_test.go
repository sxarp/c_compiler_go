package asm

import "testing"

func expects(t *testing.T, expected, expect string) {
	if expected != expect {
		t.Errorf("Expected %s, got %s", expected, expect)
	}

}

func TestIns(t *testing.T) {
	expects(t, "        ret", I().Ret().str())
	expects(t, "        mov rax, 42", I().Mov().Rax().Val(42).str())
	expects(t, "        mov rax, rax", I().Mov().Rax().Rax().str())
}
