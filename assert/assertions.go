package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func AssertEqualT(t *testing.T, ex interface{}, ac interface{}, msg string) bool {
	if ex != ac {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: %s expected '%v' actual '%v'", AssertionAt(file), line, msg, ex, ac)
		return false
	}
	return true
}

func AssertionAt(file string) string {
	var short string
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return short
}

func AssertEqual(ex interface{}, ac interface{}, msg string) error {
	if ex != ac {
		return fmt.Errorf("%s expected '%v' actual '%v'", msg, ex, ac)
	}
	return nil
}

func AssertEqualSB(ex interface{}, ac interface{}, msg string, errors *strings.Builder) bool {
	if ex != ac {
		_, file, line, _ := runtime.Caller(1)
		fmt.Fprintf(errors, "\n%s:%d: %s expected '%v' actual '%v'", AssertionAt(file), line, msg, ex, ac)
		return false
	}
	return true
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
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: %s '%s' is nil", AssertionAt(file), line, msg, a)
		return false
	}
	return true
}

func AssertNilT(t *testing.T, a interface{}, msg string) bool {
	if reflect.TypeOf(a) != nil {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: %s '%v' not nil", AssertionAt(file), line, msg, a)
		return false
	}
	return true
}
