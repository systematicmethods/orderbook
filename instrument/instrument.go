package instrument

import "math/big"

type Instrument interface {
	ID() string
	Name() string
	DecimalPlaces() int32
}

func NewInstrument(id string, name string) *instrument {
	return &instrument{id: id, name: name, dp: 2}
}

func NewInstrumentDP(id string, name string, decimalPlaces int32) *instrument {
	return &instrument{id: id, name: name, dp: decimalPlaces}
}

type instrument struct {
	id       string
	name     string
	lotsize  int
	minsize  int
	dp       int32
	rounding big.RoundingMode
}

func (i *instrument) ID() string {
	return i.id
}

func (i *instrument) Name() string {
	return i.name
}

func (i *instrument) DecimalPlaces() int32 {
	return i.dp
}
