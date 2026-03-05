package liteapi

import (
	"context"
	"testing"
	"time"

	"github.com/tonkeeper/tongo/liteapi/pool"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/ton"
)

type stubConnPool struct {
	name          string
	status        pool.Status
	masterInfoCli *pool.MasterchainInfoClient
	bestClient    *liteclient.Client
	bestHead      ton.BlockIDExt
	err           error
}

func (s *stubConnPool) BestMasterchainInfoClient() *pool.MasterchainInfoClient {
	return s.masterInfoCli
}

func (s *stubConnPool) BestMasterchainClient(ctx context.Context) (*liteclient.Client, ton.BlockIDExt, error) {
	return s.bestClient, s.bestHead, s.err
}

func (s *stubConnPool) BestClientByAccountID(ctx context.Context, accountID ton.AccountID, archiveRequired bool) (*liteclient.Client, ton.BlockIDExt, error) {
	return s.bestClient, s.bestHead, s.err
}

func (s *stubConnPool) BestClientByBlockID(ctx context.Context, blockID ton.BlockID) (*liteclient.Client, error) {
	return s.bestClient, s.err
}

func (s *stubConnPool) WaitMasterchainSeqno(ctx context.Context, seqno uint32, timeout time.Duration) error {
	return s.err
}

func (s *stubConnPool) ConnectionsNumber() int {
	return len(s.Status().Connections)
}

func (s *stubConnPool) Status() pool.Status {
	return s.status
}

func TestFallbackPool_PrefersPrimaryWhenConnected(t *testing.T) {
	primary := &stubConnPool{
		name: "primary",
		status: pool.Status{
			Connections: []pool.ConnStatus{{Connected: true}},
		},
	}
	fallback := &stubConnPool{
		name: "fallback",
		status: pool.Status{
			Connections: []pool.ConnStatus{{Connected: true}},
		},
	}

	fp := &connPoolProxy{
		primary:  primary,
		fallback: fallback,
	}

	if !fp.usePrimary() {
		t.Fatalf("expected primary to be used when it has connected nodes")
	}

	if got := fp.Status(); &got == &fallback.status {
		t.Fatalf("expected Status to reflect primary when it is healthy")
	}
}

func TestFallbackPool_UsesFallbackWhenPrimaryDown(t *testing.T) {
	primary := &stubConnPool{
		name: "primary",
		status: pool.Status{
			Connections: []pool.ConnStatus{{Connected: false}},
		},
		err: pool.ErrNoConnections,
	}
	fallback := &stubConnPool{
		name: "fallback",
		status: pool.Status{
			Connections: []pool.ConnStatus{{Connected: true}},
		},
	}

	fp := &connPoolProxy{
		primary:  primary,
		fallback: fallback,
	}

	if fp.usePrimary() {
		t.Fatalf("did not expect primary to be used when all primary connections are down")
	}

	if err := fp.WaitMasterchainSeqno(context.Background(), 1, time.Second); err != nil {
		if err == pool.ErrNoConnections {
			t.Fatalf("expected fallback to handle WaitMasterchainSeqno, got ErrNoConnections from primary")
		}
	}
}
