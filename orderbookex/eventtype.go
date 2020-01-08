package orderbookex

type EventType int

const (
	EventTypeNewOrderSingle EventType = iota
	EventTypeNewOrderAck
	EventTypeRejected
	EventTypeCancel
	EventTypeCancelled
	EventTypeReplaced
	EventTypeCancelAck
	EventTypeCancelRejected
	EventTypePartialFill
	EventTypeFill
	EventTypeDoneForDay
	EventTypeExpired
	EventTypeDone
	EventTypeRestated
	EventTypeUnknown
)

func EventTypeConv(thetype string) EventType {
	switch thetype {
	case "NewOrder":
		return EventTypeNewOrderSingle
	case "NewOrderAck":
		return EventTypeNewOrderAck
	case "PartiallyFilled":
		return EventTypePartialFill
	case "Filled":
		return EventTypeFill
	case "Cancel":
		return EventTypeCancel
	case "Cancelled":
		return EventTypeCancelled
	case "CancelRejected":
		return EventTypeCancelRejected
	case "Rejected":
		return EventTypeRejected
	case "Replaced":
		return EventTypeReplaced
	case "Done":
		return EventTypeDone
	case "Expired":
		return EventTypeExpired
	case "DoneForDay":
		return EventTypeDoneForDay
	}
	return EventTypeUnknown
}

func (it EventType) String() string {
	switch it {
	case EventTypeNewOrderSingle:
		return "EventTypeNewOrderSingle"
	case EventTypeNewOrderAck:
		return "EventTypeNewOrderAck"
	case EventTypeRejected:
		return "EventTypeRejected"
	case EventTypeCancel:
		return "EventTypeCancel"
	case EventTypeCancelled:
		return "EventTypeCancelled"
	case EventTypeReplaced:
		return "EventTypeReplaced"
	case EventTypeCancelAck:
		return "EventTypeCancelAck"
	case EventTypeCancelRejected:
		return "EventTypeCancelRejected"
	case EventTypePartialFill:
		return "EventTypePartialFill"
	case EventTypeFill:
		return "EventTypeFill"
	case EventTypeDoneForDay:
		return "EventTypeDoneForDay"
	case EventTypeExpired:
		return "EventTypeExpired"
	case EventTypeDone:
		return "EventTypeDone"
	case EventTypeRestated:
		return "EventTypeRestated"
	}
	return "EventTypeUnknown"
}
