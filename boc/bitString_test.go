package boc

import (
	"encoding/json"
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppend(t *testing.T) {
	s1 := NewBitString(8 * 10)
	s1.WriteUint(1, 80)
	s1.ReadBit()
	s2 := NewBitString(8 * 10)
	s2.WriteUint(1, 80)
	s2.ReadBit()
	s1.Append(s2)
}

func TestMinBits(t *testing.T) {
	for i := 0; i < 1000500; i++ {
		if minBitsRequired(uint64(i)) != int(math.Ceil(math.Log2(float64(i+1)))) {
			t.Fatal(i)
		}
	}
}

func BenchmarkMinbits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = minBitsRequired(uint64(i))
	}
}

func BenchmarkOldMinbits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = int(math.Ceil(math.Log2(float64(i + 1))))
	}
}

func TestBitString_WriteBit(t *testing.T) {
	bs1 := NewBitString(8)
	for i := 0; i <= 7; i++ {
		if err := bs1.WriteBit(true); err != nil {
			t.Errorf("WriteBit() failed: %v", err)
		}
	}

	tests := []struct {
		name      string
		bitstring BitString
		value     bool
		wantErr   string
	}{
		{
			name:      "overflow",
			bitstring: bs1,
			value:     true,
			wantErr:   "BitString overflow",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bitstring.WriteBit(tt.value)
			if len(tt.wantErr) > 0 {
				if err == nil {
					t.Errorf("WriteBit() must return an error")
				}
				if err.Error() != tt.wantErr {
					t.Errorf("WriteBit() error = %v, want = %v", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("WriteBit() error = %v", err)
			}
		})
	}
}

func TestNewCellWithBitsPanic(t *testing.T) {
	defer func() { recover() }()
	bs := NewBitString(2000)
	bs.len = 2000
	NewCellWithBits(bs)
	t.Errorf("should panic with cell overflow")
}

func TestNewCellWithBits(t *testing.T) {
	bs := NewBitString(0)
	if NewCellWithBits(bs).BitsAvailableForWrite() != 0 {
		t.Fatalf("should be Cell with zero bits")
	}
	bs = NewBitString(CellBits)
	for i := 0; i < CellBits; i++ {
		_ = bs.WriteBit(true)
	}
	if NewCellWithBits(bs).BitsAvailableForWrite() != 0 {
		t.Fatalf("should be full Cell")
	}
}

func Test_JSON(t *testing.T) {
	rand.Seed(24)

	b := make([]byte, CellBits)
	_, err := rand.Read(b)
	if err != nil {
		t.Errorf("rand.Read() failed: %v", err)
	}
	binaryStr := fmt.Sprintf("%08b", b)
	type testType struct {
		Value BitString
	}
	for length := 1; length < CellBits; length++ {
		bs := NewBitString(length)
		for i := 0; i < length; i++ {
			err := bs.WriteBit(binaryStr[i] == '1')
			if err != nil {
				t.Errorf("WriteBit() failed: %v", err)
			}
		}

		data, err := json.Marshal(testType{Value: bs})
		if err != nil {
			t.Errorf("json.Marshal() failed: %v", err)
		}
		var dest testType
		err = json.Unmarshal(data, &dest)
		if err != nil {
			t.Errorf("json.Unmarshal() failed: %v", err)
		}
		if bs.BinaryString() != dest.Value.BinaryString() {
			t.Errorf("\nwant: %v\n got: %v\ndata=%v", bs.BinaryString(), dest.Value.BinaryString(), string(data))
		}
	}

}

func TestBitString_ReadBits(t *testing.T) {
	for length := 1; length < 32; length++ {
		s := NewBitString(length)
		for i := 0; i < length; i++ {
			if err := s.WriteBit(true); err != nil {
				t.Fatalf("WriteBit() failed: %v", err)
			}
		}
		for startPos := 0; startPos <= length; startPos++ {
			for readbits := 1; readbits <= length; readbits++ {
				s.rCursor = startPos
				_, err := s.ReadBits(readbits)
				if startPos+readbits > length {
					if err != ErrNotEnoughBits {
						t.Fatalf("startPos: %v, readbits: %v, length: %v, err has to be ErrNotEnoughBits", startPos, readbits, length)
					}
					continue
				}
				if err != nil {
					t.Fatalf("startPos: %v, readbits: %v, length: %v, err has to nil", startPos, readbits, length)
				}
			}
		}
	}
}

func BenchmarkReadBit(b *testing.B) {
	str := NewBitString(1023)
	rand.Read(str.buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1023; j++ {
			_ = str.mustGetBit(j)
		}
	}
}

func BenchmarkReadUint(b *testing.B) {
	str := NewBitString(1023)
	rand.Read(str.buf)
	str.len = 1023
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str.ResetCounter()
		_, _ = str.ReadUint(32)
		_, _ = str.ReadUint(64)
		_, _ = str.ReadUint(10)
		_, _ = str.ReadUint(32)
		_, _ = str.ReadUint(1)
		_, _ = str.ReadUint(16)
		_, _ = str.ReadUint(64)
	}
}

func TestReadUint(t *testing.T) {
	for _, u := range []uint64{0, 1, 2, 8, 100500, (1 << 32) - 1, 1 << 32, (1 << 32) + 1, (1 << 64) - 1} {
		for offset := 0; offset <= 17; offset++ {
			t.Run(fmt.Sprintf("offset: %d, number: %d", offset, u), func(t *testing.T) {

				str := NewBitString(1023)
				for i := 0; i < offset; i++ {
					str.WriteBit(true)
				}
				str.WriteUint(u, bits.Len64(u))
				if bits.Len64(u) < 32 {
					str.WriteUint(u, bits.Len64(u)*2)
				}
				for i := 0; i < offset; i++ {
					str.ReadBit()
				}

				u2, err := str.ReadUint(bits.Len64(u))
				if assert.NoError(t, err) {
					assert.Equal(t, u, u2)
				}
			})
		}
	}
}
func TestReadUint_Prop(t *testing.T) {
	for _, pattern := range []uint64{
		0x0123456789ABCDEF,
		0xF1E2D3C4B5A69788,
		0x8040201008040201,
	} {
		for bitLen := range 65 {
			for offset := 0; offset <= 17; offset++ {
				num := pattern >> (64 - bitLen)
				outSpaceBit := false
				t.Run(fmt.Sprintf("number: %x, offset: %d, bitlen: %d", num, offset, bitLen), func(t *testing.T) {
					str := NewBitString(1023)
					for i := 0; i < offset; i++ {
						str.WriteBit(outSpaceBit)
					}
					str.WriteUint(num, bitLen)
					for i := 0; i < 64; i++ {
						str.WriteBit(outSpaceBit)
					}
					for i := 0; i < offset; i++ {
						str.ReadBit()
					}

					u2, err := str.ReadUint(bitLen)
					if assert.NoError(t, err) {
						assert.Equal(t, num, u2)
					}
				})
			}
		}
	}
}

func TestReadByte(t *testing.T) {
	for offset := range 7 {
		for _, byt := range []byte{
			0b01101011,
			0b11100001,
			0b10000111,
			0b00000001,
			0b10000000,
		} {
			t.Run(fmt.Sprintf("offset=%v, byte=%d", offset, byt), func(t *testing.T) {
				str := NewBitString(1023)
				// WRITE
				for i := 0; i < offset; i++ {
					assert.NoError(t, str.WriteBit(true))
				}
				assert.NoError(t, str.WriteByte(byt))
				// READ
				for i := 0; i < offset; i++ {
					bit, err := str.ReadBit()
					assert.NoError(t, err)
					assert.True(t, bit)
				}
				b, err := str.ReadByte()
				assert.NoError(t, err)
				assert.Equal(t, byt, b)
			})
		}
	}
}

func BenchmarkBitString_WriteByte(b *testing.B) {
	for offset := range 8 {
		b.Run(fmt.Sprintf("offset=%v OLD", offset), func(t *testing.B) {
			str := NewBitString(8 * b.N)
			for i := 0; i < b.N; i++ {
				for range 100000 {
					str.len = offset
					_ = str.writeByteOld(0xff)
				}
			}
		})
		b.Run(fmt.Sprintf("offset=%v NEW", offset), func(t *testing.B) {
			str := NewBitString(8 * b.N)
			for i := 0; i < b.N; i++ {
				for range 100000 {
					str.len = offset
					_ = str.WriteByte(0xff)
				}
			}
		})
	}
}

func BenchmarkBitString_WriteUint(b *testing.B) {
	for offset := range 8 {
		for _, size := range []int{1, 2, 4, 5, 8, 32, 64} {
			b.Run(fmt.Sprintf("offset=%v size=%v NEW", offset, size), func(t *testing.B) {
				str := NewBitString(8 * b.N)
				for i := 0; i < b.N; i++ {
					for range 100000 {
						str.len = offset
						_ = str.WriteUint(0xFFFFFFFFFFFFFFFF, size)
					}
				}
			})
			b.Run(fmt.Sprintf("offset=%v size=%v OLD", offset, size), func(t *testing.B) {
				str := NewBitString(8 * b.N)
				for i := 0; i < b.N; i++ {
					for range 100000 {
						str.len = offset
						_ = str.writeUintOld(0xFFFFFFFFFFFFFFFF, size)
					}
				}
			})
		}
	}
}
