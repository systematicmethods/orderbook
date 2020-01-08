package fixmodel

type OrdRejReason int

const (
	OrdRejReasonNotApplicable                  OrdRejReason = -1
	OrdRejReasonUnknownSymbol                  OrdRejReason = 1
	OrdRejReasonExchangeClosed                 OrdRejReason = 2
	OrdRejReasonUnknownOrder                   OrdRejReason = 5
	OrdRejReasonDuplicateOrder                 OrdRejReason = 6
	OrdRejReasonInvalidInvestorID              OrdRejReason = 10
	OrdRejReasonUnsupportedOrderCharacteristic OrdRejReason = 11 // any other reason
	OrdRejReasonOther                          OrdRejReason = 99
)

/*
0 = Broker option

1 = Unknown symbol

2 = Exchange closed

3 = Order exceeds limit

4 = Too late to enter

5 = Unknown Order

6 = Duplicate Order (e.g. dupe ClOrdID <11>)

7 = Duplicate of a verbally communicated order

8 = Stale Order
	0	=
	 Broker / Exchange option
	Added  FIX.2.7		[BrokerCredit]
	1	=
	Unknown symbol
	Added  FIX.2.7		[UnknownSymbol]
	2	=
	Exchange closed
	Added  FIX.2.7		[ExchangeClosed]
	3	=
	Order exceeds limit
	Added  FIX.2.7		[OrderExceedsLimit]
	4	=
	Too late to enter
	Added  FIX.4.0		[TooLateToEnter]
	5	=
	Unknown order
	Added  FIX.4.1		[UnknownOrder]
	6	=
	Duplicate Order (e.g. dupe ClOrdID)
	Added  FIX.4.1		[DuplicateOrder]
	7	=
	Duplicate of a verbally communicated order
	Added  FIX.4.2		[DuplicateOfAVerballyCommunicatedOrder]
	8	=
	Stale order
	Added  FIX.4.2		[StaleOrder]
	9	=
	Trade along required
	Added  FIX.4.3		[TradeAlongRequired]
	10	=
	Invalid Investor ID
	Added  FIX.4.3		[InvalidInvestorID]
	11	=
	Unsupported order characteristic
	Added  FIX.4.3		[UnsupportedOrderCharacteristic]
	12	=
	Surveillance option
	Added  FIX.4.3
	Updated  FIX.5.0SP2  EP204		[SurveillanceOption]
	13	=
	Incorrect quantity
	Added  FIX.4.4		[IncorrectQuantity]
	14	=
	Incorrect allocated quantity
	Added  FIX.4.4		[IncorrectAllocatedQuantity]
	15	=
	Unknown account(s)
	Added  FIX.4.4		[UnknownAccount]
	16	=
	Price exceeds current price band
	Added  FIX.5.0  EP-1		[PriceExceedsCurrentPriceBand]
	18	=
	Invalid price increment
	Added  FIX.4.4  EP6		[InvalidPriceIncrement]
	19	=
	Reference price not available
	Added  FIX.5.0SP2  EP134		[ReferencePriceNotAvailable]
	20	=
	Notional value exceeds threshold
	Added  FIX.5.0SP2  EP134		[NotionalValueExceedsThreshold]
	21	=
	Algorithm risk threshold breached
	A sell-side broker algorithm has detected that a risk limit has been breached which requires further communication with the client. Used in conjunction with Text(58) to convey the details of the specific event.
	Added  FIX.5.0SP2  EP149		[AlgorithRiskThresholdBreached]
	22	=
	Short sell not permitted
	Added  FIX.5.0SP2  EP164		[ShortSellNotPermitted]
	23	=
	Short sell rejected due to security pre-borrow restriction
	Added  FIX.5.0SP2  EP164		[ShortSellSecurityPreBorrowRestriction]
	24	=
	Short sell rejected due to account pre-borrow restriction
	Added  FIX.5.0SP2  EP164		[ShortSellAccountPreBorrowRestriction]
	25	=
	Insufficient credit limit
	Added  FIX.5.0SP2  EP171		[InsufficientCreditLimit]
	26	=
	Exceeded clip size limit
	Added  FIX.5.0SP2  EP171		[ExceededClipSizeLimit]
	27	=
	Exceeded maximum notional order amount
	Added  FIX.5.0SP2  EP171		[ExceededMaxNotionalOrderAmt]
	28	=
	Exceeded DV01/PV01 limit
	Added  FIX.5.0SP2  EP171		[ExceededDV01PV01Limit]
	29	=
	Exceeded CS01 limit
	Added  FIX.5.0SP2  EP171		[ExceededCS01Limit]
	99	=
	Other
*/
