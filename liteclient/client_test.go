package liteclient

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/tonkeeper/tongo/config"
)

func createTestLiteServerConnection() (*Connection, error) {
	base64Key := "wQE0MVhXNWUXpWiW5Bk8cAirIh5NNG3cZM1/fSVKIts="
	host := "135.181.140.221:46995"
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
	return NewConnection(context.Background(), pubkey, host)
}

func TestClient(t *testing.T) {
	c, err := createTestLiteServerConnection()
	if err != nil {
		t.Fatalf("NewConnection() failed: %v", err)
	}

	client := NewClient(c)
	b, _ := hex.DecodeString("df068c7978250e896bffffffff00000000000000801dcf3501db6f7082c2ea79e22a0ad7305022bd38300a3731f3dace87c2ba16b3582b89e5bb66984ce75d5f868bec2f0fb6bf3c6c3e0a506e3f15559a94ea58583251571e00000000195a6d69afdc08a11c38b729bfad789b957a50650f9bb84b1947fcc0acc9a41b000000")
	resp, err := client.Request(context.Background(), b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", resp)
}

func TestGeneratedMethod(t *testing.T) {
	c, err := createTestLiteServerConnection()
	if err != nil {
		t.Fatalf("NewConnection() failed: %v", err)
	}

	client := NewClient(c)
	resp, err := client.LiteServerGetMasterchainInfo(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Last seqno: %d\n", resp.Last.Seqno)
}

func TestGeneratedMethod2(t *testing.T) {
	c, err := createTestLiteServerConnection()
	if err != nil {
		t.Fatalf("NewConnection() failed: %v", err)
	}

	client := NewClient(c)

	r, err := client.LiteServerGetMasterchainInfo(context.Background())
	if err != nil {
		panic(err)
	}

	req := LiteServerGetBlockRequest{Id: r.Last}

	resp, err := client.LiteServerGetBlock(context.Background(), req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Req  seqno: %d\n", req.Id.Seqno)
	fmt.Printf("Resp seqno: %d\n", resp.Id.Seqno)
}

func TestGeneratedMethod3(t *testing.T) {
	c, err := createTestLiteServerConnection()
	if err != nil {
		t.Fatalf("NewConnection() failed: %v", err)
	}

	client := NewClient(c)

	r, err := client.LiteServerGetMasterchainInfo(context.Background())
	if err != nil {
		panic(err)
	}

	req := LiteServerGetValidatorStatsRequest{
		Mode:  0,
		Id:    r.Last,
		Limit: 10,
	}

	resp, err := client.LiteServerGetValidatorStats(context.Background(), req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Count: %d\n", resp.Count)
	fmt.Printf("Complete: %v\n", resp.Complete)
}

func TestGeneratedMethod4(t *testing.T) {
	c, err := createTestLiteServerConnection()
	if err != nil {
		t.Fatalf("NewConnection() failed: %v", err)
	}

	client := NewClient(c)

	r, err := client.LiteServerGetTime(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current time: %d\n", r)
}

func TestGeneratedMethod5(t *testing.T) {
	c, err := createTestLiteServerConnection()
	if err != nil {
		t.Fatalf("NewConnection() failed: %v", err)
	}

	client := NewClient(c)
	r, err := client.LiteServerGetMasterchainInfo(context.Background())
	if err != nil {
		panic(err)
	}

	req := LiteServerLookupBlockRequest{
		Mode: 1,
		Id: TonNodeBlockIdC{
			r.Last.Workchain,
			r.Last.Shard,
			r.Last.Seqno - 2,
		},
		Lt:    nil,
		Utime: nil,
	}

	r1, err := client.LiteServerLookupBlock(context.Background(), req)

	req1 := LiteServerGetBlockProofRequest{
		Mode:        0,
		KnownBlock:  r1.Id,
		TargetBlock: nil,
	}

	r2, err := client.LiteServerGetBlockProof(context.Background(), req1)
	if err != nil {
		panic(err)
	}
	_ = r2
}

func TestClient_WaitMasterchainSeqno(t *testing.T) {
	c, err := createTestLiteServerConnection()
	if err != nil {
		t.Fatalf("NewConnection() failed: %v", err)
	}
	client := NewClient(c)
	resp, err := client.LiteServerGetMasterchainInfo(context.Background())
	if err != nil {
		t.Fatalf("LiteServerGetMasterchainInfo() failed: %v", err)
	}
	seqno := resp.Last.Seqno + 2
	if err := client.WaitMasterchainSeqno(context.Background(), seqno, 15000); err != nil {
		t.Fatalf("WaitMasterchainSeqno() failed: %v", err)
	}
	resp, err = client.LiteServerGetMasterchainInfo(context.Background())
	if err != nil {
		t.Fatalf("LiteServerGetMasterchainInfo() failed: %v", err)
	}
	if resp.Last.Seqno != seqno {
		t.Fatalf("want seqno: %v, got: %v", seqno, resp.Last.Seqno)
	}
}
