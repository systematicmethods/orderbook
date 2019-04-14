package uuidext

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"testing"
	"time"
)

func Test_NewUUIDFromTimeIsTheSame(m *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	exs := []struct {
		t   time.Time
		seq int
	}{
		{time.Date(2019, 10, 11, 11, 11, 1, 0, loc), 0},
		{time.Date(2019, 11, 11, 11, 11, 59, 0, loc), 0},
		{time.Date(2019, 12, 11, 11, 11, 59, 123456700, loc), 0},
		{time.Date(2019, 12, 11, 11, 11, 59, 12345600, loc), 0},
		{time.Date(2019, 12, 11, 11, 11, 59, 12300, loc), 0},
		{time.Date(2019, 12, 11, 11, 11, 59, 12300, loc), 11},
	}

	for _, ex := range exs {
		//fmt.Printf("date %v\n", ex.t)
		id, _ := NewUUIDFromTimeSeq(ex.t, uint16(ex.seq))
		if id.Version() != 1 {
			m.Errorf("Not type 1 UUID was %d", id.Version())
		}
		t2sec, t2nano := id.Time().UnixTime()

		if t2sec != ex.t.Unix() {
			m.Errorf("Seconds different expected %d is %d", ex.t.Unix(), t2sec)
		}

		if t2nano != int64(ex.t.Nanosecond()) {
			m.Errorf("Nano different expected %d is %d", ex.t.Nanosecond(), t2nano)
		}

		if id.ClockSequence() != ex.seq {
			m.Errorf("Sequence different expected %d is %d", ex.seq, id.ClockSequence())
		}
	}
}

func Test_NewUUIDFromTimeIsEqual(m *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	exs := []struct {
		t   time.Time
		seq uint16
	}{
		{time.Date(2019, 10, 11, 11, 11, 1, 100, loc), 11},
		{time.Date(2019, 10, 11, 11, 11, 1, 100, loc), 11},
	}

	id1, _ := NewUUIDFromTimeSeq(exs[0].t, exs[0].seq)
	id2, _ := NewUUIDFromTimeSeq(exs[1].t, exs[1].seq)
	if id1.Time() != id2.Time() {
		m.Errorf("Time should be equal %v is %v", id1, id2)
	}
	if id1.ClockSequence() != id2.ClockSequence() {
		m.Errorf("Clock seq should be equal %v is %v", id1, id2)
	}
	//fmt.Printf("ex1 %v ex2 %v \n", exs[0].seq, exs[1].seq)
	//fmt.Printf("ex1 %v ex2 %v \n", id1.ClockSequence(), id2.ClockSequence())
}

func Test_MadeUUIDTimeIsOrdered(m *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	exs := []struct {
		t   time.Time
		seq uint16
	}{
		{time.Date(2019, 10, 11, 11, 11, 1, 100, loc), 10},
		{time.Date(2019, 10, 11, 11, 11, 1, 100, loc), 11},
		{time.Date(2019, 10, 11, 11, 11, 1, 100, loc), 9},
	}

	id1, _ := NewUUIDFromTimeSeq(exs[0].t, exs[0].seq)
	id2, _ := NewUUIDFromTimeSeq(exs[1].t, exs[1].seq)
	id3, _ := NewUUIDFromTimeSeq(exs[1].t, exs[1].seq)
	id4, _ := NewUUIDFromTimeSeq(exs[2].t, exs[2].seq)
	if UUIDComparator(id1, id2) >= 0 {
		m.Errorf("Time id1 should be less than id2 was %d %v is %v", UUIDComparator(id1, id2), id1.Time(), id2.Time())
	}
	if id1.ClockSequence() >= id2.ClockSequence() {
		m.Errorf("Clock seq id1 should be less than id2 was %v is %v", id1.ClockSequence(), id2.ClockSequence())
	}
	if UUIDComparator(id1, id3) >= 0 {
		m.Errorf("Time id1 should  be less than id3 was %d %v is %v", UUIDComparator(id1, id3), id1.Time(), id3.Time())
	}
	if UUIDComparator(id2, id3) != 0 {
		m.Errorf("Time id2 should equal to id3 was %d %v is %v", UUIDComparator(id2, id3), id2.Time(), id3.Time())
	}
	if UUIDComparator(id1, id4) <= 0 {
		m.Errorf("Time id1 should be greater than id4 was %d %v is %v", UUIDComparator(id1, id4), id1.Time(), id4.Time())
	}

}

func Test_NewUUIDSIsOrdered(m *testing.T) {
	id1, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()
	id3, _ := uuid.NewUUID()
	if UUIDComparator(id1, id2) >= 0 {
		m.Errorf("Time id1 should be less than id2 was %d %v is %v", UUIDComparator(id1, id2), id1.Time(), id2.Time())
	}
	if UUIDComparator(id2, id3) >= 0 {
		m.Errorf("Time id2 should be less than id3 was %d %v is %v", UUIDComparator(id2, id3), id2.Time(), id3.Time())
	}
}

func Test_NewUUIDIsEqualToNewUUIDFromTime(m *testing.T) {
	id1, _ := uuid.NewUUID()
	sec, nano := id1.Time().UnixTime()
	t1 := time.Unix(sec, nano)

	//id2, _ := NewUUIDFromTimeSeq(t1, uint16(id1.ClockSequence()))
	id2, _ := NewUUIDFromTimeSeqNode(t1, uint16(id1.ClockSequence()), id1.NodeID())
	if UUIDComparator(id1, id2) != 0 {
		m.Errorf("Time id1 should be equal id2 was %d", UUIDComparator(id1, id2))
		dumptime(m, id1, "id2 eq")
		dumptime(m, id2, "id2 eq")
	}
}

func Test_MixedNewUUIDFromTimeAndNewUIDAreOrdered(m *testing.T) {
	id1, _ := NewUUIDFromTimeSeq(time.Now(), 1)
	id2, _ := NewUUIDFromTimeSeq(time.Now(), 2)
	id3, _ := NewUUIDFromTimeSeq(time.Now(), 3)
	if UUIDComparator(id1, id2) >= 0 {
		m.Errorf("Time id1 should be less than id2 was %d %v is %v", UUIDComparator(id1, id2), id1.Time(), id2.Time())
		dumptime(m, id1, "id1 mix time")
		dumptime(m, id2, "id2 mix time")
	}
	if UUIDComparator(id2, id3) >= 0 {
		m.Errorf("Time id2 should be less than id3 was %d %v is %v", UUIDComparator(id2, id3), id2.Time(), id3.Time())
		dumptime(m, id2, "id2a mix")
		dumptime(m, id3, "id3 mix")
	}
}

func dumptime(m *testing.T, id uuid.UUID, msg string) {
	m.Errorf("Time %s %d, %d, %v %s", msg, id.Time(), id.ClockSequence(), id.Version(), hex.Dump(id[:]))
	dumpbytes(id[:])
}

func dumpbytes(b []byte) {
	for _, n := range b[:] {
		fmt.Printf(" %08b", n) // prints 00000000 11111101
	}
	fmt.Println()
}
