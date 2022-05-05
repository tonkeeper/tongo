package adnl

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
)

const (
	magicADNLQuery  = 0x7af98bb4 // crc32("adnl.message.query query_id:int256 query:bytes = adnl.Message")
	magicADNLAnswer = 0x1684ac0f // crc32("adnl.message.answer query_id:int256 answer:bytes = adnl.Message")
)

type queryID [32]byte

type Client struct {
	connection   *Connection
	queries      map[queryID]chan Message
	queriesMutex sync.Mutex
}

type Query []byte

type Message []byte

func NewClient(c *Connection) *Client {
	c2 := &Client{
		connection: c,
		queries:    make(map[queryID]chan Message),
	}
	go c2.reader()
	return c2
}

func (c *Client) Request(ctx context.Context, q Query) (Message, error) {
	var id queryID
	rand.Read(id[:])
	data := make([]byte, 4, 44+len(q)) //create with small overhead for reducing garbage collector calls
	binary.BigEndian.PutUint32(data, magicADNLQuery)
	data = append(data, id[:]...)
	data = append(data, encodeLength(len(q))...)
	data = append(data, q...)
	if len(data)%4 != 0 { //
		data = append(data, make([]byte, 4-len(data)%4)...)
	}
	p, err := NewPacket(data)
	if err != nil {
		return Message{}, err
	}
	resp := c.registerCallback(id)
	err = c.connection.Send(p)
	if err != nil {
		return Message{}, err
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case b := <-resp:
		return b, nil
	}
}

func (c *Client) registerCallback(id queryID) chan Message {
	resp := make(chan Message)
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
	resp <- data[:length] //todo: maybe copy
	return nil
}
