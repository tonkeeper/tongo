package func_testing

import (
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"math/rand"
	"time"
)

func fakeUninitAccount(addr tongo.AccountID, balance tlb.Grams) tlb.ShardAccount {
	a := tlb.ShardAccount{
		Account: tlb.Account{
			SumType: "Account",
			Account: struct {
				Addr        tlb.MsgAddress
				StorageStat tlb.StorageInfo
				Storage     tlb.AccountStorage
			}{
				Addr: addr.ToMsgAddress(),
				StorageStat: tlb.StorageInfo{
					Used:     tlb.StorageUsed{},
					LastPaid: uint32(time.Now().Unix()),
				},
				Storage: tlb.AccountStorage{
					LastTransLt: 10005001,
					Balance: tlb.CurrencyCollection{
						Grams: balance,
					},
					State: tlb.AccountState{
						SumType: "AccountUninit",
					},
				},
			},
		},
		LastTransLt: 10005000,
	}
	rand.Read(a.LastTransHash[:])
	return a
}

func externalMessageToNewAccount(code, data, body any) (tlb.Message, tongo.AccountID, error) {
	codeCell := boc.NewCell()
	dataCell := boc.NewCell()
	bodyCell := boc.NewCell()
	var err error
	switch c := code.(type) {
	case []byte:
		cells, err := boc.DeserializeBoc(c)
		if err != nil {
			return tlb.Message{}, tongo.AccountID{}, err
		}
		codeCell = cells[0]
	case string:
		codeCell, err = boc.DeserializeSinglRootBase64(c)
	default:
		err = tlb.Marshal(codeCell, code)
		if err != nil {
			return tlb.Message{}, tongo.AccountID{}, err
		}
	}
	switch c := data.(type) {
	case []byte:
		cells, err := boc.DeserializeBoc(c)
		if err != nil {
			return tlb.Message{}, tongo.AccountID{}, err
		}
		dataCell = cells[0]
	case string:
		dataCell, err = boc.DeserializeSinglRootBase64(c)
	default:
		err = tlb.Marshal(dataCell, data)
		if err != nil {
			return tlb.Message{}, tongo.AccountID{}, err
		}
	}
	switch c := body.(type) {
	case []byte:
		cells, err := boc.DeserializeBoc(c)
		if err != nil {
			return tlb.Message{}, tongo.AccountID{}, err
		}
		bodyCell = cells[0]
	case string:
		bodyCell, err = boc.DeserializeSinglRootBase64(c)
	default:
		err = tlb.Marshal(bodyCell, body)
		if err != nil {
			return tlb.Message{}, tongo.AccountID{}, err
		}
	}

	m := tlb.Message{
		Info: tlb.CommonMsgInfo{
			SumType: "ExtInMsgInfo",
			ExtInMsgInfo: &struct {
				Src       tlb.MsgAddress
				Dest      tlb.MsgAddress
				ImportFee tlb.Grams
			}{
				Src: tlb.MsgAddress{SumType: "AddrNone"},
			},
		},
	}
	m.Init.Exists = true
	m.Init.Value.Value.Code.Exists = true
	m.Init.Value.Value.Data.Exists = true
	m.Init.Value.Value.Code.Value.Value = *codeCell
	m.Init.Value.Value.Code.Value.Value = *dataCell

	initState := boc.NewCell()
	err = tlb.Marshal(initState, m.Init.Value.Value)
	if err != nil {
		return tlb.Message{}, tongo.AccountID{}, err
	}
	hash, err := initState.Hash()
	if err != nil {
		return tlb.Message{}, tongo.AccountID{}, err
	}
	a := tongo.AccountID{Workchain: 0, Address: [32]byte(hash)}
	m.Init.Value.Value.Code.Value.Value.ResetCounters()
	m.Init.Value.Value.Data.Value.Value.ResetCounters()

	m.Info.ExtInMsgInfo.Dest = a.ToMsgAddress()

	m.Body.Value = tlb.Any(*bodyCell)
	return m, a, nil
}
