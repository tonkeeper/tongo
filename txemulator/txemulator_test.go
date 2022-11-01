package txemulator

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/liteclient"
	"github.com/startfellows/tongo/tlb"
	"github.com/startfellows/tongo/tvm"
	"github.com/startfellows/tongo/wallet"
)

func TestExec(t *testing.T) {
	recipientAddr, _ := tongo.AccountIDFromRaw("0:507dea7d606f22d9e85678d3eede39bbe133a868d2a0e3e07f5502cb70b8a512")
	pk, _ := base64.StdEncoding.DecodeString("OyAWIb4FeP1bY1VhALWrU2JN9/8O1Kv8kWZ0WfXXpOM=")
	privateKey := ed25519.NewKeyFromSeed(pk)

	client, c := wallet.NewMockBlockchain(1, tongo.AccountInfo{Balance: 1000})
	w, err := wallet.NewWallet(privateKey, wallet.V4R2, 0, nil, client)
	if err != nil {
		log.Fatalf("Unable to create wallet: %v", err)
	}
	tongoClient, err := liteclient.NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}

	config, err := tongoClient.GetLastConfigAll(context.Background())
	if err != nil {
		log.Fatalf("Get account state error: %v", err)
	}

	account, err := tongoClient.GetLastRawAccount(context.Background(), w.GetAddress())
	if err != nil {
		log.Fatalf("Get account state error: %v", err)
	}

	comment := "hello"
	tonTransfer := wallet.Message{
		Amount:  10000,
		Address: *recipientAddr,
		Comment: &comment,
	}

	accountID := w.GetAddress()
	res, err := tvm.RunTvm(
		&account.Account.Storage.State.AccountActive.StateInit.Code.Value.Value,
		&account.Account.Storage.State.AccountActive.StateInit.Data.Value.Value,
		"seqno", []tvm.StackEntry{}, &accountID)
	if err != nil {
		log.Fatalf("TVM run error: %v", err)
	}
	if res.ExitCode != 0 || len(res.Stack) != 1 || !res.Stack[0].IsInt() {
		log.Fatalf("TVM execution failed")
	}

	err = w.SimpleSend(context.Background(), []wallet.Message{tonTransfer})
	if err != nil {
		log.Fatalf("Unable to generate transfer message: %v", err)
	}

	msg := <-c
	var message tongo.Message
	cell, err := boc.DeserializeBoc(msg)
	if err != nil {
		log.Fatalf("unable to deserialize transfer message: %v", err)
	}
	err = tlb.Unmarshal(cell[0], &message)
	if err != nil {
		log.Fatalf("unable to unmarshal transfer message: %v", err)
	}

	var shardAccount tongo.ShardAccount
	shardAccount.Account = account
	shardAccount.LastTransLt = account.Account.Storage.LastTransLt - 1

	e, err := NewEmulator(config, PrintsAllStackValuesForCommand)
	if err != nil {
		log.Fatalf("unable to create emulator: %v", err)
	}

	e.SetVerbosityLevel(0)

	emRes, err := e.Emulate(shardAccount, message)
	if err != nil {
		log.Fatalf("emulator error: %v", err)
	}
	if emRes.Emulation == nil {
		log.Fatalf("empty emulation")
	}
	fmt.Printf("Account last transaction hash: %x\n", emRes.Emulation.ShardAccount.LastTransHash)
	fmt.Printf("Transaction lt: %v\n", emRes.Emulation.Transaction.Lt)
}

func TestGetConfigExec(t *testing.T) {

	tongoClient, err := liteclient.NewClientWithDefaultMainnet() //
	// tongoClient, err := liteclient.NewClientWithDefaultTestnet() //
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}

	mcExtra, err := tongoClient.GetConfigAll(context.Background())
	if err != nil {
		log.Fatalf("Get account state error: %v", err)
	}

	config := mcExtra.Config
	t.Log("config addr: ", config.ConfigAddr.Hex())
	for i := range config.Config.Hashmap.Keys() {
		if binary.BigEndian.Uint32(config.Config.Hashmap.Keys()[i].Buffer()) == 34 {
			str := config.Config.Hashmap.Values()[i].Value.RawBitString()
			fmt.Printf("key: %v, value: %x\n", config.Config.Hashmap.Keys()[i].BinaryString(), str.Buffer())
			var validatorSet tongo.ValidatorsSet
			err := tlb.Unmarshal(&config.Config.Hashmap.Values()[i].Value, &validatorSet)
			if err != nil {
				t.Fatalf("Unmarshal validator set error: %v", err)
			}
			t.Log("SumType:         ", validatorSet.SumType)
			t.Log("TotalWeight:     ", validatorSet.ValidatorsExt.TotalWeight)
			t.Log("UtimeSince:      ", validatorSet.ValidatorsExt.UtimeSince)
			t.Log("UtimeUntil:      ", validatorSet.ValidatorsExt.UtimeUntil)
			t.Log("Total:           ", validatorSet.ValidatorsExt.Total)
			t.Log("Main:            ", validatorSet.ValidatorsExt.Main)
			t.Log("Validators List: ")
			var sum uint64
			for i := range validatorSet.ValidatorsExt.List.Keys() {
				t.Log("Number:    ", i)
				t.Log("Key:       ", validatorSet.ValidatorsExt.List.Keys()[i].BinaryString())
				t.Log("SumType:   ", validatorSet.ValidatorsExt.List.Values()[i].SumType)
				if validatorSet.ValidatorsExt.List.Values()[i].SumType == "ValidatorAddr" {
					t.Log("PublicKey: ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.PublicKey.PubKey.Hex())
					t.Log("Weight:    ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.Weight)
					t.Log("AdnlAddr:  ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.AdnlAddr.Hex())
					sum += validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.Weight
				} else {
					t.Log("PublicKey: ", validatorSet.ValidatorsExt.List.Values()[i].Validator.PublicKey.PubKey.Hex())
					t.Log("Weight:    ", validatorSet.ValidatorsExt.List.Values()[i].Validator.Weight)
				}
				t.Log("--------------------------------------------------------")
			}
			t.Log(validatorSet.ValidatorsExt.TotalWeight)
			t.Log(sum)
		}
	}
}

func TestValidatorLoadExec(t *testing.T) {
	ctx := context.Background()
	tongoClient, err := liteclient.NewClientWithDefaultMainnet() //
	// tongoClient, err := liteclient.NewClientWithDefaultTestnet() //
	if err != nil {
		log.Fatalf("Unable to create tongo client: %v", err)
	}

	mcInfoExtra, err := tongoClient.GetMasterchainInfoExt(ctx, 0)
	if err != nil {
		log.Fatalf("Get account state error: %v", err)
	}
	lastBlockId := tongo.TonNodeBlockId{
		Workchain: mcInfoExtra.Workchain,
		Shard:     mcInfoExtra.Shard,
		Seqno:     mcInfoExtra.Seqno,
	}

	now := time.Now().Unix()
	_, header1, err := tongoClient.LookupBlock(ctx, 4, lastBlockId, 0, uint32(now-1000))
	if err != nil {
		log.Fatalf("LookupBlock 1 error: %v", err)
	}
	_, header2, err := tongoClient.LookupBlock(ctx, 4, lastBlockId, 0, uint32(now-10))
	if err != nil {
		log.Fatalf("LookupBlock 2 error: %v", err)
	}
	parents1, err := header1.GetParents()
	if err != nil {
		log.Fatalf("GetParents 1 error: %v", err)
	}
	parents2, err := header2.GetParents()
	if err != nil {
		log.Fatalf("GetParents 2 error: %v", err)
	}

	_, err = tongoClient.GetBlockProof(ctx, 0, parents1[0], nil) //&parents2[0])
	if err != nil {
		log.Fatalf("Get account state error: %v", err)
	}

	shardState1, err := tongoClient.GetConfigById(ctx, parents1[0])
	if err != nil {
		log.Fatalf("GetConfigById 1 error: %v", err)
	}
	shardState2, err := tongoClient.GetConfigById(ctx, parents2[0])
	if err != nil {
		log.Fatalf("GetConfigById 2 error: %v", err)
	}

	block1, err := tongoClient.GetBlock(ctx, parents1[0])
	if err != nil {
		log.Fatalf("GetBlock 1 error: %v", err)
	}
	block2, err := tongoClient.GetBlock(ctx, parents2[0])
	if err != nil {
		log.Fatalf("GetBlock 2 error: %v", err)
	}
	if block1.Extra.CreatedBy.Base64() == "" && block2.Extra.CreatedBy.Base64() == "" {
		log.Fatalf("SWW")
	}

	config1 := shardState1.UnsplitState.Value.ShardStateUnsplit.Custom.Value.Value.Config
	config2 := shardState2.UnsplitState.Value.ShardStateUnsplit.Custom.Value.Value.Config

	validatorStats1, err := tongoClient.ValidatorStats(ctx, 0, parents1[0], uint32(len(config1.Config.Hashmap.Keys())), nil, nil)
	if err != nil {
		log.Fatalf("ValidatorStats 1 error: %v", err)
	}
	validatorStats2, err := tongoClient.ValidatorStats(ctx, 0, parents2[0], uint32(len(config2.Config.Hashmap.Keys())), nil, nil)
	if err != nil {
		log.Fatalf("GetParents 1 error: %v", err)
	}
	if validatorStats1.ShardStateUnsplit.GlobalID != validatorStats2.ShardStateUnsplit.GlobalID {
		log.Fatalf("SWW")
	}
	t.Log("config 1 addr: ", config1.ConfigAddr.Hex())
	t.Log("config1 len: ", len(config1.Config.Hashmap.Keys()))

	t.Log("config 2 addr: ", config2.ConfigAddr.Hex())
	t.Log("config2 len: ", len(config2.Config.Hashmap.Keys()))

	for i := range config1.Config.Hashmap.Keys() {
		if binary.BigEndian.Uint32(config1.Config.Hashmap.Keys()[i].Buffer()) == 34 &&
			binary.BigEndian.Uint32(config2.Config.Hashmap.Keys()[i].Buffer()) == 34 {
			str := config1.Config.Hashmap.Values()[i].Value.RawBitString()
			fmt.Printf("key: %v, value: %x\n", config1.Config.Hashmap.Keys()[i].BinaryString(), str.Buffer())
			var validatorSet tongo.ValidatorsSet
			err := tlb.Unmarshal(&config1.Config.Hashmap.Values()[i].Value, &validatorSet)
			if err != nil {
				log.Fatalf("Unmarshal validator set error: %v", err)
			}
			t.Log("SumType:         ", validatorSet.SumType)
			t.Log("TotalWeight:     ", validatorSet.ValidatorsExt.TotalWeight)
			t.Log("UtimeSince:      ", validatorSet.ValidatorsExt.UtimeSince)
			t.Log("UtimeUntil:      ", validatorSet.ValidatorsExt.UtimeUntil)
			t.Log("Total:           ", validatorSet.ValidatorsExt.Total)
			t.Log("Main:            ", validatorSet.ValidatorsExt.Main)
			// t.Log("Validators List: ")
			// var sum uint64
			// for i := range validatorSet.ValidatorsExt.List.Keys() {
			// 	t.Log("Number:    ", i)
			// 	t.Log("Key:       ", validatorSet.ValidatorsExt.List.Keys()[i].BinaryString())
			// 	// t.Log("SumType:   ", validatorSet.ValidatorsExt.List.Values()[i].SumType)
			// 	// if validatorSet.ValidatorsExt.List.Values()[i].SumType == "ValidatorAddr" {
			// 	// 	t.Log("PublicKey: ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.PublicKey.SigPubKey.PubKey.Hex())
			// 	// 	t.Log("Weight:    ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.Weight)
			// 	// 	t.Log("AdnlAddr:  ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.AdnlAddr.Hex())
			// 	// 	sum += validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.Weight

			// 	// } else {
			// 	// 	t.Log("PublicKey: ", validatorSet.ValidatorsExt.List.Values()[i].Validator.PublicKey.SigPubKey.PubKey.Hex())
			// 	// 	t.Log("Weight:    ", validatorSet.ValidatorsExt.List.Values()[i].Validator.Weight)
			// 	// }

			// 	t.Log("--------------------------------------------------------")
			// }
			// t.Log(validatorSet.ValidatorsExt.TotalWeight)
			// t.Log(sum)
		}
	}
}
