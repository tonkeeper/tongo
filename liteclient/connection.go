package liteclient

import (
	"context"
	"encoding/binary"
	"log/slog"
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

	// mu protects all fields below.
	mu           sync.Mutex
	status       ConnectionStatus
	econn        *encryptedConn
	pings        map[uint64]time.Time
	avgRoundTrip time.Duration
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
	c.pings = make(map[uint64]time.Time, 5)
	go c.reader(econn.handleIncomingPackets())
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
		c.registerPing(binary.LittleEndian.Uint64(ping[4:]))
		err = c.Send(p)
		if err != nil && IsNotConnectedYet(err) {
			slog.Info("ping error", "host", c.host, "error", err)
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
