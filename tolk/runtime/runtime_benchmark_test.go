package runtime

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
		typeName             string
		cell                 boc.Cell
		abiFile              string
		customUnpackResolver func(decoder *Decoder) func(alias parser.AliasRef, cell *boc.Cell, value *AliasValue) error
	}
	for _, curr := range []Case{
		{
			name:     "unmarshal small int",
			typeName: "BenchmarkSmallInt",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c72410101010005000006ff76c41616db06"),
		},
		{
			name:     "unmarshal big int",
			typeName: "BenchmarkBigInt",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101001900002dfffffffffffffffffffffffffffffffffff99bfeac6423a6f0b50c"),
		},
		{
			name:     "unmarshal small uint",
			typeName: "BenchmarkSmallUint",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000900000d00000000001d34e435eafd"),
		},
		{
			name:     "unmarshal big uint",
			typeName: "BenchmarkBigUint",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7201010101002200004000000000000000000000000000000000000000000000000000013f8842547174"),
		},
		{
			name:     "unmarshal var int 16",
			typeName: "BenchmarkVarInt16",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000600000730c98588449b6923"),
		},
		{
			name:     "unmarshal var uint 32",
			typeName: "BenchmarkVarUint32",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000800000b28119ab36b44d3a86c0f"),
		},
		{
			name:     "unmarshal bits24",
			typeName: "BenchmarkBits24",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000500000631323318854035"),
		},
		{
			name:     "unmarshal coins",
			typeName: "BenchmarkCoins",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c72410101010007000009436ec6e0189ebbd7f4"),
		},
		{
			name:     "unmarshal bool",
			typeName: "BenchmarkBool",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000300000140f6d24034"),
		},
		{
			name:     "unmarshal cell",
			typeName: "BenchmarkCell",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c724101020100090001000100080000007ba52a3292"),
		},
		{
			name:     "unmarshal remaining",
			typeName: "BenchmarkRemaining",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000900000dc0800000000ab8d04726e4"),
		},
		{
			name:     "unmarshal internal address",
			typeName: "BenchmarkAddress",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6"),
		},
		{
			name:     "unmarshal not exists optional address",
			typeName: "BenchmarkOptionalAddress",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c724101010100030000012094418655"),
		},
		{
			name:     "unmarshal exists optional address",
			typeName: "BenchmarkOptionalAddress",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6"),
		},
		{
			name:     "unmarshal external address",
			typeName: "BenchmarkExternalAddress",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000600000742082850fcbd94fd"),
		},
		{
			name:     "unmarshal none address any",
			typeName: "BenchmarkAnyAddress",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c724101010100030000012094418655"),
		},
		{
			name:     "unmarshal internal address any",
			typeName: "BenchmarkAnyAddress",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101002400004380107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6351064a3e1a6"),
		},
		{
			name:     "unmarshal external address any",
			typeName: "BenchmarkAnyAddress",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000600000742082850fcbd94fd"),
		},
		{
			name:     "unmarshal var address any",
			typeName: "BenchmarkAnyAddress",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000900000dc0800000000ab8d04726e4"),
		},
		{
			name:     "unmarshal not exists nullable",
			typeName: "BenchmarkNullableRemaining",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000300000140f6d24034"),
		},
		{
			name:     "unmarshal exists nullable",
			typeName: "BenchmarkNullableCell",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010201000b000101c001000900000c0ae007880db9"),
		},
		{
			name:     "unmarshal ref",
			typeName: "BenchmarkRefInt65",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010201000e000100010011000000000009689e40e150b4c5"),
		},
		{
			name:     "unmarshal empty tensor",
			typeName: "BenchmarkEmptyTensor",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c724101010100020000004cacb9cd"),
		},
		{
			name:     "unmarshal not empty tensor",
			typeName: "BenchmarkNotEmptyTensor",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101001f00003900000000000000000000000000021cb43b9aca00fffd550bfbaae07401a2a98117"),
		},
		{
			name:     "unmarshal small-int-key map",
			typeName: "BenchmarkSmallIntKeyMap",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010201000c000101c001000ba00000007bc09a662c32"),
		},
		{
			name:     "unmarshal small-uint-key map",
			typeName: "BenchmarkSmallUintKeyMap",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c72410104010053000101c0010202cb02030045a7400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe80045a3cff5555555555555555555555555555555555555555555555555555555555555555888440ce8"),
		},
		{
			name:     "unmarshal big-uint-key map",
			typeName: "BenchmarkBigUintKeyMap",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010301001a000101c0010115a70000000000000047550902000b000000001ab01d5bf1a9"),
		},
		{
			name:     "unmarshal bits-key map",
			typeName: "BenchmarkBitsKeyMap",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010301003b000101c0010106a0828502005ea0000000000000003e400b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe43b9aca00b89cdc86"),
		},
		{
			name:     "unmarshal address-key map",
			typeName: "BenchmarkAddressKeyMap",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010201002f000101c0010051a17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f9409502f9002016fdc16e"),
		},
		{
			name:     "unmarshal union with dec prefix",
			typeName: "BenchmarkDecUnion",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101001300002180000000000000000000000003b5577dc0660d6029"),
		},
		{
			name:     "unmarshal union with bin prefix",
			typeName: "AddressWithPrefix | MapWithPrefix | CellWithPrefix",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010201002e0001017801004fa17002c44ea652d4092859c67da44e4ca3add6565b0e2897d640a2c51bfb370d8877f900a4d89920c413c650"),
			abiFile:  "testdata/abi/bin_union.json",
		},
		{
			name:     "unmarshal union with hex prefix",
			typeName: "UInt66WithPrefix | UInt33WithPrefix | UInt4WithPrefix",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101000b000011deadbeef00000000c0d75977b9"),
			abiFile:  "testdata/abi/hex_union.json",
		},
		{
			name:     "unmarshal a-lot-refs from alias",
			typeName: "GoodNamingForMsg",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8"),
			abiFile:  "testdata/abi/refs.json",
		},
		{
			name:     "unmarshal a-lot-refs from struct",
			typeName: "ManyRefsMsg",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c724101040100b7000377deadbeef80107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e635087735940143ffffffffffffffffffffffffffff63c006010203004b80010df454cebee868f611ba8c0d4a9371fb73105396505783293a7625f75db3b9880bebc20100438006e05909e22b2e5e6087533314ee56505f85212914bd5547941a2a658ac62fe101004f801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfcc12309ce54001e09a48b8"),
			abiFile:  "testdata/abi/refs.json",
		},
		{
			name:     "unmarshal a-lot-generics from struct",
			typeName: "ManyRefsMsg<uint16>",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647"),
			abiFile:  "testdata/abi/generics.json",
		},
		{
			name:     "unmarshal a-lot-generics from alias",
			typeName: "GoodNamingForMsg<uint16>",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c72410103010043000217d017d7840343b9aca0000108010200080000007b005543b9aca001017d78402005889d4ca5a81250b38cfb489c99475bacacb61c512fac81458a37f66e1b10eff422fc7647"),
			abiFile:  "testdata/abi/generics.json",
		},
		{
			name:     "unmarshal struct with default values",
			typeName: "DefaultTest",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010101003100005d80000002414801622753296a04942ce33ed2272651d6eb2b2d87144beb2051628dfd9b86c43bfd00000156ac2c4c70811a9dde"),
			abiFile:  "testdata/abi/default_values.json",
		},
		{
			name:     "unmarshal a-lot-numbers",
			typeName: "Numbers",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c72410101010033000062000000000000000000000000000000000000000000000000000000000000000000000000000000f1106aecc4c800020926dc62f014"),
			abiFile:  "testdata/abi/numbers.json",
		},
		{
			name:     "unmarshal a-lot-random-fields",
			typeName: "RandomFields",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c7241010301007800028b79480107bfaaa5cc6e5368e5f9799188bd798cd22e04ab16d1d8ea4fc37480741e6350e038d7eb37c5e80000000ab50ee6b28000000000000016e4c000006c175300001801bc01020001c00051000000000005120041efeaa9731b94da397e5e64622f5e63348b812ac5b4763a93f0dd201d0798d4409e337ceb"),
			abiFile:  "testdata/abi/random_fields.json",
		},
		{
			name:     "unmarshal alias with custom unpack",
			typeName: "MyAlias",
			cell:     *boc.MustDeserializeSinglRootHex("b5ee9c724101010100470000890000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000043b9aca00886e91196"),
			abiFile:  "testdata/abi/custom_pack_unpack.json",
			customUnpackResolver: func(decoder *Decoder) func(alias parser.AliasRef, cell *boc.Cell, value *AliasValue) error {
				return func(alias parser.AliasRef, cell *boc.Cell, value *AliasValue) error {
					err := cell.Skip(512)
					if err != nil {
						return fmt.Errorf("failed to 512 bits from alias")
					}
					my := decoder.abiIndex.Structs["My"]
					val, err := decoder.UnmarshalTyIdx(cell, my.TyIdx)
					if err != nil {
						return fmt.Errorf("failed to unmarshal alias' coins")
					}
					*value = AliasValue(*val)
					return nil
				}
			},
		},
	} {
		abiFile := curr.abiFile
		if abiFile == "" {
			abiFile = "testdata/abi/benchmark_types.json"
		}
		contractABI := loadTestABI(b, abiFile)
		tyIdx := testTyIdx(b, contractABI, curr.typeName)
		b.Run(curr.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				curr.cell.ResetCounters()
				decoder := NewDecoder(contractABI)
				if curr.customUnpackResolver != nil {
					decoder.WithCustomUnpackResolver(curr.customUnpackResolver(decoder))
				}
				_, err := decoder.UnmarshalTyIdx(&curr.cell, tyIdx)
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
		abis := make([]parser.ContractABI, len(curr.abiFiles))
		for i, abiFile := range curr.abiFiles {
			data, err := os.ReadFile(abiFile)
			if err != nil {
				b.Fatal(err)
			}

			var abi parser.ContractABI
			err = json.Unmarshal(data, &abi)
			if err != nil {
				b.Fatal(err)
			}
			abis[i] = abi
		}
		b.Run(curr.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				curr.cell.ResetCounters()
				decoder := NewDecoder(abis[0])
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
