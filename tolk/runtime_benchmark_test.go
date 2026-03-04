package tolk

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/tolk/parser"
)

func BenchmarkRuntimeUnmarshalling(b *testing.B) {
	type Case struct {
		name                 string
		ty                   parser.Ty
		cell                 boc.Cell
		abiFiles             []string
		customUnpackResolver func(decoder *Decoder) func(alias parser.AliasRef, cell *boc.Cell, value *AliasValue) error
	}
	for _, curr := range []Case{
		{
			name: "unmarshal small int",
			ty:   parser.NewIntNType(24),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c72410101010005000006ff76c41616db06"),
		},
		{
			name: "unmarshal big int",
			ty:   parser.NewIntNType(183),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101001900002dfffffffffffffffffffffffffffffffffff99bfeac6423a6f0b50c"),
		},
		{
			name: "unmarshal small uint",
			ty:   parser.NewUIntNType(53),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000900000d00000000001d34e435eafd"),
		},
		{
			name: "unmarshal big uint",
			ty:   parser.NewUIntNType(257),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101002300004100000000000000000000000000000000000000000000000000009fc4212a38ba40b11cce12"),
		},
		{
			name: "unmarshal var int 16",
			ty:   parser.NewVarInt16Type(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000600000730c98588449b6923"),
		},
		{
			name: "unmarshal var uint 32",
			ty:   parser.NewVarUInt32Type(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000800000b28119ab36b44d3a86c0f"),
		},
		{
			name: "unmarshal bits24",
			ty:   parser.NewBitsNType(24),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000500000631323318854035"),
		},
		{
			name: "unmarshal coins",
			ty:   parser.NewCoinsType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c72410101010007000009436ec6e0189ebbd7f4"),
		},
		{
			name: "unmarshal bool",
			ty:   parser.NewBoolType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000300000140f6d24034"),
		},
		{
			name: "unmarshal cell",
			ty:   parser.NewCellType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c724101020100090001000100080000007ba52a3292"),
		},
		{
			name: "unmarshal remaining",
			ty:   parser.NewRemainingType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000900000dc0800000000ab8d04726e4"),
		},
		{
			name: "unmarshal internal address",
			ty:   parser.NewAddressType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6"),
		},
		{
			name: "unmarshal not exists optional address",
			ty:   parser.NewAddressOptType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c724101010100030000012094418655"),
		},
		{
			name: "unmarshal exists optional address",
			ty:   parser.NewAddressOptType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6"),
		},
		{
			name: "unmarshal external address",
			ty:   parser.NewAddressExtType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000600000742082850fcbd94fd"),
		},
		{
			name: "unmarshal none address any",
			ty:   parser.NewAddressAnyType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c724101010100030000012094418655"),
		},
		{
			name: "unmarshal internal address any",
			ty:   parser.NewAddressAnyType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6"),
		},
		{
			name: "unmarshal external address any",
			ty:   parser.NewAddressAnyType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000600000742082850fcbd94fd"),
		},
		{
			name: "unmarshal var address any",
			ty:   parser.NewAddressAnyType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000900000dc0800000000ab8d04726e4"),
		},
		{
			name: "unmarshal not exists nullable",
			ty:   parser.NewNullableType(parser.NewRemainingType()),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000300000140f6d24034"),
		},
		{
			name: "unmarshal exists nullable",
			ty:   parser.NewNullableType(parser.NewCellType()),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010201000b000101c001000900000c0ae007880db9"),
		},
		{
			name: "unmarshal ref",
			ty:   parser.NewCellOfType(parser.NewIntNType(65)),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010201000e000100010011000000000009689e40e150b4c5"),
		},
		{
			name: "unmarshal empty tensor",
			ty:   parser.NewTensorType(),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c724101010100020000004cacb9cd"),
		},
		{
			name: "unmarshal not empty tensor",
			ty: parser.NewTensorType(
				parser.NewUIntNType(123),
				parser.NewBoolType(),
				parser.NewCoinsType(),
				parser.NewTensorType(
					parser.NewIntNType(23),
					parser.NewNullableType(parser.NewIntNType(2)),
				),
				parser.NewVarIntType(32),
			),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101001f00003900000000000000000000000000021cb43b9aca00fffd550bfbaae07401a2a98117"),
		},
		{
			name: "unmarshal small-int-key map",
			ty:   parser.NewMapType(parser.NewIntNType(32), parser.NewBoolType()),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010201000c000101c001000ba00000007bc09a662c32"),
		},
		{
			name: "unmarshal small-uint-key map",
			ty:   parser.NewMapType(parser.NewUIntNType(16), parser.NewAddressType()),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c72410104010053000101c0010202cb02030045a7400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe80045a3cff5555555555555555555555555555555555555555555555555555555555555555888440ce8"),
		},
		{
			name: "unmarshal big-uint-key map",
			ty:   parser.NewMapType(parser.NewUIntNType(78), parser.NewCellType()),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010301001a000101c0010115a70000000000000047550902000b000000001ab01d5bf1a9"),
		},
		{
			name: "unmarshal bits-key map",
			ty: parser.NewMapType(
				parser.NewBitsNType(16),
				parser.NewMapType(
					parser.NewIntNType(64),
					parser.NewTensorType(
						parser.NewAddressType(),
						parser.NewCoinsType(),
					),
				),
			),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010301003b000101c0010106a0828502005ea0000000000000003e400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe43b9aca00b89cdc86"),
		},
		{
			name: "unmarshal address-key map",
			ty:   parser.NewMapType(parser.NewAddressType(), parser.NewCoinsType()),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010201002f000101c0010051a17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f9409502f9002016fdc16e"),
		},
		{
			name: "unmarshal union with dec prefix",
			ty: parser.NewUnionType(
				1, true,
				parser.NewUnionVariant(parser.NewIntNType(16), "0"),
				parser.NewUnionVariant(parser.NewIntNType(128), "1"),
			),
			cell: *boc.MustDeserializeSinglRootHex("b5ee9c7241010101001300002180000000000000000000000003b5577dc0660d6029"),
		},
		{
			name: "unmarshal union with bin prefix",
			ty: parser.NewUnionType(
				3, false,
				parser.NewUnionVariant(parser.NewStructType("AddressWithPrefix"), "0b001"),
				parser.NewUnionVariant(parser.NewStructType("MapWithPrefix"), "0b011"),
				parser.NewUnionVariant(parser.NewStructType("CellWithPrefix"), "0b111"),
			),
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010201002e0001017801004fa17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f900a4d89920c413c650"),
			abiFiles: []string{"testdata/bin_union.json"},
		},
		{
			name: "unmarshal union with hex prefix",
			ty: parser.NewUnionType(
				32, false,
				parser.NewUnionVariant(parser.NewStructType("UInt66WithPrefix"), "0x12345678"),
				parser.NewUnionVariant(parser.NewStructType("UInt33WithPrefix"), "0xdeadbeef"),
				parser.NewUnionVariant(parser.NewStructType("UInt4WithPrefix"), "0x89abcdef"),
			),
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000b000011deadbeef00000000c0d75977b9"),
			abiFiles: []string{"testdata/hex_union.json"},
		},
		{
			name:     "unmarshal a-lot-refs from alias",
			ty:       parser.NewAliasType("GoodNamingForMsg"),
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8"),
			abiFiles: []string{"testdata/refs.json"},
		},
		{
			name:     "unmarshal a-lot-refs from struct",
			ty:       parser.NewStructType("ManyRefsMsg"),
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8"),
			abiFiles: []string{"testdata/refs.json"},
		},
		{
			name:     "unmarshal a-lot-generics from struct",
			ty:       parser.NewStructType("ManyRefsMsg", parser.NewUIntNType(16)),
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647"),
			abiFiles: []string{"testdata/generics.json"},
		},
		{
			name:     "unmarshal a-lot-generics from alias",
			ty:       parser.NewAliasType("GoodNamingForMsg", parser.NewUIntNType(16)),
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647"),
			abiFiles: []string{"testdata/generics.json"},
		},
		{
			name:     "unmarshal struct with default values",
			ty:       parser.NewStructType("DefaultTest"),
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101003100005d80000002414801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfd00000156ac2c4c70811a9dde"),
			abiFiles: []string{"testdata/default_values.json"},
		},
		{
			name:     "unmarshal a-lot-numbers",
			ty:       parser.NewStructType("Numbers"),
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c72410101010033000062000000000000000000000000000000000000000000000000000000000000000000000000000000f1106aecc4c800020926dc62f014"),
			abiFiles: []string{"testdata/numbers.json"},
		},
		{
			name:     "unmarshal a-lot-random-fields",
			ty:       parser.NewStructType("RandomFields"),
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010301007800028b79480107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6350e038d7eb37c5e80000000ab50ee6b28000000000000016e4c000006c175300001801bc01020001c00051000000000005120041efeaa9731b94da397e5e64622f5e63348b812ac5b4763a93f0dd201d0798d4409e337ceb"),
			abiFiles: []string{"testdata/random_fields.json"},
		},
		{
			name:     "unmarshal alias with custom unpack",
			ty:       parser.NewAliasType("MyAlias"),
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c724101010100470000890000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000043b9aca00886e91196"),
			abiFiles: []string{"testdata/custom_pack_unpack.json"},
			customUnpackResolver: func(decoder *Decoder) func(alias parser.AliasRef, cell *boc.Cell, value *AliasValue) error {
				return func(alias parser.AliasRef, cell *boc.Cell, value *AliasValue) error {
					err := cell.Skip(512)
					if err != nil {
						return fmt.Errorf("failed to 512 bits from alias")
					}
					val, err := decoder.Unmarshal(cell, parser.NewStructType("My"))
					if err != nil {
						return fmt.Errorf("failed to unmarshal alias' coins")
					}
					*value = AliasValue(*val)
					return nil
				}
			},
		},
	} {
		abis := make([]parser.ABI, len(curr.abiFiles))
		for i, abiFile := range curr.abiFiles {
			data, err := os.ReadFile(abiFile)
			if err != nil {
				b.Fatal(err)
			}

			var abi parser.ABI
			err = json.Unmarshal(data, &abi)
			if err != nil {
				b.Fatal(err)
			}
			abis[i] = abi
		}
		b.Run(curr.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				curr.cell.ResetCounters()
				decoder := NewDecoder()
				if err := decoder.WithABIs(abis...); err != nil {
					b.Fatal(err)
				}
				if curr.customUnpackResolver != nil {
					decoder.WithCustomUnpackResolver(curr.customUnpackResolver(decoder))
				}
				_, err := decoder.Unmarshal(&curr.cell, curr.ty)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkRuntimeMessageUnmarshalling(b *testing.B) {
	type Case struct {
		name     string
		cell     boc.Cell
		abiFiles []string
	}
	for _, curr := range []Case{
		{
			name:     "unmarshal jetton transfer message stonfi",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c720102030100011b0001ae0f8a7ea5546de4efb35a04c230f424080125c28235ca8d125e676591513d520721b1fe99f7722f4c87723ce7ee0dfb73a300268806c2c709c47ec1b610073c38ef75cf6066e9e4368b10ddbdc015d0e59c98881c9c38010101e16664de2a801244183034d9fd59a236f71ec4271be377399056dda4cc3a5ebf5dc40967df64100268806c2c709c47ec1b610073c38ef75cf6066e9e4368b10ddbdc015d0e59c98a004d100d858e1388fd836c200e7871deeb9ec0cdd3c86d1621bb7b802ba1cb39310000000034d3e7f3c002009542ecec75480134403616384e23f60db08039e1c77bae7b03374f21b45886edee00ae872ce4c4000005400f4684b10a661eaa395f87d4a660e6dfc3bec187a8b24f6f362c0c6b1d20f1b5d8"),
			abiFiles: []string{"testdata/jetton_transfer.json", "testdata/payloads.json"},
		},
	} {
		abis := make([]parser.ABI, len(curr.abiFiles))
		for i, abiFile := range curr.abiFiles {
			data, err := os.ReadFile(abiFile)
			if err != nil {
				b.Fatal(err)
			}

			var abi parser.ABI
			err = json.Unmarshal(data, &abi)
			if err != nil {
				b.Fatal(err)
			}
			abis[i] = abi
		}
		b.Run(curr.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				curr.cell.ResetCounters()
				decoder := NewDecoder()
				if err := decoder.WithABIs(abis...); err != nil {
					b.Fatal(err)
				}
				_, err := decoder.UnmarshalMessage(&curr.cell)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkTLBMessageUnmarshalling(b *testing.B) {
	cell := boc.MustDeserializeSinglRootHex("b5ee9c720102030100011b0001ae0f8a7ea5546de4efb35a04c230f424080125c28235ca8d125e676591513d520721b1fe99f7722f4c87723ce7ee0dfb73a300268806c2c709c47ec1b610073c38ef75cf6066e9e4368b10ddbdc015d0e59c98881c9c38010101e16664de2a801244183034d9fd59a236f71ec4271be377399056dda4cc3a5ebf5dc40967df64100268806c2c709c47ec1b610073c38ef75cf6066e9e4368b10ddbdc015d0e59c98a004d100d858e1388fd836c200e7871deeb9ec0cdd3c86d1621bb7b802ba1cb39310000000034d3e7f3c002009542ecec75480134403616384e23f60db08039e1c77bae7b03374f21b45886edee00ae872ce4c4000005400f4684b10a661eaa395f87d4a660e6dfc3bec187a8b24f6f362c0c6b1d20f1b5d8")
	for i := 0; i < b.N; i++ {
		cell.ResetCounters()
		var body abi.InMsgBody
		decoder := tlb.NewDecoder()
		if err := decoder.Unmarshal(cell, &body); err != nil {
			b.Fatal(err)
		}
	}
}
