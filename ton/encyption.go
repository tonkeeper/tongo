package ton

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"github.com/oasisprotocol/curve25519-voi/curve"
	ed25519crv "github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/oasisprotocol/curve25519-voi/primitives/x25519"
)

func Encrypt(receiverPubkey ed25519.PublicKey, ourPrivateKey ed25519.PrivateKey, data []byte, salt []byte) ([]byte, error) {
	if len(receiverPubkey) != ed25519.PublicKeySize {
		return nil, errors.New("receiverPubkey size is invalid")
	}
	if len(ourPrivateKey) != ed25519.PrivateKeySize {
		return nil, errors.New("ourPrivateKey size is invalid")
	}
	ourPublicKey := ourPrivateKey.Public().(ed25519.PublicKey)
	sharedSecret, err := getSharedSecret(ourPrivateKey, receiverPubkey)
	if err != nil {
		return nil, err
	}
	prefixedData, err := addPrefix(data)
	if err != nil {
		return nil, err
	}
	dataHash := combineSecrets(salt, prefixedData)
	msgKey := dataHash[:16]
	x := combineSecrets(sharedSecret, msgKey)
	key := x[:32]
	iv := x[32:48]
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	enc := cipher.NewCBCEncrypter(c, iv)

	encryptedData := make([]byte, len(prefixedData))
	enc.CryptBlocks(encryptedData, prefixedData)

	pubXorKey := ourPublicKey
	for i := 0; i < 32; i++ {
		pubXorKey[i] ^= receiverPubkey[i]
	}

	var result []byte
	result = append(result, pubXorKey...)
	result = append(result, msgKey...)
	result = append(result, encryptedData...)

	return result, nil
}

func getSharedSecret(ourPrivateKey ed25519.PrivateKey, theirPublicKey ed25519.PublicKey) ([]byte, error) {
	comp, err := curve.NewCompressedEdwardsYFromBytes(theirPublicKey)
	if err != nil {
		return nil, err
	}

	ep, err := curve.NewEdwardsPoint().SetCompressedY(comp)
	if err != nil {
		return nil, err
	}

	mp := curve.NewMontgomeryPoint().SetEdwards(ep)
	bb := x25519.EdPrivateKeyToX25519(ed25519crv.PrivateKey(ourPrivateKey))

	key, err := x25519.X25519(bb, mp[:])
	if err != nil {
		return nil, err
	}

	return key, nil
}

func addPrefix(data []byte) ([]byte, error) {
	prefixLen := 16 + (16-len(data)%16)%16
	prefix := make([]byte, prefixLen)
	_, err := rand.Read(prefix)
	if err != nil {
		return nil, err
	}
	prefix[0] = byte(prefixLen)
	var res []byte
	res = append(res, prefix...)
	res = append(res, data...)
	if len(res)%16 != 0 {
		return nil, errors.New("data length is not a multiple of 16")
	}
	return res, nil
}

func combineSecrets(a, b []byte) []byte {
	h := hmac.New(sha512.New, a)
	h.Write(b)
	return h.Sum(nil)
}

// TODO: add decrypt
