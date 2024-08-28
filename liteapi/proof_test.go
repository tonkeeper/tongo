package liteapi

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo/ton"
	"testing"
)

func getInitBlock() (*ton.BlockIDExt, error) {
	var rootHash ton.Bits256
	err := rootHash.FromBase64("VpWyfNOLm8Rqt6CZZ9dZGqJRO3NyrlHHYN1k1oLbJ6g=")
	if err != nil {
		return nil, fmt.Errorf("incorrect root hash")
	}
	var fileHash ton.Bits256
	err = fileHash.FromBase64("8o12KX54BtJM8RERD1J97Qe1ZWk61LIIyXydlBnixK8=")
	if err != nil {
		return nil, fmt.Errorf("incorrect file hash")
	}
	return &ton.BlockIDExt{
		BlockID: ton.BlockID{
			Workchain: -1,
			Shard:     9223372036854775808,
			Seqno:     34835953,
		},
		RootHash: rootHash,
		FileHash: fileHash,
	}, nil
}

func getLastBlockInMasterchain(c *Client) (*ton.BlockIDExt, error) {
	lst, err := c.GetMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}
	blk := lst.Last.ToBlockIdExt()
	return &blk, nil
}

func TestVerifyProofChain(t *testing.T) {
	c, err := NewClientWithDefaultMainnet()
	if err != nil {
		t.Fatalf("unable to create liteclient: %v", err)
	}
	from, err := getInitBlock()
	if err != nil {
		t.Fatalf("unable to get init block: %v", err)
	}
	to, err := getLastBlockInMasterchain(c)
	if err != nil {
		t.Fatalf("unable to get last block: %v", err)
	}
	type Test struct {
		name string
		from *ton.BlockIDExt
		to   *ton.BlockIDExt
	}
	tests := []Test{
		{
			name: "test verify forward proof chain",
			from: from,
			to:   to,
		},
		{
			name: "test verify backward proof chain",
			from: to,
			to:   from,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = c.VerifyProofChain(context.Background(), *test.from, *test.to)
			if err != nil {
				t.Errorf("proof chain failed from %v, to %v: %v", test.from.Seqno, test.to.Seqno, err)
			}
		})
	}
}
