package pool

import (
	"context"
	"sync"
	"time"

	"github.com/tonkeeper/tongo/liteclient"
)

type connection struct {
	client      *liteclient.Client
	mu          sync.RWMutex
	masterSeqno uint32
}

func (c *connection) Run(ctx context.Context) {
	for {
		var seqno uint32
		for {
			res, err := c.client.LiteServerGetMasterchainInfo(ctx)
			if err != nil {
				// TODO: log error
				time.Sleep(1000 * time.Millisecond)
				continue
			}
			seqno = res.Last.Seqno
			break
		}
		c.setMasterSeqno(seqno)
		for {
			if err := c.client.WaitMasterchainSeqno(ctx, seqno+1, 15_000); err != nil {
				// TODO: log error
				time.Sleep(1000 * time.Millisecond)
				// we want to request seqno again with LiteServerGetMasterchainInfo
				// to avoid situation when this server has been offline for too long,
				// and it doesn't contain a block with the latest known seqno anymore.
				break
			}
			seqno += 1
			if ctx.Err() != nil {
				return
			}
			c.setMasterSeqno(seqno)
		}
	}
}

// IsOK returns true if there is no problems with the underlying liteclient and its connection to a lite server.
func (c *connection) IsOK() bool {
	return c.client.IsOK()
}

func (c *connection) Client() *liteclient.Client {
	return c.client
}

func (c *connection) MasterSeqno() uint32 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.masterSeqno
}

func (c *connection) setMasterSeqno(seqno uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.masterSeqno = seqno
}
