package instrument

import (
	"orderbook/assert"
	"testing"
)

func Test_NewInstrument(m *testing.T) {
	ins := MakeInstrument("ABV", "ABV Investments")
	assert.AssertEqualT(m, ins.ID(), "ABV", "NewInstrument")
	assert.AssertEqualT(m, ins.Name(), "ABV Investments", "NewInstrument")
}
