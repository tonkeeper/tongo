package liteclient

import (
	"bytes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/utils"
	"io"
	"strings"
)

type params [160]byte

func newParameters() (params, error) {
	var p params
	_, err := io.ReadFull(rand.Reader, p[:])
	return p, err
}

func (p params) rxKey() []byte {
	return p[0:32]
}

func (p params) txKey() []byte {
	return p[32:64]
}

func (p params) rxNonce() []byte {
	return p[64:80]
}

func (p params) txNonce() []byte {
	return p[80:96]
}

func (p params) padding() []byte {
	return p[96:160]
}

func (p params) hash() []byte {
	h := sha256.New()
	h.Write(p[:])
	return h.Sum(nil)
}

type Packet struct {
	Payload []byte
	nonce   [32]byte
}

func NewPacket(payload []byte) (Packet, error) {
	packet := Packet{Payload: payload}
	_, err := io.ReadFull(rand.Reader, packet.nonce[:])
	return packet, err
}

func (p Packet) hash() []byte {
	h := sha256.New()
	h.Write(p.nonce[:])
	h.Write(p.Payload)
	return h.Sum(nil)
}

func (p Packet) size() []byte {
	s := make([]byte, 4)
	binary.LittleEndian.PutUint32(s[:], uint32(len(p.Payload)+32+32))
	return s
}
func (p Packet) marshal() []byte {
	b := make([]byte, 4+32+len(p.Payload)+32)
	copy(b[:4], p.size())
	copy(b[4:36], p.nonce[:])
	copy(b[36:36+len(p.Payload)], p.Payload)
	copy(b[36+len(p.Payload):], p.hash())
	return b
}

func (p Packet) MagicType() uint32 {
	if len(p.Payload) < 4 {
		return 0
	}
	return binary.LittleEndian.Uint32(p.Payload[:4])
}

func ParsePacket(r io.Reader, decryptor cipher.Stream) (Packet, error) {
	var p Packet
	size := make([]byte, 4) //todo: reuse via sync.pool
	n, err := io.ReadFull(r, size)
	if err != nil {
		return Packet{}, err
	}
	if n < 4 {
		return p, fmt.Errorf("not enough bytes (%v) for parsing packet", n)
	}
	decryptor.XORKeyStream(size, size)
	length := int(binary.LittleEndian.Uint32(size))
	if length < 64 || length > 8<<20 {
		return Packet{}, fmt.Errorf("invalid length of data: %d bytes", length)
	}
	data := make([]byte, length)
	n, err = io.ReadFull(r, data)
	if err != nil {
		return Packet{}, err
	}
	if n != length {
		return p, fmt.Errorf("invalid packet length. should be %v by header but real length is %v", length, n)
	}
	decryptor.XORKeyStream(data, data)
	copy(p.nonce[:], data[:32])
	p.Payload = make([]byte, length-32-32) //todo: maybe remove copy
	copy(p.Payload, data[32:length-32])
	if !bytes.Equal(data[length-32:], p.hash()) {
		return p, fmt.Errorf("checksum error")
	}
	return p, nil
}

type Address struct {
	pubkey ed25519.PublicKey
}

func NewAddress(key []byte) (Address, error) {
	if len(key) != 32 {
		return Address{}, fmt.Errorf("invalid key length: %v", len(key))
	}
	var a Address
	a.pubkey = key
	return a, nil
}

func (a Address) hash() []byte {
	h := sha256.New()
	h.Write([]byte{0xc6, 0xb4, 0x13, 0x48})
	h.Write(a.pubkey[:])
	return h.Sum(nil)
}

// ADNL address is equal Address.hash()

func ADNLAddressToBase32(addr ton.Bits256) string {
	a := append([]byte{0x2d}, addr[:]...)
	crcBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(crcBytes, utils.Crc16(a))
	return strings.ToLower(base32.StdEncoding.EncodeToString(append(a, crcBytes...))[1:])
}

func ParseADNLAddress(addr string) (ton.Bits256, error) {
	addr = strings.TrimSuffix(addr, ".adnl")
	if len(addr) != 55 {
		return [32]byte{}, errors.New("wrong adnl address length")
	}
	buf, err := base32.StdEncoding.DecodeString("F" + strings.ToUpper(addr))
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to decode address: %w", err)
	}
	if buf[0] != 0x2d {
		return [32]byte{}, errors.New("invalid first byte")
	}
	hash := binary.BigEndian.Uint16(buf[33:])
	calculatedHash := utils.Crc16(buf[:33])
	if hash != calculatedHash {
		return [32]byte{}, errors.New("invalid address")
	}
	var res ton.Bits256
	copy(res[:], buf[1:33])
	return res, nil
}
