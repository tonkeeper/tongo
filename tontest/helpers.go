package tontest

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type accountBuilder struct {
	state      *tlb.AccountStatus
	balance    *tlb.Grams
	code       *boc.Cell
	data       *boc.Cell
	lastLT     *uint64
	lastHash   *tlb.Bits256
	workchain  *int32
	address    *ton.AccountID
	frozenHash *tlb.Bits256
}

func Account() accountBuilder {
	return accountBuilder{}
}

func (b accountBuilder) Address(a ton.AccountID) accountBuilder {
	b.address = &a
	return b
}

func (b accountBuilder) Balance(grams tlb.Grams) accountBuilder {
	b.balance = &grams
	return b
}

func (b accountBuilder) State(state tlb.AccountStatus) accountBuilder {
	b.state = &state
	return b
}

func (b accountBuilder) StateInit(code, data *boc.Cell) accountBuilder {
	b.code = code
	b.data = data
	return b
}

func (b accountBuilder) Last(lt uint64, hash tlb.Bits256) accountBuilder {
	b.lastLT = &lt
	b.lastHash = &hash
	return b
}

func (b accountBuilder) MustShardAccount() tlb.ShardAccount {
	a, err := b.ShardAccount()
	if err != nil {
		panic(err)
	}
	return a
}

func (b accountBuilder) ShardAccount() (tlb.ShardAccount, error) {
	// todo: conflicts checks
	// frozen without hash
	// frozen + data/code
	// uninit + data/code
	//etc

	if (b.state == nil || *b.state == tlb.AccountNone) &&
		b.balance == nil &&
		b.code == nil &&
		b.data == nil &&
		b.lastHash == nil &&
		b.lastLT == nil {
		return tlb.ShardAccount{
			Account: tlb.Account{
				SumType: "AccountNone",
			},
		}, nil
	}
	var address *ton.AccountID
	if b.address != nil {
		address = b.address
	}
	if address == nil && (b.code != nil || b.data != nil) {
		var workchain int32
		if b.workchain != nil {
			workchain = *b.workchain
		}
		var err error
		address, err = initToAddress(workchain, b.code, b.data)
		if err != nil {
			return tlb.ShardAccount{}, err
		}
	}
	if address == nil {
		return tlb.ShardAccount{}, fmt.Errorf("account should have address or stateInit")
	}
	var balance tlb.Grams
	if b.balance != nil {
		balance = *b.balance
	}
	lastLt := uint64(100500)
	if b.lastLT != nil {
		lastLt = *b.lastLT
	}
	state := tlb.AccountState{
		SumType: "AccountUninit",
	}
	if b.state != nil {
		if *b.state == tlb.AccountActive {
			state.SumType = "AccountActive"
			if b.code != nil {
				state.AccountActive.StateInit.Code.Exists = true
				state.AccountActive.StateInit.Code.Value.Value = *b.code
			}
			if b.data != nil {
				state.AccountActive.StateInit.Data.Exists = true
				state.AccountActive.StateInit.Data.Value.Value = *b.data
			}
		}
		//todo: frozen
	}
	account := tlb.ShardAccount{
		Account: tlb.Account{
			SumType: "Account",
			Account: tlb.ExistedAccount{
				Addr: address.ToMsgAddress(),
				StorageStat: tlb.StorageInfo{
					LastPaid: uint32(time.Now().Unix()),
					StorageExtra: tlb.StorageExtraInfo{
						SumType: "StorageExtraNone",
					},
				},
				Storage: tlb.AccountStorage{
					LastTransLt: lastLt,
					Balance: tlb.CurrencyCollection{
						Grams: balance,
					},
					State: state,
				},
			},
		},
		LastTransLt: lastLt - 1, //todo: investigate lt number
	}
	if b.lastHash != nil {
		account.LastTransHash = *b.lastHash
	} else {
		rand.Read(account.LastTransHash[:])
	}

	return account, nil
}

func MustAnyToCell(a any) *boc.Cell {
	c, err := AnyToCell(a)
	if err != nil {
		panic(err)
	}
	return c
}

func AnyToCell(a any) (*boc.Cell, error) {
	cell := boc.NewCell()
	var err error
	switch c := a.(type) {
	case []byte:
		cells, err := boc.DeserializeBoc(c)
		if err != nil {
			return nil, err
		}
		return cells[0], nil
	case string:
		return boc.DeserializeSinglRootBase64(c)
	default:
		err = tlb.Marshal(cell, a)
		return cell, err
	}
}

type MsgType int

const External MsgType = 0
const Internal MsgType = 1

type messageBuilder struct {
	msgType          MsgType
	to               *ton.AccountID
	from             ton.AccountID
	value            tlb.Grams
	bounce           bool
	bounced          bool
	body, data, code *boc.Cell
}

func NewMessage() messageBuilder {
	return messageBuilder{}
}

func (b messageBuilder) To(a ton.AccountID) messageBuilder {
	b.to = &a
	return b
}
func (b messageBuilder) Internal(from ton.AccountID, value tlb.Grams) messageBuilder {
	b.msgType = Internal
	b.value = value
	b.from = from
	return b
}

func (b messageBuilder) WithBody(body *boc.Cell) messageBuilder {
	b.body = body
	return b
}

func (b messageBuilder) WithInit(code, data *boc.Cell) messageBuilder {
	b.code = code
	b.data = data
	return b
}

func (b messageBuilder) MustMessage() tlb.Message {
	m, err := b.Message()
	if err != nil {
		panic(err)
	}
	return m
}
func (b messageBuilder) Message() (tlb.Message, error) {
	var m tlb.Message
	if b.code != nil {
		m.Init.Exists = true
		m.Init.Value.Value.Code.Exists = true
		m.Init.Value.Value.Code.Value.Value = *b.code
	}
	if b.data != nil {
		m.Init.Exists = true
		m.Init.Value.Value.Data.Exists = true
		m.Init.Value.Value.Data.Value.Value = *b.data
	}
	if b.to == nil && !m.Init.Exists {
		return tlb.Message{}, fmt.Errorf("can't determine address for message")
	} else if b.to == nil && m.Init.Exists {
		var err error
		b.to, err = initToAddress(0, b.code, b.data)
		if err != nil {
			return tlb.Message{}, err
		}
	}

	switch b.msgType {
	case External:
		m.Info.SumType = "ExtInMsgInfo"
		m.Info.ExtInMsgInfo = &struct {
			Src, Dest tlb.MsgAddress
			ImportFee tlb.VarUInteger16
		}{Src: tlb.MsgAddress{SumType: "AddrNone"}, Dest: b.to.ToMsgAddress()}
	case Internal:
		m.Info.SumType = "IntMsgInfo"
		m.Info.IntMsgInfo = &struct {
			IhrDisabled bool
			Bounce      bool
			Bounced     bool
			Src         tlb.MsgAddress
			Dest        tlb.MsgAddress
			Value       tlb.CurrencyCollection
			IhrFee      tlb.Grams
			FwdFee      tlb.Grams
			CreatedLt   uint64
			CreatedAt   uint32
		}{
			IhrDisabled: true,
			Bounce:      b.bounce,
			Bounced:     b.bounced,
			Src:         b.from.ToMsgAddress(),
			Dest:        b.to.ToMsgAddress(),
			Value:       tlb.CurrencyCollection{Grams: b.value},
			IhrFee:      0,
			FwdFee:      0,
			CreatedLt:   100500,
			CreatedAt:   uint32(time.Now().Unix())}
	}
	if b.body != nil {
		m.Body.Value = tlb.Any(*b.body)
	}

	return m, nil
}

func initToAddress(workchain int32, code, data *boc.Cell) (*ton.AccountID, error) {
	initState := boc.NewCell()
	s := tlb.StateInit{}
	if code != nil {
		s.Code.Exists = true
		s.Code.Value.Value = *code
	}
	if data != nil {
		s.Data.Exists = true
		s.Data.Value.Value = *data
	}

	err := tlb.Marshal(initState, s)
	if err != nil {
		return nil, err
	}
	hash, err := initState.Hash()
	if err != nil {
		return nil, err
	}
	a := ton.AccountID{Workchain: workchain}
	copy(a.Address[:], hash)
	code.ResetCounters()
	data.ResetCounters()
	return &a, nil
}
