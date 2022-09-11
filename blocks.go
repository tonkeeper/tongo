package tongo

import (
	"encoding/binary"
	"fmt"
)

type TonNodeBlockIdExt struct {
	Workchain int32
	Shard     int64
	Seqno     int32
	RootHash  Hash
	FileHash  Hash
}

func (id TonNodeBlockIdExt) MarshalTL() ([]byte, error) {
	payload := make([]byte, 80)
	binary.LittleEndian.PutUint32(payload[:4], uint32(id.Workchain))
	binary.LittleEndian.PutUint64(payload[4:12], uint64(id.Shard))
	binary.LittleEndian.PutUint32(payload[12:16], uint32(id.Seqno))
	copy(payload[16:48], id.RootHash[:])
	copy(payload[48:80], id.FileHash[:])
	return payload, nil
}

func (id *TonNodeBlockIdExt) UnmarshalTL(data []byte) error {
	if len(data) != 80 {
		return fmt.Errorf("invalid data length")
	}
	id.Workchain = int32(binary.LittleEndian.Uint32(data[:4]))
	id.Shard = int64(binary.LittleEndian.Uint64(data[4:12]))
	id.Seqno = int32(binary.LittleEndian.Uint32(data[12:16]))
	copy(id.RootHash[:], data[16:48])
	copy(id.FileHash[:], data[48:80])
	return nil
}

func (id TonNodeBlockIdExt) String() string {
	return fmt.Sprintf("(%d,%x,%d,%x,%x)", id.Workchain, id.Shard, id.Seqno, id.RootHash, id.FileHash)
}
