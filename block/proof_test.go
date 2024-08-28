package block

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/ton"
	"os"
	"testing"
)

func createTestLiteServerConnection() (*liteclient.Connection, error) {
	base64Key := ""
	host := ""
	if serversEnv, ok := os.LookupEnv("LITE_SERVERS"); ok {
		servers, err := config.ParseLiteServersEnvVar(serversEnv)
		if err != nil {
			return nil, err
		}
		if len(servers) > 0 {
			base64Key = servers[0].Key
			host = servers[0].Host
		}
	}
	pubkey, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}
	return liteclient.NewConnection(context.Background(), pubkey, host)
}

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

func Test(t *testing.T) {
	c, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		fmt.Printf("Unable to create tongo client: %v", err)
	}

	from, err := getInitBlock()
	if err != nil {
		t.Fatal("cannot get init block: %w", err)
	}

	a, err := createTestLiteServerConnection()
	if err != nil {
		t.Fatal("cannot get")
	}
	lst, err := c.GetMasterchainInfo(context.Background())
	to := blockIdExtMapper(lst.Last)

	// make backward
	//from, to = to, from
	//err = VerifyProofChain(context.Background(), liteclient.NewClient(a), from, to)

	lst, err = c.GetMasterchainInfo(context.Background())
	to = blockIdExtMapper(lst.Last)
	err = VerifyProofChain(context.Background(), liteclient.NewClient(a), from, to)
}
