package orderbook

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"log"
	"testing"
)

func printerror(err error, m *testing.T) {
	if err != nil {
		m.Errorf("Add order failed %v", err)
	}
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

func loglines() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
