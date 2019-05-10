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
)
