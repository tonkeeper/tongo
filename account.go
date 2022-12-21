package tongo

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/snksoft/crc"

	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
	"github.com/startfellows/tongo/utils"
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

type AccountInfo struct {
	Status            AccountStatus
	Balance           uint64
	Data              []byte
	Code              []byte
	FrozenHash        Hash
	LastTransactionLt uint64
}

func AccountIDFromBase64Url(s string) (*AccountID, error) {
	if len(s) == 0 {
		return nil, nil
	}
	var aa AccountID
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(b) != 36 {
		return nil, fmt.Errorf("invalid account 'user friendly' form length: %v", s)
	}
	checksum := uint64(binary.BigEndian.Uint16(b[34:36]))
	if checksum != crc.CalculateCRC(crc.XMODEM, b[0:34]) {
		return nil, fmt.Errorf("invalid checksum")
	}
	aa.Workchain = int32(int8(b[1]))
	copy(aa.Address[:], b[2:34])
	return &aa, nil
}

func AccountIDFromRaw(s string) (*AccountID, error) {
	if len(s) == 0 {
		return nil, nil
	}
	var (
		workchain int32
		address   []byte
		aa        AccountID
	)
	_, err := fmt.Sscanf(s, "%d:%x", &workchain, &address)
	if err != nil {
		return nil, err
	}
	if len(address) != 32 {
		return nil, fmt.Errorf("address len must be 32 bytes")
	}
	aa.Workchain = workchain
	copy(aa.Address[:], address)
	return &aa, nil
}

func ParseAccountID(s string) (*AccountID, error) {
	aa, err := AccountIDFromRaw(s)
	if err != nil {
		aa, err = AccountIDFromBase64Url(s)
		if err != nil {
			return nil, err
		}
	}
	return aa, nil
}

func MustParseAccountID(s string) *AccountID {
	aa, err := ParseAccountID(s)
	if err != nil {
		panic(err)
	}
	return aa
}

func AccountIDFromCell(cell *boc.Cell) (*AccountID, error) {
	var msgAddress MsgAddress
	err := tlb.Unmarshal(cell, &msgAddress)
	if err != nil {
		return nil, err
	}
	return msgAddress.AccountID()
}

// MsgAddressInt
// addr_std$10 anycast:(Maybe Anycast)
// workchain_id:int8 address:bits256  = MsgAddressInt;

func (id AccountID) MarshalTLB(c *boc.Cell, tag string) error {
	err := c.WriteUint(2, 2)
	if err != nil {
		return err
	}
	err = c.WriteBit(false)
	if err != nil {
		return err
	}
	err = c.WriteUint(uint64(id.Workchain), 8)
	if err != nil {
		return err
	}
	err = c.WriteBytes(id.Address[:])
	if err != nil {
		return err
	}
	return nil
}

func MsgAddressFromAccountID(id *AccountID) MsgAddress {
	if id == nil {
		return MsgAddress{SumType: "AddrNone"}
	}
	addr := MsgAddress{
		SumType: "AddrStd",
	}
	addr.AddrStd.WorkchainId = id.Workchain
	addr.AddrStd.Address = id.Address
	addr.AddrStd.Anycast.Null = true
	return addr
}

// AccountStatus
// acc_state_uninit$00 = AccountStatus;
// acc_state_frozen$01 = AccountStatus;
// acc_state_active$10 = AccountStatus;
// acc_state_nonexist$11 = AccountStatus;
type AccountStatus string

const (
	AccountEmpty  AccountStatus = "empty" // empty state from node
	AccountNone   AccountStatus = "nonexist"
	AccountUninit AccountStatus = "uninit"
	AccountActive AccountStatus = "active"
	AccountFrozen AccountStatus = "frozen"
)

func (a AccountStatus) MarshalTLB(c *boc.Cell, tag string) error {
	switch a {
	case AccountUninit:
		return c.WriteUint(0, 2)
	case AccountFrozen:
		return c.WriteUint(1, 2)
	case AccountActive:
		return c.WriteUint(2, 2)
	case AccountNone:
		return c.WriteUint(3, 2)
	}
	return nil
}

func (a *AccountStatus) UnmarshalTLB(c *boc.Cell, tag string) error {
	t, err := c.ReadUint(2)
	if err != nil {
		return err
	}
	switch t {
	case 0:
		*a = AccountUninit
	case 1:
		*a = AccountFrozen
	case 2:
		*a = AccountActive
	case 3:
		*a = AccountNone
	}
	return nil
}

// ShardAccount
// account_descr$_ account:^Account last_trans_hash:bits256
// last_trans_lt:uint64 = ShardAccount;
type ShardAccount struct {
	Account       Account `tlb:"^"`
	LastTransHash Hash
	LastTransLt   uint64
}

// Account
// account_none$0 = Account;
// account$1 addr:MsgAddressInt storage_stat:StorageInfo
// storage:AccountStorage = Account;
type Account struct {
	tlb.SumType
	AccountNone struct {
	} `tlbSumType:"account_none$0"`
	Account struct {
		Addr        MsgAddress
		StorageStat StorageInfo
		Storage     AccountStorage
	} `tlbSumType:"account$1"`
}

// AccountStorage
// account_storage$_ last_trans_lt:uint64
// balance:CurrencyCollection state:AccountState
// = AccountStorage;
type AccountStorage struct {
	LastTransLt uint64
	Balance     CurrencyCollection
	State       AccountState
}

// AccountState
// account_uninit$00 = AccountState;
// account_active$1 _:StateInit = AccountState;
// account_frozen$01 state_hash:bits256 = AccountState;
type AccountState struct {
	tlb.SumType
	AccountUninit struct {
	} `tlbSumType:"account_uninit$00"`
	AccountActive struct {
		StateInit StateInit
	} `tlbSumType:"account_active$1"`
	AccountFrozen struct {
		StateHash Hash
	} `tlbSumType:"account_frozen$01"`
}

// StorageInfo
// storage_info$_ used:StorageUsed last_paid:uint32
// due_payment:(Maybe Grams) = StorageInfo;
type StorageInfo struct {
	Used       StorageUsed
	LastPaid   uint32
	DuePayment tlb.Maybe[Grams]
}

// StorageUsed
// storage_used$_ cells:(VarUInteger 7) bits:(VarUInteger 7)
// public_cells:(VarUInteger 7) = StorageUsed;
type StorageUsed struct {
	Cells       tlb.VarUInteger `tlb:"7bytes"`
	Bits        tlb.VarUInteger `tlb:"7bytes"`
	PublicCells tlb.VarUInteger `tlb:"7bytes"`
}
