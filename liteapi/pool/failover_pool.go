package pool

import (
	"context"
	"sync"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/liteclient"
)

const (
	updateBestConnectionInterval = 10 * time.Second
)

// FailoverPool is a pool of connections to lite servers
// that implements a failover strategy and in case of failure
// it'll automatically switch to a working connection.
//
// The failover strategy is as follows:
// given a list of connections,
// a pool uses the first connection from the list that is
// 1. working
// 2. not more than 1 block behind the head of masterchain.
type FailoverPool struct {
	conns              []conn
	updateBestInterval time.Duration

	mu       sync.RWMutex
	bestConn conn
}

// conn contains all methods needed by a pool.
// used to implement tests.
type conn interface {
	MasterSeqno() uint32
	IsOK() bool
	Client() *liteclient.Client
	Run(ctx context.Context)
}

// NewFailoverPool returns a new instance of a failover pool.
// The given list of clients is ordered by priority and starts with a connection with the highest priority.
func NewFailoverPool(clients []*liteclient.Client) *FailoverPool {
	conns := make([]conn, 0, len(clients))
	for _, cli := range clients {
		conns = append(conns, &connection{client: cli})
	}
	return &FailoverPool{
		conns:              conns,
		updateBestInterval: updateBestConnectionInterval,
		bestConn:           conns[0],
	}
}

func (p *FailoverPool) Run(ctx context.Context) {
	for _, c := range p.conns {
		go c.Run(ctx)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(p.updateBestInterval):
			p.updateBest()
		}
	}
}

// updateBest finds the best suitable connection to work with and switches to it.
func (p *FailoverPool) updateBest() {
	var maxSeqno uint32
	for _, c := range p.conns {
		masterSeqno := c.MasterSeqno()
		if maxSeqno < masterSeqno {
			maxSeqno = masterSeqno
		}
	}

	for _, c := range p.conns {
		if !c.IsOK() {
			continue
		}
		if c.MasterSeqno()+1 >= maxSeqno {
			p.setBestConnection(c)
			return
		}
	}
}

func (p *FailoverPool) setBestConnection(conn conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.bestConn = conn
}

func (p *FailoverPool) bestConnection() conn {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.bestConn
}

func (p *FailoverPool) BestMasterchainServer() *liteclient.Client {
	return p.bestConnection().Client()
}

func (p *FailoverPool) BestServerByAccountID(tongo.AccountID) (*liteclient.Client, error) {
	return p.BestMasterchainServer(), nil
}

func (p *FailoverPool) BestServerByBlockID(tongo.BlockID) (*liteclient.Client, error) {
	return p.BestMasterchainServer(), nil
}

// ConnectionsNumber returns a number of connections in this pool.
func (p *FailoverPool) ConnectionsNumber() int {
	return len(p.conns)
}
