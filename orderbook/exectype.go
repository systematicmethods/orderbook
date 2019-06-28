package orderbook

type ExecType rune

const (
	ExecTypeNew                            ExecType = '0'
	ExecTypeDoneForDay                              = '3'
	ExecTypeCanceled                                = '4'
	ExecTypeReplaced                                = '5'
	ExecTypePendingCancel                           = '6'
	ExecTypeStopped                                 = '7'
	ExecTypeRejected                                = '8'
	ExecTypeSuspended                               = '9'
	ExecTypePendingNew                              = 'A'
	ExecTypeCalculated                              = 'B'
	ExecTypeExpired                                 = 'C'
	ExecTypeRestated                                = 'D'
	ExecTypePendingReplace                          = 'E'
	ExecTypeTrade                                   = 'F'
	ExecTypeTradeCorrect                            = 'G'
	ExecTypeTradeCancel                             = 'H'
	ExecTypeOrderStatus                             = 'I'
	ExecTypeTradeInAClearingHold                    = 'J'
	ExecTypeTradeHasBeenReleasedToClearing          = 'K'
	ExecTypeTriggeredOrActivatedBySystem            = 'L'
	ExecTypeLocked                                  = 'M'
	ExecTypeReleased                                = 'N'
	ExecTypeUnknown                                 = 'x'
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
	}
	return "Unknown"
}
