package main

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"testing"
)

func expects(t *testing.T, expected, expect string) {
	if expected != expect {
		t.Errorf("Expected %s, got %s", expected, expect)
	}

}

func expectt(t *testing.T, expected, expect bool) {
	if expected != expect {
		t.Errorf("Expected %v, got %v", expected, expect)
	}

}

func execCode(t *testing.T, code string) string {

	b := []byte(code)

	if err := ioutil.WriteFile("../tmp/src.s", b, 0644); err != nil {
		t.Errorf("Failed to put asm file.")
	}

	if err := exec.Command("gcc", "-o", "../tmp/tmp", "../tmp/src.s").Run(); err != nil {
		t.Errorf("Failed to comple: %s", err)
	}

	err := exec.Command("../tmp/tmp").Run()

	re := regexp.MustCompile("[0-9]+")
	res := re.FindString(err.Error())

	return res
}

func compare(t *testing.T, expected, code string) {

	expects(t, expected, execCode(t, compile(code)))

}

func TestComp(t *testing.T) {
	r := "42"

	expected := `
.intel_syntax noprefix
.global main

main:
        mov rax, 42
        ret
`

	expects(t, expected, compile(r))
}

func TestByCamperation(t *testing.T) {
	compare(t, "42", "42")
	compare(t, "41", "41")
	compare(t, "1", "1")

}

func ExpectToken(t *testing.T, tt TokenType, inputStr, expectedVal, expectedStr string) {
	tok, outputStr := tt.match(inputStr)
	outputVal := tok.val()

	if outputVal != expectedVal || outputStr != expectedStr {
		t.Errorf("Expected %s, %s for input %s, got %s, %s.", expectedVal, expectedStr, inputStr, outputVal, outputStr)

	}

}

func TestTokenize(t *testing.T) {
	ExpectToken(t, TPlus, "+", "+", "")
	ExpectToken(t, TPlus, "+1234", "+", "1234")
	ExpectToken(t, TPlus, "1+1234", "FAIL", "1+1234")
	ExpectToken(t, TPlus, "", "FAIL", "")

	ExpectToken(t, TMinus, "-", "-", "")
	ExpectToken(t, TMinus, "-123", "-", "123")
	ExpectToken(t, TMinus, "1-1234", "FAIL", "1-1234")
	ExpectToken(t, TMinus, "", "FAIL", "")

	ExpectToken(t, TInt, "1", "1", "")
	ExpectToken(t, TInt, "123", "123", "")
	ExpectToken(t, TInt, "123x", "123", "x")
	ExpectToken(t, TInt, "x123", "FAIL", "x123")
	ExpectToken(t, TInt, "", "FAIL", "")

}
