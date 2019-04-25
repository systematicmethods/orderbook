package instrument

type Instrument interface {
	InstrumentID() string
	InstrumentName() string
}

func MakeInstrument(id string, name string) *instrument {
	return &instrument{id: id, name: name}
}

type instrument struct {
	id      string
	name    string
	lotsize int
	minsize int
}

func (i *instrument) InstrumentID() string {
	return i.id
}

func (i *instrument) InstrumentName() string {
	return i.name
}
