package instrument

type Instrument interface {
	ID() string
	Name() string
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

func (i *instrument) ID() string {
	return i.id
}

func (i *instrument) Name() string {
	return i.name
}
