package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/tl"
	"github.com/tonkeeper/tongo/tlb"
)

func main() {
	// Deserialize lite server response
	// TL schema:
	// liteServer.transactionList ids:(vector tonNode.blockIdExt) transactions:bytes = liteServer.TransactionList;
	// liteServer.error code:int message:string = liteServer.Error;
	// TL tag for SumType: crc32("liteServer.error code:int message:string = liteServer.Error") little endian
	type tonNodeBlockIdExt struct {
		Workchain int32
		Shard     int64
		Seqno     int32
		RootHash  tongo.Hash
		FileHash  tongo.Hash
	}
	var response struct {
		tlb.SumType
		TransactionList struct {
			Ids          []tonNodeBlockIdExt // corresponds to type TL vector
			Transactions []byte              // corresponds to type TL bytes
		} `tlSumType:"0bc6266f"`
		LiteServerError struct {
		} `tlSumType:"48e1a9bb"`
	}
	data, _ := hex.DecodeString("0bc6266f0200000000000000000000000000008091769201b73fd70ebac6405823096dc667b6e42b8cec1a92f6b140a0309c4cc95f54d21847b5033bf3fbc5b4580313caba9999d25dbd7188a048cee644f0fac0d08a29b300000000000000000000008091769201b73fd70ebac6405823096dc667b6e42b8cec1a92f6b140a0309c4cc95f54d21847b5033bf3fbc5b4580313caba9999d25dbd7188a048cee644f0fac0d08a29b3fef50300b5ee9c72010211020003e8010003b57e2d41ed396a9f1ba03839d63c5650fafc6fcfb574fd03f2e67d6555b61a3acd9000019fa692bda41c0c3b53376473958a45db5c122ee5ec7b5dbd6afea1b51be8d894e85377c79ce000019fa6621a78a629f85660003469da278802030403b57e2d41ed396a9f1ba03839d63c5650fafc6fcfb574fd03f2e67d6555b61a3acd9000019fa692bda4abf09a2c354f195ae4ccc0e0823a64be24a61ac5de7b117144d0a25cf0b7284ca000019fa692bda41629f85660001461e3e3080b0c0d0201e005060082728734d5189f3298fd505cec1e2ece027b7361f9a47a3d0221c0d556e023dc4d5f351fdfe28574ae29fc95c113a819fd8d4609355948efd527f8a97c4f1f883f9a020f0c4b061993cf0440090a01e18801c5a83da72d53e37407073ac78aca1f5f8df9f6ae9fa07e5ccfacaab6c34759b2047002a1c28bcd96496121965c341237afd99c4ec31f5098d0ededa5d307b2eba42dc340251441d16c3a069bfe20751c08b7ddc83067daecb7d9071a260cee68314d4d18bb14fc2cf800000020001c070101df08006a6200574ba8a53890bf135a88d761900ee44f40247050298824d3617c3fffb89e46eaa812a05f20000000000000000000000000000000b36801c5a83da72d53e37407073ac78aca1f5f8df9f6ae9fa07e5ccfacaab6c34759b3002ba5d4529c485f89ad446bb0c8077227a012382814c41269b0be1fffdc4f23755409502f900006145860000033f4d257b484c53f0acc40009d419d8313880000000000000000110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020006fc987a1204c14584000000000000200000000000357d3046a89bb7b0f886158490114f4ac5e3b1ff7fb5f6be1cde59b14d6c616044050164c0101a00e008272351fdfe28574ae29fc95c113a819fd8d4609355948efd527f8a97c4f1f883f9ab9fb184ccac0504709880bef9df17e51f44d62273b1886fe1f45d4268a31439f02150c090e8bfc2b5861e3e3110f1000c948008dbe435819ec7bfae0721aa85a4d01bc6414619b03a6faec7af93c2fe48234030038b507b4e5aa7c6e80e0e758f15943ebf1bf3ed5d3f40fcb99f59556d868eb3650e8bfc2b406145860000033f4d257b492c53f0acc6a993b6d800000000000000040009e407bec3b957000000000000000001d00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005bc00000000000000000000000012d452da449e50b8cf7dd27861f146122afe1b546bb8b70fc8216f0c614139f8e04000000")
	reader := bytes.NewReader(data)
	err := tl.Unmarshal(reader, &response)
	if err != nil {
		panic(err)
	}
	if response.SumType == "TransactionList" {
		fmt.Printf("TL deserialization: \n")
		fmt.Printf("response.TransactionList.Ids[0].Seqno: %v\n", response.TransactionList.Ids[0].Seqno)
		fmt.Printf("response.TransactionList.Ids[0].Shard: %v\n", response.TransactionList.Ids[0].Shard)
	} else {
		panic("invalid data")
	}

	// TL serialization:
	var X struct {
		A int32
		B []byte
	}
	X.A = 123
	X.B, _ = hex.DecodeString("AABBCCDD")
	b, err := tl.Marshal(X)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TL serialization: \n")
	fmt.Printf("Result: %x\n", b)
}
