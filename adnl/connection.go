package adnl

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
	"net"
	"time"
)

const (
	magicTCPPing = 0x9a2b084d //crc32(tcp.ping random_id:long = tcp.Pong)
	magicTCPPong = 0x03fb69dc //crc32(tcp.pong random_id:long = tcp.Pong)
)

type Connection struct {
	address Address
	params  params
	keys    x25519Keys

	cipher   cipher.Stream
	decipher cipher.Stream

	conn net.Conn
	resp chan Packet

	peerPublicKey []byte
	host          string
	IsReady       bool
	lastPong      time.Time
	ioReader      *bufio.Reader
}

func NewConnection(ctx context.Context, peerPublicKey []byte, host string) (*Connection, error) {
	c, err := createConnection(ctx, peerPublicKey, host)
	if err != nil {
		return nil, err
	}
	c.IsReady = true
	go c.reader()
	go c.ping()
	return c, nil
}

func createConnection(ctx context.Context, peerPublicKey []byte, host string) (*Connection, error) {
	a, err := NewAddress(peerPublicKey)
	if err != nil {
		return nil, err
	}
	params, err := newParameters()
	if err != nil {
		return nil, err
	}
	ci, err := aes.NewCipher(params.txKey())
	if err != nil {
		return nil, err
	}
	dci, err := aes.NewCipher(params.rxKey())
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
		address:  a,
		params:   params,
		keys:     keys,
		cipher:   cipher.NewCTR(ci, params.txNonce()),
		decipher: cipher.NewCTR(dci, params.rxNonce()),
		conn:     conn,
		resp:     make(chan Packet, 1000),
	}
	err = c.handshake()
	if err != nil {
		return nil, err
	}
	c.peerPublicKey = peerPublicKey
	c.host = host
	return c, nil
}

func (c *Connection) reconnect() {
	c.IsReady = false
	err := c.conn.Close()
	if err != nil {
		fmt.Printf("Cant close connection\n")
	}
	for {
		conn, err := createConnection(context.Background(), c.peerPublicKey, c.host)
		if err == nil {
			c.address = conn.address
			c.params = conn.params
			c.keys = conn.keys
			c.cipher = conn.cipher
			c.decipher = conn.decipher
			c.conn = conn.conn
			c.IsReady = true
			c.ioReader = bufio.NewReader(c.conn)
			return
		}
		fmt.Printf("Reconnect failed: %v\n", err)
		time.Sleep(time.Second * 5)
	}
}

func (c *Connection) reader() {
	c.ioReader = bufio.NewReader(c.conn)
	c.lastPong = time.Now()
	for {
		if c.IsReady != true {
			continue
		}
		err := c.conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		if err != nil {
			fmt.Printf("set deadline error: %v\n", err)
		}
		p, err := ParsePacket(c.ioReader, c.decipher)
		if err != nil {
			//fmt.Printf("reader error: %v\n", err)
			continue
		}
		if p.MagicType() == magicTCPPong {
			c.lastPong = time.Now()
			continue
		}
		c.resp <- p
	}
}

func (c *Connection) handshake() error {
	key := append([]byte{}, c.keys.shared[:16]...)
	key = append(key, c.params.hash()[16:32]...)
	nonce := append([]byte{}, c.params.hash()[0:4]...)
	nonce = append(nonce, c.keys.shared[20:32]...)
	cipherKey, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	data := append([]byte{}, c.params[:]...)
	cipher.NewCTR(cipherKey, nonce).XORKeyStream(data, data)
	req := make([]byte, 256)
	copy(req[:32], c.address.hash())
	copy(req[32:64], c.keys.public)
	copy(req[64:96], c.params.hash())
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
	b := p.marshal()
	c.cipher.XORKeyStream(b, b)
	_, err := c.conn.Write(b)
	return err
}

func (c *Connection) Responses() chan Packet {
	return c.resp
}

func (c *Connection) ping() {
	ping := make([]byte, 12)
	binary.BigEndian.PutUint32(ping[:4], magicTCPPing)
	for {
		time.Sleep(time.Second * 10)
		mrand.Read(ping[4:])
		p, err := NewPacket(ping)
		if err != nil {
			panic(err) // impossible if NewPacket function is correct
		}
		err = c.Send(p)
		if err != nil {
			fmt.Printf("send ping error: %v\n", err)
			continue
		}
		if time.Since(c.lastPong) > time.Second*20 {
			c.reconnect()
			c.lastPong = time.Now()
		}
	}
}
