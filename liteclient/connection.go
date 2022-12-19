package liteclient

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	mrand "math/rand"
	"net"
	"sync"
	"time"
)

const (
	magicTCPPing = 0x4d082b9a //crc32(tcp.ping random_id:long = tcp.Pong)
	magicTCPPong = 0xdc69fb03 //crc32(tcp.pong random_id:long = tcp.Pong)
)

type ConnectionStatus int

const (
	NotInit ConnectionStatus = iota
	Connecting
	Connected
	Closed
)

type Connection struct {
	cipher        cipher.Stream
	decipher      cipher.Stream
	conn          net.Conn
	packetMutex   sync.Mutex
	resp          chan Packet
	peerPublicKey []byte
	host          string
	Status        ConnectionStatus
	lastPong      time.Time
}

func NewConnection(ctx context.Context, peerPublicKey []byte, host string) (*Connection, error) {
	c, err := createConnection(ctx, peerPublicKey, host)
	if err != nil {
		return nil, err
	}
	go c.reader()
	go c.ping()
	go c.watchdog(time.Second * 10)
	return c, nil
}

func createConnection(ctx context.Context, peerPublicKey []byte, host string) (*Connection, error) {
	a, err := NewAddress(peerPublicKey)
	if err != nil {
		return nil, err
	}
	param, err := newParameters()
	if err != nil {
		return nil, err
	}
	ci, err := aes.NewCipher(param.txKey())
	if err != nil {
		return nil, err
	}
	dci, err := aes.NewCipher(param.rxKey())
	if err != nil {
		return nil, err
	}
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", host)
	if err != nil {
		return nil, err
	}
	keys, err := newKeys(a.pubkey)
	if err != nil {
		return nil, err
	}
	var c = &Connection{
		cipher:        cipher.NewCTR(ci, param.txNonce()),
		decipher:      cipher.NewCTR(dci, param.rxNonce()),
		conn:          conn,
		resp:          make(chan Packet, 1000),
		peerPublicKey: peerPublicKey,
		host:          host,
	}
	err = c.handshake(a, param, keys)
	if err != nil {
		return nil, err
	}
	c.Status = Connected
	c.lastPong = time.Now()
	return c, nil
}

func (c *Connection) watchdog(timeout time.Duration) {
	for {
		if time.Since(c.lastPong) > timeout && c.Status != Connecting {
			c.reconnect()
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func (c *Connection) reconnect() {
	c.Status = Connecting
	err := c.conn.Close()
	if err != nil {
		fmt.Printf("Cant close connection\n")
	}
	for {
		conn, err := createConnection(context.Background(), c.peerPublicKey, c.host)
		if err == nil {
			c.cipher = conn.cipher
			c.decipher = conn.decipher
			c.conn = conn.conn
			c.packetMutex = sync.Mutex{}
			c.lastPong = time.Now()
			go c.reader()
			c.Status = Connected
			return
		}
		fmt.Printf("Reconnect failed: %v\n", err)
		time.Sleep(time.Second * 5)
	}
}

func (c *Connection) reader() {
	ioReader := bufio.NewReader(c.conn)
	for {
		p, err := ParsePacket(ioReader, c.decipher)
		if errors.Is(err, net.ErrClosed) {
			fmt.Printf("Old connection closed. Drop reader.: %v\n", err)
			break
		}
		if err != nil {
			fmt.Printf("reader error: %v\n", err)
			time.Sleep(time.Millisecond * 200)
			continue
		}
		if p.MagicType() == magicTCPPong {
			c.lastPong = time.Now()
			continue
		}
		c.resp <- p
	}
}

func (c *Connection) handshake(address Address, params params, keys x25519Keys) error {
	key := append([]byte{}, keys.shared[:16]...)
	key = append(key, params.hash()[16:32]...)
	nonce := append([]byte{}, params.hash()[0:4]...)
	nonce = append(nonce, keys.shared[20:32]...)
	cipherKey, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	data := append([]byte{}, params[:]...)
	cipher.NewCTR(cipherKey, nonce).XORKeyStream(data, data)
	req := make([]byte, 256)
	copy(req[:32], address.hash())
	copy(req[32:64], keys.public)
	copy(req[64:96], params.hash())
	copy(req[96:], data)
	_, err = c.conn.Write(req)
	if err != nil {
		return err
	}
	_, err = ParsePacket(c.conn, c.decipher)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) Send(p Packet) error {
	if c.Status != Connected {
		return nil
	}
	b := p.marshal()
	c.packetMutex.Lock()
	c.cipher.XORKeyStream(b, b)
	_, err := c.conn.Write(b)
	c.packetMutex.Unlock()
	if err != nil {
		fmt.Printf("Sending error: %v\n", err)
		c.reconnect()
	}
	return err
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
		if err != nil {
			fmt.Printf("ping error: %v\n", err)
			continue
		}
	}
}
