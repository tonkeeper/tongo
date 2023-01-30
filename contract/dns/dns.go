package dns

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"math/big"
	"strings"
)

type blockchain interface {
	DnsResolve(ctx context.Context, address tongo.AccountID, domain string, category *big.Int) (int, *boc.Cell, error)
	GetRootDNS(ctx context.Context) (tongo.AccountID, error)
}

var (
	ErrNotResolved = fmt.Errorf("not resolved")
)

type DNS struct {
	Root       tongo.AccountID
	blockchain blockchain
}

// NewDNS
// If root == nil then use root from network config
func NewDNS(root *tongo.AccountID, blockchain blockchain) (*DNS, error) {
	var (
		err error
		r   tongo.AccountID
	)
	if root == nil {
		if blockchain == nil {
			return nil, tongo.BlockchainInterfaceIsNil
		}
		r, err = blockchain.GetRootDNS(context.Background())
		if err != nil {
			return nil, err
		}
	} else {
		r = *root
	}
	return &DNS{
		Root:       r,
		blockchain: blockchain,
	}, nil
}

func (d *DNS) Resolve(ctx context.Context, domain string) ([]tlb.DNSRecord, error) {
	if d.blockchain == nil {
		return nil, tongo.BlockchainInterfaceIsNil
	}
	if domain == "" {
		domain = "."
	}
	dom := convertDomain(domain)
	r, err := d.resolve(ctx, d.Root, dom)
	if err != nil && strings.Contains(err.Error(), "method execution failed") {
		return nil, ErrNotResolved
	}
	return r, err
}

func (d *DNS) resolve(ctx context.Context, resolver tongo.AccountID, dom string) ([]tlb.DNSRecord, error) {
	n := len(dom)
	i, res, err := d.blockchain.DnsResolve(ctx, resolver, dom, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	if i%8 != 0 {
		return nil, fmt.Errorf("invalid qty of resolved bits")
	}
	if i/8 == 0 { // not resolved
		return nil, nil
	} else if i/8 == n { // resolved
		return parseDnsRecords(res)
	} // m < n partial resolved
	rec, err := parseDnsRecords(res)
	if err != nil {
		return nil, err
	}
	if len(rec) != 1 {
		return nil, fmt.Errorf("must be only one record for partial resolved")
	}
	if rec[0].SumType != "DNSNextResolver" {
		return nil, fmt.Errorf("must be only next resolver record for partial resolved")
	}
	account, err := tongo.AccountIDFromTlb(rec[0].DNSNextResolver)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("nil account id")
	}
	return d.resolve(ctx, *account, string([]byte(dom)[i/8:]))
}

func convertDomain(domain string) string {
	domains := strings.Split(domain, ".")
	for i, j := 0, len(domains)-1; i < j; i, j = i+1, j-1 { // reverse array
		domains[i], domains[j] = domains[j], domains[i]
	}
	return strings.Join(domains, "\x00") + "\x00"
}

func parseDnsRecords(c *boc.Cell) ([]tlb.DNSRecord, error) {
	var record tlb.DNSRecord
	err := tlb.Unmarshal(c, &record)
	if err == nil && record.SumType == "DNSNextResolver" {
		return []tlb.DNSRecord{record}, nil
	}
	c.ResetCounters()
	var records tlb.DNSRecordSet
	c2 := boc.NewCell()
	_ = c2.WriteBit(true)
	_ = c2.AddRef(c)
	err = tlb.Unmarshal(c2, &records)
	if err != nil {
		return nil, err
	}
	var res []tlb.DNSRecord
	for _, r := range records.Records.Values() {
		res = append(res, r.Value)
	}
	return res, nil
}
