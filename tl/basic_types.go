package tl

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

type Int256 [32]byte

// Scan implements Scanner for database/sql.
func (i *Int256) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Hash", src)
	}
	if len(srcB) != 32 {
		return fmt.Errorf("can't scan []byte of len %d into Hash, want %d", len(srcB), 32)
	}
	copy(i[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (i Int256) Value() (driver.Value, error) {
	return i[:], nil
}

func (i Int256) MarshalTL() ([]byte, error) {
	return i[:], nil
}

func (i *Int256) UnmarshalTL(r io.Reader) error {
	var b [32]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	*i = b
	return nil
}

func (i Int256) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(i[:]))
}

func (i *Int256) UnmarshalJSON(data []byte) error {
	var h string
	err := json.Unmarshal(data, &h)
	if err != nil {
		return err
	}
	b, err := hex.DecodeString(h)
	if err != nil {
		return err
	}
	if len(b) != len(i) {
		return fmt.Errorf("invalid int256 len")
	}
	copy(i[:], b)
	return nil
}
