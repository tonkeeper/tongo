package pool

import (
	"context"
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
	return 0
}

func (m *mockConn) MasterHead() ton.BlockIDExt {
	return ton.BlockIDExt{BlockID: ton.BlockID{Seqno: m.seqno}}
}

func (m *mockConn) SetMasterHead(ext ton.BlockIDExt) {
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
				t.Fatalf("want connection id: %v, got: %v", tt.wantID, c.id)
			}
		})
	}
}
