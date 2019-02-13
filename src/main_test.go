package main

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"testing"

	"github.com/sxarp/c_compiler_go/src/h"
)

func execCode(t *testing.T, code string) string {

	b := []byte(code)

	if err := ioutil.WriteFile("../tmp/src.s", b, 0644); err != nil {
		t.Errorf("Failed to put asm file.")
	}

	if err := exec.Command("gcc", "-o", "../tmp/tmp", "../tmp/src.s").Run(); err != nil {
		t.Errorf("Failed to comple: %s", err)
	}

	err := exec.Command("../tmp/tmp").Run()

	// Run returns nil when exit code is 0.
	if err == nil {
		return "0"

	}

	re := regexp.MustCompile("[0-9]+")
	res := re.FindString(err.Error())

	return res
}

func compare(t *testing.T, expected, code string) {

	h.Expects(t, expected, execCode(t, compile(code)))

}

func TestComp(t *testing.T) {
	r := "42"

	expected := `
.intel_syntax noprefix
.global main

main:
        push 42
        pop rax
        ret
`

	h.Expects(t, expected, compile(r))
}

func TestByCamperation(t *testing.T) {
	compare(t, "42", "42")
	compare(t, "41", "41")
	compare(t, "1", "1")
	compare(t, "41", " 41 ")
	compare(t, "2", "1 + 1")
	compare(t, "0", "1 - 1")
	compare(t, "2", "1 - 5 + 6")
	compare(t, "3", "1 - 2 + 3 -4 + 5")
	compare(t, "3", "7 - (2 + 3) -4 + 5")
	compare(t, "9", "1 + (2 - 1) - (1 - (3 + 5) )")

	// Only 8bits are available for the parent processes, then exit codes are in 0 ~ 255.
	// https://unix.stackexchange.com/questions/418784/what-is-the-min-and-max-values-of-exit-codes-in-linux
	compare(t, "249", "1 - 2 + 3 -4 + 5 - 10")

}

func TestPreprocess(t *testing.T) {
	h.Expects(t, preprocess(""), "")
	h.Expects(t, preprocess("3"), "3")
	h.Expects(t, preprocess("3 "), "3")
	h.Expects(t, preprocess(" 12 3"), "123")

}
