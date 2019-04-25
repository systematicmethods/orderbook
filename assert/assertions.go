package assert

import (
	"fmt"
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}, msg string) {
	if a != b {
		t.Errorf("%s was %v != %v", msg, a, b)
	}
}

func AssertEqualD(t *testing.T, a interface{}, b interface{}, msg string) {
	fmt.Printf("debug a=%v b=%v", a, b)
	AssertEqual(t, a, b, msg)
}

func AssertNotEqual(t *testing.T, a interface{}, b interface{}, msg string) {
	if a == b {
		t.Errorf("%s %v == %v", msg, a, b)
	}
}

func AssertNotNil(t *testing.T, a interface{}, msg string) bool {
	if reflect.ValueOf(a).IsNil() {
		t.Errorf("%s %s is nil", msg, a)
		return false
	}
	return true
}

func AssertNil(t *testing.T, a interface{}, msg string) bool {
	if !reflect.ValueOf(a).IsNil() {
		t.Errorf("%s %v not nil", msg, a)
		return false
	}
	return true
}
