package assert

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
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

func AssertTrueT(t *testing.T, value bool, msg string) bool {
	if !value {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: %s", AssertionAt(file), line, msg)
		return false
	}
	return true
}

func Fail(t *testing.T, msg string) {
	_, file, line, _ := runtime.Caller(1)
	t.Errorf("\n%s:%d: %s", AssertionAt(file), line, msg)
}

func AssertEqualTfloat64(t *testing.T, ex float64, ac float64, epsilon float64, msg string) bool {
	if (math.Abs(ex - ac)) > epsilon {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: %s expected '%v' actual '%v' abs '%v'", AssertionAt(file), line, msg, ex, ac, math.Abs(ex-ac))
		return false
	}
	return true
}

func AssertEqualTdecimal(t *testing.T, ex decimal.Decimal, ac decimal.Decimal, epsilon float64, msg string) bool {
	if ex.Sub(ac).Abs().GreaterThan(decimal.NewFromFloat(epsilon)) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: %s expected '%v' actual '%v' abs '%v'", AssertionAt(file), line, msg, ex, ac, ex.Sub(ac).Abs().String())
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
	ty := reflect.TypeOf(a)
	_, file, line, _ := runtime.Caller(1)
	if ty != nil && ty.Kind() == reflect.Ptr {
		x := reflect.ValueOf(a)
		//fmt.Printf("\n*a='%v'\n", x.Pointer())
		if x.Pointer() != 0 {
			t.Errorf("\n%s:%d: %s '%v' not nil %v", AssertionAt(file), line, msg, a, reflect.TypeOf(a))
			return false
		}
		return true
	}
	if reflect.TypeOf(a) != nil {
		t.Errorf("\n%s:%d: %s '%v' not nil %v", AssertionAt(file), line, msg, a, reflect.TypeOf(a))
		return false
	}
	return true
}
