package orderbook

import (
	"github.com/google/uuid"
	"time"
)

type ExecutionReport interface {
	InstrumentID() string
	ClientID() string
	ClOrdID() string
	Side() Side

	LastQty() int64
	LastPrice() float64
	ExecType() ExecType

	LeavesQty() int64
	CumQty() int64
	OrdStatus() OrdStatus

	OrderID() string
	ExecID() string
	OrderQty() int64
	TransactTime() time.Time
}

type executionReport struct {
	instrumentID string
	clientID     string
	clOrdID      string
	side         Side

	lastQty   int64
	lastPrice float64
	execType  ExecType

	leavesQty int64
	cumQty    int64
	ordStatus OrdStatus

	orderID      string
	execID       string
	orderQty     int64
	transactTime time.Time

	eventType EventType
}

func MakeNewOrderAckExecutionReport(ord OrderState) ExecutionReport {
	theExecID, _ := uuid.NewUUID()
	return ExecutionReport(&executionReport{
		ord.InstrumentID(),
		ord.ClientID(),
		ord.ClOrdID(),
		ord.Side(),
		0,
		0,
		ExecTypeNew,
		ord.LeavesQty(),
		ord.CumQty(),
		OrdStatusNew,
		ord.OrderID(),
		theExecID.String(),
		ord.OrderQty(),
		ord.TransactTime(),
		EventTypeNewOrderAck,
	})
}

func MakeCancelOrderExecutionReport(ord OrderState, order OrderCancelRequest) ExecutionReport {
	theExecID, _ := uuid.NewUUID()
	return ExecutionReport(&executionReport{
		ord.InstrumentID(),
		ord.ClientID(),
		order.ClOrdID(),
		ord.Side(),
		0,
		0,
		ExecTypeCanceled,
		ord.LeavesQty(),
		ord.CumQty(),
		OrdStatusCanceled,
		ord.OrderID(),
		theExecID.String(),
		ord.OrderQty(),
		ord.TransactTime(),
		EventTypeCancelAck,
	})
}
func MakeExecutionReport(
	eventType EventType,
	instrumentID string,
	clientID string,
	clOrdID string,
	side Side,
	lastQty int64,
	lastPrice float64,
	execType ExecType,
	leavesQty int64,
	cumQty int64,
	ordStatus OrdStatus,
	orderID string,
	execID string,
	orderQty int64,
	transactTime time.Time,
) ExecutionReport {
	return ExecutionReport(&executionReport{
		instrumentID,
		clientID,
		clOrdID,
		side,
		lastQty,
		lastPrice,
		execType,
		leavesQty,
		cumQty,
		ordStatus,
		orderID,
		execID,
		orderQty,
		transactTime,
		eventType,
	})
}

func (e *executionReport) InstrumentID() string {
	return e.instrumentID
}

func (e *executionReport) ClientID() string {
	return e.clientID
}

func (e *executionReport) ClOrdID() string {
	return e.clOrdID
}

func (e *executionReport) Side() Side {
	return e.side
}

func (e *executionReport) LastQty() int64 {
	return e.lastQty
}

func (e *executionReport) LastPrice() float64 {
	return e.lastPrice
}

func (e *executionReport) ExecType() ExecType {
	return e.execType
}

func (e *executionReport) LeavesQty() int64 {
	return e.leavesQty
}

func (e *executionReport) CumQty() int64 {
	return e.cumQty
}

func (e *executionReport) OrdStatus() OrdStatus {
	return e.ordStatus
}

func (e *executionReport) OrderID() string {
	return e.orderID
}

func (e *executionReport) ExecID() string {
	return e.execID
}

func (e *executionReport) OrderQty() int64 {
	return e.orderQty
}

func (e *executionReport) TransactTime() time.Time {
	return e.transactTime
}
