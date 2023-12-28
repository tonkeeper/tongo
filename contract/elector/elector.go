package elector

import (
	"context"
	"fmt"

	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/utils"
)

type executor interface {
	RunSmcMethodByID(context.Context, ton.AccountID, int, tlb.VmStack) (uint32, tlb.VmStack, error)
}

type Validator struct {
	Pubkey    tlb.Bits256
	Stake     int64
	MaxFactor int64
	Address   ton.AccountID
	AdnlAddr  string
}

type ParticipantList struct {
	ElectAt    int64
	ElectClose int64
	MinStake   int64
	TotalStake int64
	Validators []Validator
}

// participantList describes a result type of participant_list_extended method.
// https://github.com/ton-blockchain/governance-contract/blob/master/elector-code.fc#L1461
type participantList struct {
	ElectAt    int64
	ElectClose int64
	MinStake   int64
	TotalStake int64
	Validators []struct {
		ID        tlb.Bits256
		Validator struct {
			Stake     tlb.Coins
			MaxFactor uint32
			Address   tlb.Bits256
			AdnlAddr  tlb.Bits256
		}
	}
}

func GetParticipantListExtended(ctx context.Context, electorAddr ton.AccountID, e executor) (*ParticipantList, error) {
	stack := tlb.VmStack{}
	status, result, err := e.RunSmcMethodByID(ctx, electorAddr, utils.MethodIdFromName("participant_list_extended"), stack)
	if err != nil {
		return nil, err
	}
	if status != 0 {
		return nil, fmt.Errorf("emulator status: %d", status)
	}
	var list participantList
	if err = result.Unmarshal(&list); err != nil {
		return nil, err
	}
	validators := make([]Validator, 0, len(list.Validators))
	for _, v := range list.Validators {
		accountID := ton.AccountID{
			Workchain: -1,
			Address:   v.Validator.Address,
		}
		v.Validator.AdnlAddr.Hex()
		validators = append(validators, Validator{
			Pubkey:    v.ID,
			Stake:     int64(v.Validator.Stake),
			MaxFactor: int64(v.Validator.MaxFactor),
			Address:   accountID,
			AdnlAddr:  v.Validator.AdnlAddr.Hex(),
		})
	}
	return &ParticipantList{
		ElectAt:    list.ElectAt,
		ElectClose: list.ElectClose,
		MinStake:   list.MinStake,
		TotalStake: list.TotalStake,
		Validators: validators,
	}, nil
}
