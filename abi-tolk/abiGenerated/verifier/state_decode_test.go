package abiVerifier

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func TestVerifierRegistryStorageDecode(t *testing.T) {
	// 0:9f49427c96d43e10fa493f8bc21f1b8fef7171cdb8fc6fc2d50024eae0407da60
	cell := boc.MustDeserializeSinglRootHex("b5ee9c7241020d0100017300010381400102027302030387bf2abc83fb562e979b6da4a4279e39717194978a3473ea2f48b70859ca6d2106ce0040d68bad60a978f5d65b1ddd4c1f5562ccefb1a8106da3e36bc617136a96536101600405060387bf1252d302b2a608380a94092a4026a10c86f4de519f22f8272749689bc6342e1e002a7d17860f7ccf3a61297505f13879fd8a38225c08224ed8cb61c756429d7ad680e00b0c0c020120070800106f7262732e636f6d002068747470733a2f2f6f7262732e636f6d004bbffceb46647a5b9688aeede28abbdd47f4c15ab1f92b56b71f429fd93a61de268b01818181c002037aa0090a0049be9f8f251513cf1f25ad88816a4d8e3cfa3bd4c2f1fb45a12c1bf832351ef14cf0303030480049be9f59c646a30388d204c76724cf8535f6b378819cc49daac36a258096363d00e030303058004ba00395e7bf1b80f84a71a23ff0fecd9ccabefc460311e1cf75d77850236a1cb32ee000000010002076657269666965722e746f6e2e6f72677dab31ca")
	var storage VerifierRegistryStorage
	err := tlb.Unmarshal(cell, &storage)
	require.NoError(t, err, "Unmarshal VerifierRegistryStorage")

	keys := storage.Verifiers.Keys()
	require.Len(t, keys, 2)

	assert.Equal(t, mustSha256("orbs.com"), keys[0].HexString())
	assert.Equal(t, mustSha256("verifier.ton.org"), keys[1].HexString())

	t.Run("orbs.com", func(t *testing.T) {
		v, ok := storage.Verifiers.Get(keys[0])
		require.True(t, ok)
		assert.Equal(t, tlb.Uint8(2), v.Settings.MultiSigThreshold)
		assert.Equal(t, "orbs.com", v.Settings.Name)
		assert.Equal(t, "https://orbs.com", v.Settings.MarketingUrl)

		pubKeys := v.Settings.PubKeyEndpoints.Keys()
		require.NotEmpty(t, pubKeys)
		endpoint, ok := v.Settings.PubKeyEndpoints.Get(pubKeys[0])
		require.True(t, ok)
		assert.Equal(t, tlb.Uint32(50529027), endpoint)
	})

	t.Run("verifier.ton.org", func(t *testing.T) {
		v2, ok := storage.Verifiers.Get(keys[1])
		require.True(t, ok)
		assert.Equal(t, tlb.Uint8(1), v2.Settings.MultiSigThreshold)
		assert.Equal(t, "verifier.ton.org", v2.Settings.Name)
		assert.Equal(t, "verifier.ton.org", v2.Settings.MarketingUrl)

		v2PubKeys := v2.Settings.PubKeyEndpoints.Keys()
		require.Len(t, v2PubKeys, 1)
		assert.Equal(t, "1caf3df8dc07c2538d11ff87f66ce655f7e230188f0e7baebbc2811b50e59977", v2PubKeys[0].HexString())
		v2Endpoint, ok := v2.Settings.PubKeyEndpoints.Get(v2PubKeys[0])
		require.True(t, ok)
		assert.Equal(t, tlb.Uint32(0), v2Endpoint)
	})
}

func TestSourcesRegistryStorageDecode(t *testing.T) {
	// 0:fe049495509be2b9dfd0bfb6267dddeacd776b995e76e5edd80f3590a0720088
	cell := boc.MustDeserializeSinglRootHex("b5ee9c724101090100c5000197403ef148044190ab00801035a2eb582a5e3d7596c7775307d558b33bec6a041b68f8daf185c4daa594d850039f49427c96d43e10fa493f8bc21f1b8fef7171cdb8fc6fc2d50024eae0407da6010114ff00f4a413f4bcf2c80b0202016203040202ce05060009a0ef6de0030201200708001d403c8cbff12cbff01cf16ccc9ed54800531b088831c02456f8007434c0c05c6c2456f83e900c3c004c148131c17cb86472140133c5b250cc3c00a000273b513434fff4fffe900831c0a48c1b64b50c38a0b3fd9a60")
	var storage SourcesRegistryStorage
	err := tlb.Unmarshal(cell, &storage)
	require.NoError(t, err, "Unmarshal SourcesRegistryStorage")

	assert.Equal(t, tlb.Coins(66000000), storage.MinGram)
	assert.Equal(t, tlb.Coins(1100000000), storage.MaxGram)
	assert.Equal(t,
		ton.MustParseAccountID("EQDn0lCfJbUPhD6ST-Lwh8bj-9xcc24_G_C1QAk6uBAfaeMN").ToInternal(),
		storage.VerifierRegistryAddress)
	assert.Greater(t, storage.SourceItemCode.BitSize(), 0)
}

func TestSourceItemStorageDecode(t *testing.T) {
	// 0:1fb6e874e84c3bdf923872840a682f6661f66f212b3b6b5e789829ebe85a282c
	cell := boc.MustDeserializeSinglRootHex("b5ee9c724101020100aa0001c37494b4c0aca9820e02a5024a9009a84321bd379467c8be09c9d25a26f18d0b8788db7990749293ec0493a9b04ea5b65203e7152783a0c9702e89b950c28a4fca801fc09292aa137c573bfa17f6c4cfbbbd59aeed732bcedcbdbb01e6b2140e40111001008601697066733a2f2f6261666b72656965696a7265737a6e6c7464326779656a6e6336706668373772736c667333616b68797767347937763479347971796e687333657114c67da7")
	var storage SourceItemStorage
	err := tlb.Unmarshal(cell, &storage)
	require.NoError(t, err, "Unmarshal SourceItemStorage")

	assert.Equal(t, mustSha256("verifier.ton.org"), storage.VerifierId.HexString())
	assert.Equal(t, ton.MustParseAccountID("0:fe049495509be2b9dfd0bfb6267dddeacd776b995e76e5edd80f3590a0720088").ToInternal(), storage.SourceItemRegistry)
	assert.Equal(t, "88db7990749293ec0493a9b04ea5b65203e7152783a0c9702e89b950c28a4fca", storage.VerifiedCodeCellHash.HexString())

	if assert.NotNil(t, storage.Content.Value, "content should not be nil") {
		assert.Equal(t, tlb.Uint8(1), storage.Content.Value.Version)
		assert.Equal(t, "ipfs://bafkreieijresznltd2gyejnc6pfh77rslfs3akhywg4y7v4y4yqynhs3eq", string(storage.Content.Value.Url))
	}
}

func TestSourceItemMessageDecode(t *testing.T) {
	// https://tonviewer.com/transaction/8aaca8cc73259bcd684ab738d54d13c7df42efa9270adf7418437afc0541d6bf
	cell := boc.MustDeserializeSinglRootHex("b5ee9c7201010101004500008601697066733a2f2f6261666b72656967716b666d65626a636274686b61626c6b656b767463783571776a327978627468786c62647374683274626c3566627371646e79")
	var msg SourceContent
	require.NoError(t, tlb.Unmarshal(cell, &msg))
	assert.Equal(t, SourceContent{Version: 1, Url: "ipfs://bafkreigqkfmebjcbthkablkekvtcx5qwj2yxbthxlbdsth2tbl5fbsqdny"}, msg)
}

func TestVerifierRegistryMessageDecode(t *testing.T) {
	cases := []struct {
		name     string
		bocHex   string
		expected VerifierRegistryInternalMessage
	}{
		{
			name: "ForwardMessage",
			// https://tonviewer.com/transaction/4e3b8606f9fb7a2a054eb520c0b781121f67295811252c10f0b86c8ef9dbaaf2
			bocHex: "b5ee9c7201020501000170000218752177585abb8cab18ac4eb7010201cd7494b4c0aca9820e02a5024a9009a84321bd379467c8be09c9d25a26f18d0b876a3d096f800680fb2e470984a28b15b395f0a35fb3e26189abe4873abb24213c581ebe75b2d003f8125255426f8ae77f42fed899f777ab35ddae6579db97b7603cd64281c802220300c06fc8cc5649c3ff435f79e65537010e286086a32756ea46aacee1f03c7660f7ab6ecd677477fa83e304d34c5382a8ecab11b45eb98413905440a38c3e89c781031caf3df8dc07c2538d11ff87f66ce655f7e230188f0e7baebbc2811b50e599770198000003ea5abb8cab18ac4eb77494b4c0aca9820e02a5024a9009a84321bd379467c8be09c9d25a26f18d0b87385a5f0ddb841c17dd518a327c322d7d098e91c795eb9a152fe9be5cd20e40c904008601697066733a2f2f6261666b72656967716b666d65626a636274686b61626c6b656b767463783571776a327978627468786c62647374683274626c3566627371646e79",
			expected: func() VerifierRegistryInternalMessage {
				payloadCell := boc.MustDeserializeSinglRootHex("b5ee9c72010102010094000198000003ea5abb8cab18ac4eb77494b4c0aca9820e02a5024a9009a84321bd379467c8be09c9d25a26f18d0b87385a5f0ddb841c17dd518a327c322d7d098e91c795eb9a152fe9be5cd20e40c901008601697066733a2f2f6261666b72656967716b666d65626a636274686b61626c6b656b767463783571776a327978627468786c62647374683274626c3566627371646e79")
				sigCell := boc.MustDeserializeSinglRootHex("b5ee9c720101010100620000c06fc8cc5649c3ff435f79e65537010e286086a32756ea46aacee1f03c7660f7ab6ecd677477fa83e304d34c5382a8ecab11b45eb98413905440a38c3e89c781031caf3df8dc07c2538d11ff87f66ce655f7e230188f0e7baebbc2811b50e59977")
				return VerifierRegistryInternalMessage{
					SumType: VerifierRegistryInternalMessageKind_ForwardMessage,
					ForwardMessage: &ForwardMessage{
						QueryId: 6537973950539648695,
						Msg: tlb.RefT[*MessageDescription]{
							Value: &MessageDescription{
								VerifierId: mustUint256("7494b4c0aca9820e02a5024a9009a84321bd379467c8be09c9d25a26f18d0b87"),
								ValidUntil: 1782385007,
								SourceAddr: mustMsgAddress("0:3407d972384c251458ad9caf851afd9f130c4d5f2439d5d92109e2c0f5f3ad96"),
								TargetAddr: mustMsgAddress("0:fe049495509be2b9dfd0bfb6267dddeacd776b995e76e5edd80f3590a0720088"),
								Msg:        *payloadCell,
							},
						},
						Signatures: *sigCell,
					},
				}
			}(),
		},
		{
			name: "UpdateVerifier",
			// https://tonviewer.com/transaction/707e01131907becf2e98f25ee17c8ddb448d3d7b45450b8f09ef0925ae9c844e
			bocHex: "b5ee9c720101080100d000035b6002d61a00000000000000006aaf20fed58ba5e6db692909e78e5c5c6525e28d1cfa8bd22dc216729b4841b302c0010203020120040500106f7262732e636f6d002068747470733a2f2f6f7262732e636f6d004bbffceb46647a5b9688aeede28abbdd47f4c15ab1f92b56b71f429fd93a61de268b01818181c002037aa006070049be9f8f251513cf1f25ad88816a4d8e3cfa3bd4c2f1fb45a12c1bf832351ef14cf0303030480049be9f59c646a30388d204c76724cf8535f6b378819cc49daac36a258096363d00e030303058",
			expected: VerifierRegistryInternalMessage{
				SumType: VerifierRegistryInternalMessageKind_UpdateVerifier,
				UpdateVerifier: &UpdateVerifier{
					QueryId:    0,
					VerifierId: mustUint256("6aaf20fed58ba5e6db692909e78e5c5c6525e28d1cfa8bd22dc216729b4841b3"),
					Settings: VerifierSettings{
						MultiSigThreshold: 2,
						Name:              "orbs.com",
						MarketingUrl:      "https://orbs.com",
						PubKeyEndpoints: tlb.NewHashmapE(
							[]tlb.Uint256{
								mustUint256("79d68cc8f4b72d115ddbc51577ba8fe982b563f256ad6e3e853fb274c3bc4d16"),
								mustUint256("d1f8f251513cf1f25ad88816a4d8e3cfa3bd4c2f1fb45a12c1bf832351ef14cf"),
								mustUint256("d5f59c646a30388d204c76724cf8535f6b378819cc49daac36a258096363d00e"),
							},
							[]tlb.Uint32{
								tlb.Uint32(50529027),
								tlb.Uint32(50529028),
								tlb.Uint32(50529029),
							},
						),
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cell := boc.MustDeserializeSinglRootHex(c.bocHex)
			var msg VerifierRegistryInternalMessage
			require.NoError(t, tlb.Unmarshal(cell, &msg))
			assert.Equal(t, c.expected, msg)
		})
	}
}

func TestSourcesRegistryMessageDecode(t *testing.T) {
	cases := []struct {
		name     string
		bocHex   string
		expected SourcesRegistryInternalMessage
	}{
		{
			name: "ChangeAdmin",
			// https://tonviewer.com/transaction/0c2c2b41938d4a700b9d4bafd04f4f988a9adfcec9b75991264019a7955bba8d
			bocHex: "b5ee9c7201010101003000005b00000bbc0000000000000000800816ce77bc10a265365ec38613b27599642d0ed9aa3fe79fae45aa297d87e06eb0",
			expected: SourcesRegistryInternalMessage{
				SumType: SourcesRegistryInternalMessageKind_ChangeAdmin,
				ChangeAdmin: &ChangeAdmin{
					QueryId:  0,
					NewAdmin: ton.MustParseAccountID("0:40b673bde0851329b2f61c309d93accb216876cd51ff3cfd722d514bec3f0375").ToInternal(),
				},
			},
		},
		{
			name: "ChangeVerifierRegistry",
			// https://tonviewer.com/transaction/fbba4b34fc9a22088778986e70730402c82673bc4539f8e93f5c1643989fdd9d
			bocHex: "b5ee9c7201010101003000005b000007d3000000000000000080000000000000000000000000000000000000000000000000000000000000000010",
			expected: SourcesRegistryInternalMessage{
				SumType: SourcesRegistryInternalMessageKind_ChangeVerifierRegistry,
				ChangeVerifierRegistry: &ChangeVerifierRegistry{
					QueryId:             0,
					NewVerifierRegistry: ton.MustParseAccountID("0:0000000000000000000000000000000000000000000000000000000000000000").ToInternal(),
				},
			},
		},
		{
			name: "SetCode",
			// https://tonviewer.com/transaction/5cfe7e9385de55710dc8cb1e90ea26fcbcf91c561eb8a775fb10f4e12c07f559
			bocHex: "b5ee9c720101050100500001180000138e0000000000000000010114ff00f4a413f4bcf2c80b020201620304004ad03221c700915be0d0d3033071b09130e0d31fd33f3101812704ba94d4d1fb04e030f2c0cb000ba0467b0205cd",
			expected: SourcesRegistryInternalMessage{
				SumType: SourcesRegistryInternalMessageKind_SetCode,
				SetCode: &SetCode{
					QueryId: 0,
					NewCode: *boc.MustDeserializeSinglRootHex("b5ee9c72010104010041000114ff00f4a413f4bcf2c80b010201620203004ad03221c700915be0d0d3033071b09130e0d31fd33f3101812704ba94d4d1fb04e030f2c0cb000ba0467b0205cd"),
				},
			},
		},
		{
			name: "SetSourceItemCode",
			// https://tonviewer.com/transaction/bd1c25cdf67b6e8e46ede9befc8e95b63c2bad41e1c8fbf0f7c8a1b9af0db2f6
			bocHex: "b5ee9c7201010901009600011800000fa50000000000000000010114ff00f4a413f4bcf2c80b0202016203040202ce05060009a14123e0030201200708001d403c8cbff12cbff01cf16ccc9ed54800631b088831c02456f8007434c0c05c6c2456f83e900c3c004c00ece7d48931c17cb86472140133c5b27c00b817c16103fcbc2000393b513434fff4fffe900835d2b08025dfc0750c0510cc380c1c15481b60",
			expected: SourcesRegistryInternalMessage{
				SumType: SourcesRegistryInternalMessageKind_SetSourceItemCode,
				SetSourceItemCode: &SetSourceItemCode{
					QueryId:           0,
					NewSourceItemCode: *boc.MustDeserializeSinglRootHex("b5ee9c72010108010087000114ff00f4a413f4bcf2c80b0102016202030202ce04050009a14123e0030201200607001d403c8cbff12cbff01cf16ccc9ed54800631b088831c02456f8007434c0c05c6c2456f83e900c3c004c00ece7d48931c17cb86472140133c5b27c00b817c16103fcbc2000393b513434fff4fffe900835d2b08025dfc0750c0510cc380c1c15481b60"),
				},
			},
		},
		{
			name: "SetDeploymentCosts",
			// https://tonviewer.com/transaction/b7f809e22871bda156c92a265fd453e0fce441079f87c19746f95208b4eb4fe0
			bocHex: "b5ee9c7201010101001700002a000017770000000000000000403dfd24044190ab00",
			expected: SourcesRegistryInternalMessage{
				SumType: SourcesRegistryInternalMessageKind_SetDeploymentCosts,
				SetDeploymentCosts: &SetDeploymentCosts{
					QueryId:    0,
					NewMinGram: 65000000,
					NewMaxGram: 1100000000,
				},
			},
		},
		{
			name: "DeploySourceItemPayload",
			// https://tonviewer.com/transaction/c731c333d83ba1443218768bda56016493c1a39d43a9ca2d128eaeba1f03df1c
			bocHex: "b5ee9c72010102010087000198000003ea00000000000000006aaf20fed58ba5e6db692909e78e5c5c6525e28d1cfa8bd22dc216729b4841b3feb5ff6820e2ff0d9483e7e0d62c817d846789fb4ae580c878866d959dabd5c001006c01697066733a2f2f516d574145594e36446448726e6936413655476f55613976364c5077784270646e39574c5a38356e52456d345a42",
			expected: SourcesRegistryInternalMessage{
				SumType: SourcesRegistryInternalMessageKind_DeploySourceItemPayload,
				DeploySourceItemPayload: &DeploySourceItemPayload{
					QueryId:              0,
					VerifierId:           mustUint256("6aaf20fed58ba5e6db692909e78e5c5c6525e28d1cfa8bd22dc216729b4841b3"),
					VerifiedCodeCellHash: mustUint256("feb5ff6820e2ff0d9483e7e0d62c817d846789fb4ae580c878866d959dabd5c0"),
					SourceContent: tlb.RefT[*SourceContent]{
						Value: &SourceContent{Version: 1, Url: "ipfs://QmWAEYN6DdHrni6A6UGoUa9v6LPwxBpdn9WLZ85nREm4ZB"},
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cell := boc.MustDeserializeSinglRootHex(c.bocHex)
			var msg SourcesRegistryInternalMessage
			require.NoError(t, tlb.Unmarshal(cell, &msg))
			assert.Equal(t, c.expected, msg)
		})
	}
}

func mustSha256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

func mustUint256(hex string) tlb.Uint256 {
	var v tlb.Uint256
	(*big.Int)(&v).SetString(hex, 16)
	return v
}

func mustMsgAddress(addr string) tlb.MsgAddress {
	return ton.MustParseAccountID(addr).ToInternal().ToMsgAddress()
}
