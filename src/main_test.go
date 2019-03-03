package main

import (
	"fmt"
	"testing"

	"github.com/sxarp/c_compiler_go/src/h"
)

func compare(t *testing.T, expected, code string) {
	t.Helper()
	h.ExpectEq(t, expected, h.ExecCode(t,
		compile(fmt.Sprintf("main(){%s}", code)), "../tmp", "src"))
}

func compareMF(t *testing.T, expected, code string) {
	t.Helper()
	h.ExpectEq(t, expected, h.ExecCode(t,
		compile(code), "../tmp", "src"))
}

func TestComp(t *testing.T) {
	r := "main(){return 42;}"

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
	compare(t, "42", "return 42;")
	compare(t, "2", "return 1 + 1;")
	compare(t, "0", "return 1 - 1;")
	compare(t, "3", "return 1 - 2 + 3 -4 + 5;")
	compare(t, "9", "return 1 + (2 - 1) - (1 - (3 + 5) );")
	compare(t, "2", "return 1 * (2 - 1) - (1 - (10 / 5) );")
	compare(t, "28", "return a = 28;")
	compare(t, "15", "return z = 28 + 13 - 13 * 2;")

	compare(t, "42", "2;return 42;")
	compare(t, "15", "z = 28 + 13 - 13 * 2; 5; return z;")

	compare(t, "24", "a = 1; b = a+1; c = b+1; return 8*c;")

	compare(t, "2", "return a = b = c = 1+1;")

	compare(t, "1", "a = b = 1;return a == b == 1;")
	compare(t, "1", "a = b = 1;return a == b + 1 == 0;")
	compare(t, "1", "a = b = 1;return a != b + 1 == 1;")
	compare(t, "1", "a = b = 1;return a != b == 0;")

	compare(t, "6", "a=5;return (1+a);return 10;")

	// Only 8bits are available for the parent processes, then exit codes are in 0 ~ 255.
	// https://unix.stackexchange.com/questions/418784/what-is-the-min-and-max-values-of-exit-codes-in-linux
	compare(t, "13",
		`
a = b = 1;
c = a; a = a + b; b = c;
c = a; a = a + b; b = c;
c = a; a = a + b; b = c;
c = a; a = a + b; b = c;
c = a; a = a + b; b = c;
return a;
`)
	compare(t, "0",
		`
alpha = 32;
beta = 11;
gamma = 28;
lhs = (alpha + beta + gamma)*(alpha * beta + beta * gamma + gamma * alpha) - alpha * beta * gamma;
rhs = (alpha + beta)*(beta + gamma)*(gamma + alpha);
return lhs - rhs;
`)

}

func TestByMF(t *testing.T) {
	compareMF(t, "9",
		`
main(){
a = 2;
return add(2, 3)+a;
}
add(a, b){
c = 1;
return a+b + sub(a, b) -c;
}
sub(a,b){
d = 4;
return a-b + 4;}
`)

	compareMF(t, "89",
		`
main(){
x = 10;
return fib(x);
}
fib(x){
if (x == 1) { return 1;}
if (x== 2) { return 2;}
return fib(x-1) + fib(x-2);
}
`)

	compareMF(t, "10",
		`
main(){
a = 10;
b = 0;
while (a) {
a = a -1;
b = b + 1;
}
return b;
}
`)

	compareMF(t, "89",
		`
main(){
n = 0;
a = 1;
b = 1;

while (n != 9) {
n = n+1;
c = a;
a = a + b;
b = c;
}
return a;
}
`)

	compareMF(t, "89",
		`
main(){
a = 1;
b = 1;
for (n = 0; n != 9; n=n+1) {
c = a;
a = a + b;
b = c;
}
return a;
}
`)
}
