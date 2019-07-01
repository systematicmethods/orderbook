package assert

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func AssertEqualT(t *testing.T, a interface{}, b interface{}, msg string) bool {
	if a != b {
		t.Errorf("%s was '%v' != '%v'", msg, a, b)
		return false
	}
	return true
}

func AssertEqual(a interface{}, b interface{}, msg string) error {
	if a != b {
		return fmt.Errorf("%s was '%v' != '%v'", msg, a, b)
	}
	return nil
}

func AssertEqualSB(ex interface{}, ac interface{}, msg string, errors *strings.Builder) {
	if ex != ac {
		fmt.Fprintf(errors, "%s expected '%v' actual '%v'", msg, ex, ac)
	}
}

func AssertEqualTD(t *testing.T, a interface{}, b interface{}, msg string) {
	fmt.Printf("debug a='%v' b='%v'", a, b)
	AssertEqualT(t, a, b, msg)
}

func AssertNotEqualT(t *testing.T, a interface{}, b interface{}, msg string) {
	if a == b {
		t.Errorf("%s %v == %v", msg, a, b)
	}
}

func AssertNotNilT(t *testing.T, a interface{}, msg string) bool {
	if reflect.TypeOf(a) == nil {
		t.Errorf("%s '%s' is nil", msg, a)
		return false
	}
	return true
}

func AssertNilT(t *testing.T, a interface{}, msg string) bool {
	if reflect.TypeOf(a) != nil {
		t.Errorf("%s '%v' not nil", msg, a)
		return false
	}
	return true
}
