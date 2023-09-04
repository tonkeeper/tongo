package nft

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
)

type Item struct {
	Address    ton.AccountID
	Collection *ton.AccountID
	Owner      *ton.AccountID
}

type sender interface {
	Send(context.Context, ...wallet.Sendable) error
	GetAddress() ton.AccountID
}

func (item Item) Transfer(ctx context.Context, sender sender, destination ton.AccountID) error {
	if item.Owner != nil && sender.GetAddress() != *item.Owner {
		return errors.New("sender is not the item owner")
	}
	transfer := ItemTransferMessage{
		ItemAddress:         item.Address,
		Destination:         destination,
		ResponseDestination: sender.GetAddress(),
		AttachedTon:         ton.OneTON / 20,
		ForwardTon:          0,
	}
	return sender.Send(ctx, transfer)
}

type ItemTransferMessage struct {
	ItemAddress         ton.AccountID
	Destination         ton.AccountID
	ResponseDestination ton.AccountID
	AttachedTon         tlb.Grams
	ForwardTon          tlb.Grams
	ForwardPayload      *boc.Cell
	CustomPayload       *boc.Cell
}

func (itm ItemTransferMessage) ToInternal() (tlb.Message, byte, error) {
	c := boc.NewCell()
	forwardTon := big.NewInt(int64(itm.ForwardTon))
	msgBody := abi.NftTransferMsgBody{
		QueryId:             uint64(time.Now().UnixNano()),
		NewOwner:            itm.Destination.ToMsgAddress(),
		ResponseDestination: itm.ResponseDestination.ToMsgAddress(),
		ForwardAmount:       tlb.VarUInteger16(*forwardTon),
	}
	if itm.CustomPayload != nil {
		msgBody.CustomPayload.Exists = true
		msgBody.CustomPayload.Value.Value = tlb.Any(*itm.CustomPayload)
	}
	if itm.ForwardPayload != nil {
		msgBody.ForwardPayload.IsRight = true
		msgBody.ForwardPayload.Value = tlb.Any(*itm.ForwardPayload)
	}
	if err := c.WriteUint(0x5fcc3d14, 32); err != nil {
		return tlb.Message{}, 0, err
	}
	if err := tlb.Marshal(c, msgBody); err != nil {
		return tlb.Message{}, 0, err
	}
	m := wallet.Message{
		Amount:  itm.AttachedTon,
		Address: itm.ItemAddress,
		Bounce:  true,
		Mode:    wallet.DefaultMessageMode,
		Body:    c,
	}
	return m.ToInternal()
}
