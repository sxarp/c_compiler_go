package main

import (
	"testing"

	"github.com/sxarp/c_compiler_go/src/h"
)

func compare(t *testing.T, expected, code string) {
	t.Helper()
	h.ExpectEq(t, expected, h.ExecCode(t, compile(code), "../tmp", "src"))
}

func TestComp(t *testing.T) {
	r := "42;"

	expected := `
.intel_syntax noprefix
.global main

main:
        push rbp
        mov rbp, rsp
        sub rsp, 0
        push 42
        pop rax
        mov rsp, rbp
        pop rbp
        ret
`

	h.ExpectEq(t, expected, compile(r))
}

func TestByCamperation(t *testing.T) {
	compare(t, "42", "42;")
	compare(t, "41", "41;")
	compare(t, "1", "1;")
	compare(t, "41", " 41 ;")
	compare(t, "2", "1 + 1;")
	compare(t, "0", "1 - 1;")
	compare(t, "2", "1 - 5 + 6;")
	compare(t, "3", "1 - 2 + 3 -4 + 5;")
	compare(t, "3", "7 - (2 + 3) -4 + 5;")
	compare(t, "9", "1 + (2 - 1) - (1 - (3 + 5) );")
	compare(t, "2", "1 * (2 - 1) - (1 - (10 / 5) );")
	compare(t, "10", "1 * (2 / 1) * 8 - 6;")
	compare(t, "28", "a = 28;")
	compare(t, "15", "z = 28 + 13 - 13 * 2;")

	compare(t, "42", "2;42;")
	compare(t, "15", "z = 28 + 13 - 13 * 2; 5; c=15;")
	compare(t, "15", "z = 28 + 13 - 13 * 2; 5; z;")

	compare(t, "24", "a = 1; b = a+1; c = b+1; 8*c;")

	compare(t, "2", "a = b = c = 1+1;")

	compare(t, "1", "a = b = 1;a == b == 1;")
	compare(t, "1", "a = b = 1;a == b + 1 == 0;")
	compare(t, "1", "a = b = 1;a != b + 1 == 1;")
	compare(t, "1", "a = b = 1;a != b == 0;")

	compare(t, "6", "a=5;return (1+a);10;")

	// Only 8bits are available for the parent processes, then exit codes are in 0 ~ 255.
	// https://unix.stackexchange.com/questions/418784/what-is-the-min-and-max-values-of-exit-codes-in-linux
	compare(t, "249", "1 - 2 + 3 -4 + 5 - 10;")
	compare(t, "13",
		`
a = b = 1;
c = a; a = a + b; b = c;
c = a; a = a + b; b = c;
c = a; a = a + b; b = c;
c = a; a = a + b; b = c;
c = a; a = a + b; b = c;
a;
`)
	compare(t, "0",
		`
alpha = 32;
beta = 11;
gamma = 28;
lhs = (alpha + beta + gamma)*(alpha * beta + beta * gamma + gamma * alpha) - alpha * beta * gamma;
rhs = (alpha + beta)*(beta + gamma)*(gamma + alpha);
lhs - rhs;
`)

}
