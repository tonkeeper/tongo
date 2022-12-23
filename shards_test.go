package tongo

import (
	"testing"
)

func TestParseShardID(t *testing.T) {
	shard := uint64(0xfb80000000000000)
	shardID, err := ParseShardID(int64(shard))
	if err != nil {
		panic(err)
	}
	t.Logf("Shard  : %b\n", shard)
	t.Logf("Mask   : %b\n", uint64(shardID.mask))
	t.Logf("Prefix : %b\n", uint64(shardID.prefix))
	if int64(shard) != shardID.Encode() {
		t.Logf("%b\n", shard)
		t.Logf("%b\n", uint64(shardID.Encode()))
		t.Fatal(shard)
	}
}

func TestShardID_MatchAccountID(t *testing.T) {
	for _, c := range []struct {
		account string
		shard   uint64
		match   bool
	}{
		{"0:FB65906EB4EC4803550C0842105667E7270B6C22C32A4FB7D2B3C49B96C15773", 0xfb80000000000000, true},
		{"0:FE65906EB4EC4803550C0842105667E7270B6C22C32A4FB7D2B3C49B96C15773", 0xfb80000000000000, false},
		{"0:FE65906EB4EC4803550C0842105667E7270B6C22C32A4FB7D2B3C49B96C15773", 0x8000000000000000, true},
	} {
		a := MustParseAccountID(c.account)
		shardID := MustParseShardID(int64(c.shard))
		if shardID.MatchAccountID(a) != c.match {
			t.Errorf("%v %x", c.account, c.shard)
		}
	}
}

func TestShardID_MatchBlockID(t *testing.T) {
	for _, c := range []struct {
		shard uint64
		block uint64
		match bool
	}{
		{shard: 0xfb80000000000000, block: 0xfb80000000000000, match: true},  // equal prefixes
		{shard: 0xfb80000000000000, block: 0x8000000000000000, match: true},  // only one real shard. node tracks it
		{shard: 0xfb80000000000000, block: 0xc000000000000000, match: true},  // two real shards. node tracks this shard
		{shard: 0xfb80000000000000, block: 0x4000000000000000, match: false}, // two real shards. node not tracks this shard
		{shard: 0x8000000000000000, block: 0xfb80000000000000, match: true},  // node tracks all shards
		{shard: 0xfc80000000000000, block: 0xfb80000000000000, match: false}, // not equal prefixes
		{shard: 0xfb80000000000000, block: 0xffc0000000000000, match: false},
		{shard: 0xfb80000000000000, block: 0xff00000000000000, match: false},
	} {
		shard := MustParseShardID(int64(c.shard))
		t.Logf("Shard       : %b\n", c.shard)
		t.Logf("Block shard : %b\n\n", c.block)
		if shard.MatchBlockID(BlockID{Shard: c.block}) != c.match {
			t.Errorf("shard %v block %v", c.shard, c.block)
		}
	}
}
