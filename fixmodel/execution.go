package fixmodel

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ExecutionReport struct {
	instrumentID string
	clientID     string
	clOrdID      string
	side         Side

	lastQty   int64
	lastPrice float64
	execType  ExecType

	leavesQty   int64
	cumQty      int64
	ordStatus   OrdStatus
	origClOrdID string

	orderID               string
	execID                string
	orderQty              int64
	transactTime          time.Time
	execRestatementReason ExecRestatementReason
	ordRejReason          OrdRejReason
	rejectText            string

	eventType EventType
}

func (b *ExecutionReport) String() string {
	str1 := fmt.Sprintf("ExecutionReport: instrumentID:%s, clientID:%s, clOrdID:%s, side:%v, lastQty:%d, lastPrice:%f, execType:%v, leavesQty:%d, cumQty:%d, ordStatus;%v",
		b.instrumentID,
		b.clientID,
		b.clOrdID,
		b.side,
		b.lastQty,
		b.lastPrice,
		b.execType,
		b.leavesQty,
		b.cumQty,
		OrdStatusToString(b.ordStatus))
	str2 := fmt.Sprintf("origClOrdID:%s, orderID:%s, execID:%s, orderQty:%d, transactTime:%s, ordRejReason:%d, rejectText:%s",
		b.origClOrdID,
		b.orderID,
		b.execID,
		b.orderQty,
		b.transactTime,
		b.ordRejReason,
		b.rejectText,
	)
	return fmt.Sprintf("%s %s", str1, str2)
}

func NewRejectExecutionReport(ord *NewOrderSingle, reason OrdRejReason, rejecttext string) *ExecutionReport {
	theExecID, _ := uuid.NewUUID()
	return &ExecutionReport{
		ord.InstrumentID(),
		ord.ClientID(),
		ord.ClOrdID(),
		ord.Side(),
		0,
		0,
		ExecTypeRejected,
		ord.OrderQty(),
		0,
		OrdStatusRejected,
		"",
		ord.OrderID(),
		theExecID.String(),
		ord.OrderQty(),
		ord.TransactTime(),
		ExecRestatementReasonNone,
		reason,
		rejecttext,
		EventTypeRejected,
	}
}

func NewExecutionReport(
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
	origClOrdID string,
	orderID string,
	execID string,
	orderQty int64,
	transactTime time.Time,
	reason ExecRestatementReason) *ExecutionReport {
	return &ExecutionReport{
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
		origClOrdID,
		orderID,
		execID,
		orderQty,
		transactTime,
		reason,
		OrdRejReasonNotApplicable,
		"",
		eventType,
	}
}

func (e *ExecutionReport) InstrumentID() string {
	return e.instrumentID
}

func (e *ExecutionReport) ClientID() string {
	return e.clientID
}

func (e *ExecutionReport) ClOrdID() string {
	return e.clOrdID
}

func (e *ExecutionReport) Side() Side {
	return e.side
}

func (e *ExecutionReport) LastQty() int64 {
	return e.lastQty
}

func (e *ExecutionReport) LastPrice() float64 {
	return e.lastPrice
}

func (e *ExecutionReport) ExecType() ExecType {
	return e.execType
}

func (e *ExecutionReport) LeavesQty() int64 {
	return e.leavesQty
}

func (e *ExecutionReport) CumQty() int64 {
	return e.cumQty
}

func (e *ExecutionReport) OrdStatus() OrdStatus {
	return e.ordStatus
}

func (e *ExecutionReport) OrigClOrdID() string {
	return e.origClOrdID
}

func (e *ExecutionReport) OrderID() string {
	return e.orderID
}

func (e *ExecutionReport) ExecID() string {
	return e.execID
}

func (e *ExecutionReport) OrderQty() int64 {
	return e.orderQty
}

func (e *ExecutionReport) TransactTime() time.Time {
	return e.transactTime
}

func (e *ExecutionReport) RejectText() string {
	return e.rejectText
}

func (e *ExecutionReport) RejectReason() OrdRejReason {
	return e.ordRejReason
}
