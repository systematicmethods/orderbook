package orderstate

import (
	"github.com/google/uuid"
	"orderbook/fixmodel"
)

func NewExecutionReport(
	ord *OrderState,
	event fixmodel.EventType,
	exec fixmodel.ExecType,
	status fixmodel.OrdStatus,
	lastQty int64,
	lastPrice float64,
	reason fixmodel.ExecRestatementReason,
) *fixmodel.ExecutionReport {
	theExecID, _ := uuid.NewUUID()
	return fixmodel.NewExecutionReport(
		event,
		ord.InstrumentID(),
		ord.ClientID(),
		ord.ClOrdID(),
		ord.Side(),
		lastQty,
		lastPrice,
		exec,
		ord.LeavesQty(),
		ord.CumQty(),
		status,
		"",
		ord.OrderID(),
		theExecID.String(),
		ord.OrderQty(),
		ord.TransactTime(),
		reason,
	)
}

func NewExecutionReportOrigClOrdID(
	ord *OrderState,
	event fixmodel.EventType,
	exec fixmodel.ExecType,
	status fixmodel.OrdStatus,
	lastQty int64,
	lastPrice float64,
	clOrdID string,
	reason fixmodel.ExecRestatementReason,
) *fixmodel.ExecutionReport {
	theExecID, _ := uuid.NewUUID()
	return fixmodel.NewExecutionReport(
		event,
		ord.InstrumentID(),
		ord.ClientID(),
		clOrdID,
		ord.Side(),
		lastQty,
		lastPrice,
		exec,
		ord.LeavesQty(),
		ord.CumQty(),
		status,
		ord.ClOrdID(),
		ord.OrderID(),
		theExecID.String(),
		ord.OrderQty(),
		ord.TransactTime(),
		reason,
	)
}
func NewOrderCancelledExecutionReport(ord *OrderState) *fixmodel.ExecutionReport {
	return NewExecutionReport(ord, fixmodel.EventTypeCancel, fixmodel.ExecTypeCanceled, fixmodel.OrdStatusCanceled, 0, 0, fixmodel.ExecRestatementReasonNone)
}

func NewFillExecutionReport(ord *OrderState, fillPrice float64, qty int64) *fixmodel.ExecutionReport {
	var etype fixmodel.EventType
	if ord.OrdStatus() == fixmodel.OrdStatusFilled {
		etype = fixmodel.EventTypeFill
	} else {
		etype = fixmodel.EventTypePartialFill
	}
	return NewExecutionReport(ord, etype, fixmodel.ExecTypeTrade, ord.ordStatus, qty, fillPrice, fixmodel.ExecRestatementReasonNone)
}

func NewNewOrderAckExecutionReport(ord *OrderState) *fixmodel.ExecutionReport {
	return NewExecutionReport(ord, fixmodel.EventTypeNewOrderAck, fixmodel.ExecTypeNew, ord.ordStatus, 0, 0, fixmodel.ExecRestatementReasonNone)
}

func NewRejectExecutionRepor(ord *OrderState) *fixmodel.ExecutionReport {
	return NewExecutionReport(ord, fixmodel.EventTypeRejected, fixmodel.ExecTypeRejected, fixmodel.OrdStatusRejected, 0, 0, fixmodel.ExecRestatementReasonNone)
}

func NewCancelOrderExecutionReport(ord *OrderState, order *fixmodel.OrderCancelRequest) *fixmodel.ExecutionReport {
	return NewExecutionReportOrigClOrdID(ord, fixmodel.EventTypeCancelAck, fixmodel.ExecTypeCanceled, fixmodel.OrdStatusCanceled, 0, 0, order.ClOrdID(), fixmodel.ExecRestatementReasonNone)
}
