package orderbookex

import (
	"encoding/hex"
	"fmt"
	"github.com/andres-erbsen/clock"
	"github.com/google/uuid"
	"log"
	"orderbook/test"
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

func makeAuctionLimitOrder(clientID string, clOrdID string, side Side, qty int64, price float64, clock clock.Clock) NewOrderSingle {
	dt := test.NewTime(11, 11, 1)
	return NewNewOrder(
		inst,
		clientID,
		clOrdID,
		side,
		price,
		qty,
		TimeInForceGoodForAuction,
		dt,
		clock.Now(),
		OrderTypeLimit)
}
