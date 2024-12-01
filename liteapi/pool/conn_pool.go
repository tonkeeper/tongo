package pool

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/ton"
)

const (
	updateBestConnectionInterval = 10 * time.Second
)

var (
	ErrNoConnections = fmt.Errorf("no connections available")
)

type Strategy string

const (
	BestPingStrategy       = "best-ping"
	FirstWorkingConnection = "first-working"
)

// ConnPool is a pool of connections to lite servers
// that implements two different strategy:
//  1. BestPingStrategy - it'll switch to a connection with the best ping.
//  2. FirstWorkingConnection - it'll switch to the first working connection.
//
// For both strategies, a connection has to be not more than 1 block behind the head of masterchain to be considered as working.
type ConnPool struct {
	strategy           Strategy
	updateBestInterval time.Duration

	masterHeadUpdatedCh chan masterHeadUpdated

	mu         sync.RWMutex
	conns      []conn
	bestConn   conn
	waitListID uint64
	waitList   map[uint64]chan ton.BlockIDExt
}

// conn contains all methods needed by a pool.
// used to implement tests.
type conn interface {
	ID() int
	MasterHead() ton.BlockIDExt
	SetMasterHead(ton.BlockIDExt)
	IsOK() bool
	Client() *liteclient.Client
	Run(ctx context.Context, detectArchive bool)
	IsArchiveNode() bool
	AverageRoundTrip() time.Duration
	Status() ConnStatus
}

// New returns a new instance of a connections pool.
func New(strategy Strategy) *ConnPool {
	return &ConnPool{
		strategy:            strategy,
		updateBestInterval:  updateBestConnectionInterval,
		waitList:            map[uint64]chan ton.BlockIDExt{},
		masterHeadUpdatedCh: make(chan masterHeadUpdated, 10),
	}
}

func (p *ConnPool) InitializeConnections(ctx context.Context, timeout time.Duration, maxConnections int, detectArchiveNodes bool, servers []config.LiteServer) chan error {
	ch := make(chan error, 1)
	go func() {
		clientsCh := make(chan clientWrapper, len(servers))
		for connID, server := range servers {
			go func(connID int, server config.LiteServer) {
				cli, _ := connect(ctx, timeout, server)
				// TODO: log error
				clientsCh <- clientWrapper{
					connID:     connID,
					cli:        cli,
					serverHost: server.Host,
				}
			}(connID, server)
		}

		processedConnections := 0
		for {
			if processedConnections == len(servers) {
				break
			}
			select {
			case <-ctx.Done():
				break
			case wrapper := <-clientsCh:
				processedConnections += 1
				if wrapper.cli == nil {
					continue
				}
				if p.ConnectionsNumber() < maxConnections {
					c := p.addConnection(wrapper.connID, wrapper.cli, wrapper.serverHost)
					go c.Run(context.TODO(), detectArchiveNodes)
				}
				if p.ConnectionsNumber() == maxConnections {
					processedConnections = len(servers)
					break
				}
			}
		}
		if p.ConnectionsNumber() == 0 {
			ch <- fmt.Errorf("all liteservers are unavailable")
			return
		}
		ch <- nil
	}()
	return ch
}

func connect(ctx context.Context, timeout time.Duration, server config.LiteServer) (*liteclient.Client, error) {
	serverPubkey, err := base64.StdEncoding.DecodeString(server.Key)
	if err != nil {
		return nil, err
	}
	c, err := liteclient.NewConnection(ctx, serverPubkey, server.Host)
	if err != nil {
		return nil, err
	}
	cli := liteclient.NewClient(c, liteclient.OptionTimeout(timeout))
	if _, err := cli.LiteServerGetMasterchainInfo(ctx); err != nil {
		return nil, err
	}
	return cli, nil
}

type clientWrapper struct {
	connID     int
	cli        *liteclient.Client
	serverHost string
}

func (p *ConnPool) addConnection(connID int, cli *liteclient.Client, serverHost string) *connection {
	p.mu.Lock()
	defer p.mu.Unlock()
	c := &connection{
		id:                  connID,
		serverHost:          serverHost,
		client:              cli,
		masterHeadUpdatedCh: p.masterHeadUpdatedCh,
	}
	p.conns = append(p.conns, c)
	sort.Slice(p.conns, func(i, j int) bool {
		return p.conns[i].ID() < p.conns[j].ID()
	})
	if len(p.conns) == 1 {
		p.bestConn = c
	}
	return c
}

func (p *ConnPool) Run(ctx context.Context) {
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
func (p *ConnPool) updateBest() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.conns) == 0 {
		return
	}

	var maxSeqno uint32
	for _, c := range p.conns {
		masterSeqno := c.MasterHead().Seqno
		if maxSeqno < masterSeqno {
			maxSeqno = masterSeqno
		}
	}

	switch p.strategy {
	case BestPingStrategy:
		if bestConn := p.findBestPingConnection(maxSeqno); bestConn != nil {
			p.bestConn = bestConn
		}
	case FirstWorkingConnection:
		if bestConn := p.findFirstWorkingConnection(maxSeqno); bestConn != nil {
			p.bestConn = bestConn
		}
	}
}

func (p *ConnPool) findFirstWorkingConnection(maxSeqno uint32) conn {
	for _, c := range p.conns {
		if !c.IsOK() {
			continue
		}
		if c.MasterHead().Seqno+1 >= maxSeqno {
			return c
		}
	}
	return nil
}

func (p *ConnPool) findBestPingConnection(maxSeqno uint32) conn {
	var bestConn conn
	for _, c := range p.conns {
		if !c.IsOK() {
			continue
		}
		if c.MasterHead().Seqno+1 < maxSeqno {
			continue
		}
		if bestConn == nil || c.AverageRoundTrip() < bestConn.AverageRoundTrip() {
			bestConn = c
		}
	}
	return bestConn
}

func (p *ConnPool) bestConnection() conn {
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

func (p *ConnPool) BestMasterchainInfoClient() *MasterchainInfoClient {
	return &MasterchainInfoClient{
		conn: p.bestConnection(),
	}
}

// BestMasterchainClient returns a liteclient and its known masterchain head.
func (p *ConnPool) BestMasterchainClient(ctx context.Context) (*liteclient.Client, ton.BlockIDExt, error) {
	bestConnection := p.bestConnection()
	if bestConnection == nil {
		return nil, ton.BlockIDExt{}, ErrNoConnections
	}
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
		return nil, ton.BlockIDExt{}, ctx.Err()
	case head := <-ch:
		return bestConnection.Client(), head, nil
	}
}
func (p *ConnPool) BestArchiveClient(ctx context.Context) (*liteclient.Client, ton.BlockIDExt, error) {
	for _, c := range p.conns {
		if c.IsOK() && c.IsArchiveNode() {
			return c.Client(), c.MasterHead(), nil
		}
	}
	return nil, ton.BlockIDExt{}, fmt.Errorf("no archive nodes available")
}

// BestClientByAccountID returns a liteclient and its known masterchain head.
func (p *ConnPool) BestClientByAccountID(ctx context.Context, accountID ton.AccountID, archiveRequired bool) (*liteclient.Client, ton.BlockIDExt, error) {
	if archiveRequired {
		return p.BestArchiveClient(ctx)
	}
	return p.BestMasterchainClient(ctx)
}

// BestClientByBlockID returns a liteclient and its known masterchain head.
func (p *ConnPool) BestClientByBlockID(ctx context.Context, blockID ton.BlockID) (*liteclient.Client, error) {
	server, _, err := p.BestMasterchainClient(ctx)
	return server, err
}

// ConnectionsNumber returns a number of connections in this pool.
func (p *ConnPool) ConnectionsNumber() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.conns)
}

func (p *ConnPool) notifySubscribers(update masterHeadUpdated) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.bestConn == nil {
		return
	}
	if update.Conn.ID() != p.bestConn.ID() {
		return
	}
	for _, ch := range p.waitList {
		ch <- update.Head
	}
}

func (p *ConnPool) subscribe(seqno uint32) (uint64, chan ton.BlockIDExt) {
	ch := make(chan ton.BlockIDExt, 1)

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

func (p *ConnPool) unsubscribe(waitID uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.waitList, waitID)
}

func (p *ConnPool) WaitMasterchainSeqno(ctx context.Context, seqno uint32, timeout time.Duration) error {
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

type Status struct {
	Connections []ConnStatus
}

func (p *ConnPool) Status() Status {
	p.mu.Lock()
	defer p.mu.Unlock()

	connStatuses := make([]ConnStatus, 0, len(p.conns))
	for _, c := range p.conns {
		connStatuses = append(connStatuses, c.Status())
	}
	return Status{
		Connections: connStatuses,
	}
}
