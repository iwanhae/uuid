package uuid

import (
	"encoding/json"
	"fmt"

	"github.com/iwanhae/uuid/base58"
)

type UUID [16]byte

type Version int

const (
	Version4 Version = Version(4)
	Version7 Version = Version(7)
)

var (
	Nil UUID
)

func Parse(uuid string) (UUID, error) {
	decoded := base58.Decode(uuid)
	if len(decoded) != 16 {
		return Nil, fmt.Errorf("expected base58 encoded uuid string to have 16 bytes, but got %d bytes", len(decoded))
	}
	return UUID(decoded), nil
}

func (u UUID) String() string {
	return base58.Encode(u[:])
}

func (u UUID) Version() Version {
	// Version is stored in bits 4-7 of the 6th byte
	switch u[6] >> 4 {
	case 4:
		return Version4
	case 7:
		return Version7
	default:
		return Version(0)
	}
}

// MarshalJSON implements json.Marshaler.
func (v UUID) MarshalJSON() ([]byte, error) {
	return []byte(
		fmt.Sprintf(`"%s"`, v.String()),
	), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (u *UUID) UnmarshalJSON(data []byte) error {
	var encoded string
	if err := json.Unmarshal(data, &encoded); err != nil {
		return err
	}
	parsed, err := Parse(encoded)
	if err != nil {
		return fmt.Errorf("failed to parse uuid %q, %w", encoded, err)
	}
	copy(u[:], parsed[:])
	return nil
}
