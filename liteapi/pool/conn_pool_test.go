package pool

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/ton"
)

type mockConn struct {
	id           int
	seqno        uint32
	isOK         bool
	avgRoundTrip time.Duration
}

func (m *mockConn) AverageRoundTrip() time.Duration {
	return m.avgRoundTrip
}

func (m *mockConn) IsArchiveNode() bool {
	return false
}

func (m *mockConn) ID() int {
	return m.id
}

func (m *mockConn) MasterHead() ton.BlockIDExt {
	return ton.BlockIDExt{BlockID: ton.BlockID{Seqno: m.seqno}}
}

func (m *mockConn) SetMasterHead(ext ton.BlockIDExt) {
	m.seqno = ext.Seqno
}

func (m *mockConn) MasterSeqno() uint32 {
	return m.seqno
}

func (m *mockConn) IsOK() bool {
	return m.isOK
}

func (m *mockConn) Client() *liteclient.Client {
	panic("implement me")
}

func (m *mockConn) Status() ConnStatus {
	panic("implement me")
}

func (m *mockConn) Run(ctx context.Context, detectArchiveNodes bool) {
}

var _ conn = &mockConn{}

func TestConnPool_updateBest(t *testing.T) {
	tests := []struct {
		name     string
		conns    []conn
		strategy Strategy
		wantID   int
	}{
		{
			name:     "pick up the first connection when it works",
			strategy: FirstWorkingConnection,
			conns: []conn{
				&mockConn{seqno: 100, isOK: true, id: 0},
				&mockConn{seqno: 100, isOK: true, id: 1},
				&mockConn{seqno: 100, isOK: true, id: 2},
			},
			wantID: 0,
		},
		{
			name:     "pick up the first connection when it works but is slightly slow",
			strategy: FirstWorkingConnection,
			conns: []conn{
				&mockConn{seqno: 99, isOK: true, id: 0},
				&mockConn{seqno: 100, isOK: true, id: 1},
				&mockConn{seqno: 100, isOK: true, id: 2},
			},
			wantID: 0,
		},
		{
			name:     "pick up the second connection when the first one doesnt work",
			strategy: FirstWorkingConnection,
			conns: []conn{
				&mockConn{seqno: 100, isOK: false, id: 0},
				&mockConn{seqno: 100, isOK: true, id: 1},
				&mockConn{seqno: 100, isOK: true, id: 2},
			},
			wantID: 1,
		},
		{
			name:     "pick up the second connection when the first one is too slow",
			strategy: FirstWorkingConnection,
			conns: []conn{
				&mockConn{seqno: 98, isOK: true, id: 0},
				&mockConn{seqno: 100, isOK: true, id: 1},
				&mockConn{seqno: 100, isOK: true, id: 2},
			},
			wantID: 1,
		},
		{
			name:     "0 - slow, 1 - broken, 2 - OK",
			strategy: FirstWorkingConnection,
			conns: []conn{
				&mockConn{seqno: 98, isOK: true, id: 0},
				&mockConn{seqno: 100, isOK: false, id: 1},
				&mockConn{seqno: 100, isOK: true, id: 2},
			},
			wantID: 2,
		},
		{
			name:     "3 - is the best",
			strategy: BestPingStrategy,
			conns: []conn{
				&mockConn{seqno: 99, isOK: true, id: 0, avgRoundTrip: 10 * time.Millisecond},
				&mockConn{seqno: 100, isOK: true, id: 1, avgRoundTrip: 20 * time.Millisecond},
				&mockConn{seqno: 100, isOK: true, id: 2, avgRoundTrip: 5 * time.Millisecond},
			},
			wantID: 2,
		},
		{
			name:     "3 - is the best but doesn't work, 1 - OK",
			strategy: BestPingStrategy,
			conns: []conn{
				&mockConn{seqno: 99, isOK: true, id: 0, avgRoundTrip: 10 * time.Millisecond},
				&mockConn{seqno: 100, isOK: true, id: 1, avgRoundTrip: 20 * time.Millisecond},
				&mockConn{seqno: 100, isOK: false, id: 2, avgRoundTrip: 5 * time.Millisecond},
			},
			wantID: 0,
		},
		{
			name:     "3 - is the best but doesn't work, 1 - slow, 2 - ok",
			strategy: BestPingStrategy,
			conns: []conn{
				&mockConn{seqno: 98, isOK: true, id: 0, avgRoundTrip: 10 * time.Millisecond},
				&mockConn{seqno: 100, isOK: true, id: 1, avgRoundTrip: 20 * time.Millisecond},
				&mockConn{seqno: 100, isOK: false, id: 2, avgRoundTrip: 5 * time.Millisecond},
			},
			wantID: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ConnPool{
				conns:              tt.conns,
				updateBestInterval: time.Second,
				strategy:           tt.strategy,
			}
			ctx := context.Background()
			go p.Run(ctx)
			p.updateBest()
			c := p.bestConnection().(*mockConn)
			if tt.wantID != c.id {
				t.Fatalf("expected connection id %d, got %d", tt.wantID, c.id)
			}
		})
	}
}

type mockConnWithClient struct {
	*connection
	mockIsOK bool
}

func (m *mockConnWithClient) IsOK() bool {
	return m.mockIsOK
}

func TestWaitMasterchainSeqno(t *testing.T) {
	t.Run("timer reuse prevents memory leak", func(t *testing.T) {
		pool := New(BestPingStrategy)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go pool.Run(ctx)

		testConn := &connection{
			id:                  1,
			masterHeadUpdatedCh: pool.masterHeadUpdatedCh,
		}
		testConn.SetMasterHead(ton.BlockIDExt{BlockID: ton.BlockID{Seqno: 1}})

		wrappedConn := &mockConnWithClient{
			connection: testConn,
			mockIsOK:   true,
		}

		pool.mu.Lock()
		pool.conns = []conn{wrappedConn}
		pool.bestConn = wrappedConn
		pool.mu.Unlock()

		go func() {
			for i := uint32(2); i < 50; i++ {
				time.Sleep(10 * time.Millisecond)
				testConn.SetMasterHead(ton.BlockIDExt{BlockID: ton.BlockID{Seqno: i}})
			}
		}()

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		for i := uint32(2); i < 45; i++ {
			err := pool.WaitMasterchainSeqno(ctx, i, 500*time.Millisecond)
			if err != nil {
				t.Fatalf("failed to wait for seqno %d: %v", i, err)
			}
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		allocatedKB := (m2.TotalAlloc - m1.TotalAlloc) / 1024
		if allocatedKB > 500 {
			t.Errorf("excessive memory allocation: %d KB (expected < 500 KB with timer reuse)", allocatedKB)
		}
	})

	t.Run("timeout when seqno not reached", func(t *testing.T) {
		pool := New(BestPingStrategy)

		ctx := context.Background()
		go pool.Run(ctx)

		testConn := &connection{
			id:                  1,
			masterHeadUpdatedCh: pool.masterHeadUpdatedCh,
		}
		testConn.SetMasterHead(ton.BlockIDExt{BlockID: ton.BlockID{Seqno: 1}})

		wrappedConn := &mockConnWithClient{
			connection: testConn,
			mockIsOK:   true,
		}

		pool.mu.Lock()
		pool.conns = []conn{wrappedConn}
		pool.bestConn = wrappedConn
		pool.mu.Unlock()

		err := pool.WaitMasterchainSeqno(ctx, 9999, 100*time.Millisecond)
		if err == nil {
			t.Fatal("expected timeout error, got nil")
		}
		if err.Error() != "timeout" {
			t.Fatalf("expected 'timeout' error, got %q", err.Error())
		}
	})
}

func TestMultipleConnectionsChannelHandling(t *testing.T) {
	pool := New(BestPingStrategy)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pool.Run(ctx)

	var wrappedConns []*mockConnWithClient
	for i := 0; i < 5; i++ {
		c := &connection{
			id:                  i,
			masterHeadUpdatedCh: pool.masterHeadUpdatedCh,
		}
		wrapped := &mockConnWithClient{
			connection: c,
			mockIsOK:   true,
		}
		wrappedConns = append(wrappedConns, wrapped)
	}

	pool.mu.Lock()
	pool.conns = make([]conn, len(wrappedConns))
	for i, c := range wrappedConns {
		pool.conns[i] = c
	}
	pool.bestConn = wrappedConns[0]
	pool.mu.Unlock()

	var blocked atomic.Int32
	var wg sync.WaitGroup

	for idx, wrapped := range wrappedConns {
		wg.Add(1)
		go func(connID int, conn *connection) {
			defer wg.Done()

			for seqno := uint32(1); seqno <= 20; seqno++ {
				done := make(chan bool, 1)

				go func(s uint32) {
					conn.SetMasterHead(ton.BlockIDExt{BlockID: ton.BlockID{Seqno: s}})
					done <- true
				}(seqno)

				select {
				case <-done:
				case <-time.After(100 * time.Millisecond):
					blocked.Add(1)
					return
				}

				time.Sleep(5 * time.Millisecond)
			}
		}(idx, wrapped.connection)
	}

	wg.Wait()

	if blocked.Load() > 0 {
		t.Errorf("failed to send updates: %d connection(s) blocked (channel overflow)", blocked.Load())
	}
}

func TestChannelBufferCapacity(t *testing.T) {
	pool := New(BestPingStrategy)

	capacity := cap(pool.masterHeadUpdatedCh)
	if capacity < 10 {
		t.Errorf("insufficient channel buffer: got %d, expected >= 10", capacity)
	}
}
