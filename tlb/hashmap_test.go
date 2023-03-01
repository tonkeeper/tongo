package tlb

import (
	"os"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

func TestHashmapAug_UnmarshalTLB(t *testing.T) {
	tests := []struct {
		name        string
		bocFilename string
		wantKeys    map[string]struct{}
	}{
		{
			bocFilename: "testdata/hashmap_aug.hex",
			wantKeys: map[string]struct{}{
				"824b9bca026a8ebaa17c278ab6fe342ea88766b54a8b22a1d61d5fcd1b3636f2": struct{}{},
				"85bec9755b8e384c44375b2b341efc87771f9d42feef0dd0395ce1dc2fc4eb40": struct{}{},
				"86e9617280e747235574d07d7b0202cf5a2b3c6d70f583a82dda3e8bb9b1b0c6": struct{}{},
				"8860d68280f629d34206817f7ba34d1536530d0ee260ba803ad20a59dd394013": struct{}{},
				"89678983e42fb0511cc593a0eff90403059051e3b442d0ff3d789452d8c5aefe": struct{}{},
				"91eb8789b4bdb4363bab92b747456896f254c0ba234177d54b4b43ee728d399d": struct{}{},
				"a349c3a670f6a84a105f7bc6cdb90954174ff3273e37e2da0b6a867cf689cd3b": struct{}{},
				"a8a817e2d33bbcf09204a841b3ac4d459c7c07a040d437005e7048fa27825958": struct{}{},
				"a955339e6a64d90ba5e7a8e19c2b00b54927f7de59e94f91c74ddcaabfa68b83": struct{}{},
				"ad9c21c5659d2775c6a303b43178c9c609c6799f59f72164a1eb57e69bc9f78d": struct{}{},
				"b10f11f48b2a929078096475365e9d0c1dc3d91ad1c7968ccb3f82e0d372100c": struct{}{},
				"b1ba2b9ccbec3c1375c77dbef5e422bd72c1ddc37a983f1212c77e3ba494ef62": struct{}{},
				"b370d166683c4e348b02b2113617d52f69bb8cf9fd4dbb1132c37c83683ae0f5": struct{}{},
				"b7499a24e15f5b02bb19c619c085fa6a92b4a326069ff012adb01d940af5d083": struct{}{},
				"b7fe6895380805f3821c7927877f87834fa86350a38473ab78057d2ef716fbb3": struct{}{},
				"b9eb34054eb95cf017a75267e70ac6efa63d95d5f8c0cebcf3f808d6c93ed4cb": struct{}{},
				"bc20148610d9f889e6a0df32b55952980b96b5fed7dae9ab3929e575a23d923d": struct{}{},
				"c26e9670833e413f739704ab3e3369e4351b72362fa60fca4eb225b83a9eeaf0": struct{}{},
				"c2b8b1d6496e507866870b2fa57e61280b17b7bb4ce0a2a5a8a2e4605ec0b616": struct{}{},
				"c7de8b1dccec0c8164f468b4ae4a4d2b85311bb2a9b023ae9af44840fc2f3d53": struct{}{},
				"c92fa5cfaa88ca3ea41e79f73f6c52d99b608e86085a00310c17758c436cf756": struct{}{},
				"d414f4e21e63eba9322d60660ad14c0f39e7eb5e0269bc6a8b756ff0ddef069f": struct{}{},
				"d425610dfbc21f02bf44a70b1c90d5d47d99bd4c8ec98286638058d42adbe03e": struct{}{},
				"d7b3df3705de80741a262fe293b703446913030a17ac3c27cd6cc61b94e3396e": struct{}{},
				"d97699ecebb6ab2f3b5749b5ecba7e68a65438c8ec4862a387be5a8dc2487f69": struct{}{},
				"dd6a08146ef8258edbf54a25e18b6db32664b8175b2f3b3687b26ccc594947f6": struct{}{},
				"de7f0b5e66d681b8cdb80af274216af6c5789e5aa4facf08d748308987da6d2f": struct{}{},
				"df82c045053fa76b6369fa17486f4dcd135710341b5d90f97b9a5ff4f1ee073d": struct{}{},
				"eedf414be3b4198d03bc7a65f272708f89716d852ea9aed1482e9b2a03594cc6": struct{}{},
				"f262d3f7ad776258f0d7c48b6af8f525923dccb7506a6003524b955072f27278": struct{}{},
				"f5c69091848e84f1fd5cd836d662469637a08d20f6928ce744af81f165daa184": struct{}{},
				"f935ce9f5c8655570e86e18c92eab7c14e66b5e1e270c31d21bba2cdc040e754": struct{}{},
				"fb52eec56cf1aeb319a6baf3f89f318ff08c234dbc7d40b29a0f99ee1d26821c": struct{}{},
				"fec7ca7301cbee6013be7c796cf19e0eaf8666c85c6e9915831c753a135ece91": struct{}{},
				"ffc0e3801749eaf340e407381ba50d2043082c9aceb7a93e84f447ce886bf9af": struct{}{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, err := os.ReadFile(tt.bocFilename)
			if err != nil {
				t.Fatalf("ReadFile() failed: %v", err)
			}
			cell, err := boc.DeserializeBocHex(string(bs))
			if err != nil {
				t.Fatalf("DeserializeBocHex() failed: %v", err)
			}
			var m HashmapAugE[Bits256, AccountBlock, CurrencyCollection]
			err = m.UnmarshalTLB(cell[0], &Decoder{})
			if err != nil {
				t.Fatalf("UnmarshalTLB() failed: %v", err)
			}
			values := m.Values()
			keys := map[string]struct{}{}
			for i, key := range m.Keys() {
				value := values[i]
				_ = value
				keys[key.Hex()] = struct{}{}
			}
			if !reflect.DeepEqual(tt.wantKeys, keys) {
				t.Fatalf("want: %v, got: %v", tt.wantKeys, keys)
			}
		})
	}
}

func Benchmark_HashmapAug_Unmarshal(b *testing.B) {
	data, err := os.ReadFile("testdata/hashmap_aug.hex")
	if err != nil {
		b.Errorf("ReadFile() failed: %v", err)
	}
	cell, err := boc.DeserializeBocHex(string(data))
	if err != nil {
		b.Errorf("boc.DeserializeBoc() failed: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cell[0].ResetCounters()
		var m HashmapAugE[Bits256, AccountBlock, CurrencyCollection]
		err = m.UnmarshalTLB(cell[0], &Decoder{})
		if err != nil {
			b.Errorf("UnmarshalTLB() failed: %v", err)
		}
	}
}
