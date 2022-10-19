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

	w, err := wallet.NewWallet(privateKey, wallet.V4R2, 0, nil)
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

	tonTransfer := wallet.TonTransfer{
		Recipient: *recipientAddr,
		Amount:    10000,
		Comment:   "hello",
		Bounce:    false,
		Mode:      1,
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

	msg, err := w.GenerateTonTransferMessage(uint32(res.Stack[0].Int64()), 0xFFFFFFFF, []wallet.TonTransfer{tonTransfer})
	if err != nil {
		log.Fatalf("Unable to generate transfer message: %v", err)
	}

	var message tongo.Message[tlb.Any]
	c, err := boc.DeserializeBoc(msg)
	if err != nil {
		log.Fatalf("unable to deserialize transfer message: %v", err)
	}
	err = tlb.Unmarshal(c[0], &message)
	if err != nil {
		log.Fatalf("unable to unmarshal transfer message: %v", err)
	}

	var shardAccount tongo.ShardAccount
	shardAccount.Account = account
	shardAccount.LastTransLt = account.Account.Storage.LastTransLt - 1

	e, err := NewEmulator(config)
	if err != nil {
		log.Fatalf("unable to create emulator: %v", err)
	}

	acc, tx, err := e.Emulate(shardAccount, message)
	if err != nil {
		log.Fatalf("emulator error: %v", err)
	}
	fmt.Printf("Account last transaction hash: %x\n", acc.LastTransHash)
	fmt.Printf("Prev transaction hash: %x\n", tx.PrevTransHash)

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
	fmt.Println("config addr: ", config.ConfigAddr.Hex())
	for i := range config.Config.Hashmap.Keys() {
		if binary.BigEndian.Uint32(config.Config.Hashmap.Keys()[i].Buffer()) == 34 {
			str := config.Config.Hashmap.Values()[i].Value.RawBitString()
			fmt.Printf("key: %v, value: %x\n", config.Config.Hashmap.Keys()[i].BinaryString(), str.Buffer())
			var validatorSet tongo.ValidatorsSet
			err := tlb.Unmarshal(&config.Config.Hashmap.Values()[i].Value, &validatorSet)
			if err != nil {
				log.Fatalf("Unmarshal validator set error: %v", err)
			}
			fmt.Println("SumType:         ", validatorSet.SumType)
			fmt.Println("TotalWeight:     ", validatorSet.ValidatorsExt.TotalWeight)
			fmt.Println("UtimeSince:      ", validatorSet.ValidatorsExt.UtimeSince)
			fmt.Println("UtimeUntil:      ", validatorSet.ValidatorsExt.UtimeUntil)
			fmt.Println("Total:           ", validatorSet.ValidatorsExt.Total)
			fmt.Println("Main:            ", validatorSet.ValidatorsExt.Main)
			fmt.Println("Validators List: ")
			var sum uint64
			for i := range validatorSet.ValidatorsExt.List.Keys() {
				fmt.Println("Number:    ", i)
				fmt.Println("Key:       ", validatorSet.ValidatorsExt.List.Keys()[i].BinaryString())
				fmt.Println("SumType:   ", validatorSet.ValidatorsExt.List.Values()[i].SumType)
				if validatorSet.ValidatorsExt.List.Values()[i].SumType == "ValidatorAddr" {
					fmt.Println("PublicKey: ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.PublicKey.PubKey.Hex())
					fmt.Println("Weight:    ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.Weight)
					fmt.Println("AdnlAddr:  ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.AdnlAddr.Hex())
					sum += validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.Weight
				} else {
					fmt.Println("PublicKey: ", validatorSet.ValidatorsExt.List.Values()[i].Validator.PublicKey.PubKey.Hex())
					fmt.Println("Weight:    ", validatorSet.ValidatorsExt.List.Values()[i].Validator.Weight)
				}
				fmt.Println("--------------------------------------------------------")
			}
			fmt.Println(validatorSet.ValidatorsExt.TotalWeight)
			fmt.Println(sum)
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
	header1, err := tongoClient.LookupBlock(ctx, 4, lastBlockId, 0, uint32(now-1000))
	if err != nil {
		log.Fatalf("LookupBlock 1 error: %v", err)
	}
	header2, err := tongoClient.LookupBlock(ctx, 4, lastBlockId, 0, uint32(now-10))
	if err != nil {
		log.Fatalf("LookupBlock 2 error: %v", err)
	}
	parents1, err := header1.Info.GetParents()
	if err != nil {
		log.Fatalf("GetParents 1 error: %v", err)
	}
	parents2, err := header2.Info.GetParents()
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
	fmt.Println("config 1 addr: ", config1.ConfigAddr.Hex())
	fmt.Println("config1 len: ", len(config1.Config.Hashmap.Keys()))

	fmt.Println("config 2 addr: ", config2.ConfigAddr.Hex())
	fmt.Println("config2 len: ", len(config2.Config.Hashmap.Keys()))

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
			fmt.Println("SumType:         ", validatorSet.SumType)
			fmt.Println("TotalWeight:     ", validatorSet.ValidatorsExt.TotalWeight)
			fmt.Println("UtimeSince:      ", validatorSet.ValidatorsExt.UtimeSince)
			fmt.Println("UtimeUntil:      ", validatorSet.ValidatorsExt.UtimeUntil)
			fmt.Println("Total:           ", validatorSet.ValidatorsExt.Total)
			fmt.Println("Main:            ", validatorSet.ValidatorsExt.Main)
			// fmt.Println("Validators List: ")
			// var sum uint64
			// for i := range validatorSet.ValidatorsExt.List.Keys() {
			// 	fmt.Println("Number:    ", i)
			// 	fmt.Println("Key:       ", validatorSet.ValidatorsExt.List.Keys()[i].BinaryString())
			// 	// fmt.Println("SumType:   ", validatorSet.ValidatorsExt.List.Values()[i].SumType)
			// 	// if validatorSet.ValidatorsExt.List.Values()[i].SumType == "ValidatorAddr" {
			// 	// 	fmt.Println("PublicKey: ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.PublicKey.SigPubKey.PubKey.Hex())
			// 	// 	fmt.Println("Weight:    ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.Weight)
			// 	// 	fmt.Println("AdnlAddr:  ", validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.AdnlAddr.Hex())
			// 	// 	sum += validatorSet.ValidatorsExt.List.Values()[i].ValidatorAddr.Weight

			// 	// } else {
			// 	// 	fmt.Println("PublicKey: ", validatorSet.ValidatorsExt.List.Values()[i].Validator.PublicKey.SigPubKey.PubKey.Hex())
			// 	// 	fmt.Println("Weight:    ", validatorSet.ValidatorsExt.List.Values()[i].Validator.Weight)
			// 	// }

			// 	fmt.Println("--------------------------------------------------------")
			// }
			// fmt.Println(validatorSet.ValidatorsExt.TotalWeight)
			// fmt.Println(sum)
		}
	}
}
