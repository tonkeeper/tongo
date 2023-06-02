package account

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/snksoft/crc"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/contract/dns"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
)

type executor interface {
	RunSmcMethodByID(context.Context, ID, int, tlb.VmStack) (uint32, tlb.VmStack, error)
}

type parser struct {
	root     tongo.AccountID
	executor executor
}

type ID = tongo.AccountID

type Address struct {
	ID
	Bounce    bool
	StateInit *tlb.StateInit
}

func (p *parser) Root(root ID) *parser {
	p.root = root
	return p
}

func (p *parser) Executor(executor executor) *parser {
	p.executor = executor
	return p
}

var DefaultParser *parser

func init() {
	DefaultParser = NewAccountParser()
}

func Parser(a string) (Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return DefaultParser.ParseAddress(ctx, a)
}

func NewAccountParser() *parser {
	defaultRoot := tongo.MustParseAccountID("-1:e56754f83426f69b09267bd876ac97c44821345b7e266bd956a7bfbfb98df35c")
	return &parser{
		root: defaultRoot,
	}
}

func (p *parser) ParseAddress(ctx context.Context, address string) (Address, error) {

	accountID, err := tongo.AccountIDFromRaw(address)
	if err == nil {
		return Address{ID: accountID, Bounce: true, StateInit: nil}, nil
	}
	bytesAddress, _ := base64.URLEncoding.DecodeString(address)
	if len(bytesAddress) == 36 {
		checksum := uint64(binary.BigEndian.Uint16(bytesAddress[34:36]))
		if checksum == crc.CalculateCRC(crc.XMODEM, bytesAddress[0:34]) {
			bounce := bytesAddress[0]&0x11 == 0x11
			accountID.Workchain = int32(int8(bytesAddress[1]))
			copy(accountID.Address[:], bytesAddress[2:34])
			return Address{ID: accountID, Bounce: bounce, StateInit: nil}, nil
		}
	}
	if !strings.Contains(address, ".") {
		return Address{}, fmt.Errorf("can't decode address %v", address)
	}
	if p.executor == nil {
		var err error
		p.executor, err = liteapi.NewClientWithDefaultMainnet()
		if err != nil {
			return Address{}, err
		}
	}
	newDns := dns.NewDNS(p.root, p.executor) // import cycle in package/dns.go
	result, err := newDns.Resolve(ctx, address)
	if err != nil {
		return Address{}, err
	}
	account := Address{Bounce: true}
	for _, r := range result {
		if r.SumType == "DNSSmcAddress" {
			accountID, err := tongo.AccountIDFromTlb(r.DNSSmcAddress.Address)
			if err != nil {
				return Address{}, err
			}
			if accountID == nil {
				return Address{}, fmt.Errorf("destination account is null")
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
	return Address{}, fmt.Errorf("address not found")
}
