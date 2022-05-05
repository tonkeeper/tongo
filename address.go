package tongo

import (
	"encoding/binary"
	"fmt"
)

type AccountID struct {
	Workchain int32
	Address   []byte
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
