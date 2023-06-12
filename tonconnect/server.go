package tonconnect

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/wallet"
)

type Server struct {
	executor        abi.Executor
	secret          string
	lifeTimeProof   int64
	lifeTimePayload int64
}

type Options struct {
	lifeTimePayload int64
	lifeTimeProof   int64
}

type Option func(o *Options)

func WithLifeTimePayload(lifeTimePayload int64) Option {
	return func(o *Options) {
		o.lifeTimePayload = lifeTimePayload
	}
}

func WithLifeTimeProof(lifeTimeProof int64) Option {
	return func(o *Options) {
		o.lifeTimeProof = lifeTimeProof
	}
}

const (
	defaultLifeTimeProof   = 300
	defaultLifeTimePayload = 300
)

const (
	tonProofPrefix   = "ton-proof-item-v2/"
	tonConnectPrefix = "ton-connect"
)

func NewTonConnect(executor abi.Executor, secret string, opts ...Option) (*Server, error) {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}

	if executor == nil {
		return nil, fmt.Errorf("executor is not configured")
	}

	if options.lifeTimePayload == 0 {
		options.lifeTimePayload = defaultLifeTimePayload
	}

	if options.lifeTimeProof == 0 {
		options.lifeTimeProof = defaultLifeTimeProof
	}

	return &Server{
		executor:        executor,
		secret:          secret,
		lifeTimeProof:   options.lifeTimeProof,
		lifeTimePayload: options.lifeTimePayload,
	}, nil
}

var knownHashes = make(map[string]wallet.Version)

func init() {
	for i := wallet.Version(0); i <= wallet.V4R2; i++ {
		ver := wallet.GetCodeHashByVer(i)
		knownHashes[hex.EncodeToString(ver[:])] = i
	}
}

type Proof struct {
	Address string    `json:"address"`
	Proof   ProofData `json:"proof"`
}

type ProofData struct {
	Timestamp int64  `json:"timestamp"`
	Domain    string `json:"domain"`
	Signature string `json:"signature"`
	Payload   string `json:"payload"`
	StateInit string `json:"state_init"`
}

type parsedMessage struct {
	workChain int32
	address   []byte
	ts        int64
	domain    string
	signature []byte
	payload   string
	stateInit string
}

// GeneratePayload this is the first stage of the authorization process. Client fetches payload to be signed by wallet
func (s *Server) GeneratePayload() (string, error) {
	payload := make([]byte, 16, 48)
	_, err := rand.Read(payload[:8])
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}
	binary.BigEndian.PutUint64(payload[8:16], uint64(time.Now().Add(time.Duration(s.lifeTimePayload)).Unix()))
	hmacHash := hmac.New(sha256.New, []byte(s.secret))
	hmacHash.Write(payload)
	payload = hmacHash.Sum(payload)
	return hex.EncodeToString(payload[:32]), nil
}

// CheckProof aggregated with ready-made data from TonConnect 2.0
// 1) Client fetches payload from GeneratePayload to be signed by wallet
// 2) Client connects to the wallet via TonConnect 2.0 and passes ton_proof request with specified payload, for more
// details see the frontend SDK: https://github.com/ton-connect/sdk/tree/main/packages/sdk
// 3) User approves connection and client receives signed payload with additional prefixes.
// 4) Client sends signed result (Proof) to the backend and CheckProof checks correctness of the all prefixes and signature correctness
func (s *Server) CheckProof(ctx context.Context, tp *Proof) (bool, ed25519.PublicKey, error) {
	verified, err := s.checkPayload(tp.Proof.Payload)
	if !verified {
		return false, nil, fmt.Errorf("failed verify payload")
	}

	parsed, err := s.convertTonProofMessage(tp)
	if err != nil {
		return false, nil, err
	}

	if time.Since(time.Unix(parsed.ts, 0)) > time.Duration(s.lifeTimeProof)*time.Second {
		return false, nil, fmt.Errorf("proof has been expired")
	}

	accountID, err := tongo.ParseAccountID(tp.Address)
	if err != nil {
		return false, nil, err
	}

	pubKey, err := s.getWalletPubKey(ctx, accountID)
	if err != nil {
		if tp.Proof.StateInit == "" {
			return false, nil, fmt.Errorf("failed get public key")
		}
		if ok, err := compareStateInitWithAddress(accountID, tp.Proof.StateInit); err != nil || !ok {
			return false, nil, fmt.Errorf("failed compare state init with address")
		}
		pubKey, err = ParseStateInit(tp.Proof.StateInit)
		if err != nil {
			return false, nil, fmt.Errorf("failed get public key")
		}
	}

	mes, err := createMessage(parsed)
	if err != nil {
		return false, nil, err
	}

	check := signatureVerify(pubKey, mes, parsed.signature)
	if !check {
		return false, nil, fmt.Errorf("failed proof")
	}

	return true, pubKey, nil
}

func (s *Server) GetSecret() string {
	return s.secret
}

func (s *Server) checkPayload(payload string) (bool, error) {
	bytesPayload, err := hex.DecodeString(payload)
	if err != nil {
		return false, err
	}
	if len(bytesPayload) != 32 {
		return false, fmt.Errorf("invalid payload length")
	}
	hmacHash := hmac.New(sha256.New, []byte(s.secret))
	hmacHash.Write(bytesPayload[:16])
	computedSignature := hmacHash.Sum(nil)

	if subtle.ConstantTimeCompare(bytesPayload[16:], computedSignature[:16]) != 1 {
		return false, fmt.Errorf("invalid payload signature")
	}

	if time.Since(time.Unix(int64(binary.BigEndian.Uint64(bytesPayload[8:16])), 0)) > time.Duration(s.lifeTimePayload)*time.Second {
		return false, fmt.Errorf("payload expired")
	}

	return true, nil
}

func (s *Server) convertTonProofMessage(tp *Proof) (*parsedMessage, error) {
	addr := strings.Split(tp.Address, ":")
	if len(addr) != 2 {
		return nil, fmt.Errorf("invalid address param: %v", tp.Address)
	}

	workChain, err := strconv.ParseInt(addr[0], 10, 32)
	if err != nil {
		return nil, err
	}

	walletAddr, err := hex.DecodeString(addr[1])
	if err != nil {
		return nil, err
	}

	sig, err := base64.StdEncoding.DecodeString(tp.Proof.Signature)
	if err != nil {
		return nil, err
	}

	return &parsedMessage{
		workChain: int32(workChain),
		address:   walletAddr,
		domain:    tp.Proof.Domain,
		ts:        tp.Proof.Timestamp,
		signature: sig,
		payload:   tp.Proof.Payload,
		stateInit: tp.Proof.StateInit,
	}, nil
}

func (s *Server) getWalletPubKey(ctx context.Context, address tongo.AccountID) (ed25519.PublicKey, error) {
	_, result, err := abi.GetPublicKey(ctx, s.executor, address)
	if err != nil {
		return nil, err
	}
	if r, ok := result.(abi.GetPublicKeyResult); ok {
		i := big.Int(r.PublicKey)
		b := i.Bytes()
		if len(b) < 24 || len(b) > 32 {
			return nil, fmt.Errorf("invalid publock key")
		}
		return append(make([]byte, 32-len(b)), b...), nil
	}

	return nil, fmt.Errorf("can't get publick key")
}

func createMessage(message *parsedMessage) ([]byte, error) {
	wc := make([]byte, 4)
	binary.BigEndian.PutUint32(wc, uint32(message.workChain))

	ts := make([]byte, 8)
	binary.LittleEndian.PutUint64(ts, uint64(message.ts))

	dl := make([]byte, 4)
	binary.LittleEndian.PutUint32(dl, uint32(len(message.domain)))

	m := []byte(tonProofPrefix)
	m = append(m, wc...)
	m = append(m, message.address...)
	m = append(m, dl...)
	m = append(m, []byte(message.domain)...)
	m = append(m, ts...)
	m = append(m, []byte(message.payload)...)

	messageHash := sha256.Sum256(m)
	fullMes := []byte{0xff, 0xff}
	fullMes = append(fullMes, []byte(tonConnectPrefix)...)
	fullMes = append(fullMes, messageHash[:]...)

	res := sha256.Sum256(fullMes)
	return res[:], nil
}

func signatureVerify(pubKey ed25519.PublicKey, message, signature []byte) bool {
	return ed25519.Verify(pubKey, message, signature)
}

func ParseStateInit(stateInit string) ([]byte, error) {
	cells, err := boc.DeserializeBocBase64(stateInit)
	if err != nil || len(cells) != 1 {
		return nil, err
	}

	var state tlb.StateInit
	err = tlb.Unmarshal(cells[0], &state)
	if err != nil {
		return nil, err
	}

	if !state.Data.Exists || !state.Code.Exists {
		return nil, err
	}

	hash, err := state.Code.Value.Value.HashString()
	if err != nil {
		return nil, err
	}

	version, prs := knownHashes[hash]
	if !prs {
		return nil, fmt.Errorf("unknown hash")
	}

	var pubKey tlb.Bits256
	switch version {
	case wallet.V1R1, wallet.V1R2, wallet.V1R3, wallet.V2R1, wallet.V2R2:
		var data wallet.DataV1V2
		err = tlb.Unmarshal(&state.Data.Value.Value, &data)
		if err != nil {
			return nil, err
		}
		pubKey = data.PublicKey

	case wallet.V3R1, wallet.V3R2, wallet.V4R1, wallet.V4R2:
		var data wallet.DataV3
		err = tlb.Unmarshal(&state.Data.Value.Value, &data)
		if err != nil {
			return nil, err
		}
		pubKey = data.PublicKey
	}

	return pubKey[:], nil
}

func compareStateInitWithAddress(a tongo.AccountID, stateInit string) (bool, error) {
	cells, err := boc.DeserializeBocBase64(stateInit)
	if err != nil || len(cells) != 1 {
		return false, err
	}
	h, err := cells[0].Hash()
	if err != nil {
		return false, err
	}
	return bytes.Equal(h, a.Address[:]), nil
}
