package fixmodel

import (
	"github.com/google/uuid"
	"time"
)

type OrderCancelRequest struct {
	instrumentID string
	clientID     string
	clOrdID      string

	side         Side
	origClOrdID  string
	transactTime time.Time

	eventType EventType
	orderID   string
}

func NewOrderCancelRequest(
	instrumentID string,
	clientID string,
	clOrdID string,
	side Side,
	origClOrdID string,
	transactTime time.Time) *OrderCancelRequest {
	theOrderID, _ := uuid.NewUUID()
	return &OrderCancelRequest{
		instrumentID,
		clientID,
		clOrdID,
		side,
		origClOrdID,
		transactTime,
		EventTypeCancel,
		theOrderID.String(),
	}
}

func (p *OrderCancelRequest) InstrumentID() string {
	return p.instrumentID
}

func (p *OrderCancelRequest) ClientID() string {
	return p.clientID
}

func (p *OrderCancelRequest) ClOrdID() string {
	return p.clOrdID
}

func (p *OrderCancelRequest) Side() Side {
	return p.side
}

func (p *OrderCancelRequest) OrigClOrdID() string {
	return p.origClOrdID
}

func (p *OrderCancelRequest) TransactTime() time.Time {
	return p.transactTime
}

func (p *OrderCancelRequest) OrderID() string {
	return p.orderID
}

func (p *OrderCancelRequest) isBuy() bool {
	return p.Side() == SideBuy
}
