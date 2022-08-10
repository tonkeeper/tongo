package tongo

import (
	"fmt"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
)

// DNSRecordSet
// _ (HashmapE 256 DNSRecord) = DNS_RecordSet;
type DNSRecordSet struct {
	Records tlb.HashmapE[tlb.Ref[DNSRecord]] `tlb:"256bits"`
}

// DNSRecord
// dns_text#1eda _:Text = DNSRecord;
// dns_next_resolver#ba93 resolver:MsgAddressInt = DNSRecord;  // usually in record #-1
// dns_adnl_address#ad01 adnl_addr:bits256 flags:(## 8) { flags <= 1 }
// proto_list:flags . 0?ProtoList = DNSRecord;  // often in record #2
// dns_smc_address#9fd3 smc_addr:MsgAddressInt flags:(## 8) { flags <= 1 }
// cap_list:flags . 0?SmcCapList = DNSRecord;   // often in record #1
type DNSRecord struct {
	tlb.SumType
	DNSText         DNSText   `tlbSumType:"dns_text#1eda"`
	DNSNextResolver AccountID `tlbSumType:"dns_next_resolver#ba93"`
	DNSAdnlAddress  struct {
		Address   [32]byte
		ProtoList []string
	} `tlbSumType:"dns_adnl_address#ad01"`
	DNSSmcAddress struct {
		Address       AccountID
		SmcCapability SmcCapabilities
	} `tlbSumType:"dns_smc_address#9fd3"`
}

func (r *DNSRecord) UnmarshalTLB(c *boc.Cell, tag string) error {
	t, err := c.ReadUint(16)
	if err != nil {
		return err
	}
	switch t {
	case 0x1eda: // dns_text#1eda _:Text = DNSRecord;
		var text DNSText
		err := tlb.Unmarshal(c, &text)
		if err != nil {
			return err
		}
		r.SumType = "DNSText"
		r.DNSText = text
		return nil
	case 0xba93: // dns_next_resolver#ba93 resolver:MsgAddressInt = DNSRecord;  // usually in record #-1
		var msgAddr MsgAddress
		err := tlb.Unmarshal(c, &msgAddr)
		if err != nil {
			return err
		}
		addr, err := msgAddr.AccountId()
		if err != nil {
			return err
		}
		if addr == nil {
			return fmt.Errorf("nil next resolver address")
		}
		r.SumType = "DNSNextResolver"
		r.DNSNextResolver = *addr
		return nil
	case 0xad01:
		res, err := readDnsAdnlAddress(c)
		if err != nil {
			return err
		}
		*r = res
		return nil
	case 0x9fd3: //dns_smc_address#9fd3 smc_addr:MsgAddressInt flags:(## 8) { flags <= 1 }
		// cap_list:flags . 0?SmcCapList = DNSRecord;   // often in record #1
		res, err := readDNSSmcAddress(c)
		if err != nil {
			return err
		}
		*r = res
		return nil
	}
	return nil
}

// DNSText
// text$_ chunks:(## 8) rest:(TextChunks chunks) = Text;
type DNSText string

func (t *DNSText) UnmarshalTLB(c *boc.Cell, tag string) error {
	chunksQty, err := c.ReadUint(8)
	if err != nil {
		return err
	}
	res, err := readChunks(c, int(chunksQty))
	if err != nil {
		return err
	}
	*t = DNSText(res)
	return nil
}

// chunk_ref$_ {n:#} ref:^(TextChunks (n + 1)) = TextChunkRef (n + 1);
// chunk_ref_empty$_ = TextChunkRef 0;
// text_chunk$_ {n:#} len:(## 8) data:(bits (len * 8)) next:(TextChunkRef n) = TextChunks (n + 1);
// text_chunk_empty$_ = TextChunks 0;
func readChunks(c *boc.Cell, chunksQty int) (string, error) {
	if chunksQty == 0 {
		return "", nil
	}
	ln, err := c.ReadUint(8)
	if err != nil {
		return "", err
	}
	data, err := c.ReadBytes(int(ln))
	if err != nil {
		return "", err
	}
	res := string(data)
	if chunksQty > 1 {
		next, err := c.NextRef()
		if err != nil {
			return "", err
		}
		nextChunk, err := readChunks(next, chunksQty-1)
		if err != nil {
			return "", err
		}
		res = res + nextChunk
	}
	return res, nil
}

// dns_adnl_address#ad01 adnl_addr:bits256 flags:(## 8) { flags <= 1 } // ad01 tag reading at upper level
// proto_list:flags . 0?ProtoList = DNSRecord;  // often in record #2
// proto_list_nil$0 = ProtoList;
// proto_list_next$1 head:Protocol tail:ProtoList = ProtoList;
// proto_http#4854 = Protocol;
func readDnsAdnlAddress(c *boc.Cell) (DNSRecord, error) {
	addr, err := c.ReadBytes(32)
	if err != nil {
		return DNSRecord{}, err
	}
	flags, err := c.ReadUint(8)
	if err != nil {
		return DNSRecord{}, err
	}
	if flags > 2 {
		return DNSRecord{}, fmt.Errorf("invalid dns_adnl_address flags")
	}
	var protoList []string
	if flags > 0 {
		next, err := c.ReadBit()
		if err != nil {
			return DNSRecord{}, err
		}
		for next {
			t, err := c.ReadUint(16)
			if err != nil {
				return DNSRecord{}, err
			}
			switch t {
			case 0x4854:
				protoList = append(protoList, "http")
			}
			next, err = c.ReadBit()
			if err != nil {
				return DNSRecord{}, err
			}
		}
	}
	var res DNSRecord
	res.SumType = "DNSAdnlAddress"
	copy(res.DNSAdnlAddress.Address[:], addr[:])
	res.DNSAdnlAddress.ProtoList = protoList
	return res, nil
}

// SmcCapabilities
// Reorganized SmcCapList type
// cap_list_nil$0 = SmcCapList;
// cap_list_next$1 head:SmcCapability tail:SmcCapList = SmcCapList;
// cap_method_seqno#5371 = SmcCapability;
// cap_method_pubkey#71f4 = SmcCapability;
// cap_is_wallet#2177 = SmcCapability;
// cap_name#ff name:Text = SmcCapability;
type SmcCapabilities struct {
	Name       []string
	Interfaces []string
}

func readDNSSmcAddress(c *boc.Cell) (DNSRecord, error) {
	var a MsgAddress
	err := tlb.Unmarshal(c, &a)
	if err != nil {
		return DNSRecord{}, err
	}
	addr, err := a.AccountId()
	if err != nil {
		return DNSRecord{}, err
	}
	if addr == nil {
		return DNSRecord{}, fmt.Errorf("nil smc_addr")
	}
	flags, err := c.ReadUint(8)
	if err != nil {
		return DNSRecord{}, err
	}
	if flags > 2 {
		return DNSRecord{}, fmt.Errorf("invalid smc_addr flags")
	}

	var capabilities SmcCapabilities
	var capability struct {
		tlb.SumType
		CapMethodSeqno  struct{} `tlbSumType:"cap_method_seqno#5371"`
		CapMethodPubkey struct{} `tlbSumType:"cap_method_pubkey#71f4"`
		CapIsWallet     struct{} `tlbSumType:"cap_is_wallet#2177"`
		CapName         DNSText  `tlbSumType:"cap_name#ff"`
	}

	if flags > 0 {
		next, err := c.ReadBit()
		if err != nil {
			return DNSRecord{}, err
		}
		for next {
			err = tlb.Unmarshal(c, &capability)
			if err != nil {
				return DNSRecord{}, err
			}
			switch capability.SumType {
			case "CapMethodSeqno":
				capabilities.Interfaces = append(capabilities.Interfaces, "seqno")
			case "CapMethodPubkey":
				capabilities.Interfaces = append(capabilities.Interfaces, "pubkey")
			case "CapIsWallet":
				capabilities.Interfaces = append(capabilities.Interfaces, "wallet")
			case "CapName":
				capabilities.Name = append(capabilities.Name, string(capability.CapName))
			}
			next, err = c.ReadBit()
			if err != nil {
				return DNSRecord{}, err
			}
		}
	}
	var res DNSRecord
	res.SumType = "DNSSmcAddress"
	res.DNSSmcAddress.Address = *addr
	res.DNSSmcAddress.SmcCapability = capabilities
	return res, nil
}
