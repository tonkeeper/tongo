package pool

import (
	"context"
	"testing"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/liteclient"
)

type mockConn struct {
	id    int
	seqno uint32
	isOK  bool
}

func (m *mockConn) ID() int {
	return 0
}

func (m *mockConn) MasterHead() tongo.BlockIDExt {
	return tongo.BlockIDExt{BlockID: tongo.BlockID{Seqno: m.seqno}}
}

func (m *mockConn) SetMasterHead(ext tongo.BlockIDExt) {
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

func (m *mockConn) Run(ctx context.Context) {
}

var _ conn = &mockConn{}

func TestFailoverPool_updateBest(t *testing.T) {
	tests := []struct {
		name   string
		conns  []conn
		wantID int
	}{
		{
			name: "pick up the first connection when it works",
			conns: []conn{
				&mockConn{seqno: 100, isOK: true, id: 0},
				&mockConn{seqno: 100, isOK: true, id: 1},
				&mockConn{seqno: 100, isOK: true, id: 2},
			},
			wantID: 0,
		},
		{
			name: "pick up the first connection when it works but is slightly slow",
			conns: []conn{
				&mockConn{seqno: 99, isOK: true, id: 0},
				&mockConn{seqno: 100, isOK: true, id: 1},
				&mockConn{seqno: 100, isOK: true, id: 2},
			},
			wantID: 0,
		},
		{
			name: "pick up the second connection when the first one doesnt work",
			conns: []conn{
				&mockConn{seqno: 100, isOK: false, id: 0},
				&mockConn{seqno: 100, isOK: true, id: 1},
				&mockConn{seqno: 100, isOK: true, id: 2},
			},
			wantID: 1,
		},
		{
			name: "pick up the second connection when the first one is too slow",
			conns: []conn{
				&mockConn{seqno: 98, isOK: true, id: 0},
				&mockConn{seqno: 100, isOK: true, id: 1},
				&mockConn{seqno: 100, isOK: true, id: 2},
			},
			wantID: 1,
		},
		{
			name: "0 - slow, 1 - broken, 2 - OK",
			conns: []conn{
				&mockConn{seqno: 98, isOK: true, id: 0},
				&mockConn{seqno: 100, isOK: false, id: 1},
				&mockConn{seqno: 100, isOK: true, id: 2},
			},
			wantID: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &FailoverPool{
				conns:              tt.conns,
				updateBestInterval: time.Second,
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
