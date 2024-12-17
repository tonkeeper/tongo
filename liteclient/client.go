package liteclient

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/tonkeeper/tongo/tl"
)

const (
	magicADNLQuery                      = 0xb48bf97a // crc32("adnl.message.query query_id:int256 query:bytes = adnl.Message")
	magicADNLAnswer                     = 0x0fac8416 // crc32("adnl.message.answer query_id:int256 answer:bytes = adnl.Message")
	magicLiteServerQuery                = 0x798c06df // crc32("liteServer.query#df068c79 data:bytes = Object")
	magicLiteServerWaitMasterchainSeqno = 0xbaeab892
)

type queryID [32]byte

const (
	defaultTimeout = time.Minute
)

type Client struct {
	// timeout configures a timeout of a lite client method.
	// if such a method makes several calls to a lite server,
	// the total time is bounded by the timeout.
	timeout      time.Duration
	connections  []*Connection
	nextConn     int
	connMutex    sync.Mutex
	queries      map[queryID]chan []byte
	queriesMutex sync.Mutex
}

type Options func(connection *Client)

func OptionTimeout(t time.Duration) Options {
	return func(c *Client) {
		c.timeout = t
	}
}

func OptionWorkersPerConnection(n int) Options {
	return func(c *Client) {
		if n < 1 {
			n = 1
		}
		connFirst := c.connections[0]
		for i := 0; i < n-1; i++ {
			conn, err := NewConnection(context.Background(), connFirst.peerPublicKey, connFirst.host)
			if err != nil {
				slog.Warn("liteclient clone connection error", err)
				continue
			}
			c.connections = append(c.connections, conn)
		}
	}
}

func NewClient(c *Connection, opts ...Options) *Client {
	c2 := &Client{
		timeout:     defaultTimeout,
		connections: []*Connection{c},
		queries:     make(map[queryID]chan []byte),
	}
	for _, f := range opts {
		f(c2)
	}

	for _, conn := range c2.connections {
		go c2.reader(conn)
	}
	return c2
}

// IsOK returns true if there is no problems with this client and its underlying connection to a lite server.
func (c *Client) IsOK() bool {
	for _, conn := range c.connections {
		if conn.Status() == Connected {
			return true
		}
	}
	return false
}

// Request sends q as query in adnl.message.query and receives answer from adnl.message.answer
// adnl.message.query query_id:int256 query:bytes = adnl.Message
// adnl.message.answer query_id:int256 answer:bytes = adnl.Message
func (c *Client) Request(ctx context.Context, q []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	var id queryID
	rand.Read(id[:])
	data := make([]byte, 4, 44+len(q)) //create with small overhead for reducing garbage collector calls
	binary.LittleEndian.PutUint32(data, magicADNLQuery)
	data = append(data, id[:]...)
	data = append(data, encodeLength(len(q))...)
	data = append(data, q...)
	data = alignBytes(data)
	p, err := NewPacket(data)
	if err != nil {
		return nil, newClientError("NewPacket() failed: %v", err)
	}
	resp := c.registerCallback(id)
	defer c.unregisterCallback(id)

	c.connMutex.Lock()
	conn := c.connections[c.nextConn]
	c.nextConn = (c.nextConn + 1) % len(c.connections)
	c.connMutex.Unlock()

	err = conn.Send(p)
	if err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, newClientError("request timeout: %v", ctx.Err())
	case b := <-resp:
		return b, nil
	}
}

func (c *Client) registerCallback(id queryID) chan []byte {
	resp := make(chan []byte, 1)
	c.queriesMutex.Lock()
	c.queries[id] = resp
	c.queriesMutex.Unlock()
	return resp
}

func (c *Client) unregisterCallback(id queryID) {
	c.queriesMutex.Lock()
	delete(c.queries, id)
	c.queriesMutex.Unlock()
}

func encodeLength(i int) []byte {
	if i >= 254 {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(i<<8))
		b[0] = 254
		return b
	} else {
		return []byte{byte(i)}
	}
}

func decodeLength(b []byte) (int, []byte, error) {
	if len(b) == 0 {
		return 0, nil, fmt.Errorf("size should contains at least one byte")
	}
	if b[0] == 255 {
		return 0, nil, fmt.Errorf("invalid first byte value %x", b[0])
	}
	if b[0] < 254 {
		return int(b[0]), b[1:], nil
	}
	if b[0] != 254 {
		panic("how it cat be possible? you are fucking wizard!")
	}
	if len(b) < 4 {
		return 0, nil, fmt.Errorf("not enought bytes for decoding size")
	}
	b[0] = 0
	i := binary.LittleEndian.Uint32(b[:4])
	b[0] = 254
	return int(i) >> 8, b[4:], nil
}

func (c *Client) reader(conn *Connection) {
	for p := range conn.Responses() {
		if p.MagicType() != magicADNLAnswer {
			continue
		}
		err := c.processQueryAnswer(p)
		if err != nil {
			slog.Info("liteclient.reader() error", "err", err)
		}
	}
}

func (c *Client) processQueryAnswer(p Packet) error {
	if len(p.Payload) < 37 {
		return fmt.Errorf("too short payload")
	}
	var id queryID
	copy(id[:], p.Payload[4:36])
	c.queriesMutex.Lock()
	resp, prs := c.queries[id]
	delete(c.queries, id)
	c.queriesMutex.Unlock()
	if !prs {
		return fmt.Errorf("unknown query %x with id %x", p.Payload[:4], id)
	}
	length, data, err := decodeLength(p.Payload[36:])
	if err != nil {
		return err
	}
	if len(data) < length {
		return fmt.Errorf("payload is smaller than should be according to length")
	}
	resp <- data[:length]
	return nil
}

// liteServerRequest sends q as liteServer.query data:bytes = Object;
func (c *Client) liteServerRequest(ctx context.Context, q []byte) ([]byte, error) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, magicLiteServerQuery)
	data = append(data, tl.EncodeLength(len(q))...)
	data = append(data, q...)
	data = alignBytes(data)
	return c.Request(ctx, data)
}

func alignBytes(data []byte) []byte {
	left := len(data) % 4
	if left != 0 {
		data = append(data, make([]byte, 4-left)...)
	}
	return data
}

// WaitMasterchainSeqno waits for the given block to become committed.
// If timeout happens, it returns an error.
func (c *Client) WaitMasterchainSeqno(ctx context.Context, seqno uint32, timeout uint32) error {
	data := make([]byte, 0, 12)
	data = binary.LittleEndian.AppendUint32(data, magicLiteServerWaitMasterchainSeqno)
	data = binary.LittleEndian.AppendUint32(data, seqno)
	data = binary.LittleEndian.AppendUint32(data, timeout)
	resp, err := c.liteServerRequest(ctx, data)
	if err != nil {
		return err
	}
	if len(resp) < 4 {
		return fmt.Errorf("not enough bytes for tag")
	}
	tag := binary.LittleEndian.Uint32(resp[:4])
	if tag == 0xbba9e148 {
		var errRes LiteServerErrorC
		if err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes); err != nil {
			return err
		}
		if errRes.Code == 0 {
			return nil
		}
		return errRes
	}
	return fmt.Errorf("invalid tag")
}

func (c *Client) AverageRoundTrip() time.Duration {
	var total time.Duration
	for _, conn := range c.connections {
		total += conn.AverageRoundTrip()
	}
	return total / time.Duration(len(c.connections))
}

func (c *Client) WaitMasterchainBlock(ctx context.Context, seqno uint32, timeout uint32) (res LiteServerBlockHeaderC, err error) {
	var (
		mc     int    = -1
		uintMc uint32 = uint32(mc)
	)
	request := LiteServerLookupBlockRequest{
		Mode: 1,
		Id: TonNodeBlockIdC{
			Workchain: uintMc,
			Shard:     0x8000000000000000,
			Seqno:     seqno,
		},
	}
	data := make([]byte, 0, 38)
	data = binary.LittleEndian.AppendUint32(data, magicLiteServerWaitMasterchainSeqno)
	data = binary.LittleEndian.AppendUint32(data, seqno)
	data = binary.LittleEndian.AppendUint32(data, timeout)
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerLookupBlockRequest `tlSumType:"fac8f71e"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return res, err
	}
	data = append(data, payload...)
	resp, err := c.liteServerRequest(ctx, data)
	if err != nil {
		return res, err
	}
	if len(resp) < 4 {
		return res, fmt.Errorf("not enough bytes for tag")
	}
	tag := binary.LittleEndian.Uint32(resp[:4])
	if tag == 0xbba9e148 {
		var errRes LiteServerErrorC
		if err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes); err != nil {
			return res, err
		}
		return res, errRes
	}
	if tag == 0x752d8219 {
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return res, fmt.Errorf("invalid tag")
}
