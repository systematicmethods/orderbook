package orderbook

type ExecType rune

const (
	ExecTypeNew                            ExecType = '0'
	DoneForDay                                      = '3'
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
)
