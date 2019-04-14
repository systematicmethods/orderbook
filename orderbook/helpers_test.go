package orderbook

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func printerror(err error, m *testing.T) {
	if err != nil {
		m.Errorf("Add order failed %v", err)
	}
}

func assertequal(m *testing.T, a interface{}, b interface{}, msg string) {
	if a != b {
		m.Errorf("%s was %s != %s", msg, a, b)
	}
}

func assertNotEqual(m *testing.T, a interface{}, b interface{}, msg string) {
	if a == b {
		m.Errorf("%s %s == %s", msg, a, b)
	}
}

func assertNotNil(m *testing.T, a interface{}, msg string) bool {
	if reflect.ValueOf(a).IsNil() {
		m.Errorf("%s %s is nil", msg, a)
		return false
	}
	return true
}

func assertNil(m *testing.T, a interface{}, msg string) bool {
	if !reflect.ValueOf(a).IsNil() {
		m.Errorf("%s %s not nil", msg, a)
		return false
	}
	return true
}

func newID(uuid uuid.UUID, _ error) uuid.UUID {
	return uuid
}

func dumptime(m *testing.T, id uuid.UUID, msg string) {
	m.Errorf("Time %s %d, %d, %v %s", msg, id.Time(), id.ClockSequence(), id.Version(), hex.Dump(id[:]))
	dumpbytes(id[:])
}

func dumpbytes(b []byte) {
	for _, n := range b[:] {
		fmt.Printf(" %08b", n) // prints 00000000 11111101
	}
	fmt.Println()
}
