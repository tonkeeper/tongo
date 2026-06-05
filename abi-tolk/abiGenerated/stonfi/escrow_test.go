package abiStonfi

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func TestStonfiEscrowFactory_IncomingMessages(t *testing.T) {
	tests := []struct {
		name     string
		txHash   string
		boc      string
		expected *FactoryIncomingMessage
	}{
		{
			// tx 2b9edeb3fe1b351fa6477b6c9cec341c37f6eec05d8421cf4d5acf3f1001d283
			name:   "MinterInitTransfer",
			txHash: "2b9edeb3fe1b351fa6477b6c9cec341c37f6eec05d8421cf4d5acf3f1001d283",
			boc:    "b5ee9c72010204010001170003a74716dd1600000000000000003b009af132ea2233d20f1f5523f244d9045112fcccd0a1dbb551bb5c046bbb4b5616c57200002000152b9e3e45dcf22f286053256bd622b81df6415fb20423ce5940588726466294010203009340002a573c7c8bb9e45e50c0a64ad7ac45703bec82bf6408479cb280b10e4c8cc52800f73db504c3bab6a5803055097b6bf35f18b8d5f0043d76323b1079f94b6c1e78c12309ce54000800010100d9800054ae78f91773c8bca1814c95af588ae077d9057ec8108f396501621c99198a5002e53f017b44348a6e123c3cfae17da7967f1958e16bcf64fcabeda5cba319ba500800054ae78f91773c8bca1814c95af588ae077d9057ec8108f396501621c99198a40000000011e1a301",
			expected: &FactoryIncomingMessage{
				SumType: FactoryIncomingMessageKind_MinterInitTransfer,
				MinterInitTransfer: &MinterInitTransfer{
					QuoteId:      u256("26687527438059201005361057571915321181268910921497390286994359456668219849547"),
					TransferType: tlb.Uint32(1444332914),
					RefFee:       addr(""),
					Excesses:     addr("0:02a573c7c8bb9e45e50c0a64ad7ac45703bec82bf6408479cb280b10e4c8cc52"),
					BidPaymentData: tlb.RefT[*BidPaymentData]{Value: &BidPaymentData{
						Recipient:       addr("0:02a573c7c8bb9e45e50c0a64ad7ac45703bec82bf6408479cb280b10e4c8cc52"),
						BidJettonWallet: addr("0:7b9eda8261dd5b52c0182a84bdb5f9af8c5c6af8021ebb191d883cfca5b60f3c"),
						Amount:          tlb.Grams(10000000000000),
					}},
					AskRefundData: tlb.RefT[*AskRefundData]{Value: &AskRefundData{RefundTo: addr("")}},
					UserPaymentData: tlb.RefT[*UserPaymentData]{Value: &UserPaymentData{
						User:                           addr("0:02a573c7c8bb9e45e50c0a64ad7ac45703bec82bf6408479cb280b10e4c8cc52"),
						UserReceiveJettonWallet:        addr("0:b94fc05ed10d229b848f0f3eb85f69e59fc656385af3d93f2afb6972e8c66e94"),
						UserExcesses:                   addr("0:02a573c7c8bb9e45e50c0a64ad7ac45703bec82bf6408479cb280b10e4c8cc52"),
						UserSafeDepositAndForwardValue: tlb.Uint64(150000000),
					}},
				},
			},
		},
		{
			// tx 49c61518c02d219ec3f86d9421057e07b05a9bf632e9ae08981c6279268e627c
			name:   "MinterRefundRequest",
			txHash: "49c61518c02d219ec3f86d9421057e07b05a9bf632e9ae08981c6279268e627c",
			boc:    "b5ee9c720101020100a20001b2d714500d00000000000000005dad97beac94017719d6e689ae83652b37053e3c4c6e7f55ed4b106818eeb3fa3079d06400cc707e0ceff5545feaae6c9265376d75f9e780ca3268d93737e6af26d83c06380000cb2a696fa2d8010087800c3e2c8eff4bc7e72c6c7cfe0b1ffb35897c497080d3020da8100d82d8ecaef9900331c1f833bfd5517faab9b24994ddb5d7e79e0328c9a364dcdf9abc9b60f018e008",
			expected: &FactoryIncomingMessage{
				SumType: FactoryIncomingMessageKind_MinterRefundRequest,
				MinterRefundRequest: &MinterRefundRequest{
					QuoteId:     u256("42371806764713266931649974104434782121910926036884287515841385037919863878650"),
					Amount:      tlb.Grams(498950),
					Recipient:   addr("0:cc707e0ceff5545feaae6c9265376d75f9e780ca3268d93737e6af26d83c0638"),
					ExitCode:    tlb.Uint32(52010),
					PrevMessage: tlb.Uint32(1768923864),
					ExtraFields: tlb.RefT[*MinterRefundRequestExtraFields]{Value: &MinterRefundRequestExtraFields{
						JettonWallet: addr("0:61f16477fa5e3f396363e7f058ffd9ac4be24b840698106d40806c16c76577cc"),
						Excesses:     addr("0:cc707e0ceff5545feaae6c9265376d75f9e780ca3268d93737e6af26d83c0638"),
					}},
				},
			},
		},
		{
			// tx 5c1ea5badc3bc539b833b61611c58632654c1b9697a34c907ee89683db867aad
			name:   "MinterGiveProtocolOwnership",
			txHash: "5c1ea5badc3bc539b833b61611c58632654c1b9697a34c907ee89683db867aad",
			boc:    "b5ee9c7201010101005100009d5df8fa350000000000000000800811ffec4c5344f9238a101b7575cfc5bda86a5a3483d6bc3b1da9f0ff0bfae9f000e4e7ec18bdc7ee40cdf23de021486145a0e19cc140c19e4618844828f064425a",
			expected: &FactoryIncomingMessage{
				SumType: FactoryIncomingMessageKind_MinterGiveProtocolOwnership,
				MinterGiveProtocolOwnership: &MinterGiveProtocolOwnership{
					NewProtocolAddress: addr("0:408fff62629a27c91c5080dbabae7e2ded4352d1a41eb5e1d8ed4f87f85fd74f"),
					Excesses:           addr("0:3939fb062f71fb90337c8f780852185168386730503067918621120a3c191096"),
				},
			},
		},
		{
			// tx bd842ab3bc3e0c73621b012aa7f0cdf2713875b56f6aab1bea6515d8804cd67d
			name:   "MinterTakeProtocolOwnership",
			txHash: "bd842ab3bc3e0c73621b012aa7f0cdf2713875b56f6aab1bea6515d8804cd67d",
			boc:    "b5ee9c7201010101003000005b881246140000000000000000800811ffec4c5344f9238a101b7575cfc5bda86a5a3483d6bc3b1da9f0ff0bfae9f0",
			expected: &FactoryIncomingMessage{
				SumType: FactoryIncomingMessageKind_MinterTakeProtocolOwnership,
				MinterTakeProtocolOwnership: &MinterTakeProtocolOwnership{
					Excesses: addr("0:408fff62629a27c91c5080dbabae7e2ded4352d1a41eb5e1d8ed4f87f85fd74f"),
				},
			},
		},
		{
			// minter_lock_payload, jetton transfer forward payload (msg 4a72c32c1e77b55fa0d6cb84d5fa3a13a8a1bfb6208270aaa5acc71eb3ab058d)
			name:   "MinterLockPayload",
			txHash: "4a72c32c1e77b55fa0d6cb84d5fa3a13a8a1bfb6208270aaa5acc71eb3ab058d",
			boc:    "b5ee9c72010206010001140003993e2466840000000000000000000000000000000000000000000000000000000000000001000010000a95cf1f22ee791794302992b5eb115c0efb20afd90211e72ca02c4393233149023c346010010203015c80198e0fc19dfeaa8bfd55cd924ca6edaebf3cf019464d1b26e6fcd5e4db0780c70000000000001d4c0817d7840004084202183f24bedaa94c422d8f1ad3a1a9750afbb0d7b0f53098e11bb25941000ebb7901878005739b3a8a3e5e1a458ccd2800e304af00583725925b08958736475e1933e29cd0000a95cf1f22ee791794302992b5eb115c0efb20afd90211e72ca02c4393233148400500010800432006314e53e499d8c39d73172af958240205ca87ba481bb2533f1d83f66e28cd777c",
			expected: &FactoryIncomingMessage{
				SumType: FactoryIncomingMessageKind_MinterLockPayload,
				MinterLockPayload: &MinterLockPayload{
					QuoteId:        u256("1"),
					RefFee:         addr(""),
					Excesses:       addr("0:02a573c7c8bb9e45e50c0a64ad7ac45703bec82bf6408479cb280b10e4c8cc52"),
					TonSafeDeposit: tlb.Grams(150000000),
					LockArgs: tlb.RefT[*BilateralLockArgs]{Value: &BilateralLockArgs{
						Resolver:             addr("0:cc707e0ceff5545feaae6c9265376d75f9e780ca3268d93737e6af26d83c0638"),
						ResolverTimeoutDelta: tlb.Uint64(60000),
						ResolverAskAmount:    tlb.Grams(200000000),
						DutchSegments: tlb.RefT[*DutchSegments]{Value: &DutchSegments{
							Segments: *boc.NewCell().CopyRemaining(),
						}},
					}},
					UnlockCondition: *boc.MustDeserializeSinglRootBase64("te6ccgEBAQEAIwAIQgIYPyS+2qlMQi2PGtOhqXUK+7DXsPUwmOEbsllBAA67eQ=="),
					ExtraFields: tlb.RefT[*MinterLockPayloadExtraFields]{Value: &MinterLockPayloadExtraFields{
						AskJettonWallet: addr("0:2b9cd9d451f2f0d22c6669400718257802c1b92c92d844ac39b23af0c99f14e6"),
						RefundTo:        addr("0:02a573c7c8bb9e45e50c0a64ad7ac45703bec82bf6408479cb280b10e4c8cc52"),
						More: tlb.RefT[*MinterLockPayloadExtraFieldsMore]{Value: &MinterLockPayloadExtraFieldsMore{
							OrderOwner:      addr(""),
							AskJettonMinter: addr("0:c629ca7c933b1873ae62e55f2b048040b950f74903764a67e3b07ecdc519aeef"),
						}},
					}},
				},
			},
		},
		{
			// minter_unlock_payload, jetton transfer forward payload (msg 5dddb574cef9f20d42b72c3b8eca5ca0e55e1f0f601258745eef82914c9ba25d)
			name:   "MinterUnlockPayload",
			txHash: "5dddb574cef9f20d42b72c3b8eca5ca0e55e1f0f601258745eef82914c9ba25d",
			boc:    "b5ee9c7201010301005100028b6a58f85c46d12544e4d4c4f6b39609f9dbff33d95c267eaeece0c8fc278a475651b90ada80198e0fc19dfeaa8bfd55cd924ca6edaebf3cf019464d1b26e6fcd5e4db0780c708010200010200030040",
			expected: &FactoryIncomingMessage{
				SumType: FactoryIncomingMessageKind_MinterUnlockPayload,
				MinterUnlockPayload: &MinterUnlockPayload{
					QuoteId:     u256("32031427659357264878515443776074461437146733909049316968113216565124114483930"),
					Excesses:    addr("0:cc707e0ceff5545feaae6c9265376d75f9e780ca3268d93737e6af26d83c0638"),
					ExtraFields: tlb.RefT[*UnlockPayloadExtraFields]{Value: &UnlockPayloadExtraFields{Recipient: addr(""), RefundTo: addr("")}},
					UnlockArgs:  tlb.RefT[*BilateralUnlockArgs]{Value: &BilateralUnlockArgs{}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := boc.MustDeserializeSinglRootHex(tt.boc)
			var msg FactoryIncomingMessage
			require.NoError(t, tlb.Unmarshal(cell, &msg))
			assertEqualDecoded(t, tt.expected, &msg)
			assert.True(t, cell.IsEmpty(), "cell should be fully consumed")
		})
	}
}

func TestStonfiEscrowPosition_IncomingMessages(t *testing.T) {
	tests := []struct {
		name     string
		txHash   string
		boc      string
		expected *PositionIncomingMessage
	}{
		{
			// tx c19b0d1b52c035a4640d13d9f9aab0ae67f540948c01461b300b2d04ae46ba24
			name:   "ItemInternalLock",
			txHash: "c19b0d1b52c035a4640d13d9f9aab0ae67f540948c01461b300b2d04ae46ba24",
			boc:    "b5ee9c72010206010001490003a77ab06181635862ac4d440252800e4c77ee3b482c706882be33e1f27a10c9d47e3e73c97107ab9923ba219fadf10806146581000d31d1ac6942099ccc803f9c2511616c22050af29b5c81babbbfb3caae1d22ab1001020301518007e28fa60d3977b85cb977ddb6a73dde1734ca86119ec08040a7bf3580403370600052047868c00404015c801ebdde4c025f4a056cec4660383e1e4d8587d9f2def99b5f3e0db375a725fab5800000000000000f08049319e205084202183f24bedaa94c422d8f1ad3a1a9750afbb0d7b0f53098e11bb25941000ebb7900c9800b7616c0fcd81f7c40c1dc067511e88583b4697e6a6961ebb80ec324c1e8a53f7000d31d1ac6942099ccc803f9c2511616c22050af29b5c81babbbfb3caae1d22ab2001a63a358d284133999007f384a22c2d8440a15e536b90375777f67955c3a455640000e1402498cf1001e",
			expected: &PositionIncomingMessage{
				SumType: PositionIncomingMessageKind_ItemInternalLock,
				ItemInternalLock: &ItemInternalLock{
					QueryId:      tlb.Uint64(7158580099875603026),
					TokenAddress: addr("0:7263bf71da4163834415f19f0f93d0864ea3f1f39e4b883d5cc91dd10cfd6f88"),
					Amount:       tlb.Grams(51000000),
					Excesses:     addr("0:34c746b1a50826733200fe70944585b088142bca6d7206eaeefecf2ab8748aac"),
					AdditionalField: tlb.RefT[*ItemAdditionalField]{Value: &ItemAdditionalField{
						AskJettonWallet:            addr("0:3f147d3069cbbdc2e5cbbeedb539eef0b9a654308cf60402053df9ac02019b83"),
						RefFee:                     addr(""),
						RefFeeTier:                 tlb.Uint16(10),
						SafeDepositAndForwardValue: tlb.Grams(150000000),
						More: tlb.RefT[*ItemAdditionalFieldMore]{Value: &ItemAdditionalFieldMore{
							AskJettonMinter: addr("0:5bb0b607e6c0fbe2060ee033a88f442c1da34bf3534b0f5dc07619260f4529fb"),
							OrderOwner:      addr("0:34c746b1a50826733200fe70944585b088142bca6d7206eaeefecf2ab8748aac"),
							RefundTo:        addr("0:34c746b1a50826733200fe70944585b088142bca6d7206eaeefecf2ab8748aac"),
						}},
					}},
					LockArgs: tlb.RefT[*BilateralLockArgs]{Value: &BilateralLockArgs{
						Resolver:             addr("0:f5eef26012fa502b67623301c1f0f26c2c3ecf96f7ccdaf9f06d9bad392fd5ac"),
						ResolverTimeoutDelta: tlb.Uint64(120),
						ResolverAskAmount:    tlb.Grams(38374641),
						DutchSegments: tlb.RefT[*DutchSegments]{Value: &DutchSegments{
							Num:      tlb.Uint4(1),
							Segments: *boc.MustDeserializeSinglRootBase64("te6ccgEBAQEACQAADUAkmM8QAeg="),
						}},
					}},
					UnlockCondition: *boc.MustDeserializeSinglRootBase64("te6ccgEBAQEAIwAIQgIYPyS+2qlMQi2PGtOhqXUK+7DXsPUwmOEbsllBAA67eQ=="),
				},
			},
		},
		{
			// tx b6ec3d9b17f1b251bc89b0105a779db3cbc50abf32446da04f1ea6ca03d28ed7
			name:   "ItemInternalUnlock",
			txHash: "b6ec3d9b17f1b251bc89b0105a779db3cbc50abf32446da04f1ea6ca03d28ed7",
			boc:    "b5ee9c720101030100c70002a7696fa2d800000000699ade292007af77930097d2815b3b11980e0f87936161f67cb7be66d7cf836cdd69c97ead64003f147d3069cbbdc2e5cbbeedb539eef0b9a654308cf60402053df9ac02019b83402498cf140102000b4030a32c104000c9801ebdde4c025f4a056cec4660383e1e4d8587d9f2def99b5f3e0db375a725fab59003d7bbc9804be940ad9d88cc0707c3c9b0b0fb3e5bdf336be7c1b66eb4e4bf56b2007af77930097d2815b3b11980e0f87936161f67cb7be66d7cf836cdd69c97ead640",
			expected: &PositionIncomingMessage{
				SumType: PositionIncomingMessageKind_ItemInternalUnlock,
				ItemInternalUnlock: &ItemInternalUnlock{
					QueryId:                  tlb.Uint64(1771757097),
					Resolver:                 addr("0:f5eef26012fa502b67623301c1f0f26c2c3ecf96f7ccdaf9f06d9bad392fd5ac"),
					ResolverSentJettonWallet: addr("0:3f147d3069cbbdc2e5cbbeedb539eef0b9a654308cf60402053df9ac02019b83"),
					ResolverSentAmount:       tlb.Grams(38374641),
					UnlockArgs:               tlb.RefT[*BilateralUnlockArgs]{Value: &BilateralUnlockArgs{MinOut: tlb.Grams(51000001)}},
					ExtraFields: tlb.RefT[*ItemInternalUnlockExtraFields]{Value: &ItemInternalUnlockExtraFields{
						Recipient: addr("0:f5eef26012fa502b67623301c1f0f26c2c3ecf96f7ccdaf9f06d9bad392fd5ac"),
						RefundTo:  addr("0:f5eef26012fa502b67623301c1f0f26c2c3ecf96f7ccdaf9f06d9bad392fd5ac"),
						Excesses:  addr("0:f5eef26012fa502b67623301c1f0f26c2c3ecf96f7ccdaf9f06d9bad392fd5ac"),
					}},
				},
			},
		},
		{
			// tx 8f463491f6c91df8a69a17e5fa27f57914c728fbb1b6e6bcd2772c331d8ab317
			name:   "ItemWithdraw",
			txHash: "8f463491f6c91df8a69a17e5fa27f57914c728fbb1b6e6bcd2772c331d8ab317",
			boc:    "b5ee9c7201010201003300015b12b8a98799c0431cc4cc062d801bf7ab7d092ffb8193f97bf8ced0c9081bbf15ac847adaacafb361d04e1ad91588010000",
			expected: &PositionIncomingMessage{
				SumType: PositionIncomingMessageKind_ItemWithdraw,
				ItemWithdraw: &ItemWithdraw{
					QueryId:      tlb.Uint64(11078928874171270701),
					WithdrawArgs: *boc.MustDeserializeSinglRootBase64("te6ccgEBAQEAAgAAAA=="),
					Excesses:     addr("0:dfbd5be8497fdc0c9fcbdfc676864840ddf8ad6423d6d5657d9b0e8270d6c8ac"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := boc.MustDeserializeSinglRootHex(tt.boc)
			var msg PositionIncomingMessage
			require.NoError(t, tlb.Unmarshal(cell, &msg))
			assertEqualDecoded(t, tt.expected, &msg)
			assert.True(t, cell.IsEmpty(), "cell should be fully consumed")
		})
	}
}

func TestStonfiEscrowPosition_ExternalMessages(t *testing.T) {
	tests := []struct {
		name     string
		txHash   string
		boc      string
		expected *PositionExternalMessage
	}{
		{
			// msg 306bd66765d0097fe6919403092e6de13cd3246d38e9d04a4f5c88ffceba6a14
			name:   "ExternalCronTrigger",
			txHash: "306bd66765d0097fe6919403092e6de13cd3246d38e9d04a4f5c88ffceba6a14",
			boc:    "b5ee9c7201010101002c0000532114702d800c4c7a1efeef93a245170deab302f80465e43fa30f24b63db5e5d8dd704a87cfd34e343150",
			expected: &PositionExternalMessage{
				SumType: PositionExternalMessageKind_ExternalCronTrigger,
				ExternalCronTrigger: &ExternalCronTrigger{
					RewardAddress: addr("0:6263d0f7f77c9d1228b86f559817c0232f21fd187925b1edaf2ec6eb82543e7e"),
					Salt:          tlb.Uint32(2591138186),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := boc.MustDeserializeSinglRootHex(tt.boc)
			var msg PositionExternalMessage
			require.NoError(t, tlb.Unmarshal(cell, &msg))
			assertEqualDecoded(t, tt.expected, &msg)
			assert.True(t, cell.IsEmpty(), "cell should be fully consumed")
		})
	}
}

func TestStonfiEscrowVault_IncomingMessages(t *testing.T) {
	tests := []struct {
		name     string
		txHash   string
		boc      string
		expected *VaultIncomingMessage
	}{
		{
			// tx 3587b08f3fcd4c5af1c713a0f09c142d66745a01ac0fbdefa188ce9733825802
			name:   "VaultDepositTokens",
			txHash: "3587b08f3fcd4c5af1c713a0f09c142d66745a01ac0fbdefa188ce9733825802",
			boc:    "b5ee9c72010101010033000061555edf4b00000000695ba8b222b2d801ebdde4c025f4a056cec4660383e1e4d8587d9f2def99b5f3e0db375a725fab5808",
			expected: &VaultIncomingMessage{
				SumType: VaultIncomingMessageKind_VaultDepositTokens,
				VaultDepositTokens: &VaultDepositTokens{
					QueryId:          tlb.Uint64(1767614642),
					Amount:           tlb.Grams(11053),
					Excesses:         addr("0:f5eef26012fa502b67623301c1f0f26c2c3ecf96f7ccdaf9f06d9bad392fd5ac"),
					ForwardTonAmount: tlb.Grams(0),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := boc.MustDeserializeSinglRootHex(tt.boc)
			var msg VaultIncomingMessage
			require.NoError(t, tlb.Unmarshal(cell, &msg))
			assertEqualDecoded(t, tt.expected, &msg)
			assert.True(t, cell.IsEmpty(), "cell should be fully consumed")
		})
	}
}

func addr(s string) tlb.MsgAddress {
	if s == "" {
		return tlb.MsgAddress{SumType: "AddrNone"}
	}
	var ia tlb.InternalAddress
	if err := ia.UnmarshalJSON([]byte(s)); err != nil {
		panic(err)
	}
	return ia.ToMsgAddress()
}

func u256(s string) tlb.Uint256 {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("invalid uint256: " + s)
	}
	return tlb.Uint256(*v)
}

func assertEqualDecoded(t *testing.T, expected, actual any) {
	t.Helper()

	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)
	actualJSON, err := json.Marshal(actual)
	require.NoError(t, err)
	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}
