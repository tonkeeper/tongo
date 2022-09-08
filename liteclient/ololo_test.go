package liteclient

import (
	"context"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	api, err := NewClient(nil)
	if err != nil {
		panic(err)
	}
	info, err := api.GetMasterchainInfo(context.Background())
	if err != nil {
		panic(err)
	}
	shards, err := api.BlocksGetShards(context.TODO(), info)
	fmt.Println(len(shards), err)
}
