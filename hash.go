package tongo

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tonkeeper/tongo/boc"
)

type Hash [32]byte

var ErrEntityNotFound = errors.New("entity not found")

func (h Hash) Base64() string {
	return base64.StdEncoding.EncodeToString(h[:])
}

func (h Hash) Hex() string {
	return fmt.Sprintf("%x", h[:])
}
func (h *Hash) FromBase64(s string) error {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	if len(b) != 32 {
		return errors.New("invalid hash length")
	}
	copy(h[:], b)
	return nil
}

func (h *Hash) FromBase64URL(s string) error {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	if len(b) != 32 {
		return errors.New("invalid hash length")
	}
	copy(h[:], b)
	return nil
}

func (h *Hash) FromHex(s string) error {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	if len(b) != 32 {
		return errors.New("invalid hash length")
	}
	copy(h[:], b)
	return nil
}

func (h *Hash) FromUnknownString(s string) error {
	err := h.FromBase64(s)
	if err == nil {
		return nil
	}
	err = h.FromBase64URL(s)
	if err == nil {
		return nil
	}
	err = h.FromHex(s)
	if err == nil {
		return nil
	}
	return err
}

func (h *Hash) FromBytes(b []byte) error {
	if len(b) != 32 {
		return fmt.Errorf("can't scan []byte of len %d into Hash, want %d", len(b), 32)
	}
	copy(h[:], b)
	return nil
}

func ParseHash(s string) (Hash, error) {
	var h Hash
	err := h.FromUnknownString(s)
	return h, err
}

func MustParseHash(s string) Hash {
	h, err := ParseHash(s)
	if err != nil {
		panic(err)
	}
	return h
}

// Scan implements Scanner for database/sql.
func (h *Hash) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Hash", src)
	}
	if len(srcB) != 32 {
		return fmt.Errorf("can't scan []byte of len %d into Hash, want %d", len(srcB), 32)
	}
	copy(h[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (h Hash) Value() (driver.Value, error) {
	return h[:], nil
}

func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%x", h))
}

func (h Hash) MarshalTLB(c *boc.Cell, tag string) error {
	err := c.WriteBytes(h[:])
	if err != nil {
		return err
	}
	return nil
}

func (h *Hash) UnmarshalTLB(c *boc.Cell, tag string) error {
	b, err := c.ReadBytes(32)
	if err != nil {
		return err
	}
	copy(h[:], b[:])
	return nil
}

func (h Hash) MarshalTL() ([]byte, error) {
	return h[:], nil
}

func (h *Hash) UnmarshalTL(r io.Reader) error {
	var b [32]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	*h = b
	return nil
}
