package pool

import (
	"context"
	"fmt"
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

	masterHeadUpdatedCh chan masterHeadUpdated

	mu         sync.RWMutex
	bestConn   conn
	waitListID uint64
	waitList   map[uint64]chan tongo.BlockIDExt
}

// conn contains all methods needed by a pool.
// used to implement tests.
type conn interface {
	ID() int
	MasterHead() tongo.BlockIDExt
	SetMasterHead(tongo.BlockIDExt)
	IsOK() bool
	Client() *liteclient.Client
	Run(ctx context.Context)
}

// NewFailoverPool returns a new instance of a failover pool.
// The given list of clients is ordered by priority and starts with a connection with the highest priority.
func NewFailoverPool(clients []*liteclient.Client) *FailoverPool {
	if len(clients) == 0 {
		panic("empty list of clients")
	}
	conns := make([]conn, 0, len(clients))
	masterHeadUpdatedCh := make(chan masterHeadUpdated, 10)

	for connID, cli := range clients {
		conns = append(conns, &connection{
			id:                  connID,
			client:              cli,
			masterHeadUpdatedCh: masterHeadUpdatedCh,
		})
	}
	return &FailoverPool{
		conns:               conns,
		updateBestInterval:  updateBestConnectionInterval,
		bestConn:            conns[0],
		waitList:            map[uint64]chan tongo.BlockIDExt{},
		masterHeadUpdatedCh: masterHeadUpdatedCh,
	}
}

func (p *FailoverPool) Run(ctx context.Context) {
	for _, c := range p.conns {
		go c.Run(ctx)
	}
	tickTock := time.NewTicker(p.updateBestInterval)
	defer tickTock.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tickTock.C:
			p.updateBest()
		case update := <-p.masterHeadUpdatedCh:
			p.notifySubscribers(update)
		}
	}
}

// updateBest finds the best suitable connection to work with and switches to it.
func (p *FailoverPool) updateBest() {
	var maxSeqno uint32
	for _, c := range p.conns {
		masterSeqno := c.MasterHead().Seqno
		if maxSeqno < masterSeqno {
			maxSeqno = masterSeqno
		}
	}

	for _, c := range p.conns {
		if !c.IsOK() {
			continue
		}
		if c.MasterHead().Seqno+1 >= maxSeqno {
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
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.bestConn
}

type MasterchainInfoClient struct {
	conn conn
}

func (s *MasterchainInfoClient) LiteServerGetMasterchainInfoExt(ctx context.Context, request liteclient.LiteServerGetMasterchainInfoExtRequest) (res liteclient.LiteServerMasterchainInfoExtC, err error) {
	info, err := s.conn.Client().LiteServerGetMasterchainInfoExt(ctx, request)
	if err != nil {
		return liteclient.LiteServerMasterchainInfoExtC{}, err
	}
	s.conn.SetMasterHead(info.Last.ToBlockIdExt())
	return info, err
}

func (s *MasterchainInfoClient) LiteServerGetMasterchainInfo(ctx context.Context) (liteclient.LiteServerMasterchainInfoC, error) {
	info, err := s.conn.Client().LiteServerGetMasterchainInfo(ctx)
	if err != nil {
		return liteclient.LiteServerMasterchainInfoC{}, err
	}
	s.conn.SetMasterHead(info.Last.ToBlockIdExt())
	return info, err
}

func (p *FailoverPool) BestMasterchainInfoClient() *MasterchainInfoClient {
	return &MasterchainInfoClient{
		conn: p.bestConnection(),
	}
}

// BestMasterchainClient returns a liteclient and its known masterchain head.
func (p *FailoverPool) BestMasterchainClient(ctx context.Context) (*liteclient.Client, tongo.BlockIDExt, error) {
	bestConnection := p.bestConnection()
	masterHead := bestConnection.MasterHead()
	if masterHead.Seqno > 0 {
		return bestConnection.Client(), bestConnection.MasterHead(), nil
	}
	// so this client is not initialized yet,
	// let's wait for it to be initialized.
	waitID, ch := p.subscribe(1)
	defer p.unsubscribe(waitID)

	select {
	case <-ctx.Done():
		return nil, tongo.BlockIDExt{}, ctx.Err()
	case head := <-ch:
		return bestConnection.Client(), head, nil
	}
}

// BestClientByAccountID returns a liteclient and its known masterchain head.
func (p *FailoverPool) BestClientByAccountID(ctx context.Context, accountID tongo.AccountID) (*liteclient.Client, tongo.BlockIDExt, error) {
	return p.BestMasterchainClient(ctx)
}

// BestClientByBlockID returns a liteclient and its known masterchain head.
func (p *FailoverPool) BestClientByBlockID(ctx context.Context, blockID tongo.BlockID) (*liteclient.Client, error) {
	server, _, err := p.BestMasterchainClient(ctx)
	return server, err
}

// ConnectionsNumber returns a number of connections in this pool.
func (p *FailoverPool) ConnectionsNumber() int {
	return len(p.conns)
}

func (p *FailoverPool) notifySubscribers(update masterHeadUpdated) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if update.Conn.ID() != p.bestConn.ID() {
		return
	}
	for _, ch := range p.waitList {
		ch <- update.Head
	}
}

func (p *FailoverPool) subscribe(seqno uint32) (uint64, chan tongo.BlockIDExt) {
	ch := make(chan tongo.BlockIDExt, 1)

	p.mu.Lock()
	defer p.mu.Unlock()

	head := p.bestConn.MasterHead()
	if head.Seqno >= seqno {
		ch <- head
		return 0, ch
	}
	p.waitListID++
	p.waitList[p.waitListID] = ch
	return p.waitListID, ch
}

func (p *FailoverPool) unsubscribe(waitID uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.waitList, waitID)
}

func (p *FailoverPool) WaitMasterchainSeqno(ctx context.Context, seqno uint32, timeout time.Duration) error {
	waitID, ch := p.subscribe(seqno)
	defer p.unsubscribe(waitID)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(timeout):
			return fmt.Errorf("timeout")
		case head := <-ch:
			if head.Seqno >= seqno {
				return nil
			}
		}
	}
}
