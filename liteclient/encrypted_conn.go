package liteclient

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"net"
)

// encryptedConn is a wrapper for a tcp connection that additionally encrypts and decrypts messages.
type encryptedConn struct {
	cipher   cipher.Stream
	decipher cipher.Stream
	conn     net.Conn
}

func newEncryptedConnection(ctx context.Context, peerPublicKey []byte, host string) (*encryptedConn, error) {
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
	var c = &encryptedConn{
		cipher:   cipher.NewCTR(ci, param.txNonce()),
		decipher: cipher.NewCTR(dci, param.rxNonce()),
		conn:     conn,
	}
	err = c.handshake(a, param, keys)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (econn *encryptedConn) handshake(address Address, params params, keys x25519Keys) error {
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
	_, err = econn.conn.Write(req)
	if err != nil {
		return err
	}
	_, err = ParsePacket(econn.conn, econn.decipher)
	if err != nil {
		return err
	}
	return nil
}

func (econn *encryptedConn) handleIncomingPackages() chan Packet {
	ch := make(chan Packet, 0)
	go func() {
		ioReader := bufio.NewReader(econn.conn)
		for {
			p, err := ParsePacket(ioReader, econn.decipher)
			if err != nil {
				// if there is an error, we are done with this connection.
				//
				// if somebody calls encryptedConn.close(),
				// ParsePacket() will return an error, and we will close the channel.
				close(ch)
				return
			}
			ch <- p
		}
	}()
	return ch
}

// send encrypts a message and sends it to a lite server.
// This method doesn't work properly with several goroutines sending a message at the same moment.
// Please, serialize access to it.
func (econn *encryptedConn) send(b []byte) error {
	econn.cipher.XORKeyStream(b, b)
	_, err := econn.conn.Write(b)
	return err
}

func (econn *encryptedConn) close() {
	if err := econn.conn.Close(); err != nil {
		fmt.Printf("Cant close connection\n")
	}
}
