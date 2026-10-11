package liteclient

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"testing"
)

// Regression for #463: sendAuthComplete built the payload with
//
//	payload := make([]byte, 4)
//	binary.LittleEndian.PutUint32(payload, magicTcpAuthentificationComplete)
//	binary.LittleEndian.PutUint32(payload, magicPubKey)
//
// so the second PutUint32 overwrote the first at offset 0. The message type
// magic (magicTcpAuthentificationComplete) was dropped and the payload was
// 4 bytes short. The serialized tcp.authentificationComplete message must be:
//
//	magicTcpAuthentificationComplete (4) | magicPubKey (4) | pubKey (32) |
//	len(signature) prefix | signature | zero padding to a 4-byte boundary.
func TestBuildAuthCompletePayloadLayout(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signature := ed25519.Sign(priv, []byte("nonce"))

	payload := buildAuthCompletePayload(pub, signature)

	// Both magics must be present in their own 4-byte slots.
	gotMsgMagic := binary.LittleEndian.Uint32(payload[0:4])
	if gotMsgMagic != magicTcpAuthentificationComplete {
		t.Fatalf("message magic: got %#x, want %#x (it was overwritten before the fix)",
			gotMsgMagic, uint32(magicTcpAuthentificationComplete))
	}
	gotPubMagic := binary.LittleEndian.Uint32(payload[4:8])
	if gotPubMagic != magicPubKey {
		t.Fatalf("pubkey magic: got %#x, want %#x", gotPubMagic, uint32(magicPubKey))
	}

	// The 32-byte public key follows the two magics.
	if !bytes.Equal(payload[8:8+ed25519.PublicKeySize], pub) {
		t.Fatalf("public key not serialized at offset 8")
	}

	// The signature must survive round-trip: 1-byte length prefix (64 < 254)
	// then the signature bytes.
	sigStart := 8 + ed25519.PublicKeySize
	if int(payload[sigStart]) != len(signature) {
		t.Fatalf("signature length prefix: got %d, want %d", payload[sigStart], len(signature))
	}
	gotSig := payload[sigStart+1 : sigStart+1+len(signature)]
	if !bytes.Equal(gotSig, signature) {
		t.Fatalf("signature not serialized correctly")
	}

	// Payload must be padded to a 4-byte boundary.
	if len(payload)%4 != 0 {
		t.Fatalf("payload not aligned to 4 bytes: len=%d", len(payload))
	}
}
