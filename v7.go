package uuid

import (
	"sync"
	"time"
)

type V7 [16]byte

func NewV7() V7 {
	v4 := NewV4()
	makeV7(v4[:])
	return V7(v4)
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
