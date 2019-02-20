package h

import (
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
