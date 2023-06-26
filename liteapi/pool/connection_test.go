package pool

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/tonkeeper/tongo/liteclient"
)

func Test_connection_Run(t *testing.T) {
	pubkey, err := base64.StdEncoding.DecodeString("wQE0MVhXNWUXpWiW5Bk8cAirIh5NNG3cZM1/fSVKIts=")
	if err != nil {
		t.Fatalf("DecodeString() failed: %v", err)
	}
	c, err := liteclient.NewConnection(context.Background(), pubkey, "135.181.140.221:46995")
	if err != nil {
		t.Fatalf("NewConnection() failed: %v", err)
	}
	conn := &connection{
		client: liteclient.NewClient(c),
	}
	go conn.Run(context.Background())

	time.Sleep(1 * time.Second)
	res, err := conn.Client().LiteServerGetMasterchainInfo(context.Background())
	if err != nil {
		t.Fatalf("LiteServerGetMasterchainInfo() failed: %v", err)
	}
	if res.Last.Seqno != conn.MasterSeqno() {
		t.Fatalf("want seqno: %v, got: %v", res.Last.Seqno, conn.MasterSeqno())
	}
	if err := conn.Client().WaitMasterchainSeqno(context.Background(), res.Last.Seqno+1, 15_000); err != nil {
		t.Fatalf("WaitMasterchainSeqno() failed: %v", err)
	}
	// give a few milliseconds to the connection's goroutine
	time.Sleep(150 * time.Millisecond)

	if res.Last.Seqno+1 != conn.MasterSeqno() {
		t.Fatalf("want seqno: %v, got: %v", res.Last.Seqno, conn.MasterSeqno())
	}
}
