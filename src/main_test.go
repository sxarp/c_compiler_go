package main

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"testing"
)

func expects(t *testing.T, expected, expect string) {
	if expected != expect {
		t.Errorf("Expected [%s], got [%s]", expected, expect)
	}

}

func expectt(t *testing.T, expected, expect bool) {
	if expected != expect {
		t.Errorf("Expected [%v], got [%v].", expected, expect)
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
	compare(t, "41", " 41 ")

}

func TestPreprocess(t *testing.T) {
	expects(t, preprocess(""), "")
	expects(t, preprocess("3"), "3")
	expects(t, preprocess("3 "), "3")
	expects(t, preprocess(" 12 3"), "123")

}
