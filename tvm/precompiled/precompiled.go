package precompiled

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type MethodCode struct {
	MethodID int
	CodeHash [32]byte
}

// todo: exit code
type tvmPrecompiled func(data *boc.Cell, args tlb.VmStack) (tlb.VmStack, error)

var KnownMethods = map[MethodCode]tvmPrecompiled{
	//get_pow_params gram miner
	{MethodID: 101616, CodeHash: ton.MustParseHash("ccae6ffb603c7d3e779ab59ec267ffc22dc1ebe0af9839902289a7a83e4c00f1")}: getPowParamsGram,

	//wallet v3r2 seqno
	{MethodID: 85143, CodeHash: ton.MustParseHash("84dafa449f98a6987789ba232358072bc0f76dc4524002a5d0918b9a75d2d599")}: walletv3r2seqno,
	{MethodID: 78748, CodeHash: ton.MustParseHash("84dafa449f98a6987789ba232358072bc0f76dc4524002a5d0918b9a75d2d599")}: walletv3r2publicKey,

	//wallet v4r2 seqno
	{MethodID: 85143, CodeHash: ton.MustParseHash("feb5ff6820e2ff0d9483e7e0d62c817d846789fb4ae580c878866d959dabd5c0")}:  walletv4r2seqno,
	{MethodID: 81467, CodeHash: ton.MustParseHash("feb5ff6820e2ff0d9483e7e0d62c817d846789fb4ae580c878866d959dabd5c0")}:  walletv4r2SubwalletID,
	{MethodID: 78748, CodeHash: ton.MustParseHash("feb5ff6820e2ff0d9483e7e0d62c817d846789fb4ae580c878866d959dabd5c0")}:  walletv4r2publicKey,
	{MethodID: 107653, CodeHash: ton.MustParseHash("feb5ff6820e2ff0d9483e7e0d62c817d846789fb4ae580c878866d959dabd5c0")}: walletv4r2getPluginList,

	//jetton v1 get_wallet_data
	{MethodID: 97026, CodeHash: ton.MustParseHash("beb0683ebeb8927fe9fc8ec0a18bc7dd17899689825a121eab46c5a3a860d0ce")}: jettonV1getWalletData,

	//nft_item simple
	{MethodID: 102351, CodeHash: ton.MustParseHash("4c9123828682fa6f43797ab41732bca890cae01766e0674100250516e0bf8d42")}: nftV1getNftData,
}
