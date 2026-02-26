package abi

import (
	"bytes"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func mustToAny(t *testing.T, s string) *tlb.Any {
	b, err := boc.DeserializeBocHex(s)
	if err != nil {
		t.Fatal(err)
	}
	var r tlb.Any
	err = tlb.Unmarshal(b[0], &r)
	if err != nil {
		t.Fatal(err)
	}
	return &r
}

func TestJettonPayloadJSMarshaling(t *testing.T) {
	tests := []struct {
		name string
		boc  string
		want string
	}{
		{
			name: "Empty",
			boc:  "b5ee9c72010101010002000000",
			want: `{}`,
		},
		{
			name: "TextComment",
			boc:  "b5ee9c7201010101000b0000120000000048656c6c6f",
			want: `{"SumType":"TextComment","OpCode":0,"Value":{"Text":"Hello"}}`,
		},
		{
			name: "StonfiSwap",
			boc:  "b5ee9c7201010101006f0000d92593856180022a16a3164c4d5aa3133f3110ff10496e00ca8ac8abeffc5027e024d33480c3e8227f36ff000b3ced6e32247945f26b76f68700e1a82732b2ec38febd8c2d473f20548c96b2f001b23bbd527c00fd3d0874026bb7941f8fdd4d599ed8e7a63f426d5723f0388f2e",
			want: `{"SumType":"StonfiSwap","OpCode":630424929,"Value":{"TokenWallet":"0:1150b518b2626ad51899f98887f8824b70065456455f7fe2813f012699a4061f","MinOut":"289381247","ToAddress":"0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb","ReferralAddress":"0:6c8eef549f003f4f421d009aede507e3f7535667b639e98fd09b55c8fc0e23cb"}}`,
		},
		{
			name: "incomplete cell",
			boc:  "b5ee9c72010101010004000004053d",
			want: `{"SumType":"Cell","Value":"b5ee9c72010101010004000004053d"}`,
		},
		{
			name: "unknown payload",
			boc:  "b5ee9c720101010100120000206f6c6f6c6f6c6f6c6f74726f6c6f6c6f",
			want: `{"SumType":"Cell","OpCode":1869377388,"Value":"b5ee9c720101010100120000206f6c6f6c6f6c6f6c6f74726f6c6f6c6f"}`,
		},
		{
			name: "dedust swap",
			boc:  "b5ee9c72010102010034000153e3a0d482801c2bfd75b91e3b5a9402036fdb733e5dc096f4c773a328dd33586324e0cdb4bd2307b6a7400100090000000002",
			want: `{"SumType":"DedustSwap","OpCode":3818968194,"Value":{"Step":{"PoolAddr":"0:e15febadc8f1dad4a0101b7edb99f2ee04b7a63b9d1946e99ac31927066da5e9","Params":{"KindOut":false,"Limit":"505511","Next":null}},"SwapParams":{"Deadline":0,"RecipientAddr":"","ReferralAddr":"","FulfillPayload":null,"RejectPayload":null}}}`,
		},
		{
			name: "stonfi success",
			boc:  "b5ee9c72010101010006000008c64370e5",
			want: `{"SumType":"StonfiSwapOk","OpCode":3326308581,"Value":{}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := boc.DeserializeBocHex(tt.boc)
			assert(t, err)
			var j, j2 JettonPayload
			err = tlb.Unmarshal(c[0], &j)
			assert(t, err)
			b, err := json.Marshal(j)
			assert(t, err)
			if string(b) != tt.want {
				t.Errorf("JettonPayload.MarshalJSON() = %v, want %v", string(b), tt.want)
			}
			err = json.Unmarshal(b, &j2)
			assert(t, err)
			if j.SumType != UnknownJettonOp && !reflect.DeepEqual(j, j2) {
				t.Errorf("JettonPayload.UnmarshalJSON() = %v, want %v", j2, j)
			}
			c2 := boc.NewCell()
			err = tlb.Marshal(c2, j2)
			assert(t, err)
			s, err := c2.ToBocString()
			assert(t, err)
			if s != tt.boc {
				t.Errorf("JettonPayload.MarshalTLB() = %v, want %v", s, tt.boc)
			}
		})
	}
}

func TestJettonCustomUnmarshalling(t *testing.T) {
	tests := []struct {
		name string
		boc  string
		want any
	}{
		{
			name: "Valid jetton transfer with forward and without custom payload",
			boc:  "b5ee9c720101020100690001ac0f8a7ea5546de4ef59be1a6b5cdf061db67801465aa59db01447fc9fd217528b27ded0dbc07f3f7b540d3cc8504d52a46973050037ef56fa125ff70327f2f7f19da19210377e2b5908f5b5595f66c3a09c35b22b020301001c0000000031383437333938303832",
			want: JettonTransferMsgBody{
				QueryId:             6083770388301355627,
				Amount:              tlb.VarUInteger16(*big.NewInt(884501240679)),
				Destination:         mustToMsgAddress("0:a32d52ced80a23fe4fe90ba94593ef686de03f9fbdaa069e642826a95234b982"),
				ResponseDestination: mustToMsgAddress("0:dfbd5be8497fdc0c9fcbdfc676864840ddf8ad6423d6d5657d9b0e8270d6c8ac"),
				CustomPayload:       nil,
				ForwardTonAmount:    tlb.VarUInteger16(*big.NewInt(1)),
				ForwardPayload: tlb.EitherRef[JettonPayload]{
					IsRight: true,
					Value: JettonPayload{
						OpCode:  pointer(uint32(0)),
						SumType: TextCommentJettonOp,
						Value: TextCommentJettonPayload{
							Text: "1847398082",
						},
					},
				},
			},
		},
		{
			name: "Valid jetton transfer with forward and custom payload",
			boc:  "b5ee9c72010237010004fb0002ac0f8a7ea57361fe02fbd836f2539a1681a7a801f37b81bc8c38345c005f5b52aa177612b91f3ee6cd1ae18a2f5398fc0a1f1cb500333c7998a8b9d8e669d81830ad3daa01a1837f72c28c01cbc78cd0ee6ea9a054a203010201080df602d603001600000000393832303132310946034234ad7214de4fd3e58a483bad658e2daa61aa004ef0fa541c844732e283a311001e042205817002050628480101ac7196a16449b7019e14749fdd3eb33bf15d1d45abc74292fa781f7b013ea576001d22012007082848010113ee0db49bc066f2fcdf5ebc32a1b6fae201a50d9c4912d027ce0343f55a9e41001b220120090a2201200b0c284801017397e54bdf4c77f32d712cc2b102320a00bfcc5cb955319b812ab38a7b2a2db6001a2201200d0e28480101dce12c77b9d7c1eff024a3373a7fa1f7191bcba2405879b8cd82510b5fdf81060018284801017df47110b15a482d75e1465a7557c7f929a2a878cc179263ea4fba65baa9296700192201200f10284801016fafe2d5e2580903c9fdd59968e2a71941c528f90676069c334186cc5466ea7100162201201112220120131428480101f95a8c9a4b45ea9e1640eb1f01c0a903e2f08550233cccf5b8a6ffc00bf20c39001522012015162848010135c7b7aece7eb6d8746cbc5efe957cf2f58f20edccb0e1e87ae4899cba13f698001428480101d87d5ec47172689a0dea012d00897d083fc62e3f2dc79f5d9f30d08b9a6ef91b001322012017182848010162ecbcf1dc67b60938d74ad0623f31d23c243c939ea599cb649563e9f3fe850a0012220120191a284801017a494e49cad5216a39a3a47aea711055dfb18aee23cf30ab7d298a88a59b53fc00112201201b1c28480101bbd8d2fe0e489c6f749992d17892ec54a3f7ab659f063fb774866ee3f2efab1200102201201d1e2201201f2028480101c89f69f387d22673062227be63d7661d49ee6169bd512ffe50957f6c8ffcb4fc000e22012021222848010109c1df2f152428c52898acc302528493b8120c498f1069ad66fe810c5a708a05000d220120232428480101c8a2a5686729b0d23761a3bc1f8443c25a98c981b50e50e8952b03c185aa9eb3000c284801013ee950c82bb9203ca2b32db01512c74d238d0637b941dc282d5e96b7aea119bc000a220120252628480101fa7278136d1be8945413d1976eaefe2993175082fb2d1c847477ce4e240bb11d0008220120272828480101e49d0e7dd1cc094082c9f7937a7680ba9c4937412f1e06973fa1f8c3b2a2213b0007220120292a28480101e55e41515fc9ec90f8c5b0c29d3bcccbe8afadfcfc74e59458a82c5823cfb35800072201202b2c2201202d2e2848010190e8a0b6b7552c01212fcf1c2efde23961a411f538f0ceb7880c43ac96483c2300052201202f3028480101f0223ccf937da2f7102a14a0e452dccfd0971c99e481d141c685ec820d2cc8b600042848010179d4ddfaf086a46dae7a3e44674dce72629e213253440e99185326da781dbcd30002220120313228480101bc7bf637739ebf54b7d80b4473b7f7783a6e5f544cb72eedff9b5a4ebf09d48a0000220120333422047f3135362848010174b89492389e48103c1b14bc73b9ae6f86b1a4d7b3e235f11489a9cdaf7472060000005db828b9d8e669d81830ad3daa01a1837f72c28c01cbc78cd0ee6ea9a05494e685a069e800019bd4c2800001b2632c8228480101d9e142e6d1c8cd21136dafd313171a0b04de94f314a6f381f0bcd50e323fbd870000",
			want: JettonTransferMsgBody{
				QueryId:             8314205675871287026,
				Amount:              tlb.VarUInteger16(*big.NewInt(247521090170)),
				Destination:         mustToMsgAddress("0:f9bdc0de461c1a2e002fada9550bbb095c8f9f73668d70c517a9cc7e050f8e5a"),
				ResponseDestination: mustToMsgAddress("0:ccf1e662a2e76399a76060c2b4f6a806860dfdcb0a30072f1e3343b9baa68152"),
				CustomPayload:       mustToAny(t, "b5ee9c72010235010004940001080df602d6010946034234ad7214de4fd3e58a483bad658e2daa61aa004ef0fa541c844732e283a311001e022205817002030428480101ac7196a16449b7019e14749fdd3eb33bf15d1d45abc74292fa781f7b013ea576001d22012005062848010113ee0db49bc066f2fcdf5ebc32a1b6fae201a50d9c4912d027ce0343f55a9e41001b2201200708220120090a284801017397e54bdf4c77f32d712cc2b102320a00bfcc5cb955319b812ab38a7b2a2db6001a2201200b0c28480101dce12c77b9d7c1eff024a3373a7fa1f7191bcba2405879b8cd82510b5fdf81060018284801017df47110b15a482d75e1465a7557c7f929a2a878cc179263ea4fba65baa9296700192201200d0e284801016fafe2d5e2580903c9fdd59968e2a71941c528f90676069c334186cc5466ea7100162201200f10220120111228480101f95a8c9a4b45ea9e1640eb1f01c0a903e2f08550233cccf5b8a6ffc00bf20c39001522012013142848010135c7b7aece7eb6d8746cbc5efe957cf2f58f20edccb0e1e87ae4899cba13f698001428480101d87d5ec47172689a0dea012d00897d083fc62e3f2dc79f5d9f30d08b9a6ef91b001322012015162848010162ecbcf1dc67b60938d74ad0623f31d23c243c939ea599cb649563e9f3fe850a00122201201718284801017a494e49cad5216a39a3a47aea711055dfb18aee23cf30ab7d298a88a59b53fc0011220120191a28480101bbd8d2fe0e489c6f749992d17892ec54a3f7ab659f063fb774866ee3f2efab1200102201201b1c2201201d1e28480101c89f69f387d22673062227be63d7661d49ee6169bd512ffe50957f6c8ffcb4fc000e2201201f202848010109c1df2f152428c52898acc302528493b8120c498f1069ad66fe810c5a708a05000d220120212228480101c8a2a5686729b0d23761a3bc1f8443c25a98c981b50e50e8952b03c185aa9eb3000c284801013ee950c82bb9203ca2b32db01512c74d238d0637b941dc282d5e96b7aea119bc000a220120232428480101fa7278136d1be8945413d1976eaefe2993175082fb2d1c847477ce4e240bb11d0008220120252628480101e49d0e7dd1cc094082c9f7937a7680ba9c4937412f1e06973fa1f8c3b2a2213b0007220120272828480101e55e41515fc9ec90f8c5b0c29d3bcccbe8afadfcfc74e59458a82c5823cfb3580007220120292a2201202b2c2848010190e8a0b6b7552c01212fcf1c2efde23961a411f538f0ceb7880c43ac96483c2300052201202d2e28480101f0223ccf937da2f7102a14a0e452dccfd0971c99e481d141c685ec820d2cc8b600042848010179d4ddfaf086a46dae7a3e44674dce72629e213253440e99185326da781dbcd300022201202f3028480101bc7bf637739ebf54b7d80b4473b7f7783a6e5f544cb72eedff9b5a4ebf09d48a0000220120313222047f3133342848010174b89492389e48103c1b14bc73b9ae6f86b1a4d7b3e235f11489a9cdaf7472060000005db828b9d8e669d81830ad3daa01a1837f72c28c01cbc78cd0ee6ea9a05494e685a069e800019bd4c2800001b2632c8228480101d9e142e6d1c8cd21136dafd313171a0b04de94f314a6f381f0bcd50e323fbd870000"),
				ForwardTonAmount:    tlb.VarUInteger16(*big.NewInt(1)),
				ForwardPayload: tlb.EitherRef[JettonPayload]{
					IsRight: true,
					Value: JettonPayload{
						OpCode:  pointer(uint32(0)),
						SumType: TextCommentJettonOp,
						Value: TextCommentJettonPayload{
							Text: "9820121",
						},
					},
				},
			},
		},
		{
			name: "Jetton transfer with invalid custom payload",
			boc:  "b5ee9c720101010100630000c20f8a7ea57361fe02fbd836f2539a1681a7a801f37b81bc8c38345c005f5b52aa177612b91f3ee6cd1ae18a2f5398fc0a1f1cb500333c7998a8b9d8e669d81830ad3daa01a1837f72c28c01cbc78cd0ee6ea9a054a2020000000039383230313231",
			want: JettonTransferMsgBody{
				QueryId:             8314205675871287026,
				Amount:              tlb.VarUInteger16(*big.NewInt(247521090170)),
				Destination:         mustToMsgAddress("0:f9bdc0de461c1a2e002fada9550bbb095c8f9f73668d70c517a9cc7e050f8e5a"),
				ResponseDestination: mustToMsgAddress("0:ccf1e662a2e76399a76060c2b4f6a806860dfdcb0a30072f1e3343b9baa68152"),
				CustomPayload:       nil,
				ForwardTonAmount:    tlb.VarUInteger16(*big.NewInt(1)),
				ForwardPayload: tlb.EitherRef[JettonPayload]{
					IsRight: false,
					Value: JettonPayload{
						OpCode:  pointer(uint32(0)),
						SumType: TextCommentJettonOp,
						Value: TextCommentJettonPayload{
							Text: "9820121",
						},
					},
				},
			},
		},
		{
			name: "Jetton transfer with invalid forward payload",
			boc:  "b5ee9c720101010100580000ac0f8a7ea57361fe02fbd836f2539a1681a7a801f37b81bc8c38345c005f5b52aa177612b91f3ee6cd1ae18a2f5398fc0a1f1cb500333c7998a8b9d8e669d81830ad3daa01a1837f72c28c01cbc78cd0ee6ea9a0548203",
			want: JettonTransferMsgBody{
				QueryId:             8314205675871287026,
				Amount:              tlb.VarUInteger16(*big.NewInt(247521090170)),
				Destination:         mustToMsgAddress("0:f9bdc0de461c1a2e002fada9550bbb095c8f9f73668d70c517a9cc7e050f8e5a"),
				ResponseDestination: mustToMsgAddress("0:ccf1e662a2e76399a76060c2b4f6a806860dfdcb0a30072f1e3343b9baa68152"),
				CustomPayload:       nil,
				ForwardTonAmount:    tlb.VarUInteger16(*big.NewInt(1)),
				ForwardPayload: tlb.EitherRef[JettonPayload]{
					IsRight: true,
					Value:   JettonPayload{},
				},
			},
		},
		{
			name: "Jetton transfer without forward payload",
			boc:  "b5ee9c720101010100580000ab0f8a7ea57361fe02fbd836f2539a1681a7a801f37b81bc8c38345c005f5b52aa177612b91f3ee6cd1ae18a2f5398fc0a1f1cb500333c7998a8b9d8e669d81830ad3daa01a1837f72c28c01cbc78cd0ee6ea9a0548203",
			want: JettonTransferMsgBody{
				QueryId:             8314205675871287026,
				Amount:              tlb.VarUInteger16(*big.NewInt(247521090170)),
				Destination:         mustToMsgAddress("0:f9bdc0de461c1a2e002fada9550bbb095c8f9f73668d70c517a9cc7e050f8e5a"),
				ResponseDestination: mustToMsgAddress("0:ccf1e662a2e76399a76060c2b4f6a806860dfdcb0a30072f1e3343b9baa68152"),
				CustomPayload:       nil,
				ForwardTonAmount:    tlb.VarUInteger16(*big.NewInt(1)),
				ForwardPayload:      tlb.EitherRef[JettonPayload]{},
			},
		},
		{
			name: "Jetton transfer with empty payload",
			boc:  "b5ee9c720101010100580000ac0f8a7ea57361fe02fbd836f2539a1681a7a801f37b81bc8c38345c005f5b52aa177612b91f3ee6cd1ae18a2f5398fc0a1f1cb500333c7998a8b9d8e669d81830ad3daa01a1837f72c28c01cbc78cd0ee6ea9a0548202",
			want: JettonTransferMsgBody{
				QueryId:             8314205675871287026,
				Amount:              tlb.VarUInteger16(*big.NewInt(247521090170)),
				Destination:         mustToMsgAddress("0:f9bdc0de461c1a2e002fada9550bbb095c8f9f73668d70c517a9cc7e050f8e5a"),
				ResponseDestination: mustToMsgAddress("0:ccf1e662a2e76399a76060c2b4f6a806860dfdcb0a30072f1e3343b9baa68152"),
				CustomPayload:       nil,
				ForwardTonAmount:    tlb.VarUInteger16(*big.NewInt(1)),
				ForwardPayload:      tlb.EitherRef[JettonPayload]{},
			},
		},
		{
			name: "Jetton notify",
			boc:  "b5ee9c720101020100470001687362d09c546de4ef0759cac2601fe55ff623c80013d1ae5e0177ce7efcc4d75d6b2140dfcc27b135b340d93453031837bd2f8a5f01001c0000000034373035303031333230",
			want: JettonNotifyMsgBody{
				QueryId: 6083770386919049922,
				Sender:  mustToMsgAddress("0:09e8d72f00bbe73f7e626baeb590a06fe613d89ad9a06c9a29818c1bde97c52f"),
				Amount:  tlb.VarUInteger16(*big.NewInt(2191876121148)),
				ForwardPayload: tlb.EitherRef[JettonPayload]{
					IsRight: true,
					Value: JettonPayload{
						OpCode:  pointer(uint32(0)),
						SumType: TextCommentJettonOp,
						Value: TextCommentJettonPayload{
							Text: "4705001320",
						},
					},
				},
			},
		},
		//
		{
			name: "Jetton internal transfer",
			boc:  "b5ee9c720101020100600001ad178d45190000000067065c7e66532ecb20c5f80187e3b51d9b51f19a85759c449d47e36d2e3d54fea98925a9a320a4f69eae7ef90030fc76a3b36a3e3350aeb38893a8fc6da5c7aa9fd53124b53464149ed3d5cfdf040701000800000000",
			want: JettonInternalTransferMsgBody{
				QueryId:          1728470142,
				From:             mustToMsgAddress("0:c3f1da8ecda8f8cd42bace224ea3f1b6971eaa7f54c492d4d190527b4f573f7c"),
				ResponseAddress:  mustToMsgAddress("0:c3f1da8ecda8f8cd42bace224ea3f1b6971eaa7f54c492d4d190527b4f573f7c"),
				Amount:           tlb.VarUInteger16(*big.NewInt(111269393861727)),
				ForwardTonAmount: tlb.VarUInteger16(*big.NewInt(1)),
				ForwardPayload: tlb.EitherRef[JettonPayload]{
					IsRight: true,
					Value: JettonPayload{
						OpCode:  pointer(uint32(0)),
						SumType: TextCommentJettonOp,
						Value: TextCommentJettonPayload{
							Text: "",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := boc.DeserializeBocHex(tt.boc)
			assert(t, err)
			switch tt.want.(type) {
			case JettonTransferMsgBody:
				var j JettonTransferMsgBody
				err = c[0].Skip(32)
				assert(t, err)
				err = tlb.Unmarshal(c[0], &j)
				assert(t, err)
				r1, _ := json.Marshal(j)
				r2, _ := json.Marshal(tt.want)
				if !bytes.Equal(r1, r2) {
					t.Errorf("\nMsg = %v\n want %v", string(r1), string(r2))
				}
			case JettonNotifyMsgBody:
				var j JettonNotifyMsgBody
				err = c[0].Skip(32)
				assert(t, err)
				err = tlb.Unmarshal(c[0], &j)
				assert(t, err)
				r1, _ := json.Marshal(j)
				r2, _ := json.Marshal(tt.want)
				if !bytes.Equal(r1, r2) {
					t.Errorf("\nMsg = %v\n want %v", string(r1), string(r2))
				}
			case JettonInternalTransferMsgBody:
				var j JettonInternalTransferMsgBody
				err = c[0].Skip(32)
				assert(t, err)
				err = tlb.Unmarshal(c[0], &j)
				assert(t, err)
				r1, _ := json.Marshal(j)
				r2, _ := json.Marshal(tt.want)
				if !bytes.Equal(r1, r2) {
					t.Errorf("\nMsg = %v\n want %v", string(r1), string(r2))
				}
			}
		})
	}
}
