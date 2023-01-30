package tongo

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/snksoft/crc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/utils"
	"io"
	"strings"
)

type AccountID struct {
	Workchain int32
	Address   [32]byte
}

func NewAccountId(id int32, addr [32]byte) *AccountID {
	return &AccountID{Workchain: id, Address: addr}
}

func (id AccountID) String() string {
	return id.ToRaw()
}

func (id AccountID) IsZero() bool {
	for i := range id.Address {
		if id.Address[i] != 0 {
			return false
		}
	}
	return true
}

func (id AccountID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.ToRaw())
}

func (id *AccountID) UnmarshalJSON(data []byte) error {
	a, err := ParseAccountID(strings.Trim(string(data), "\"\n "))
	if err != nil {
		return err
	}
	id.Workchain = a.Workchain
	id.Address = a.Address
	return nil
}

func (id AccountID) ToRaw() string {
	return fmt.Sprintf("%v:%x", id.Workchain, id.Address)
}

func (id AccountID) ToHuman(bounce, testnet bool) string {
	prefix := byte(0b00010001)
	if testnet {
		prefix |= 0b10000000
	}
	if !bounce {
		prefix |= 0b01000000
	}
	buf := make([]byte, 36)
	buf[0] = prefix
	buf[1] = byte(id.Workchain)
	copy(buf[2:34], id.Address[:])
	binary.BigEndian.PutUint16(buf[34:36], utils.Crc16(buf[:34]))
	return base64.URLEncoding.EncodeToString(buf)
}

func (id AccountID) MarshalTL() ([]byte, error) {
	payload := make([]byte, 36)
	binary.LittleEndian.PutUint32(payload[:4], uint32(id.Workchain))
	copy(payload[4:36], id.Address[:])
	return payload, nil
}

func (id *AccountID) UnmarshalTL(r io.Reader) error {
	var b [4]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	id.Workchain = int32(binary.LittleEndian.Uint32(b[:]))
	_, err = io.ReadFull(r, id.Address[:])
	return err
}

func (id *AccountID) ToMsgAddress() tlb.MsgAddress {
	if id == nil {
		return tlb.MsgAddress{
			SumType: "AddrNone",
		}
	}
	return tlb.MsgAddress{
		SumType: "AddrStd",
		AddrStd: struct {
			Anycast     tlb.Maybe[tlb.Anycast]
			WorkchainId int8
			Address     tlb.Bits256
		}{
			WorkchainId: int8(id.Workchain),
			Address:     id.Address,
		},
	}
}

type AccountInfo struct {
	Status            tlb.AccountStatus
	Balance           uint64
	Data              []byte
	Code              []byte
	FrozenHash        Bits256
	LastTransactionLt uint64
}

func AccountIDFromBase64Url(s string) (AccountID, error) {
	var aa AccountID
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return AccountID{}, err
	}
	if len(b) != 36 {
		return AccountID{}, fmt.Errorf("invalid account 'user friendly' form length: %v", s)
	}
	checksum := uint64(binary.BigEndian.Uint16(b[34:36]))
	if checksum != crc.CalculateCRC(crc.XMODEM, b[0:34]) {
		return AccountID{}, fmt.Errorf("invalid checksum")
	}
	aa.Workchain = int32(int8(b[1]))
	copy(aa.Address[:], b[2:34])
	return aa, nil
}

func AccountIDFromRaw(s string) (AccountID, error) {
	var (
		workchain int32
		address   []byte
		aa        AccountID
	)
	_, err := fmt.Sscanf(s, "%d:%x", &workchain, &address)
	if err != nil {
		return AccountID{}, err
	}
	if len(address) != 32 {
		return AccountID{}, fmt.Errorf("address len must be 32 bytes")
	}
	aa.Workchain = workchain
	copy(aa.Address[:], address)
	return aa, nil
}

func ParseAccountID(s string) (AccountID, error) {
	aa, err := AccountIDFromRaw(s)
	if err != nil {
		aa, err = AccountIDFromBase64Url(s)
		if err != nil {
			return AccountID{}, err
		}
	}
	return aa, nil
}

func MustParseAccountID(s string) AccountID {
	aa, err := ParseAccountID(s)
	if err != nil {
		panic(err)
	}
	return aa
}

func MsgAddressFromAccountID(id *AccountID) tlb.MsgAddress {
	if id == nil {
		return tlb.MsgAddress{SumType: "AddrNone"}
	}
	addr := tlb.MsgAddress{
		SumType: "AddrStd",
	}
	addr.AddrStd.WorkchainId = int8(id.Workchain)
	addr.AddrStd.Address = id.Address
	return addr
}

func GetAccountInfo(a tlb.Account) (AccountInfo, error) {
	if a.SumType == "AccountNone" {
		return AccountInfo{Status: tlb.AccountNone}, nil
	}
	res := AccountInfo{
		Balance:           uint64(a.Account.Storage.Balance.Grams),
		LastTransactionLt: a.Account.Storage.LastTransLt,
	}
	if a.Account.Storage.State.SumType == "AccountUninit" {
		res.Status = tlb.AccountUninit
		return res, nil
	}
	if a.Account.Storage.State.SumType == "AccountFrozen" {
		res.FrozenHash = Bits256(a.Account.Storage.State.AccountFrozen.StateHash)
		res.Status = tlb.AccountFrozen
		return res, nil
	}
	res.Status = tlb.AccountActive
	if a.Account.Storage.State.AccountActive.StateInit.Data.Exists {
		data, err := a.Account.Storage.State.AccountActive.StateInit.Data.Value.Value.ToBoc()
		if err != nil {
			return AccountInfo{}, err
		}
		res.Data = data
	}
	if a.Account.Storage.State.AccountActive.StateInit.Code.Exists {
		code, err := a.Account.Storage.State.AccountActive.StateInit.Code.Value.Value.ToBoc()
		if err != nil {
			return AccountInfo{}, err
		}
		res.Code = code
	}
	return res, nil
}

// TODO: replace pointer with nullable type
func AccountIDFromTlb(a tlb.MsgAddress) (*AccountID, error) {
	switch a.SumType {
	case "AddrNone":
		return nil, nil
	case "AddrStd":
		return &AccountID{Workchain: int32(a.AddrStd.WorkchainId), Address: a.AddrStd.Address}, nil
	}
	return nil, fmt.Errorf("can not convert not std address to AccountId")
}
