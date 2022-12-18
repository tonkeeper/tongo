package code

import (
	"encoding/hex"
	"testing"

	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
	"github.com/startfellows/tongo/utils"
)

func TestCodeInspect(t *testing.T) {
	code, err := hex.DecodeString("b5ee9c724102140100021f000114ff00f4a413f4bcf2c80b0102016202030202cd04050201200e0f04e7d10638048adf000e8698180b8d848adf07d201800e98fe99ff6a2687d20699fea6a6a184108349e9ca829405d47141baf8280e8410854658056b84008646582a802e78b127d010a65b509e58fe59f80e78b64c0207d80701b28b9e382f970c892e000f18112e001718112e001f181181981e0024060708090201200a0b00603502d33f5313bbf2e1925313ba01fa00d43028103459f0068e1201a44343c85005cf1613cb3fccccccc9ed54925f05e200a6357003d4308e378040f4966fa5208e2906a4208100fabe93f2c18fde81019321a05325bbf2f402fa00d43022544b30f00623ba9302a402de04926c21e2b3e6303250444313c85005cf1613cb3fccccccc9ed54002c323401fa40304144c85005cf1613cb3fccccccc9ed54003c8e15d4d43010344130c85005cf1613cb3fccccccc9ed54e05f04840ff2f00201200c0d003d45af0047021f005778018c8cb0558cf165004fa0213cb6b12ccccc971fb008002d007232cffe0a33c5b25c083232c044fd003d0032c03260001b3e401d3232c084b281f2fff2742002012010110025bc82df6a2687d20699fea6a6a182de86a182c40043b8b5d31ed44d0fa40d33fd4d4d43010245f04d0d431d430d071c8cb0701cf16ccc980201201213002fb5dafda89a1f481a67fa9a9a860d883a1a61fa61ff480610002db4f47da89a1f481a67fa9a9a86028be09e008e003e00b01a500c6e")
	if err != nil {
		t.Fatal(err)
	}
	cell, err := boc.DeserializeBoc(code)
	if err != nil {
		t.Fatal(err)
	}
	c, err := cell[0].NextRef()
	if err != nil {
		t.Fatal(err)
	}
	type GetMethods struct {
		Hashmap tlb.Hashmap[tlb.Size19, boc.Cell]
	}
	var d GetMethods
	err = tlb.Unmarshal(c, &d)
	if err != nil {
		t.Fatal(err)
	}
	for i := range d.Hashmap.Keys() {
		num, err := d.Hashmap.Keys()[i].ReadInt(19)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(num)
	}
	t.Log(uint64(utils.Crc16String("get_nft_data")&0xffff) | 0x10000)
	t.Log(uint64(utils.Crc16String("get_collection_data")&0xffff) | 0x10000)

}
