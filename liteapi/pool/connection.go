package pool

import (
	"context"
	"sync"
	"time"

	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/ton"
)

type connection struct {
	id         int
	serverHost string
	client     *liteclient.Client

	// masterHeadUpdatedCh is used to send a notification when a known master head is changed.
	masterHeadUpdatedCh chan masterHeadUpdated

	mu sync.RWMutex
	// masterHead is the latest known masterchain head.
	masterHead ton.BlockIDExt
	isArchive  bool
}

type masterHeadUpdated struct {
	Head ton.BlockIDExt
	Conn *connection
}

func (c *connection) Run(ctx context.Context, detectArchive bool) {
	if detectArchive {
		go func() {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
			defer cancel()
			// TODO: retry several times on error
			seqno, err := c.FindMinAvailableMasterchainSeqno(ctx)
			if err != nil {
				return
			}
			if seqno == 2 {
				c.setArchive(true)
			}
		}()
	}
	for {
		var head ton.BlockIDExt
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			res, err := c.client.LiteServerGetMasterchainInfo(ctx)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case <-time.After(1000 * time.Millisecond):
					continue
				}
			}
			head = res.Last.ToBlockIdExt()
			break
		}
		c.SetMasterHead(head)
		for {
			res, err := c.client.WaitMasterchainBlock(ctx, head.Seqno+1, 15_000)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case <-time.After(1000 * time.Millisecond):
					// we want to request seqno again with LiteServerGetMasterchainInfo
					// to avoid situation when this server has been offline for too long,
					// and it doesn't contain a block with the latest known seqno anymore.
					break
				}
				break
			}
			if ctx.Err() != nil {
				return
			}
			head = res.Id.ToBlockIdExt()
			c.SetMasterHead(head)
		}
	}
}

// IsOK returns true if there is no problems with the underlying liteclient and its connection to a lite server.
func (c *connection) IsOK() bool {
	return c.client.IsOK()
}

func (c *connection) ID() int {
	return c.id
}

func (c *connection) Client() *liteclient.Client {
	return c.client
}

func (c *connection) MasterHead() ton.BlockIDExt {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.masterHead
}

func (c *connection) SetMasterHead(head ton.BlockIDExt) {
	c.mu.Lock()
	if head.Seqno <= c.masterHead.Seqno {
		c.mu.Unlock()
		return
	}
	c.masterHead = head
	c.mu.Unlock()

	select {
	case c.masterHeadUpdatedCh <- masterHeadUpdated{
		Head: head,
		Conn: c,
	}:
	default:
		// Channel full - skip notification, pool will catch up on next update
	}
}

func (c *connection) FindMinAvailableMasterchainSeqno(ctx context.Context) (uint32, error) {
	info, err := c.client.LiteServerGetMasterchainInfo(ctx)
	if err != nil {
		return 0, err
	}
	max := info.Last.Seqno
	min := uint32(2)

	next := min
	workchain := -1
	for min+1 < max {
		request := liteclient.LiteServerLookupBlockRequest{
			Mode: 1,
			Id: liteclient.TonNodeBlockIdC{
				Workchain: int32(workchain),
				Shard:     0x8000000000000000,
				Seqno:     next,
			},
		}
		_, err := c.client.LiteServerLookupBlock(ctx, request)
		if err != nil {
			if e, ok := err.(liteclient.LiteServerErrorC); ok && e.Code == 651 {
				min = next + 1
				next = (min + max) / 2
				continue
			}
			return 0, err
		}
		max = next - 1
		next = (min + max) / 2
	}
	return min, nil
}

func (c *connection) IsArchiveNode() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isArchive
}

func (c *connection) setArchive(archive bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.isArchive = archive
}

func (c *connection) AverageRoundTrip() time.Duration {
	return c.client.AverageRoundTrip()
}

type ConnStatus struct {
	ServerHost string
	Connected  bool
	Archive    bool
}

func (c *connection) Status() ConnStatus {
	return ConnStatus{
		ServerHost: c.serverHost,
		Connected:  c.IsOK(),
		Archive:    c.IsArchiveNode(),
	}
}
