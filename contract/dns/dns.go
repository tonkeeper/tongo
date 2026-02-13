package dns

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type executor interface {
	RunSmcMethodByID(context.Context, ton.AccountID, int, tlb.VmStack) (uint32, tlb.VmStack, error)
}

var (
	ErrNotResolved = errors.New("not resolved")
)

func sha256str(value string) DNSCategory {
	sum := sha256.Sum256([]byte(value))
	return DNSCategory(hex.EncodeToString(sum[:]))
}

type DNSCategory string

var (
	DNSCategoryDNSNextResolver = sha256str("dns_next_resolver")
	DNSCategoryWallet          = sha256str("wallet")
	DNSCategorySite            = sha256str("site")
	DNSCategoryStorage         = sha256str("storage")
	DNSCategoryPicture         = sha256str("picture")
	DNSCategoryLinks           = sha256str("links")
	DNSCategoryDescription     = sha256str("description")
)

type DNS struct {
	root     ton.AccountID
	executor executor
}

// NewDNS
// If root == nil then use root from network config
func NewDNS(root ton.AccountID, e executor) *DNS {
	return &DNS{
		root:     root,
		executor: e,
	}
}

func (d *DNS) Resolve(ctx context.Context, domain string) (map[DNSCategory]tlb.DNSRecord, error) {
	if d.executor == nil {
		return nil, errors.New("blockchain interface is nil")
	}
	if domain == "" {
		domain = "."
	}
	dom := convertDomain(domain)
	return d.resolve(ctx, d.root, []byte(dom))
}

func (d *DNS) resolve(ctx context.Context, resolver ton.AccountID, dom []byte) (map[DNSCategory]tlb.DNSRecord, error) {
	n := int64(len(dom))
	stack := tlb.VmStack{}
	val, err := tlb.TlbStructToVmCellSlice(dom)
	if err != nil {
		return nil, err
	}
	stack.Put(val)
	stack.Put(tlb.VmStackValue{SumType: "VmStkInt", VmStkInt: tlb.Int257{}})
	exitCode, stack, err := d.executor.RunSmcMethodByID(ctx, resolver, 123660, stack)
	if err != nil && strings.Contains(err.Error(), "method execution failed") {
		return nil, fmt.Errorf("%w: %v", ErrNotResolved, err)
	}
	if err != nil {
		return nil, err
	}
	if !(exitCode == 0 || exitCode == 1) {
		return nil, fmt.Errorf("%w: invalid exit code %v", ErrNotResolved, exitCode)
	}
	var result struct {
		ResolvedBits int64
		Result       boc.Cell
	}
	if len(stack) == 2 && stack[0].SumType == "VmStkTinyInt" && stack[0].VmStkTinyInt == 0 && stack[1].SumType == "VmStkNull" {
		return nil, ErrNotResolved
	}
	err = stack.Unmarshal(&result)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNotResolved, err)
	}
	if result.ResolvedBits&0b111 != 0 {
		return nil, fmt.Errorf("%w: invalid qty of resolved bits", ErrNotResolved)
	}
	if result.ResolvedBits == 0 {
		return nil, ErrNotResolved
	}
	if result.ResolvedBits/8 == n {
		var recordSet tlb.DNSRecordSet
		err = tlb.Unmarshal(&result.Result, &recordSet)
		if err != nil {
			return nil, err
		}
		records := make(map[DNSCategory]tlb.DNSRecord, len(recordSet.Records.Values()))
		for _, key := range recordSet.Records.Keys() {
			value, ok := recordSet.Records.Get(key)
			if !ok {
				continue
			}
			category := DNSCategory(hex.EncodeToString(key[:]))
			records[category] = value.Value
		}
		return records, nil
	}
	var record tlb.DNSRecord
	err = tlb.Unmarshal(&result.Result, &record)
	if err != nil {
		return nil, err
	}
	if record.SumType != "DNSNextResolver" {
		return nil, fmt.Errorf("should be next resolver")
	}
	account, err := ton.AccountIDFromTlb(record.DNSNextResolver)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("invalid next resolver")
	}
	return d.resolve(ctx, *account, dom[result.ResolvedBits/8:])

}

func convertDomain(domain string) string {
	domains := strings.Split(domain, ".")
	for i, j := 0, len(domains)-1; i < j; i, j = i+1, j-1 { // reverse array
		domains[i], domains[j] = domains[j], domains[i]
	}
	return strings.Join(domains, "\x00") + "\x00"
}
