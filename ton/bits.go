package ton

import (
	"bytes"
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/tonkeeper/tongo/tlb"
)

type Bits256 tlb.Bits256

func (h Bits256) Base64() string {
	return base64.StdEncoding.EncodeToString(h[:])
}

func (h Bits256) Hex() string {
	return fmt.Sprintf("%x", h[:])
}
func (h *Bits256) FromBase64(s string) error {
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

func (h *Bits256) FromBase64URL(s string) error {
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

func (h *Bits256) FromHex(s string) error {
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

func (h *Bits256) FromUnknownString(s string) error {
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

func (h *Bits256) FromBytes(b []byte) error {
	if len(b) != 32 {
		return fmt.Errorf("can't scan []byte of len %d into Bits256, want %d", len(b), 32)
	}
	copy(h[:], b)
	return nil
}

func ParseHash(s string) (Bits256, error) {
	var h Bits256
	err := h.FromUnknownString(s)
	return h, err
}

func MustParseHash(s string) Bits256 {
	h, err := ParseHash(s)
	if err != nil {
		panic(err)
	}
	return h
}

// Scan implements Scanner for database/sql.
func (h *Bits256) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Bits256", src)
	}
	if len(srcB) != 32 {
		return fmt.Errorf("can't scan []byte of len %d into Bits256, want %d", len(srcB), 32)
	}
	copy(h[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (h Bits256) Value() (driver.Value, error) {
	return h[:], nil
}

func (b Bits256) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%x\"", b), nil
}

func (b *Bits256) UnmarshalJSON(buf []byte) error {
	var sl []byte
	_, err := fmt.Fscanf(bytes.NewReader(buf), "\"%x\"", &sl)
	if len(sl) != 32 {
		return fmt.Errorf("can't parse 256bits %v", string(buf))
	}
	copy(b[:], sl)
	return err
}
