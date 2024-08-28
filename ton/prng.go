package ton

import (
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"math/big"
)

type ValidatorPRNGDescr struct {
	seed          [32]byte // seed for validator set computation, set to zero if none
	shard         uint64
	workchain     int32
	catchainSeqno uint32
	hash          []byte
}

// ValidatorPRNG is a pseudorandom number generator to randomize validators order
type ValidatorPRNG struct {
	descr ValidatorPRNGDescr
	pos   int
	limit int
}

func NewValidatorPRNG(seed []byte, shard uint64, workchain int32, catchainSeqno uint32) (*ValidatorPRNG, error) {
	if len(seed) != 0 && len(seed) != 32 {
		return nil, fmt.Errorf("invalid seed")
	}
	descr := ValidatorPRNGDescr{
		shard:         shard,
		workchain:     workchain,
		catchainSeqno: catchainSeqno,
	}
	copy(descr.seed[:], seed)
	return &ValidatorPRNG{
		descr: descr,
	}, nil
}

func (v *ValidatorPRNG) NextUInt64() uint64 {
	if v.pos < v.limit {
		temp := v.pos
		v.pos++
		return binary.BigEndian.Uint64(v.descr.hash[temp*8:])
	}
	v.rebuildHash()
	v.increaseSeed()
	v.pos = 1
	v.limit = 8
	return binary.BigEndian.Uint64(v.descr.hash)
}

func (v *ValidatorPRNG) NextRanged(rng uint64) uint64 {
	y := new(big.Int).SetUint64(v.NextUInt64())
	bigRange := new(big.Int).SetUint64(rng)
	// return (y * rng) >> 64.
	// Use big int to avoid uint64 overflow
	return new(big.Int).Rsh(new(big.Int).Mul(y, bigRange), 64).Uint64()
}

func (v *ValidatorPRNG) increaseSeed() {
	for i := 31; i >= 0; i-- {
		v.descr.seed[i]++
		if v.descr.seed[i] != 0 {
			break
		}
	}
}

func (v *ValidatorPRNG) rebuildHash() {
	h := sha512.New()
	h.Write(v.descr.seed[:])
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf, v.descr.shard)
	binary.BigEndian.PutUint32(buf[8:], uint32(v.descr.workchain))
	binary.BigEndian.PutUint32(buf[12:], v.descr.catchainSeqno)
	h.Write(buf)
	v.descr.hash = h.Sum(nil)
}
