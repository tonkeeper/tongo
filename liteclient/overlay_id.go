package liteclient

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/tonkeeper/tongo/tl"
	"github.com/tonkeeper/tongo/tlb"
)

type OverlayID struct {
	Workchain         int32
	Shard             int64
	ZeroStateFileHash tlb.Bits256
}

// FullID computes an ID of the overlay that is used to represent overlay in the TON overlay network.
func (o OverlayID) FullID() ([]byte, error) {
	overlayID := TonNodeShardPublicOverlayIdC{
		Workchain:         o.Workchain,
		Shard:             uint64(o.Shard),
		ZeroStateFileHash: tl.Int256(o.ZeroStateFileHash),
	}
	m, err := overlayID.MarshalTL()
	if err != nil {
		return nil, err
	}
	var id [4]byte
	binary.LittleEndian.PutUint32(id[:], uint32(1302254377))
	hasher := sha256.New()
	hasher.Write(id[:])
	hasher.Write(m)
	return hasher.Sum(nil), nil
}

func (o OverlayID) ComputeShortID() (tl.Int256, error) {
	m, err := o.FullID()
	if err != nil {
		return tl.Int256{}, err
	}
	var id [4]byte
	binary.LittleEndian.PutUint32(id[:], uint32(884622795)) // pubkey_overlay
	hasher := sha256.New()
	hasher.Write(id[:])
	hasher.Write([]byte{32})
	hasher.Write(m)
	hasher.Write([]byte{0, 0, 0})
	hash := hasher.Sum(nil)
	var h tl.Int256
	copy(h[:], hash[:])
	return h, nil
}
