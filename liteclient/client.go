package liteclient

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/tonkeeper/tongo/tl"
	"math/rand"
	"sync"
	"time"
)

const (
	magicADNLQuery       = 0xb48bf97a // crc32("adnl.message.query query_id:int256 query:bytes = adnl.Message")
	magicADNLAnswer      = 0x0fac8416 // crc32("adnl.message.answer query_id:int256 answer:bytes = adnl.Message")
	magicLiteServerQuery = 0x798c06df // crc32("liteServer.query#df068c79 data:bytes = Object")
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
	connection   *Connection
	queries      map[queryID]chan []byte
	queriesMutex sync.Mutex
}

type Options func(connection *Client)

func OptionTimeout(t time.Duration) Options {
	return func(c *Client) {
		c.timeout = t
	}
}

func NewClient(c *Connection, opts ...Options) *Client {
	c2 := &Client{
		timeout:    defaultTimeout,
		connection: c,
		queries:    make(map[queryID]chan []byte),
	}
	for _, f := range opts {
		f(c2)
	}
	go c2.reader()
	return c2
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
		return nil, err
	}
	resp := c.registerCallback(id)
	err = c.connection.Send(p)
	if err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		c.unregisterCallback(id)
		return nil, ctx.Err()
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

func (c *Client) reader() {
	for p := range c.connection.Responses() {
		if p.MagicType() != magicADNLAnswer {
			fmt.Println("unknown type", p.MagicType()) //todo: remove
			continue
		}
		err := c.processQueryAnswer(p)
		if err != nil {
			fmt.Println(err) //todo: switch to debug logger
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
		return fmt.Errorf("unknow query %x", id)
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
