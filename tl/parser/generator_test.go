package parser

import (
	"fmt"
	"testing"
)

var SOURCE = `
tonNode.blockId workchain:int shard:long seqno:int = tonNode.BlockId;
tonNode.blockIdExt workchain:int shard:long seqno:int root_hash:int256 file_hash:int256 = tonNode.BlockIdExt;
tonNode.zeroStateIdExt workchain:int root_hash:int256 file_hash:int256 = tonNode.ZeroStateIdExt;

adnl.message.query#7af98bb4 query_id:int256 query:bytes = adnl.Message;
adnl.message.answer#1684ac0f query_id:int256 answer:bytes = adnl.Message;

liteServer.error code:int message:string = liteServer.Error; 

liteServer.accountId workchain:int id:int256 = liteServer.AccountId;
liteServer.libraryEntry hash:int256 data:bytes = liteServer.LibraryEntry;

liteServer.masterchainInfo last:tonNode.blockIdExt state_root_hash:int256 init:tonNode.zeroStateIdExt = liteServer.MasterchainInfo;
liteServer.masterchainInfoExt mode:# version:int capabilities:long last:tonNode.blockIdExt last_utime:int now:int state_root_hash:int256 init:tonNode.zeroStateIdExt = liteServer.MasterchainInfoExt;
liteServer.currentTime now:int = liteServer.CurrentTime;
liteServer.version mode:# version:int capabilities:long now:int = liteServer.Version;
liteServer.blockData id:tonNode.blockIdExt data:bytes = liteServer.BlockData;
liteServer.blockState id:tonNode.blockIdExt root_hash:int256 file_hash:int256 data:bytes = liteServer.BlockState;
liteServer.blockHeader id:tonNode.blockIdExt mode:# header_proof:bytes = liteServer.BlockHeader;
liteServer.sendMsgStatus status:int = liteServer.SendMsgStatus;
liteServer.accountState id:tonNode.blockIdExt shardblk:tonNode.blockIdExt shard_proof:bytes proof:bytes state:bytes = liteServer.AccountState;
liteServer.runMethodResult mode:# id:tonNode.blockIdExt shardblk:tonNode.blockIdExt shard_proof:mode.0?bytes proof:mode.0?bytes state_proof:mode.1?bytes init_c7:mode.3?bytes lib_extras:mode.4?bytes exit_code:int result:mode.2?bytes = liteServer.RunMethodResult;
liteServer.shardInfo id:tonNode.blockIdExt shardblk:tonNode.blockIdExt shard_proof:bytes shard_descr:bytes = liteServer.ShardInfo;
liteServer.allShardsInfo id:tonNode.blockIdExt proof:bytes data:bytes = liteServer.AllShardsInfo;
liteServer.transactionInfo id:tonNode.blockIdExt proof:bytes transaction:bytes = liteServer.TransactionInfo;
liteServer.transactionList ids:(vector tonNode.blockIdExt) transactions:bytes = liteServer.TransactionList;
liteServer.transactionId mode:# account:mode.0?int256 lt:mode.1?long hash:mode.2?int256 = liteServer.TransactionId;
liteServer.transactionId3 account:int256 lt:long = liteServer.TransactionId3;
liteServer.blockTransactions id:tonNode.blockIdExt req_count:# incomplete:Bool ids:(vector liteServer.transactionId) proof:bytes = liteServer.BlockTransactions;
liteServer.signature node_id_short:int256 signature:bytes = liteServer.Signature;
liteServer.signatureSet validator_set_hash:int catchain_seqno:int signatures:(vector liteServer.signature) = liteServer.SignatureSet;
liteServer.blockLinkBack#ef1b7eef to_key_block:Bool from:tonNode.blockIdExt to:tonNode.blockIdExt dest_proof:bytes proof:bytes state_proof:bytes = liteServer.BlockLink;
liteServer.blockLinkForward#1cce0f52 to_key_block:Bool from:tonNode.blockIdExt to:tonNode.blockIdExt dest_proof:bytes config_proof:bytes signatures:liteServer.SignatureSet = liteServer.BlockLink;
liteServer.partialBlockProof complete:Bool from:tonNode.blockIdExt to:tonNode.blockIdExt steps:(vector liteServer.BlockLink) = liteServer.PartialBlockProof;
liteServer.configInfo mode:# id:tonNode.blockIdExt state_proof:bytes config_proof:bytes = liteServer.ConfigInfo;
liteServer.validatorStats mode:# id:tonNode.blockIdExt count:int complete:Bool state_proof:bytes data_proof:bytes = liteServer.ValidatorStats;
liteServer.libraryResult result:(vector liteServer.libraryEntry) = liteServer.LibraryResult;
liteServer.shardBlockLink id:tonNode.blockIdExt proof:bytes = liteServer.ShardBlockLink;
liteServer.shardBlockProof masterchain_id:tonNode.blockIdExt links:(vector liteServer.shardBlockLink) = liteServer.ShardBlockProof;

liteServer.debug.verbosity value:int = liteServer.debug.Verbosity;

---functions---

liteServer.getMasterchainInfo = liteServer.MasterchainInfo;
liteServer.getMasterchainInfoExt mode:# = liteServer.MasterchainInfoExt;
`

func TestGenerateGolangTypes(t *testing.T) {
	parsed, err := Parse(SOURCE)
	if err != nil {
		panic(err)
	}
	s, err := GenerateGolangTypes(*parsed)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", s)
}

func TestCheckBits(t *testing.T) {
	mode := 53
	fmt.Printf("Mode: %b\n", mode)
	for i := 0; i < 6; i++ {
		fmt.Printf("%v bit: %v\n", i, (mode>>i)&1 == 1)
	}
}
