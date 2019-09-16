package orderbook

type ExecType rune

const (
	ExecTypeNew                            ExecType = '0'
	ExecTypeDoneForDay                     ExecType = '3'
	ExecTypeCanceled                       ExecType = '4'
	ExecTypeReplaced                       ExecType = '5'
	ExecTypePendingCancel                  ExecType = '6'
	ExecTypeStopped                        ExecType = '7'
	ExecTypeRejected                       ExecType = '8'
	ExecTypeSuspended                      ExecType = '9'
	ExecTypePendingNew                     ExecType = 'A'
	ExecTypeCalculated                     ExecType = 'B'
	ExecTypeExpired                        ExecType = 'C'
	ExecTypeRestated                       ExecType = 'D'
	ExecTypePendingReplace                 ExecType = 'E'
	ExecTypeTrade                          ExecType = 'F'
	ExecTypeTradeCorrect                   ExecType = 'G'
	ExecTypeTradeCancel                    ExecType = 'H'
	ExecTypeOrderStatus                    ExecType = 'I'
	ExecTypeTradeInAClearingHold           ExecType = 'J'
	ExecTypeTradeHasBeenReleasedToClearing ExecType = 'K'
	ExecTypeTriggeredOrActivatedBySystem   ExecType = 'L'
	ExecTypeLocked                         ExecType = 'M'
	ExecTypeReleased                       ExecType = 'N'
	ExecTypeUnknown                        ExecType = 'x'
)

func ExecTypeConv(thetype string) ExecType {
	switch thetype {
	case "New":
		return ExecTypeNew
	case "Trade":
		return ExecTypeTrade
	case "Cancelled":
		return ExecTypeCanceled
	case "Rejected":
		return ExecTypeRejected
	case "Replaced":
		return ExecTypeReplaced
	case "DoneForDay":
		return ExecTypeDoneForDay
	case "Expired":
		return ExecTypeExpired
	case "Restated":
		return ExecTypeRestated
	}
	return ExecTypeUnknown
}

func ExecTypeToString(thetype ExecType) string {
	switch thetype {
	case ExecTypeNew:
		return "New"
	case ExecTypeTrade:
		return "Trade"
	case ExecTypeCanceled:
		return "Cancelled"
	case ExecTypeRejected:
		return "Rejected"
	case ExecTypeReplaced:
		return "Replaced"
	case ExecTypeDoneForDay:
		return "DoneForDay"
	case ExecTypeExpired:
		return "Expired"
	case ExecTypeRestated:
		return "Restated"
	}
	return "Unknown"
}
