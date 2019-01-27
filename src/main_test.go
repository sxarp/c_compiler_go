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

func TestComp(t *testing.T) {
	r := "4"

	expected := `
.intel_syntax noprefix
.global main

main:
        mov rax, 42
        ret
`

	expects(t, expected, compile(r))
}

func TestStatusCode(t *testing.T) {
	asm := compile("42")

	expects(t, execCode(t, asm), "42")

}
