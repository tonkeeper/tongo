package tongo

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/snksoft/crc"
	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/contract/dns"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

const (
	// DefaultRoot is the default DNS root address used by the addressParser.
	DefaultRoot = "-1:e56754f83426f69b09267bd876ac97c44821345b7e266bd956a7bfbfb98df35c"
)

type AccountID = ton.AccountID

// Deprecated use ton.NewAccountID instead
var NewAccountId = ton.NewAccountID

// Deprecated: use ParseAddress instead.
var ParseAccountID = ton.ParseAccountID

// Deprecated: use MustParseAddress instead.
var MustParseAccountID = ton.MustParseAccountID

var AccountIDFromTlb = ton.AccountIDFromTlb

// mu protects defaultParser.
var mu sync.RWMutex
var defaultParser *addressParser

// DefaultAddressParser returns a default address parser that works in the mainnet.
// For other networks, use SetDefaultExecutor(testnetLiteapiClient).
// Take a look at NewAccountAddressParser to create a parser for a different network or with a different root address.
func DefaultAddressParser() *addressParser {
	mu.RLock()
	defer mu.RUnlock()
	return defaultParser
}

// SetDefaultExecutor sets the default executor for the default address parser.
// The executor is used to resolve DNS records.
func SetDefaultExecutor(executor abi.Executor) {
	mu.Lock()
	defer mu.Unlock()
	resolver := dns.NewDNS(MustParseAccountID(DefaultRoot), executor)
	defaultParser = NewAccountAddressParser(resolver)
}

func init() {
	defaultParser = NewAccountAddressParser(&lazyResolver{})
}

// addressParser converts a string of different formats to a ton.Address.
type addressParser struct {
	resolver dnsResolver
}

// ParseAddress parses a string of different formats to a ton.Address.
func ParseAddress(a string) (ton.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return DefaultAddressParser().ParseAddress(ctx, a)
}

func MustParseAddress(a string) ton.Address {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	addr, err := DefaultAddressParser().ParseAddress(ctx, a)
	if err != nil {
		panic(err)
	}
	return addr
}

// dnsResolver provides a method to resolve a domain name to a list of DNS records.
type dnsResolver interface {
	Resolve(context.Context, string) ([]tlb.DNSRecord, error)
}

func NewAccountAddressParser(resolver dnsResolver) *addressParser {
	return &addressParser{
		resolver: resolver,
	}
}

func (p *addressParser) ParseAddress(ctx context.Context, address string) (ton.Address, error) {
	accountID, err := ton.AccountIDFromRaw(address)
	if err == nil {
		return ton.Address{ID: accountID, Bounce: true, StateInit: nil}, nil
	}
	// ignore the error because we'll try dns in case of error
	bytesAddress, _ := base64.URLEncoding.DecodeString(address)
	if len(bytesAddress) == 36 {
		checksum := uint64(binary.BigEndian.Uint16(bytesAddress[34:36]))
		if checksum == crc.CalculateCRC(crc.XMODEM, bytesAddress[0:34]) {
			testnetOnly := bytesAddress[0]&0b10000000 != 0
			bounce := bytesAddress[0]&0x11 == 0x11
			accountID.Workchain = int32(int8(bytesAddress[1]))
			copy(accountID.Address[:], bytesAddress[2:34])
			return ton.Address{ID: accountID, Bounce: bounce, StateInit: nil, TestnetOnly: testnetOnly}, nil
		}
	}
	if !strings.Contains(address, ".") {
		return ton.Address{}, fmt.Errorf("can't decode address %v", address)
	}
	result, err := p.resolver.Resolve(ctx, address)
	if err != nil {
		return ton.Address{}, err
	}
	account := ton.Address{Bounce: true}
	for _, r := range result {
		if r.SumType == "DNSSmcAddress" {
			accountID, err := ton.AccountIDFromTlb(r.DNSSmcAddress.Address)
			if err != nil {
				return ton.Address{}, err
			}
			if accountID == nil {
				return ton.Address{}, fmt.Errorf("destination account is null")
			}
			account.ID = *accountID
			for _, c := range r.DNSSmcAddress.SmcCapability.Interfaces {
				if c == "wallet" {
					account.Bounce = false
				}
			}
			return account, nil
		}
	}
	return ton.Address{}, fmt.Errorf("address not found")
}

// lazyResolver is a dnsResolver that creates a new dns resolver on the first call to Resolve.
type lazyResolver struct {
	mu  sync.Mutex
	dns dnsResolver
}

func (l *lazyResolver) Resolve(ctx context.Context, s string) ([]tlb.DNSRecord, error) {
	resolver, err := l.resolver()
	if err != nil {
		return nil, err
	}
	return resolver.Resolve(ctx, s)
}

func (l *lazyResolver) resolver() (dnsResolver, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.dns != nil {
		return l.dns, nil
	}
	cli, err := liteapi.NewClient(liteapi.Mainnet(), liteapi.FromEnvs())
	if err != nil {
		return nil, fmt.Errorf("failed to create liteapi client: %w", err)
	}
	l.dns = dns.NewDNS(MustParseAccountID(DefaultRoot), cli)
	return l.dns, nil
}
