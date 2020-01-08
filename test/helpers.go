package test

import (
	"github.com/andres-erbsen/clock"
	"log"
	"testing"
	"time"
)

func PrintError(err error, m *testing.T) {
	if err != nil {
		m.Errorf("Failed %v", err)
	}
}

func Loglines() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func NewMockClock(hour, min, sec int) *clock.Mock {
	aclock := clock.NewMock()
	aclock.Set(NewTime(hour, min, sec))
	return aclock
}

func NewMockClockFromTime(dt time.Time) *clock.Mock {
	aclock := clock.NewMock()
	aclock.Set(dt)
	return aclock
}

func NewTime(hour, min, sec int) time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Date(2019, 10, 11, hour, min, sec, 0, loc)
}

func NewLongTime(hour, min, sec int) time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Date(2019, 10, 21, hour, min, sec, 0, loc)
}

func NewDateTime(year, month, day, hour, min, sec int) time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Date(year, time.Month(month), day, hour, min, sec, 0, loc)
}
