package h

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"testing"
)

func ExpectEq(t *testing.T, expectedValue, gotValue interface{}) {
	t.Helper()

	switch expected := expectedValue.(type) {
	case string:
		if got := gotValue.(string); got != expected {
			t.Errorf("Expected [%s], got [%s]", expected, got)
		}
	case bool:
		if got := gotValue.(bool); got != expected {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	case int:
		if got := gotValue.(int); got != expected {
			t.Errorf("Expected %d, got %d", expected, got)
		}
	default:
		panic("invalid type value is passed")
	}
}

func ExecCode(t *testing.T, code string, path, fn string) string {
	t.Helper()

	b := []byte(code)

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s.s", path, fn), b, 0644); err != nil {
		t.Errorf("Failed to put asm file.")
	}

	if err := exec.Command("gcc", "-o", fmt.Sprintf("%s/%s.o", path, fn), fmt.Sprintf("%s/%s.s", path, fn)).Run(); err != nil {
		t.Errorf("Failed to comple: %s", err)
	}

	err := exec.Command(fmt.Sprintf("%s/%s.o", path, fn)).Run()

	// Run returns nil when exit code is 0.
	if err == nil {
		return "0"

	}

	re := regexp.MustCompile("[0-9]+")
	res := re.FindString(err.Error())

	return res
}
