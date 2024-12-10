package uuid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sync"
)

type V4 struct {
	UUID
}

func NewV4() V4 {
	return Must(newRandomFromPool())
}

func (v *V4) Parse(uuid string) error {
	u, err := Parse(uuid)
	if err != nil {
		return err
	}
	if version := u.Version(); version != Version4 {
		return fmt.Errorf("expect uuid v4, but get uuid v%d", version)
	}
	v.UUID = u
	return nil
}

// UnmarshalJSON implements json.Unmarshaler
func (u *V4) UnmarshalJSON(data []byte) error {
	if err := u.UUID.UnmarshalJSON(data); err != nil {
		return err
	}
	if u.UUID.Version() != Version4 {
		return fmt.Errorf("expect uuid v4, but get uuid v%d", u.UUID.Version())
	}
	return nil
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
	var uuid UUID
	poolMu.Lock()
	if poolPos == randPoolSize {
		_, err := io.ReadFull(rander, pool[:])
		if err != nil {
			poolMu.Unlock()
			return V4{Nil}, err
		}
		poolPos = 0
	}
	copy(uuid[:], pool[poolPos:(poolPos+16)])
	poolPos += 16
	poolMu.Unlock()

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return V4{uuid}, nil
}
