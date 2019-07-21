package main

import (
	"fmt"
	"testing"

	"github.com/sxarp/c_compiler_go/src/h"
)

func compare(t *testing.T, expected, code string) {
	t.Helper()
	h.ExpectEq(t, expected, h.ExecCode(t,
		compile(fmt.Sprintf("int main(){%s}", code)), "../tmp", "src"))
}

func compareMF(t *testing.T, expected, code string) {
	t.Helper()
	h.ExpectEq(t, expected, h.ExecCode(t,
		compile(code), "../tmp", "src"))
}

func TestComp(t *testing.T) {
	r := "int main(){return 42;}"

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
	compare(t, "28", "int a; return a = 28;")
	compare(t, "15", "int z; return z = 28 + 13 - 13 * 2;")

	compare(t, "42", "2;return 42;")
	compare(t, "15", "int z; z = 28 + 13 - 13 * 2; 5; return z;")

	compare(t, "24", "int a; int b; int c; a = 1; b = a+1; c = b+1; return 8*c;")

	compare(t, "2", "int a; int b; int c;return a = b = c = 1+1;")

	compare(t, "1", "int a; int b; a = b = 1;return a == b == 1;")
	compare(t, "1", "int a; int b; a = b = 1;return a == b + 1 == 0;")
	compare(t, "1", "int a; int b; a = b = 1;return a != b + 1 == 1;")
	compare(t, "1", "int a; int b;a = b = 1;return a != b == 0;")

	compare(t, "6", "int a; a=5;return (1+a);return 10;")

	// Only 8bits are available for the parent processes, then exit codes are in 0 ~ 255.
	// https://unix.stackexchange.com/questions/418784/what-is-the-min-and-max-values-of-exit-codes-in-linux
	compare(t, "13",
		`
int a; int b; int c;
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
int alpha; int beta; int gamma; int lhs; int rhs;
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
int main(){
int a;
a = 2;
return add(2, 3)+a;
}
int add(int a, int b){
int c;
c = 1;
return a+b + sub(a, b) -c;
}
int sub(int a, int b){
int d;
d = 4;
return a-b + 4;}
`)

	compareMF(t, "89",
		`
int main(){
int x;
x = 10;
return fib(x);
}
int fib(int x){
if (x == 1) { return 1;}
if (x== 2) { return 2;}
return fib(x-1) + fib(x-2);
}
`)

	compareMF(t, "10",
		`
int main(){
int a;
int b;
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
int main(){
int n; int a;int b;int c;
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
int main(){
int a; int b; int c; int n;
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

	compareMF(t, "9",
		`
int main(){
int a; int *b; int **c;
a = 3;
b = &a;
c = &b;
return a + *b + **c;
}
`)

	compareMF(t, "2",
		`
int main(){
int a;
a = 0;
inc(&a);
inc(&a);
return a;
}
int inc(int *x) { *x = *x + 1; return 0;}
`)

	compareMF(t, "22",
		`
int main(){
int b;
int a;

int* ap;
ap = &a;

a = 1; b= 22;

return *(ap + 1);
}
`)

	compareMF(t, "22",
		`
int main(){
int b;
int a;
int c; c = 1;

int* ap;
ap = &a;

a = 1; b= 22;

return *(ap + c);
}
`)

	compareMF(t, "2",
		`
int main(){
int ret;
int val; val = 156;
int size; size = 2;
ret = syscall 1 1 val size;
return ret;
}
`)

	compareMF(t, "199",
		`
int main(){
int a; a = 199;
int array[10];
int *ap;
ap = &ap;
return *(ap+11);
}
`)
	compareMF(t, "199",
		`
int main(){
int array[10];
*(array + 9) = 199;
return *(array+9);
}`)
}
