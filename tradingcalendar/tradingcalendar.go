package tradingcalendar

import "time"

type TradingScheduleType string
type Date time.Time

const (
	TradingScheduleTypeOrderEntry TradingScheduleType = "OrderEntry"
	TradingScheduleTypeAuction    TradingScheduleType = "Auction"
	TradingScheduleTypeTrading    TradingScheduleType = "Trading"
)

type TradingSchedule interface {
	InTradingSchedule(time.Time) (bool, TradingScheduleType)
}

type tradingSchedule struct {
	startSchedule       time.Time
	endSchedule         time.Time
	tradingScheduleType TradingScheduleType
}

type DayTradingSchedule interface {
	InTradingSchedule(time.Time) TradingSchedule
}

type dayTradingSchedule struct {
	tradingSchedule []TradingSchedule
}

type tradingCalendar struct {
	calendar map[Date]DayTradingSchedule
}

type TradingCalendar interface {
	DayTradingScheduleFor(date Date) DayTradingSchedule
}
