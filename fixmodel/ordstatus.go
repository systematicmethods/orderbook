package fixmodel

type OrdStatus rune

const (
	OrdStatusNew                OrdStatus = '0'
	OrdStatusPartiallyFilled    OrdStatus = '1'
	OrdStatusFilled             OrdStatus = '2'
	OrdStatusDoneForDay         OrdStatus = '3'
	OrdStatusCanceled           OrdStatus = '4'
	OrdStatusReplaced           OrdStatus = '5'
	OrdStatusPendingCancel      OrdStatus = '6'
	OrdStatusStopped            OrdStatus = '7'
	OrdStatusRejected           OrdStatus = '8'
	OrdStatusSuspended          OrdStatus = '9'
	OrdStatusPendingNew         OrdStatus = 'A'
	OrdStatusCalculated         OrdStatus = 'B'
	OrdStatusExpired            OrdStatus = 'C'
	OrdStatusAcceptedForBidding OrdStatus = 'D'
	OrdStatusPendingReplace     OrdStatus = 'E'
	OrdStatusUnknown            OrdStatus = 'x'
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

func (it OrdStatus) String() string {
	switch it {
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

func OrdStatusToString(thetype OrdStatus) string {
	return thetype.String()
}
