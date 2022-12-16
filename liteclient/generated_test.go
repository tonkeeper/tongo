package liteclient

import (
	"bytes"
	"github.com/startfellows/tongo/tl"
	"reflect"
	"testing"
)

var blk = TonNodeBlockIdExt{
	Workchain: 0,
	Shard:     123,
	Seqno:     321,
	RootHash:  tl.Int256([32]byte{1, 2, 3}),
	FileHash:  tl.Int256([32]byte{3, 2, 1}),
}

func TestSimpleType(t *testing.T) {
	a := TonNodeBlockIdExt{
		Workchain: 0,
		Shard:     123,
		Seqno:     321,
		RootHash:  tl.Int256([32]byte{1, 2, 3}),
		FileHash:  tl.Int256([32]byte{3, 2, 1}),
	}
	b, err := tl.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var a1 TonNodeBlockIdExt
	err = tl.Unmarshal(bytes.NewReader(b), &a1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, a1) {
		t.Fatal("not equal")
	}
}

func TestSimpleTypeWithMode(t *testing.T) {
	account := tl.Int256([32]byte{3, 2, 1})
	var lt int64 = 123
	a := LiteServerTransactionId{
		Mode:    3,
		Account: &account,
		Lt:      &lt,
		Hash:    nil,
	}
	b, err := tl.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var a1 LiteServerTransactionId
	err = tl.Unmarshal(bytes.NewReader(b), &a1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, a1) {
		t.Fatal("not equal")
	}
}

func TestSimpleTypeInvalidMode(t *testing.T) {
	account := tl.Int256([32]byte{3, 2, 1})
	var lt int64 = 123
	a := LiteServerTransactionId{
		Mode:    2,
		Account: &account,
		Lt:      &lt,
		Hash:    nil,
	}
	b, err := tl.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var a1 LiteServerTransactionId
	err = tl.Unmarshal(bytes.NewReader(b), &a1)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(a, a1) {
		t.Fatal("must be not equal")
	}
}

func TestSimpleTypeWithModeAndSlice(t *testing.T) {
	// liteServer.runMethodResult mode:# id:tonNode.blockIdExt shardblk:tonNode.blockIdExt shard_proof:mode.0?bytes
	// proof:mode.0?bytes state_proof:mode.1?bytes init_c7:mode.3?bytes lib_extras:mode.4?bytes exit_code:int
	// result:mode.2?bytes = liteServer.RunMethodResult;
	a := LiteServerRunMethodResult{
		Mode:       5, // 101
		Id:         blk,
		Shardblk:   blk,
		ShardProof: []byte{1, 2, 3}, // mode.0 == 1
		Proof:      []byte{4, 5, 6}, // mode.0 == 1
		//StateProof: []byte, // mode.1 == 0
		//InitC7:     []byte, // mode.3 == 0
		//LibExtras:  []byte, // mode.4 == 0
		ExitCode: 1,
		Result:   []byte{7, 8, 9}, // mode.2 == 1
	}
	b, err := tl.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var a1 LiteServerRunMethodResult
	err = tl.Unmarshal(bytes.NewReader(b), &a1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, a1) {
		t.Fatal("not equal")
	}
}

func TestSumType(t *testing.T) {
	a := AdnlMessage{
		SumType: "AdnlMessageAnswer",
	}
	a.AdnlMessageAnswer.Answer = []byte{1, 2, 3}
	a.AdnlMessageAnswer.QueryId = [32]byte{4, 5, 6}

	b, err := tl.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var a1 AdnlMessage
	err = tl.Unmarshal(bytes.NewReader(b), &a1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, a1) {
		t.Fatal("not equal")
	}
}

func TestSimpleTypeWithVector(t *testing.T) {
	account := tl.Int256([32]byte{3, 2, 1})
	var lt int64 = 123
	id1 := LiteServerTransactionId{
		Mode:    3,
		Account: &account,
		Lt:      &lt,
		Hash:    nil,
	}
	id2 := LiteServerTransactionId{
		Mode:    1,
		Account: &account,
		Lt:      nil,
		Hash:    nil,
	}
	a := LiteServerBlockTransactions{
		Id:         blk,
		ReqCount:   123,
		Incomplete: true,
		Ids:        []LiteServerTransactionId{id1, id2},
		Proof:      []byte{1, 2, 3},
	}
	b, err := tl.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var a1 LiteServerBlockTransactions
	err = tl.Unmarshal(bytes.NewReader(b), &a1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, a1) {
		t.Fatal("not equal")
	}
}
