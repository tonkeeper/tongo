package pool

import (
	"context"
	"encoding/base64"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/ton"
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
		t.Skipf("cannot connect to lite server: %v", err)
	}
	conn := &connection{
		client:              liteclient.NewClient(c),
		masterHeadUpdatedCh: make(chan masterHeadUpdated, 100),
	}
	go conn.Run(context.Background(), false)

	time.Sleep(1 * time.Second)
	res, err := conn.Client().LiteServerGetMasterchainInfo(context.Background())
	if err != nil {
		t.Fatalf("failed to get masterchain info: %v", err)
	}
	masterHead := conn.MasterHead()
	if res.Last.Seqno > masterHead.Seqno {
		t.Fatalf("expected seqno >= %d, got %d", res.Last.Seqno, masterHead.Seqno)
	}
	err = conn.Client().WaitMasterchainSeqno(context.Background(), masterHead.Seqno+1, 15_000)
	if err != nil {
		t.Fatalf("failed to wait for next seqno: %v", err)
	}
	// give a few milliseconds to the connection's goroutine
	time.Sleep(1 * time.Second)

	newMasterHead := conn.MasterHead()
	if masterHead.Seqno+1 != newMasterHead.Seqno {
		t.Fatalf("expected seqno %d, got %d", masterHead.Seqno+1, newMasterHead.Seqno)
	}
}

func Test_connection_FindMinAvailableMasterchainSeqno(t *testing.T) {
	tests := []struct {
		name         string
		host         string
		key          string
		wantMinSeqno uint32
	}{
		{
			name:         "regular node",
			host:         "135.181.140.221:46995",
			key:          "wQE0MVhXNWUXpWiW5Bk8cAirIh5NNG3cZM1/fSVKIts=",
			wantMinSeqno: 36283540,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubkey, err := base64.StdEncoding.DecodeString(tt.key)
			if err != nil {
				t.Fatalf("failed to decode pubkey: %v", err)
			}
			c, err := liteclient.NewConnection(context.Background(), pubkey, tt.host)
			if err != nil {
				t.Skipf("cannot connect to %s: %v", tt.host, err)
			}
			conn := &connection{
				client:              liteclient.NewClient(c),
				masterHeadUpdatedCh: make(chan masterHeadUpdated, 100),
			}
			seqno, err := conn.FindMinAvailableMasterchainSeqno(context.Background())
			if err != nil {
				t.Fatalf("failed to find min seqno: %v", err)
			}
			if seqno < tt.wantMinSeqno {
				t.Fatalf("expected seqno >= %d, got %d", tt.wantMinSeqno, seqno)
			}
		})
	}
}

func Test_connection_SetMasterHead(t *testing.T) {
	t.Run("non-blocking when channel full", func(t *testing.T) {
		ch := make(chan masterHeadUpdated, 2)

		conn := &connection{
			id:                  1,
			masterHeadUpdatedCh: ch,
		}

		conn.SetMasterHead(ton.BlockIDExt{BlockID: ton.BlockID{Seqno: 1}})
		conn.SetMasterHead(ton.BlockIDExt{BlockID: ton.BlockID{Seqno: 2}})

		if len(ch) != 2 {
			t.Fatalf("expected channel to be full (2), got length %d", len(ch))
		}

		done := make(chan bool, 1)
		go func() {
			conn.SetMasterHead(ton.BlockIDExt{BlockID: ton.BlockID{Seqno: 3}})
			done <- true
		}()

		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("SetMasterHead blocked when channel was full")
		}

		readDone := make(chan ton.BlockIDExt, 1)
		go func() {
			readDone <- conn.MasterHead()
		}()

		select {
		case <-readDone:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("MasterHead() blocked, possible mutex deadlock")
		}
	})

	t.Run("concurrent access safe", func(t *testing.T) {
		conn := &connection{
			id:                  1,
			masterHeadUpdatedCh: make(chan masterHeadUpdated, 100),
		}

		var wg sync.WaitGroup

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(base uint32) {
				defer wg.Done()
				for j := uint32(0); j < 20; j++ {
					conn.SetMasterHead(ton.BlockIDExt{
						BlockID: ton.BlockID{Seqno: base*100 + j},
					})
					time.Sleep(time.Millisecond)
				}
			}(uint32(i))
		}

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 20; j++ {
					_ = conn.MasterHead()
					time.Sleep(time.Millisecond)
				}
			}()
		}

		done := make(chan bool)
		go func() {
			wg.Wait()
			done <- true
		}()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("deadlock detected during concurrent access")
		}
	})

	t.Run("ignores older seqno", func(t *testing.T) {
		conn := &connection{
			id:                  1,
			masterHeadUpdatedCh: make(chan masterHeadUpdated, 10),
		}

		conn.SetMasterHead(ton.BlockIDExt{BlockID: ton.BlockID{Seqno: 100}})

		head := conn.MasterHead()
		if head.Seqno != 100 {
			t.Fatalf("expected seqno 100, got %d", head.Seqno)
		}

		conn.SetMasterHead(ton.BlockIDExt{BlockID: ton.BlockID{Seqno: 50}})

		head = conn.MasterHead()
		if head.Seqno != 100 {
			t.Fatalf("older seqno was not ignored, expected 100, got %d", head.Seqno)
		}

		conn.SetMasterHead(ton.BlockIDExt{BlockID: ton.BlockID{Seqno: 101}})

		head = conn.MasterHead()
		if head.Seqno != 101 {
			t.Fatalf("newer seqno was not accepted, expected 101, got %d", head.Seqno)
		}
	})
}
