package abi

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/tonkeeper/tongo/liteapi"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/tvm"
)

func mustToAddress(x string) tongo.AccountID {
	accountID, err := tongo.AccountIDFromRaw(x)
	if err != nil {
		panic(err)
	}
	return accountID
}

func TestGetPluginList(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		data    string
		account string
		want    GetPluginListResult
	}{
		{
			code:    "b5ee9c72010214010002d4000114ff00f4a413f4bcf2c80b01020120020f020148030602e6d001d0d3032171b0925f04e022d749c120925f04e002d31f218210706c7567bd22821064737472bdb0925f05e003fa403020fa4401c8ca07cbffc9d0ed44d0810140d721f404305c810108f40a6fa131b3925f07e005d33fc8258210706c7567ba923830e30d03821064737472ba925f06e30d0405007801fa00f40430f8276f2230500aa121bef2e0508210706c7567831eb17080185004cb0526cf1658fa0219f400cb6917cb1f5260cb3f20c98040fb0006008a5004810108f45930ed44d0810140d720c801cf16f400c9ed540172b08e23821064737472831eb17080185005cb055003cf1623fa0213cb6acb1fcb3fc98040fb00925f03e2020120070e020120080d020158090a003db29dfb513420405035c87d010c00b23281f2fff274006040423d029be84c600201200b0c0019adce76a26840206b90eb85ffc00019af1df6a26840106b90eb858fc00011b8c97ed44d0d70b1f80059bd242b6f6a2684080a06b90fa0218470d4080847a4937d29910ce6903e9ff9837812801b7810148987159f318404f8f28308d71820d31fd31fd31f02f823bbf264ed44d0d31fd31fd3fff404d15143baf2a15151baf2a205f901541064f910f2a3f80024a4c8cb1f5240cb1f5230cbff5210f400c9ed54f80f01d30721c0009f6c519320d74a96d307d402fb00e830e021c001e30021c002e30001c0039130e30d03a4c8cb1f12cb1fcbff10111213006ed207fa00d4d422f90005c8ca0715cbffc9d077748018c8cb05cb0222cf165005fa0214cb6b12ccccc973fb00c84014810108f451f2a7020070810108d718fa00d33fc8542047810108f451f2a782106e6f746570748018c8cb05cb025006cf165004fa0214cb6a12cb1fcb3fc973fb0002006c810108d718fa00d33f305224810108f459f2a782106473747270748018c8cb05cb025005cf165003fa0213cb6acb1f12cb3fc973fb00000af400c9ed54",
			data:    "b5ee9c720101060100a100015100001e8f29a9a3179190b2045703bc50371dd6897707b1791b8c4cbdd775e41b1ab7b3ebc8ae467bc001020581400c020502027303040041be5d7739ed1b853cab49680427bd0c78a375531c7a203092a115f1f9033cb7f6d00041be4d31620ef663745b654fb7bd023166f67fc8097fbfb34a5f6d60eb4480a303900041bf618db178d663998eaf849ac43cde4965aeac444373c92a65a10f10a2e6f406af",
			account: "0:9a7752cab755c829967b33e7f2692f9bdb81a47415168bb89c10b74ee0defc6b",
			want: GetPluginListResult{
				Plugins: []struct {
					Workchain int32
					Address   tlb.Bits256
				}{
					{Workchain: 0, Address: mustToAddress("0:70c6d8bc6b31ccc757c24d621e6f24b2d7562221b9e49532d0878851737a0357").Address},
					{Workchain: 0, Address: mustToAddress("0:4e698b1077b31ba2db2a7dbde8118b37b3fe404bfdfd9a52fb6b075a2405181c").Address},
					{Workchain: 0, Address: mustToAddress("0:4cebb9cf68dc29e55a4b40213de863c51baa98e3d101849508af8fc819e5bfb6").Address},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mainnetConfig, _ := boc.DeserializeBocBase64(mainnetConfig)
			code, err := hex.DecodeString(tt.code)
			if err != nil {
				t.Fatalf("DecodeString() failed: %v", err)
			}
			data, err := hex.DecodeString(tt.data)
			if err != nil {
				t.Fatalf("DecodeString() failed: %v", err)
			}
			accountID, err := tongo.AccountIDFromRaw(tt.account)
			if err != nil {
				t.Fatalf("AccountIDFromRaw() failed: %v", err)
			}
			codeCell, _ := boc.DeserializeBoc(code)
			dataCell, _ := boc.DeserializeBoc(data)
			emulator, err := tvm.NewEmulator(codeCell[0], dataCell[0], mainnetConfig[0], 1_000_000_000, -1)
			_, got, err := GetPluginList(context.Background(), emulator, accountID)
			if err != nil {
				t.Fatalf("GetPluginList() failed: %v", err)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestWhalesNominators(t *testing.T) {
	address := tongo.MustParseAccountID("EQBI-wGVp_x0VFEjd7m9cEUD3tJ_bnxMSp0Tb9qz757ATEAM")
	client, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		t.Fatal(err)
	}
	_, v, err := GetMembers(context.Background(), client, address)
	if err != nil {
		t.Fatal(err)
	}
	members := v.(GetMembers_WhalesNominatorResult).Members
	if len(members) == 0 || members[0].Address.SumType != "AddrStd" {
		t.Fatal(len(members))
	}
	_, v, err = GetPoolStatus(context.Background(), client, address)
	if err != nil {
		t.Fatal(err)
	}
	status := v.(GetPoolStatusResult)
	fmt.Printf("%+v\n", status)
	_, v, err = GetStakingStatus(context.Background(), client, address)
	if err != nil {
		t.Fatal(err)
	}
	stakingStatus := v.(GetStakingStatusResult)
	fmt.Printf("%+v\n", stakingStatus)
}
