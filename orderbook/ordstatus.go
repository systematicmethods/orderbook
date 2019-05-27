package orderbook

type OrdStatus rune

const (
	OrdStatusNew                OrdStatus = '0'
	OrdStatusPartiallyFilled              = '1'
	OrdStatusFilled                       = '2'
	OrdStatusDoneForDay                   = '3'
	OrdStatusCanceled                     = '4'
	OrdStatusReplaced                     = '5'
	OrdStatusPendingCancel                = '6'
	OrdStatusStopped                      = '7'
	OrdStatusRejected                     = '8'
	OrdStatusSuspended                    = '9'
	OrdStatusPendingNew                   = 'A'
	OrdStatusCalculated                   = 'B'
	OrdStatusExpired                      = 'C'
	OrdStatusAcceptedForBidding           = 'D'
	OrdStatusPendingReplace               = 'E'
	OrdStatusUnknown                      = 'x'
)

func OrdStatusConv(thetype string) OrdStatus {
	switch thetype {
	case "New":
		return OrdStatusNew
	case "PartiallyFilled":
		return OrdStatusPartiallyFilled
	case "Filled":
		return OrdStatusFilled
	case "Cancelled":
		return OrdStatusCanceled
	case "Rejected":
		return OrdStatusRejected
	case "Replaced":
		return OrdStatusReplaced
	case "DoneForDay":
		return OrdStatusDoneForDay
	case "Expired":
		return OrdStatusExpired
	}
	return OrdStatusUnknown
}

func OrdStatusToString(thetype OrdStatus) string {
	switch thetype {
	case OrdStatusNew:
		return "New"
	case OrdStatusPartiallyFilled:
		return "PartiallyFilled"
	case OrdStatusFilled:
		return "Filled"
	case OrdStatusCanceled:
		return "Cancelled"
	case OrdStatusRejected:
		return "Rejected"
	case OrdStatusReplaced:
		return "Replaced"
	case OrdStatusDoneForDay:
		return "DoneForDay"
	case OrdStatusExpired:
		return "Expired"
	}
	return "OrdStatusUnknown"
}
