package tlb

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
)

// DNSRecordSet
// _ (HashmapE 256 DNSRecord) = DNS_RecordSet;
type DNSRecordSet struct {
	Records Hashmap[Bits256, Ref[DNSRecord]]
}

// DNSRecord
// dns_text#1eda _:Text = DNSRecord;
// dns_next_resolver#ba93 resolver:MsgAddressInt = DNSRecord;  // usually in record #-1
// dns_adnl_address#ad01 adnl_addr:bits256 flags:(## 8) { flags <= 1 }
// proto_list:flags . 0?ProtoList = DNSRecord;  // often in record #2
// dns_smc_address#9fd3 smc_addr:MsgAddressInt flags:(## 8) { flags <= 1 }
// cap_list:flags . 0?SmcCapList = DNSRecord;   // often in record #1
// dns_storage_address#7473 bag_id:bits256 = DNSRecord;
type DNSRecord struct {
	SumType
	DNSText         DNSText    `tlbSumType:"dns_text#1eda"`
	DNSNextResolver MsgAddress `tlbSumType:"dns_next_resolver#ba93"`
	DNSAdnlAddress  struct {
		Address   Bits256
		ProtoList []string
	} `tlbSumType:"dns_adnl_address#ad01"`
	DNSSmcAddress struct {
		Address       MsgAddress
		SmcCapability SmcCapabilities
	} `tlbSumType:"dns_smc_address#9fd3"`
	DNSStorageAddress Bits256   `tlbSumType:"dns_storage_address#7473"`
	NotStandard       *boc.Cell // only for custom unmarshaling
}

func (r DNSRecord) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	switch r.SumType {
	case "DNSText":
		err := c.WriteUint(0x1eda, 16) // dns_text#1eda
		if err != nil {
			return err
		}
		return Marshal(c, r.DNSText)
	case "DNSNextResolver":
		err := c.WriteUint(0xba93, 16) // dns_next_resolver#ba93
		if err != nil {
			return err
		}
		return Marshal(c, r.DNSNextResolver)
	case "DNSAdnlAddress":
		err := c.WriteUint(0xad01, 16) // dns_adnl_address#ad01
		if err != nil {
			return err
		}
		err = c.WriteBytes(r.DNSAdnlAddress.Address[:])
		if err != nil {
			return err
		}
		flags := uint8(len(r.DNSAdnlAddress.ProtoList))
		err = c.WriteUint(uint64(flags), 8)
		if err != nil {
			return err
		}
		for _, proto := range r.DNSAdnlAddress.ProtoList {
			switch proto {
			case "http":
				err = c.WriteUint(0x4854, 16) // proto_http#4854
				if err != nil {
					return err
				}
			}
		}
	case "DNSSmcAddress":
		err := c.WriteUint(0x9fd3, 16) // dns_smc_address#9fd3
		if err != nil {
			return err
		}
		err = Marshal(c, r.DNSSmcAddress.Address)
		if err != nil {
			return err
		}
		flags := uint8(len(r.DNSSmcAddress.SmcCapability.Interfaces))
		err = c.WriteUint(uint64(flags), 8)
		if err != nil {
			return err
		}
		for _, iface := range r.DNSSmcAddress.SmcCapability.Interfaces {
			switch iface {
			case "seqno":
				err = c.WriteUint(0x5371, 16) // cap_method_seqno#5371
				if err != nil {
					return err
				}
			case "pubkey":
				err = c.WriteUint(0x71f4, 16) // cap_method_pubkey#71f4
				if err != nil {
					return err
				}
			case "wallet":
				err = c.WriteUint(0x2177, 16) // cap_is_wallet#2177
				if err != nil {
					return err
				}
			}
		}
		for _, name := range r.DNSSmcAddress.SmcCapability.Name {
			err = c.WriteUint(0xff, 8) // cap_name#ff
			if err != nil {
				return err
			}
			err = Marshal(c, DNSText(name))
			if err != nil {
				return err
			}
		}
	case "DNSStorageAddress":
		err := c.WriteUint(0x7473, 16) // dns_storage_address#7473
		if err != nil {
			return err
		}
		return c.WriteBytes(r.DNSStorageAddress[:])
	case "NotStandard":
		c = r.NotStandard.CopyRemaining()
	}
	return nil
}

func (t DNSText) MarshalTLB(c *boc.Cell, encoder *Encoder) error {
	text := string(t)
	if len(text) == 0 {
		return c.WriteUint(0, 8) // chunks:(## 8) = 0
	}

	chunks := (len(text) + 7) / 8
	err := c.WriteUint(uint64(chunks), 8)
	if err != nil {
		return err
	}

	chunkSize := min(len(text), 8)
	err = c.WriteUint(uint64(chunkSize), 8)
	if err != nil {
		return err
	}
	err = c.WriteBytes([]byte(text[:chunkSize]))
	if err != nil {
		return err
	}

	remaining := text[chunkSize:]
	for len(remaining) > 0 {
		newCell := boc.NewCell()
		chunkSize = min(len(remaining), 8)
		err = newCell.WriteUint(uint64(chunkSize), 8)
		if err != nil {
			return err
		}
		err = newCell.WriteBytes([]byte(remaining[:chunkSize]))
		if err != nil {
			return err
		}
		err = c.AddRef(newCell)
		if err != nil {
			return err
		}
		remaining = remaining[chunkSize:]
		c = newCell
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r *DNSRecord) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
	t, err := c.ReadUint(16)
	if err != nil {
		return err
	}
	switch t {
	case 0x1eda: // dns_text#1eda _:Text = DNSRecord;
		var text DNSText
		err := Unmarshal(c, &text)
		if err != nil {
			return err
		}
		r.SumType = "DNSText"
		r.DNSText = text
		return nil
	case 0xba93: // dns_next_resolver#ba93 resolver:MsgAddressInt = DNSRecord;  // usually in record #-1
		var msgAddr MsgAddress
		err := Unmarshal(c, &msgAddr)
		if err != nil {
			return err
		}
		r.SumType = "DNSNextResolver"
		r.DNSNextResolver = msgAddr
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
	case 0x7473: // dns_storage_address#7473 bag_id:bits256 = DNSRecord;
		addr, err := c.ReadBytes(32)
		if err != nil {
			return err
		}
		r.SumType = "DNSStorageAddress"
		copy(r.DNSStorageAddress[:], addr[:])
		return nil
	}
	c.ResetCounters()
	r.SumType = "NotStandard"
	r.NotStandard = c
	return nil
}

// DNSText
// text$_ chunks:(## 8) rest:(TextChunks chunks) = Text;
type DNSText string

func (t *DNSText) UnmarshalTLB(c *boc.Cell, decoder *Decoder) error {
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
	var addr MsgAddress
	err := Unmarshal(c, &addr)
	if err != nil {
		return DNSRecord{}, err
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
		SumType
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
			err = Unmarshal(c, &capability)
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
	res.DNSSmcAddress.Address = addr
	res.DNSSmcAddress.SmcCapability = capabilities
	return res, nil
}
