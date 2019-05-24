package orderbook

import (
	"github.com/google/uuid"
	"time"
)

type OrderCancelRequest interface {
	InstrumentID() string
	ClientID() string
	ClOrdID() string

	Side() Side
	OrigClOrdID() string
	TransactTime() time.Time

	OrderID() string
}

type orderCancelRequest struct {
	instrumentID string
	clientID     string
	clOrdID      string

	side         Side
	origClOrdID  string
	transactTime time.Time

	eventType EventType
	orderID   string
}

func MakeOrderCancelRequest(
	instrumentID string,
	clientID string,
	clOrdID string,
	side Side,
	origClOrdID string,
	transactTime time.Time) OrderCancelRequest {
	theOrderID, _ := uuid.NewUUID()
	return OrderCancelRequest(&orderCancelRequest{
		instrumentID,
		clientID,
		clOrdID,
		side,
		origClOrdID,
		transactTime,
		EventTypeCancel,
		theOrderID.String(),
	})
}

func (p *orderCancelRequest) InstrumentID() string {
	return p.instrumentID
}

func (p *orderCancelRequest) ClientID() string {
	return p.clientID
}

func (p *orderCancelRequest) ClOrdID() string {
	return p.clOrdID
}

func (p *orderCancelRequest) Side() Side {
	return p.side
}

func (p *orderCancelRequest) OrigClOrdID() string {
	return p.origClOrdID
}

func (p *orderCancelRequest) TransactTime() time.Time {
	return p.transactTime
}

func (p *orderCancelRequest) OrderID() string {
	return p.orderID
}
