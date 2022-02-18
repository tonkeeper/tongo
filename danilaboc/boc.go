package danilaboc

import (
	"encoding/binary"
	"hash/crc32"
)

var bocHeader = []byte{
	0xb5, 0xee, 0x9c, 0x72, // boc header
	0x41, // has_idx:(## 1):0 has_crc32c:(## 1):1 has_cache_bits:(## 1):0 flags:(## 2):00 size:(## 3):001 = 01000001
	0x01, // offset_byte_size
	0x02, // cell_count (size * 8)
	0x01, // root_count (size * 8)
	0x00, // absent:(size * 8) { roots + absent <= cells }
	0xa1, // total_cells_size
	0x00, // root index
}

func (c *Cell) bocReprWithoutRefs() []byte {
	d1 := byte(c.RefsSize())
	d2 := byte((c.BitSize()+7)/8 + c.BitSize()/8)

	res := make([]byte, ((c.BitSize()+7)/8)+2)
	res[0] = d1
	res[1] = d2
	copy(res[2:], c.contents)

	// add padding
	if c.BitSize()%8 != 0 {
		res[len(res)-1] |= 0x01 << (7 - c.BitSize()%8)
	}

	return res
}

func (c *Cell) BocRepr(startIndex int) []byte {
	res := c.bocReprWithoutRefs()

	for i := startIndex + 1; i < startIndex+c.RefsSize()+1; i++ {
		res = append(res, byte(i))
	}

	return res
}

func (c *Cell) ToBOC() []byte {
	cells := c.GetAllTree()

	cellsSerialized := make([]byte, 0)
	for i := 0; i < len(cells); i++ {
		cellsSerialized = append(cellsSerialized, cells[i].BocRepr(i)...)
	}
	totalSize := len(cellsSerialized)

	contents := make([]byte, len(bocHeader)+totalSize+4)
	copy(contents[:len(bocHeader)], bocHeader)
	contents[9] = byte(totalSize)
	contents[6] = byte(len(cells)) // cell count

	copy(contents[len(bocHeader):len(bocHeader)+totalSize], cellsSerialized)

	binary.LittleEndian.PutUint32(contents[len(bocHeader)+totalSize:], crc32.Checksum(contents[:len(bocHeader)+totalSize], crc32.MakeTable(crc32.Castagnoli)))

	return contents
}
