package uuid

import (
	"fmt"
	"sync"
	"time"
)

type V7 struct {
	UUID `json:",inline"`
}

func NewV7() V7 {
	v4 := NewV4().UUID
	makeV7(v4[:])
	return V7{UUID: v4}
}

func (v *V7) Parse(uuid string) error {
	u, err := Parse(uuid)
	if err != nil {
		return err
	}
	if version := u.Version(); version != Version7 {
		return fmt.Errorf("expect uuid v7, but get uuid v%d", version)
	}
	v.UUID = u
	return nil
}

func (v *V7) Timestamp() time.Time {
	// Extract milliseconds from the first 48 bits (6 bytes)
	msec := int64(v.UUID[0])<<40 |
		int64(v.UUID[1])<<32 |
		int64(v.UUID[2])<<24 |
		int64(v.UUID[3])<<16 |
		int64(v.UUID[4])<<8 |
		int64(v.UUID[5])

	// Convert milliseconds to time.Time
	return time.UnixMilli(msec)
}

// UnmarshalJSON implements json.Unmarshaler
func (u *V7) UnmarshalJSON(data []byte) error {
	if err := u.UUID.UnmarshalJSON(data); err != nil {
		return err
	}
	if u.UUID.Version() != Version7 {
		return fmt.Errorf("expect uuid v7, but get uuid v%d", u.UUID.Version())
	}
	return nil
}

// Copied from
// https://github.com/google/uuid/blob/d55c313874fe007c6aaecc68211b6c7c7fc84aad/version7.go#L45-L75
func makeV7(uuid []byte) {
	/*
		 0                   1                   2                   3
		 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                           unix_ts_ms                          |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|          unix_ts_ms           |  ver  |  rand_a (12 bit seq)  |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|var|                        rand_b                             |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                            rand_b                             |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	*/
	_ = uuid[15] // bounds check

	t, s := getV7Time()

	uuid[0] = byte(t >> 40)
	uuid[1] = byte(t >> 32)
	uuid[2] = byte(t >> 24)
	uuid[3] = byte(t >> 16)
	uuid[4] = byte(t >> 8)
	uuid[5] = byte(t)

	uuid[6] = 0x70 | (0x0F & byte(s>>8))
	uuid[7] = byte(s)
}

var lastV7time int64

const nanoPerMilli = 1000000

var timeMu sync.Mutex

func getV7Time() (milli, seq int64) {
	timeMu.Lock()
	defer timeMu.Unlock()

	nano := time.Now().UnixNano()
	milli = nano / nanoPerMilli
	// Sequence number is between 0 and 3906 (nanoPerMilli>>8)
	seq = (nano - milli*nanoPerMilli) >> 8
	now := milli<<12 + seq
	if now <= lastV7time {
		now = lastV7time + 1
		milli = now >> 12
		seq = now & 0xfff
	}
	lastV7time = now
	return milli, seq
}
