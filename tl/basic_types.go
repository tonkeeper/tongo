package tl

import (
	"database/sql/driver"
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
