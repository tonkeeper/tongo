package liteclient

import (
	"context"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
	"sync"
	"time"
)

const (
	magicTCPPing = 0x4d082b9a //crc32(tcp.ping random_id:long = tcp.Pong)
	magicTCPPong = 0xdc69fb03 //crc32(tcp.pong random_id:long = tcp.Pong)
)

type ConnectionStatus int

const (
	Connecting ConnectionStatus = iota
	Connected
)
const (
	reconnectTimeout = 10 * time.Second
)

type Connection struct {
	peerPublicKey []byte
	host          string
	resp          chan Packet

	// mu protects status and econn.
	mu     sync.Mutex
	status ConnectionStatus
	econn  *encryptedConn
}

func NewConnection(ctx context.Context, peerPublicKey []byte, host string) (*Connection, error) {
	c := Connection{
		host:          host,
		peerPublicKey: peerPublicKey,
		resp:          make(chan Packet),
		status:        Connecting,
	}
	if err := c.setupEncryptedConnection(ctx); err != nil {
		return nil, err
	}
	go c.ping()
	return &c, nil
}

func (c *Connection) setupEncryptedConnection(ctx context.Context) error {
	econn, err := newEncryptedConnection(ctx, c.peerPublicKey, c.host)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.econn = econn
	c.status = Connected
	go c.reader(econn.handleIncomingPackages())
	return nil
}

func (c *Connection) reconnect() {
	c.mu.Lock()
	if c.status == Connecting {
		c.mu.Unlock()
		return
	}
	c.status = Connecting
	c.econn.close()
	c.mu.Unlock()

	for {
		if err := c.setupEncryptedConnection(context.Background()); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
}

func (c *Connection) reader(packetCh chan Packet) {
	for {
		select {
		case p, ok := <-packetCh:
			if !ok {
				// packetCh is closed, we are done.
				// setupEncryptedConnection will run another reader.
				return
			}
			if p.MagicType() == magicTCPPong {
				// no need to do anything,
				// the next iteration of this for loop will restart reconnect timeout.
				continue
			}
			c.resp <- p

		case <-time.After(reconnectTimeout):
			c.reconnect()
			// setupEncryptedConnection will run another reader.
			return
		}
	}
}

func (c *Connection) Send(p Packet) error {
	b := p.marshal()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status != Connected {
		return newClientError("not connected yet")
	}
	if err := c.econn.send(b); err != nil {
		go c.reconnect()
		return newClientError("net.Conn.send() failed: %v", err)
	}
	return nil
}

func (c *Connection) Responses() chan Packet {
	return c.resp
}

func (c *Connection) ping() {
	ping := make([]byte, 12)
	binary.LittleEndian.PutUint32(ping[:4], magicTCPPing)
	for {
		time.Sleep(time.Second * 3)
		mrand.Read(ping[4:])
		p, err := NewPacket(ping)
		if err != nil {
			panic(err) // impossible if NewPacket function is correct
		}
		err = c.Send(p)
		if err != nil && IsNotConnectedYet(err) {
			fmt.Printf("ping error: %v\n", err)
			continue
		}
	}
}

func (c *Connection) Status() ConnectionStatus {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status
}
