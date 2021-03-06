package uuidext

import (
	"bytes"
	"encoding/binary"
	"github.com/google/uuid"
	"time"
)

const (
	lillian    = 2299160          // Julian day of 15 Oct 1582
	unix       = 2440587          // Julian day of 1 Jan 1970
	epoch      = unix - lillian   // Days between epochs
	g1582      = epoch * 86400    // seconds between epochs
	g1582ns100 = g1582 * 10000000 // 100s of a nanoseconds between epochs
)

// NewUUID returns a Version 1 UUID based on given time and seq
// There are use cases for generating timeUUIDs for example back filling data
// The resolution is to 100th of micro second or 7 places
// See uuid.NewUUID and and uuid.GetTime
func NewUUIDFromTimeSeq(atime time.Time, seq uint16) (uuid.UUID, error) {
	now := uuid.Time(uint64(atime.UnixNano()/100) + g1582ns100)
	auuid := generate(now, seq)
	return auuid, nil
}

func NewUUIDFromTime(atime time.Time) (uuid.UUID, error) {
	nano := atime.UnixNano()
	now := uuid.Time(uint64(nano/100) + g1582ns100)
	auuid := generate(now, 0)
	return auuid, nil
}

func NewUUIDFromTimeSeqNode(atime time.Time, seq uint16, nodeid []byte) (uuid.UUID, error) {
	now := uuid.Time(uint64(atime.UnixNano()/100) + g1582ns100)
	auuid := generate(now, seq)
	// not sure why but this fixes tests
	auuid[8] |= 0x80
	copy(auuid[10:], nodeid[:])
	return auuid, nil
}

// The ordering is intentionally set up so that the UUIDs
// can simply be numerically compared as a set of bytes
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func UUIDComparator(a, b interface{}) int {
	as := a.(uuid.UUID)
	bs := b.(uuid.UUID)
	return bytes.Compare(as[:], bs[:])
}

func generate(atime uuid.Time, seq uint16) uuid.UUID {
	timeLow := uint32(atime & 0xffffffff)
	timeMid := uint16((atime >> 32) & 0xffff)
	timeHi := uint16((atime >> 48) & 0x0fff)
	timeHi |= 0x1000 // Version 1

	var auuid uuid.UUID
	binary.BigEndian.PutUint32(auuid[0:], timeLow)
	binary.BigEndian.PutUint16(auuid[4:], timeMid)
	binary.BigEndian.PutUint16(auuid[6:], timeHi)
	binary.BigEndian.PutUint16(auuid[8:], seq)
	return auuid
}
