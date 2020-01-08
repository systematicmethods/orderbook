package orderbookex

type ExecRestatementReason int

const (
	ExecRestatementReasonCancelOnTradingHalt ExecRestatementReason = 6
	ExecRestatementReasonNone                ExecRestatementReason = -1
)

/*
0 = GT corporate action
1 = GT renewal / restatement (no corporate action)
10 = Warehouse Recap
2 = Verbal change
3 = Repricing of order
4 = Broker option
5 = Partial decline of OrderQty (e.g. exchange initiated partial cancel)
6 = Cancel on Trading Halt
7 = Cancel on System Failure
8 = Market (Exchange) option
9 = Canceled, not best
99 = Other
11 = Peg Refresh
*/
