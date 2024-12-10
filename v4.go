package uuid

import (
	"crypto/rand"
	"errors"
	"io"
	"sync"
)

type V4 UUID

func NewV4() V4 {
	return Must(newRandomFromPool())
}

const randPoolSize = 16 * 16

var (
	rander = rand.Reader // random function

	poolMu  sync.Mutex
	poolPos = randPoolSize     // protected with poolMu
	pool    [randPoolSize]byte // protected with poolMu

	ErrInvalidUUIDFormat      = errors.New("invalid UUID format")
	ErrInvalidBracketedFormat = errors.New("invalid bracketed UUID format")
)

func newRandomFromPool() (V4, error) {
	var uuid V4
	poolMu.Lock()
	if poolPos == randPoolSize {
		_, err := io.ReadFull(rander, pool[:])
		if err != nil {
			poolMu.Unlock()
			return V4(Nil), err
		}
		poolPos = 0
	}
	copy(uuid[:], pool[poolPos:(poolPos+16)])
	poolPos += 16
	poolMu.Unlock()

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid, nil
}
