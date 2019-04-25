package instrument

import (
	"orderbook/assert"
	"testing"
)

func Test_NewInstrument(m *testing.T) {
	ins := MakeInstrument("ABV", "ABV Investments")
	assert.AssertEqual(m, ins.InstrumentID(), "ABV", "NewInstrument")
	assert.AssertEqual(m, ins.InstrumentName(), "ABV Investments", "NewInstrument")
}
