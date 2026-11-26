package liteclient

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/tonkeeper/tongo/tl"
)

const (
	magicTCPPing                     = 0x4d082b9a //crc32(tcp.ping random_id:long = tcp.Ping)
	magicTCPPong                     = 0xdc69fb03 //crc32(tcp.pong random_id:long = tcp.Pong)
	magicTcpAuthentificate           = 0x445bab12 //crc32(tcp.authentificate nonce:bytes = tcp.Message)
	magicTcpAuthentificationNonce    = 0xe35d4ab6 //crc32(tcp.authentificationNonce nonce:bytes = tcp.Message)
	magicTcpAuthentificationComplete = 0xf7ad9ea6 //crc32(tcp.authentificationComplete key:PublicKey signature:bytes = tcp.Message)
	magicPubKey                      = 0x4813b4c6 //crc32(pub.ed25519 key:int256 = PublicKey)
)

type ConnectionStatus int

const (
	Connecting ConnectionStatus = iota
	Connected
)
const (
	reconnectTimeout = 10 * time.Second
)

const (
	ClientNonceSize    = 32
	MaxServerNonceSize = 512
)

type Connection struct {
	peerPublicKey    []byte
	host             string
	resp             chan Packet
	authKey          ed25519.PrivateKey
	authCompleteChan chan error // Closes when auth is complete or error is sent

	// mu protects all fields below.
	mu           sync.Mutex
	status       ConnectionStatus
	econn        *encryptedConn
	pings        map[uint64]time.Time
	avgRoundTrip time.Duration
	nonce        []byte
}

func NewConnection(ctx context.Context, peerPublicKey []byte, host string, authKeys ...ed25519.PrivateKey) (*Connection, error) {
	if len(authKeys) > 1 {
		return nil, fmt.Errorf("expected 0 or 1 authentication key, got: %d", len(authKeys))
	}
	c := Connection{
		host:             host,
		peerPublicKey:    peerPublicKey,
		resp:             make(chan Packet),
		status:           Connecting,
		authCompleteChan: make(chan error),
	}
	if len(authKeys) == 1 {
		c.authKey = authKeys[0]
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
	c.econn = econn
	c.pings = make(map[uint64]time.Time, 5)

	go c.reader(econn.handleIncomingPackets())

	if c.authKey == nil {
		c.status = Connected
		c.mu.Unlock()
		return nil
	}

	c.mu.Unlock()

	if err := c.sendAuthRequest(); err != nil {
		c.econn.close()
		return fmt.Errorf("failed to send auth request: %w", err)
	}

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	select {
	case authErr := <-c.authCompleteChan:
		if authErr != nil {
			err = authErr
		}
	case <-timer.C:
		err = fmt.Errorf("authentication timeout")
	case <-ctx.Done():
		err = ctx.Err()
	}
	if err != nil {
		c.econn.close()
		return err
	}
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
			fmt.Printf("error reconnecting to %s: %s\n", c.host, err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
}

func averageRoundTrip(roundTrips []time.Duration) time.Duration {
	total := time.Duration(0)
	for _, rt := range roundTrips {
		total += rt
	}
	return total / time.Duration(len(roundTrips))
}

func (c *Connection) reader(packetCh chan Packet) {
	var roundTrips []time.Duration
	for {
		select {
		case p, ok := <-packetCh:
			if !ok {
				// packetCh is closed, we are done.
				// setupEncryptedConnection will run another reader.
				return
			}
			if p.MagicType() == magicTCPPong && len(p.Payload) == 12 {
				roundTrip, ok := c.processPong(binary.LittleEndian.Uint64(p.Payload[4:]))
				if ok {
					roundTrips = append(roundTrips, roundTrip)
					if len(roundTrips) > 5 {
						roundTrips = roundTrips[1:]
					}
					avg := averageRoundTrip(roundTrips)
					c.setAverageRoundTrip(avg)
				}
				// no need to do anything,
				// the next iteration of this for loop will restart reconnect timeout.
				continue
			}

			if p.MagicType() == magicTcpAuthentificationNonce {
				if err := c.handleAuthResponse(p); err != nil {
					fmt.Printf("failed to handle auth nonce: %s\n", err)
				}
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
		if _, err := rand.Read(ping[4:]); err != nil {
			panic(err) // impossible if source of randomness is correct
		}
		p, err := NewPacket(ping)
		if err != nil {
			panic(err) // impossible if NewPacket function is correct
		}
		c.registerPing(binary.LittleEndian.Uint64(ping[4:]))
		err = c.Send(p)
		if err != nil && IsNotConnectedYet(err) {
			continue
		}
	}
}

func (c *Connection) registerPing(randomID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.pings[randomID] = now
	for id, t := range c.pings {
		if t.Add(30 * time.Second).Before(now) {
			delete(c.pings, id)
		}
	}
}

func (c *Connection) processPong(randomID uint64) (time.Duration, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if startTime, ok := c.pings[randomID]; ok {
		delete(c.pings, randomID)
		return time.Since(startTime), true
	}
	return 0, false
}

func (c *Connection) setAverageRoundTrip(avg time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.avgRoundTrip = avg
}

// AverageRoundTrip returns an average round trip of the last several pings.
func (c *Connection) AverageRoundTrip() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.avgRoundTrip
}

func (c *Connection) Status() ConnectionStatus {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status
}

func (c *Connection) sendAuthRequest() error {
	nonce := make([]byte, ClientNonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("error generating nonce: %w", err)
	}

	c.mu.Lock()
	c.nonce = nonce
	c.mu.Unlock()

	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload, magicTcpAuthentificate)
	payload = append(payload, tl.EncodeLength(len(nonce))...)
	payload = append(payload, nonce...)
	payload = alignBytes(payload)

	p, err := NewPacket(payload)
	if err != nil {
		return fmt.Errorf("failed to create auth packet: %w", err)
	}
	b := p.marshal()
	return c.econn.send(b)
}

func (c *Connection) handleAuthResponse(p Packet) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.status != Connecting {
		return fmt.Errorf("received unexpected auth packet")
	}

	err := c.sendAuthComplete(p)
	if err == nil {
		c.status = Connected
	}
	c.authCompleteChan <- err
	return nil
}

func (c *Connection) sendAuthComplete(received Packet) error {
	if c.authKey == nil {
		return fmt.Errorf("no auth key available")
	}

	if len(received.Payload) < 37 {
		return fmt.Errorf("too short payload")
	}
	length, data, err := decodeLength(received.Payload[4:])
	if err != nil {
		return err
	}
	if len(data) < length {
		return fmt.Errorf("payload is smaller than should be according to length")
	}
	nonce := data[:length]
	if len(nonce) > MaxServerNonceSize {
		return fmt.Errorf("too long nonce")
	}

	signature := ed25519.Sign(c.authKey, append(append([]byte{}, c.nonce...), nonce...))

	payload := make([]byte, 0, 4+4+ed25519.PublicKeySize+ed25519.SignatureSize+8)
	payload = binary.LittleEndian.AppendUint32(payload, magicTcpAuthentificationComplete)

	payload = binary.LittleEndian.AppendUint32(payload, magicPubKey)

	pubKey := c.authKey.Public().(ed25519.PublicKey)
	payload = append(payload, pubKey[:]...)

	payload = append(payload, tl.EncodeLength(len(signature))...)
	payload = append(payload, signature...)
	payload = alignBytes(payload)

	p, err := NewPacket(payload)
	if err != nil {
		return fmt.Errorf("failed to create packet: %w", err)
	}
	b := p.marshal()
	if err := c.econn.send(b); err != nil {
		return fmt.Errorf("failed to send auth sign request, err: %w", err)
	}
	return nil
}
