package pool

import (
	"context"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteclient"
)

func createTestLiteServerConnection() (*liteclient.Connection, error) {
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
	return liteclient.NewConnection(context.Background(), pubkey, host)
}

func Test_connection_Run(t *testing.T) {
	c, err := createTestLiteServerConnection()
	if err != nil {
		t.Fatalf("NewConnection() failed: %v", err)
	}
	conn := &connection{
		client:              liteclient.NewClient(c),
		masterHeadUpdatedCh: make(chan masterHeadUpdated, 100),
	}
	go conn.Run(context.Background())

	time.Sleep(1 * time.Second)
	res, err := conn.Client().LiteServerGetMasterchainInfo(context.Background())
	if err != nil {
		t.Fatalf("LiteServerGetMasterchainInfo() failed: %v", err)
	}
	masterHead := conn.MasterHead()
	if res.Last.Seqno > masterHead.Seqno {
		t.Fatalf("want seqno: %v, got: %v", res.Last.Seqno, masterHead.Seqno)
	}
	if err := conn.Client().WaitMasterchainSeqno(context.Background(), masterHead.Seqno+1, 15_000); err != nil {
		t.Fatalf("WaitMasterchainSeqno() failed: %v", err)
	}
	// give a few milliseconds to the connection's goroutine
	time.Sleep(1 * time.Second)

	newMasterHead := conn.MasterHead()
	if masterHead.Seqno+1 != newMasterHead.Seqno {
		t.Fatalf("want seqno: %v, got: %v", res.Last.Seqno, newMasterHead.Seqno)
	}
}
