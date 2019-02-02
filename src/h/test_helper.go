package h

import (
	"testing"
)

func Expects(t *testing.T, expected, expect string) {
	if expected != expect {
		t.Errorf("Expected [%s], got [%s]", expected, expect)
	}

}

func Expectt(t *testing.T, expected, expect bool) {
	if expected != expect {
		t.Errorf("Expected [%v], got [%v].", expected, expect)
	}

}
func Expecti(t *testing.T, expected, expect int) {
	if expected != expect {
		t.Errorf("Expected %d, got %d.", expected, expect)
	}

}
