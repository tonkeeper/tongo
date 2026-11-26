package abiMooncx_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/abi"
	abiMooncx "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/mooncx"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tep64"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// Mainnet accounts the samples below come from. The memepad factory is
// EQDOEt2Hvs6nOFA22KHUEa3BEbS7Hy2ed-9yllPy5nk8Hjhs; every memepad token has a curve
// contract that is also the jetton master of that token.
const (
	cpmmFactory       = "0:80606139fa881a29fbe515bf933bf198290b2e89538e2db2920ef00a27d8e4bc"
	protocolRecipient = "0:24606bd29c7a1546e7c66e6e04335adeac90f0af9773e05f217d3012b14309f1"
	graduatedCurve    = "0:b0e8286b645919ab2b4fd3b1fa87555d84b2d4b07908d64a8d4cc8a741ac66a2"
	graduatedPool     = "0:d86a49d6cc0fe87bf736a241099847ffed659764eea92a1a17b91a171603db0f"
	graduatedCreator  = "0:df36ec602af55dc78c2f2ffeac2e0fb3c01318967b51384367a7acd1a7382a22"
	// code hash of the jetton wallets deployed by the curves (wallet tx 8fbf7b578bf04138bfdbd80ff1b17bb4b6dfa0ded646517d200a4870e4773dcf)
	jettonWalletCodeHash = "131028e558c8f6f5687a64629c17a241074fd19e56ac75f8538f5b3096af078a"
	// wallet that buys and then sells tokens of the curve 0:6e09c58542f3090040c564ca9fbed5fd2a0393fa10e45edec484017b0758dabc
	seller = "0:7449c721a2fb2fa40207f4d6d034ed4660b0c1be546c6b9b94da00e9a6445baa"
	// custom_payload of the jetton burn of tx 06042351bd29fa7e072f8892b11b71375da054019f09103a14f4ddec4dc5f630 (sell order, minOut 882090000)
	sellPayload = "b5ee9c7201010101000b000011b97a94b743493a0108"
)

func TestMemepadFactoryMessages(t *testing.T) {
	t.Run("MemepadFactoryCreateCurve with initial buy", func(t *testing.T) {
		// tx 5997d8651962e8a1d22971f6747a8df608e6a3cbda72ebb8b568f92c3fdb739f
		msg := decodeBody[abiMooncx.MemepadFactoryCreateCurve](t, "b5ee9c7241010601006a0001357e838abb94c9818a09eb4f46b6bae201a71003ce0a5012a05f2008010401c0050203040006505650006a697066733a2f2f516d654c4c4d4c36517a4653395a4479724a573776736e42775334537052704c417a4268725247704c596632645801014005000670767052b4deb7", abiMooncx.MooncxMemepadFactoryCreateCurveMsgOp, abiMooncx.PrefixMemepadFactoryCreateCurve)
		assert.Equal(t, tlb.Uint64(10721242817804062534), msg.QueryId)
		assert.Equal(t, tlb.Uint64(13167084957296755662), msg.Salt)
		assert.Equal(t, tlb.Uint8(10), msg.FeeBps)
		assert.Equal(t, tlb.Coins(5000000000), msg.InitialBuy)
		assertCurveContent(t, msg.Content, "pvp", "PVP", "ipfs://QmeLLML6QzFS9ZDyrJW7vsnBwS4SpRpLAzBhrRGpLYf2dX")
	})
	t.Run("MemepadFactoryCreateCurve without initial buy", func(t *testing.T) {
		// tx 96bbab8efcfc11eb6104d79fda3f208bf0bd014c1bcb35dd0d027917e03253d5
		msg := decodeBody[abiMooncx.MemepadFactoryCreateCurve](t, "b5ee9c7241010801009700012b7e838abb4594d661e5869383fcd1077bcb39fca13208010401c005020304000e524143434f4f4e006a697066733a2f2f516d553264314179704e4a6231515959466231623563453559546a72574a4e416f526936765a6753795253454e630201c0050600184d6f6f6e20526163636f6f6e01013007003c68747470733a2f2f742e6d652f2b735f4168743267587131316d4d546b79438f715a", abiMooncx.MooncxMemepadFactoryCreateCurveMsgOp, abiMooncx.PrefixMemepadFactoryCreateCurve)
		assert.Equal(t, tlb.Uint64(5013868001121178499), msg.QueryId)
		assert.Equal(t, tlb.Uint64(18217350195962313889), msg.Salt)
		assert.Equal(t, tlb.Uint8(50), msg.FeeBps)
		assert.Equal(t, tlb.Coins(0), msg.InitialBuy)
		assertCurveContent(t, msg.Content, "Moon Raccoon", "RACCOON", "ipfs://QmU2d1AypNJb1QYYFb1b5cE5YTjrWJNAoRi6vZgSyRSENc")
	})
}

func TestMemepadCurveMessages(t *testing.T) {
	t.Run("MemepadCurveInit with initial buy", func(t *testing.T) {
		// sent by the factory in tx 5997d8651962e8a1d22971f6747a8df608e6a3cbda72ebb8b568f92c3fdb739f, deploys curve 0:ab07512b38cd032f84b9d55bce9d664fde175be7b40b0af733ca48d882c05048 in tx 422b10506d9184c5ca731c21667265c995b2717bfa69208d28309d8ed23eb606
		msg := decodeBody[abiMooncx.MemepadCurveInit](t, "b5ee9c7241021e010004540003a961c9641a94c9818a09eb4f4680100c0c273f5103453f7ca2b7f2677e33052165d12a71c5b65241de0144fb1c9790009181af4a71e8551b9f19b9b810cd6b7ab243c2be5dcf817c85f4c04ac50c27c5404a817c8020010203005580b1a2bc2ec500000802c68af0bb14000059b386e0aaa80ecd8fae90669900601d1a94a20000014000a0a80114ff00f4a413f4bcf2c80b04010300c00e020162050603c4d0f8918e34d31f31d72c20bc6a28cc96d33f31fa00308e11d72c23deecbef492f23fe1d33f31fa0030e2ed44d0fa0002a0c801fa02cec9ed54e0d72c20bc6a28cce302d72c207c53f52ce302d72c22caf83de4e302d72c269b90ac6431dc840ff2f0070809001da0f605da89a1f401f491f49061f05502e6ed44d001d33ffa00fa50fa50fa0006fa0020fa48fa4830f89221c70591308e3af892f82a28c8cf8420fa5213fa52c97829541242c8cf83cb04cf85a0ccccf91684f7b013800b5004d724c8cf8a0040ce12cbf7cf50c705f2e04ae25126a0c801fa02cec9ed5421935b345be30d216e915be30e0a0b01fed33ffa00fa48fa50f401fa0020f404016e913091d1e223fa4430f2d14df897f89370f83a237271e304f839206e8118b722e304216e811d135803e3045023a825a07381032c70f83ca00170f836a00170f836a07381040282100966018070f837a0bcf2b0ed44d0fa0020fa48fa4830f89222c705f2e0495338bef2af5138a10c01c6f897820afaf080bef2b0ed44d0fa0020fa48fa4830f89222c705f2e04904d33ffa00fa50f4055362bef2af5162a1c801fa0215cec9ed5488256e91359130e2c8cf91ef765f7a12cb3f01fa02fa52fa54ccc9c8cf858812fa5271cf0b6eccc98050fb000d0052c8cf91cd8b427226cf0b3f5005fa0213fa5415cec9c8cf850813fa5201fa0271cf0b6accc98011fb000068f8276f10f897a1f82fa07381040282100966018070f837b60972fb02c8cf850812fa528210d53276dbcf0b8ecb3fc9810082fb0000c0c801fa0212cec9ed54f82a26c8cf8420fa5213fa52c978c8cf905e3514661acb3f5008fa02fa5414fa5458fa02cec9c8cf898801547425c8cf83cb04cf85a0ccccf91684f7b004800b27d7243615ce12cbf781150dcf0b79ccccccc98050fb0000000201200f100143bff082eb663b57a00192f4a6ac467288df2dfeddb9da1bee28f6521c8bebd21f1ec011020120131401020012006a697066733a2f2f516d654c4c4d4c36517a4653395a4479724a573776736e42775334537052704c417a4268725247704c59663264580201201516020120191a0141bf4546a6ffe1b79cfdd86bad3db874313dcde2fb05e6a74aa7f3552d9617c79d131b0141bf6ed4f942a7848ce2cb066b77a1128c6a1ff8c43f438a2dce24612ba9ffab8b03170102001800065056500141bf5208def46f5a1d4f9dce66ab309f4a851305f166f91ef79d923ef58e34f9a2091b0141bf5d01fa5e3c06901c45046c6b2ddcea5af764fea0eed72a10d404f2312ceb247d1d0102001c000670767000040039fb72d66f", abiMooncx.MooncxMemepadCurveInitMsgOp, abiMooncx.PrefixMemepadCurveInit)
		assert.Equal(t, tlb.Uint64(10721242817804062534), msg.QueryId) // query id of the MemepadFactoryCreateCurve above
		assert.Equal(t, addr(cpmmFactory), msg.CpmmFactory)
		assert.Equal(t, addr(protocolRecipient), msg.ProtocolRecipient)
		assert.Equal(t, tlb.Coins(5000000000), msg.InitialBuy)
		require.NotNil(t, msg.Econ.Value)
		assert.Equal(t, abiMooncx.EconState{
			TokensForSale: 800000000000000000,
			TokensForLp:   200000000000000000,
			VirtualGram:   666666666666,
			VirtualTokens: 1066666666666400000,
			TargetGram:    2000000000000,
			TradeFeeBps:   20,
			CreatorFeeBps: 10,
			PoolFeeBps:    10,
		}, *msg.Econ.Value)
		assertCellHash(t, jettonWalletCodeHash, &msg.JwCode)
		assertMetadata(t, &msg.Metadata, "pvp", "PVP", "ipfs://QmeLLML6QzFS9ZDyrJW7vsnBwS4SpRpLAzBhrRGpLYf2dX")
	})
	t.Run("MemepadCurveInit without initial buy", func(t *testing.T) {
		// sent by the factory in tx 96bbab8efcfc11eb6104d79fda3f208bf0bd014c1bcb35dd0d027917e03253d5, deploys curve 0:6e09c58542f3090040c564ca9fbed5fd2a0393fa10e45edec484017b0758dabc in tx 33cea18a08810f309f04612a673e811eed007022ba4fceaadac017fd11ce4743
		msg := decodeBody[abiMooncx.MemepadCurveInit](t, "b5ee9c72410222010004a800039f61c9641a4594d661e586938380100c0c273f5103453f7ca2b7f2677e33052165d12a71c5b65241de0144fb1c9790009181af4a71e8551b9f19b9b810cd6b7ab243c2be5dcf817c85f4c04ac50c27c420010203005580b1a2bc2ec500000802c68af0bb14000059b386e0aaa80ecd8fae90669900601d1a94a2000006400503280114ff00f4a413f4bcf2c80b04010300c00e020162050603c4d0f8918e34d31f31d72c20bc6a28cc96d33f31fa00308e11d72c23deecbef492f23fe1d33f31fa0030e2ed44d0fa0002a0c801fa02cec9ed54e0d72c20bc6a28cce302d72c207c53f52ce302d72c22caf83de4e302d72c269b90ac6431dc840ff2f0070809001da0f605da89a1f401f491f49061f05502e6ed44d001d33ffa00fa50fa50fa0006fa0020fa48fa4830f89221c70591308e3af892f82a28c8cf8420fa5213fa52c97829541242c8cf83cb04cf85a0ccccf91684f7b013800b5004d724c8cf8a0040ce12cbf7cf50c705f2e04ae25126a0c801fa02cec9ed5421935b345be30d216e915be30e0a0b01fed33ffa00fa48fa50f401fa0020f404016e913091d1e223fa4430f2d14df897f89370f83a237271e304f839206e8118b722e304216e811d135803e3045023a825a07381032c70f83ca00170f836a00170f836a07381040282100966018070f837a0bcf2b0ed44d0fa0020fa48fa4830f89222c705f2e0495338bef2af5138a10c01c6f897820afaf080bef2b0ed44d0fa0020fa48fa4830f89222c705f2e04904d33ffa00fa50f4055362bef2af5162a1c801fa0215cec9ed5488256e91359130e2c8cf91ef765f7a12cb3f01fa02fa52fa54ccc9c8cf858812fa5271cf0b6eccc98050fb000d0052c8cf91cd8b427226cf0b3f5005fa0213fa5415cec9c8cf850813fa5201fa0271cf0b6accc98011fb000068f8276f10f897a1f82fa07381040282100966018070f837b60972fb02c8cf850812fa528210d53276dbcf0b8ecb3fc9810082fb0000c0c801fa0212cec9ed54f82a26c8cf8420fa5213fa52c978c8cf905e3514661acb3f5008fa02fa5414fa5458fa02cec9c8cf898801547425c8cf83cb04cf85a0ccccf91684f7b004800b27d7243615ce12cbf781150dcf0b79ccccccc98050fb0000000201200f10020120111202012017180142bfbf40462915a3e6026a4d790127b95ded4d870f6ab18d9af2fcbc454168255237130142bfa105d6cc76af400325e94d588ce511be5bfdbb73b437dc51eca43917d7a43e3d1501020014003c68747470733a2f2f742e6d652f2b735f4168743267587131316d4d546b7901020016006a697066733a2f2f516d553264314179704e4a6231515959466231623563453559546a72574a4e416f526936765a6753795253454e63020120191a0201201d1e0141bf4546a6ffe1b79cfdd86bad3db874313dcde2fb05e6a74aa7f3552d9617c79d131f0141bf6ed4f942a7848ce2cb066b77a1128c6a1ff8c43f438a2dce24612ba9ffab8b031b0102001c000e524143434f4f4e0141bf5208def46f5a1d4f9dce66ab309f4a851305f166f91ef79d923ef58e34f9a2091f0141bf5d01fa5e3c06901c45046c6b2ddcea5af764fea0eed72a10d404f2312ceb247d210102002000184d6f6f6e20526163636f6f6e0004003944b4a623", abiMooncx.MooncxMemepadCurveInitMsgOp, abiMooncx.PrefixMemepadCurveInit)
		assert.Equal(t, tlb.Uint64(5013868001121178499), msg.QueryId)
		assert.Equal(t, addr(cpmmFactory), msg.CpmmFactory)
		assert.Equal(t, addr(protocolRecipient), msg.ProtocolRecipient)
		assert.Equal(t, tlb.Coins(0), msg.InitialBuy)
		require.NotNil(t, msg.Econ.Value)
		assert.Equal(t, abiMooncx.EconState{
			TokensForSale: 800000000000000000,
			TokensForLp:   200000000000000000,
			VirtualGram:   666666666666,
			VirtualTokens: 1066666666666400000,
			TargetGram:    2000000000000,
			TradeFeeBps:   100,
			CreatorFeeBps: 80,
			PoolFeeBps:    50,
		}, *msg.Econ.Value)
		assertCellHash(t, jettonWalletCodeHash, &msg.JwCode)
		assertMetadata(t, &msg.Metadata, "Moon Raccoon", "RACCOON", "ipfs://QmU2d1AypNJb1QYYFb1b5cE5YTjrWJNAoRi6vZgSyRSENc")
	})
	t.Run("MemepadCurveBuy with min out", func(t *testing.T) {
		// tx ffe6cd0d92fbae1525ac40196fd725cd68f82bfe70316a5fcb10007dc1674447 (curve mints 1450201061650171 tokens with the same query id)
		msg := decodeBody[abiMooncx.MemepadCurveBuy](t, "b5ee9c7241010101001600002729d1f1c1c4c885649e5e1576704a30e1f3823e18e3bb3ffa", abiMooncx.MooncxMemepadCurveBuyMsgOp, abiMooncx.PrefixMemepadCurveBuy)
		assert.Equal(t, tlb.Uint64(14179730093976196470), msg.QueryId)
		assert.Equal(t, tlb.Coins(1305180955485153), msg.MinOut)
	})
	t.Run("MemepadCurveBuy without min out", func(t *testing.T) {
		// tx 3d0a623881c54450a4f349284fa88a3937c9c44054178eb621e7d597223daffd
		msg := decodeBody[abiMooncx.MemepadCurveBuy](t, "b5ee9c7241010101000f00001929d1f1c1000001a17760a4a90840004eab", abiMooncx.MooncxMemepadCurveBuyMsgOp, abiMooncx.PrefixMemepadCurveBuy)
		assert.Equal(t, tlb.Uint64(1793004184745), msg.QueryId)
		assert.Equal(t, tlb.Coins(0), msg.MinOut)
	})
	t.Run("MemepadCurveSellPayout", func(t *testing.T) {
		// sent by curve 0:6e09c58542f3090040c564ca9fbed5fd2a0393fa10e45edec484017b0758dabc in tx bf57ed944c7a5469619c65c6adf7ab17d7641d55393c228c7e332d1c0395eba4 after a burn_notification
		// with the same query id; gramOut = message value - incoming value + 50000
		// (curve reserve) + gas fee + forward fee of that transaction.
		msg := decodeBody[abiMooncx.MemepadCurveSellPayout](t, "b5ee9c72410101010014000023dee7a887408f488047e14cfe5062c415c0f87c1b1934", abiMooncx.MooncxMemepadCurveSellPayoutMsgOp, abiMooncx.PrefixMemepadCurveSellPayout)
		assert.Equal(t, tlb.Uint64(4652016655895973118), msg.QueryId)
		assert.Equal(t, tlb.Coins(26512284687), msg.GramOut)
	})
	t.Run("MemepadCurveSellPayload in a jetton burn", func(t *testing.T) {
		// tx 06042351bd29fa7e072f8892b11b71375da054019f09103a14f4ddec4dc5f630: the buyer of the MemepadCurveBuy with min out sample (tx ffe6cd0d)
		// sells the same 1450201061650171 tokens with jetton_burn to its jetton wallet 0:74fd45dec1fd9c704677132398158fbc9c05d09c647459d7ac50c52c19202152;
		// custom_payload is the sell order, and the curve 0:6e09c58542f3090040c564ca9fbed5fd2a0393fa10e45edec484017b0758dabc
		// pays 980100000 with MemepadCurveSellPayout in tx fbe7be8bf9095408c16d70eddd187dbce0ce51d5f92e4164f8a58f7d835c3ac3
		const burn = "b5ee9c7241010201004300016a595f07bc7203779bd0b6182c70526f33f21eefb800e8938e4345f65f48040fe9ada069da8cc161837ca8d8d73729b401d34c88b755010011b97a94b743493a0108519fbcd6"
		_, gotOpName, value, err := abi.InternalMessageDecoder(boc.MustDeserializeSinglRootHex(burn), nil)
		require.NoError(t, err)
		require.NotNil(t, gotOpName)
		require.Equal(t, abi.JettonBurnMsgOp, *gotOpName)
		body, ok := value.(abi.JettonBurnMsgBody)
		require.True(t, ok, "abi decoder returned %T", value)
		assert.Equal(t, uint64(8215541656359147564), body.QueryId)
		amount := big.Int(body.Amount)
		assert.Equal(t, "1450201061650171", amount.String())
		assert.Equal(t, msgAddr(seller), body.ResponseDestination)
		require.NotNil(t, body.CustomPayload)
		require.Equal(t, abi.MooncxMemepadCurveSellPayloadJettonOp, body.CustomPayload.SumType)
		order, ok := body.CustomPayload.Value.(abiMooncx.MemepadCurveSellPayload)
		require.True(t, ok, "custom payload decoded as %T", body.CustomPayload.Value)
		assert.Equal(t, tlb.Coins(882090000), order.MinOut)
		// the same payload cell on its own, through the jetton payload table
		assert.Equal(t, bodyHash(t, sellPayload), refHash(t, burn))
		assert.Equal(t, order, decodePayload[abiMooncx.MemepadCurveSellPayload](t, sellPayload, abi.MooncxMemepadCurveSellPayloadJettonOp, abiMooncx.PrefixMemepadCurveSellPayload))
	})
	t.Run("jetton burn_notification of the same sell", func(t *testing.T) {
		// tx fbe7be8bf9095408c16d70eddd187dbce0ce51d5f92e4164f8a58f7d835c3ac3: the jetton wallet notifies the curve; tongo still
		// decodes it as the standard burn_notification, whose extra ref is the sell payload of the burn
		const notification = "b5ee9c724101020100650001ad7bdd97de7203779bd0b6182c70526f33f21eefb800e8938e4345f65f48040fe9ada069da8cc161837ca8d8d73729b401d34c88b755001d1271c868becbe90081fd35b40d3b51982c306f951b1ae6e536803a699116eaa0010011b97a94b743493a0108acd8489d"
		_, gotOpName, value, err := abi.InternalMessageDecoder(boc.MustDeserializeSinglRootHex(notification), nil)
		require.NoError(t, err)
		require.NotNil(t, gotOpName)
		require.Equal(t, abi.JettonBurnNotificationMsgOp, *gotOpName)
		body, ok := value.(abi.JettonBurnNotificationMsgBody)
		require.True(t, ok, "abi decoder returned %T", value)
		assert.Equal(t, uint64(8215541656359147564), body.QueryId)
		amount := big.Int(body.Amount)
		assert.Equal(t, "1450201061650171", amount.String())
		assert.Equal(t, msgAddr(seller), body.Sender)
		assert.Equal(t, msgAddr(seller), body.ResponseDestination)
		assert.Equal(t, bodyHash(t, sellPayload), refHash(t, notification))
	})
	t.Run("MemepadCurveCreatorClaimFees", func(t *testing.T) {
		// tx 483a17ce68c2b7552fb772c2b62e7194a3b935d9205e1f0d57ed8f1c874c8f71 (creator payout carries the same query id)
		msg := decodeBody[abiMooncx.MemepadCurveCreatorClaimFees](t, "b5ee9c7241010101000e0000186ec846e7000001a177a22e8e38f9a3d2", abiMooncx.MooncxMemepadCurveCreatorClaimFeesMsgOp, abiMooncx.PrefixMemepadCurveCreatorClaimFees)
		assert.Equal(t, tlb.Uint64(1793008479886), msg.QueryId)
	})
	t.Run("MemepadCurveCreatorFeesPayout", func(t *testing.T) {
		// sent by the curve to the creator in the same tx 483a17ce68c2b7552fb772c2b62e7194a3b935d9205e1f0d57ed8f1c874c8f71
		msg := decodeBody[abiMooncx.MemepadCurveCreatorFeesPayout](t, "b5ee9c7241010101000e000018aa871c43000001a177a22e8ebc5aec78", abiMooncx.MooncxMemepadCurveCreatorFeesPayoutMsgOp, abiMooncx.PrefixMemepadCurveCreatorFeesPayout)
		assert.Equal(t, tlb.Uint64(1793008479886), msg.QueryId)
	})
}

// Graduation of curve graduatedCurve into the CPMM pool graduatedPool. The first
// readiness check was answered "not ready", so the curve asked again.
func TestMemepadGraduationMessages(t *testing.T) {
	const queryId = tlb.Uint64(39163) // query id of the buy that completed the curve
	t.Run("CpmmFactoryCreatePool", func(t *testing.T) {
		// sent by the curve to the CPMM factory in tx d67e0d8a27a5fc9bdfed13f4c32ebe82e95fb266b4638091eddec8163c0bbb2d
		msg := decodeBody[abiMooncx.CpmmFactoryCreatePool](t, "b5ee9c7241010101003100005dec725c4f00000000000098fb20058741435b22c8cd595a7e9d8fd43aaaec2596a583c846b2546a66453a0d633510f4164db274", abiMooncx.MooncxCpmmFactoryCreatePoolMsgOp, abiMooncx.PrefixCpmmFactoryCreatePool)
		assert.Equal(t, queryId, msg.QueryId)
		assert.Equal(t, msgAddr(""), msg.XMaster)
		assert.Equal(t, msgAddr(graduatedCurve), msg.YMaster)
		assert.Equal(t, tlb.Uint8(30), msg.FeeBps)
	})
	t.Run("CpmmFactoryProvidePoolAddress", func(t *testing.T) {
		// sent by the curve to the CPMM factory in tx d67e0d8a27a5fc9bdfed13f4c32ebe82e95fb266b4638091eddec8163c0bbb2d
		msg := decodeBody[abiMooncx.CpmmFactoryProvidePoolAddress](t, "b5ee9c724101010100520000a01d42af3700000000000098fb20058741435b22c8cd595a7e9d8fd43aaaec2596a583c846b2546a66453a0d633510f400b0e8286b645919ab2b4fd3b1fa87555d84b2d4b07908d64a8d4cc8a741ac66a2790fa2a3", abiMooncx.MooncxCpmmFactoryProvidePoolAddressMsgOp, abiMooncx.PrefixCpmmFactoryProvidePoolAddress)
		assert.Equal(t, queryId, msg.QueryId)
		assert.Equal(t, msgAddr(""), msg.XMaster)
		assert.Equal(t, addr(graduatedCurve), msg.YMaster)
		assert.Equal(t, tlb.Uint8(30), msg.FeeBps)
		assert.Equal(t, addr(graduatedCurve), msg.Initiator)
	})
	t.Run("CpmmFactoryTakePoolAddress", func(t *testing.T) {
		// tx a899e5e5f87c05f6f463f787659a2b4821f12e2727a571e02901469edb1782ac
		msg := decodeBody[abiMooncx.CpmmFactoryTakePoolAddress](t, "b5ee9c724101010100520000a081a5842200000000000098fb801b0d493ad981fd0f7ee6d448213308fffdacb2ec9dd5254342f72342e2c07b61e400b0e8286b645919ab2b4fd3b1fa87555d84b2d4b07908d64a8d4cc8a741ac66a21e0aedd8f0", abiMooncx.MooncxCpmmFactoryTakePoolAddressMsgOp, abiMooncx.PrefixCpmmFactoryTakePoolAddress)
		assert.Equal(t, queryId, msg.QueryId)
		assert.Equal(t, addr(graduatedPool), msg.Pool)
		assert.Equal(t, msgAddr(""), msg.XMaster)
		assert.Equal(t, addr(graduatedCurve), msg.YMaster)
		assert.Equal(t, tlb.Uint8(30), msg.FeeBps)
	})
	t.Run("CpmmPoolProvideReady first attempt", func(t *testing.T) {
		// tx ff28115493e2be2ed4adcf678ccf885b3599bee4012c7cad5bd73ab8191984ba
		msg := decodeBody[abiMooncx.CpmmPoolProvideReady](t, "b5ee9c7241010101000f00001a44daf70f00000000000098fb007543130f", abiMooncx.MooncxCpmmPoolProvideReadyMsgOp, abiMooncx.PrefixCpmmPoolProvideReady)
		assert.Equal(t, abiMooncx.CpmmPoolProvideReady{QueryId: queryId, Attempt: 0}, msg)
	})
	t.Run("CpmmPoolTakeReady not ready", func(t *testing.T) {
		// tx fb3f0bc4939f40766d4d281965d0a03f4df6eb3ce0e9dbd35505566b5e643dc5
		msg := decodeBody[abiMooncx.CpmmPoolTakeReady](t, "b5ee9c7241010101001000001b67071f0d00000000000098fb00406860b294", abiMooncx.MooncxCpmmPoolTakeReadyMsgOp, abiMooncx.PrefixCpmmPoolTakeReady)
		assert.Equal(t, abiMooncx.CpmmPoolTakeReady{QueryId: queryId, Attempt: 0, IsReady: false}, msg)
	})
	t.Run("CpmmPoolProvideReady second attempt", func(t *testing.T) {
		// tx d3013042410d3da00122a9794fdcdeea6e82d27aab524600dd86b0b12bb7b727
		msg := decodeBody[abiMooncx.CpmmPoolProvideReady](t, "b5ee9c7241010101000f00001a44daf70f00000000000098fb0176c078fd", abiMooncx.MooncxCpmmPoolProvideReadyMsgOp, abiMooncx.PrefixCpmmPoolProvideReady)
		assert.Equal(t, abiMooncx.CpmmPoolProvideReady{QueryId: queryId, Attempt: 1}, msg)
	})
	t.Run("CpmmPoolTakeReady ready", func(t *testing.T) {
		// tx f8580193bdcb2ebf591b0517ca3f4b65de68e77ba8306f40eebc6fe7511509df (the curve then sends the liquidity to the pool)
		msg := decodeBody[abiMooncx.CpmmPoolTakeReady](t, "b5ee9c7241010101001000001b67071f0d00000000000098fb01c067c3e605", abiMooncx.MooncxCpmmPoolTakeReadyMsgOp, abiMooncx.PrefixCpmmPoolTakeReady)
		assert.Equal(t, abiMooncx.CpmmPoolTakeReady{QueryId: queryId, Attempt: 1, IsReady: true}, msg)
	})
	t.Run("MemepadCurveExcess", func(t *testing.T) {
		// change of the buy that completed the curve, tx d67e0d8a27a5fc9bdfed13f4c32ebe82e95fb266b4638091eddec8163c0bbb2d
		msg := decodeBody[abiMooncx.MemepadCurveExcess](t, "b5ee9c7241010101000e000018ca6c577800000000000098fbfe916291", abiMooncx.MooncxMemepadCurveExcessMsgOp, abiMooncx.PrefixMemepadCurveExcess)
		assert.Equal(t, queryId, msg.QueryId)
	})
	t.Run("MemepadCurveProtocolPayout", func(t *testing.T) {
		// platform fees sent by the curve to the factory in tx f8580193bdcb2ebf591b0517ca3f4b65de68e77ba8306f40eebc6fe7511509df
		msg := decodeBody[abiMooncx.MemepadCurveProtocolPayout](t, "b5ee9c7241010101000e0000180b8cca5200000000000098fb6b4d96c4", abiMooncx.MooncxMemepadCurveProtocolPayoutMsgOp, abiMooncx.PrefixMemepadCurveProtocolPayout)
		assert.Equal(t, queryId, msg.QueryId)
	})
	t.Run("CpmmPoolDepositLiquidity", func(t *testing.T) {
		// GRAM leg sent by the curve to the pool in tx f8580193bdcb2ebf591b0517ca3f4b65de68e77ba8306f40eebc6fe7511509df:
		// locked LP, 70% to the creator and the rest to the protocol
		msg := decodeBody[abiMooncx.CpmmPoolDepositLiquidity](t, "b5ee9c7241010201005700011990af0b6300000000000098fb0e010089801be6dd8c055eabb8f185e5ffd585c1f678026312cf6a27086cf4f59a34e70544436b10009181af4a71e8551b9f19b9b810cd6b7ab243c2be5dcf817c85f4c04ac50c27c6d0ba233b", abiMooncx.MooncxCpmmPoolDepositLiquidityMsgOp, abiMooncx.PrefixCpmmPoolDepositLiquidity)
		assert.Equal(t, queryId, msg.QueryId)
		assert.Equal(t, tlb.Coins(0), msg.MinLpOut)
		assert.True(t, msg.Locked)
		require.True(t, msg.Split.Exists)
		require.NotNil(t, msg.Split.Value.Value)
		assert.Equal(t, abiMooncx.LpSplit{
			ARecipient: addr(graduatedCreator),
			ABps:       7000,
			BRecipient: addr(protocolRecipient),
		}, *msg.Split.Value.Value)
	})
}

// decodeBody decodes a message body with the generated type T, checks that the
// root cell is fully consumed, that the global abi decoder resolves the body to
// the expected operation, and that both decoded values encode back to the
// original cell.
func decodeBody[T interface{ ToCell() (*boc.Cell, error) }, PT interface {
	*T
	tlb.UnmarshalerTLB
}](t *testing.T, bocHex string, opName abi.MsgOpName, opCode uint64) T {
	t.Helper()
	want, err := boc.MustDeserializeSinglRootHex(bocHex).Hash256()
	require.NoError(t, err)

	cell := boc.MustDeserializeSinglRootHex(bocHex)
	var msg T
	require.NoError(t, PT(&msg).UnmarshalTLB(cell, tlb.NewDecoder()))
	assert.Zero(t, cell.BitsAvailableForRead(), "leftover bits")
	assert.Zero(t, cell.RefsAvailableForRead(), "leftover refs")
	assertEncodesTo(t, want, msg)

	gotOpCode, gotOpName, value, err := abi.InternalMessageDecoder(boc.MustDeserializeSinglRootHex(bocHex), nil)
	require.NoError(t, err)
	require.NotNil(t, gotOpCode)
	require.NotNil(t, gotOpName)
	assert.Equal(t, uint32(opCode), *gotOpCode)
	assert.Equal(t, opName, *gotOpName)
	decoded, ok := value.(*T)
	require.True(t, ok, "abi decoder returned %T", value)
	assertEncodesTo(t, want, *decoded)
	return msg
}

func assertEncodesTo[T interface{ ToCell() (*boc.Cell, error) }](t *testing.T, want [32]byte, msg T) {
	t.Helper()
	cell, err := msg.ToCell()
	require.NoError(t, err)
	got, err := cell.Hash256()
	require.NoError(t, err)
	assert.Equal(t, hex.EncodeToString(want[:]), hex.EncodeToString(got[:]), "re-encoded body differs")
}

// assertCurveContent checks the content cell carried by MemepadFactoryCreateCurve:
// name, symbol and image as text refs, then an optional ref with extra fields.
func assertCurveContent(t *testing.T, content boc.Cell, name, symbol, image string) {
	t.Helper()
	refs := content.Refs()
	require.Len(t, refs, 4)
	assert.Equal(t, name, cellText(t, refs[0]))
	assert.Equal(t, symbol, cellText(t, refs[1]))
	assert.Equal(t, image, cellText(t, refs[2]))
}

// assertMetadata checks the on-chain jetton metadata the factory builds from the MemepadFactoryCreateCurve content.
func assertMetadata(t *testing.T, metadata *boc.Cell, name, symbol, image string) {
	t.Helper()
	metadata.ResetCounters()
	var content tlb.FullContent
	require.NoError(t, tlb.Unmarshal(metadata, &content))
	meta, err := tep64.ConvertOnchainData(content)
	require.NoError(t, err)
	assert.Equal(t, name, meta.Name)
	assert.Equal(t, symbol, meta.Symbol)
	assert.Equal(t, image, meta.Image)
	assert.Equal(t, "9", meta.Decimals)
}

// bodyHash returns the hash of the root cell of a BoC, refHash the hash of the only ref of that root.
func bodyHash(t *testing.T, bocHex string) string {
	t.Helper()
	h, err := boc.MustDeserializeSinglRootHex(bocHex).Hash256()
	require.NoError(t, err)
	return hex.EncodeToString(h[:])
}

func refHash(t *testing.T, bocHex string) string {
	t.Helper()
	refs := boc.MustDeserializeSinglRootHex(bocHex).Refs()
	require.Len(t, refs, 1)
	h, err := refs[0].Hash256()
	require.NoError(t, err)
	return hex.EncodeToString(h[:])
}

func assertCellHash(t *testing.T, want string, cell *boc.Cell) {
	t.Helper()
	got, err := cell.Hash256()
	require.NoError(t, err)
	assert.Equal(t, want, hex.EncodeToString(got[:]))
}

func cellText(t *testing.T, cell *boc.Cell) string {
	t.Helper()
	cell.ResetCounters()
	var text tlb.Text
	require.NoError(t, tlb.Unmarshal(cell, &text))
	return string(text)
}

func addr(raw string) tlb.InternalAddress {
	return ton.MustParseAccountID(raw).ToInternal()
}

func msgAddr(raw string) tlb.MsgAddress {
	if raw == "" {
		return tlb.MsgAddress{SumType: "AddrNone"}
	}
	return addr(raw).ToMsgAddress()
}
