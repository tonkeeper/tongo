package tongo

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/snksoft/crc"
)

type AccountStatus string

const (
	AccountNone   AccountStatus = "none"
	AccountEmpty  AccountStatus = "empty" // empty state from node
	AccountUninit AccountStatus = "uninit"
	AccountActive AccountStatus = "active"
	AccountFrozen AccountStatus = "frozen"
)

type AccountID struct {
	Workchain int32
	Address   []byte
}

func NewAccountId(id int32, addr []byte) *AccountID {
	return &AccountID{Workchain: id, Address: addr}
}

func (id AccountID) MarshalTL() ([]byte, error) {
	if len(id.Address) != 32 {
		return nil, fmt.Errorf("invalid account address length %v", len(id.Address))
	}
	payload := make([]byte, 36)
	binary.LittleEndian.PutUint32(payload[:4], uint32(id.Workchain))
	copy(payload[4:36], id.Address)
	return payload, nil
}

func (id *AccountID) UnmarshalTL(data []byte) error {
	if len(data) != 36 {
		return fmt.Errorf("invalid data length")
	}
	id.Workchain = int32(binary.LittleEndian.Uint32(data[:4]))
	id.Address = make([]byte, 32)
	copy(id.Address, data[4:36])
	return nil
}

func FromBase64Url(s string) (*AccountID, error) {
	if len(s) == 0 {
		return nil, nil
	}
	var aa AccountID
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(b) != 36 {
		return nil, fmt.Errorf("invalid account 'user friendly' form length: %v", s)
	}
	checksum := uint64(binary.BigEndian.Uint16(b[34:36]))
	if checksum != crc.CalculateCRC(crc.XMODEM, b[0:34]) {
		return nil, fmt.Errorf("invalid checksum")
	}
	aa.Workchain = int32(int8(b[1]))
	aa.Address = b[2:34]
	return &aa, nil
}

func FromRaw(s string) (*AccountID, error) {
	if len(s) == 0 {
		return nil, nil
	}
	var (
		workchain int32
		address   []byte
		aa        AccountID
	)
	_, err := fmt.Sscanf(s, "%d:%x", &workchain, &address)
	if err != nil {
		return nil, err
	}
	aa.Workchain = workchain
	aa.Address = address
	return &aa, nil
}

func ParseAccountId(s string) (*AccountID, error) {
	aa, err := FromRaw(s)
	if err != nil {
		aa, err = FromBase64Url(s)
		if err != nil {
			return nil, err
		}
	}
	return aa, nil
}

func MustParseAccountId(s string) *AccountID {
	aa, err := ParseAccountId(s)
	if err != nil {
		panic(err)
	}
	return aa
}
