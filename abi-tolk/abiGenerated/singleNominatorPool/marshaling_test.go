package abiSingleNominatorPool

import (
	"context"
	"math/big"
	"testing"

	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/tvm"
)

func TestSingleNominatorPoolGetMethods(t *testing.T) {
	code := "te6cckECDQEAAfAAART/APSkE/S88sgLAQIBYgIDArzQ7UTQ+kD6QNEixwCSXwbgA9DTAwFxsJJfBuD6QDAC0x9wIsAAIosXeMcFsCLXSsAAsI4TbCGDC8hTdqGCEDuaygCh+gLJ0JQw0z8S4lNDxwWRM+MNUjXHBZJfBuMNBAUCASALDAHEIYMLuo6g+gBTh6GCEDuaygChErYIgSAEIcIA8vRSQG2AGIBA2zzeIYEQAbqe+kBQRMhYzxYBzxbJ7VSRM+IggXcCupgC0wfUAvsAAt4gggCZA7qdAtSBIAIibvLyAfsEAt4JBPIjghBOc3RLuo/hAvpEMPgo+kQwgSADAsD/EvL0gwwBwP/y9IEgASLy9IEgBSSCEEeGjAC+8vT6ACDbPDAFgSAEBaGCEDuaygChUhC7FPL02zyCEE5zdEvIyx9SIMs/UAXPFslEMIAYgEDbPJQQNWxB4gGCEEdldCS6BggJBwAc0/8x0x/THzHT/zHUMdECNo8WghBHZXQkyMsfyz/J2zxwWIAYgEDbPJEw4ggJARZx+DPQ1wv/fwHbPAoASCJusyCRcZFw4gPIywVQBs8WUAT6AstqA5NYzAGRMOIByQH7AAAWdMjLAhLKB8v/ydAAJ734y5OLgqOAAqOAA2tqo5FSmItsABW/5QdqJofSB9IGjO7hyXM="
	data := "te6cckEBAQEARQAAhYARhy+Ifz/haSyza6FGBWNSde+ZjHy+uBieS1O4PxaPVnP8Dqn4VGKEh2fE8nnJMX7yTvNzoJgaUN4BrrMlylrF7L6zYsnB"
	pool := ton.MustParseAccountID("-1:ffdb952420893666069f674c14c516296e43bb4808e07f6726ef12ff98ad78bf")

	e, err := tvm.NewEmulatorFromBOCsBase64(code, data, "")
	if err != nil {
		t.Fatalf("NewEmulatorFromBOCsBase64() failed: %v", err)
	}

	roles, err := GetRoles(context.Background(), e, pool)
	if err != nil {
		t.Fatalf("GetRoles() failed: %v", err)
	}
	assertAddress(t, roles.Owner, "0:8c397c43f9ff0b49659b5d0a302b1a93af7ccc63e5f5c0c4f25a9dc1f8b47ab3")
	assertAddress(t, roles.Validator, "-1:03aa7e1518a121d9f13c9e724c5fbc93bcdce826069437806bacc97296b17b2f")

	poolData, err := GetPoolData(context.Background(), e, pool)
	if err != nil {
		t.Fatalf("GetPoolData() failed: %v", err)
	}
	assertInt257(t, "state", poolData.State, 2)
	assertInt257(t, "nominators_count", poolData.NominatorsCount, 1)
	assertInt257(t, "stake_amount_sent", poolData.StakeAmountSent, 0)
	assertInt257(t, "validator_amount", poolData.ValidatorAmount, 0)
	assertInt257(t, "min_stake", poolData.PoolConfig.MinStake, 0)
	assertInt257(t, "deposit_fee", poolData.PoolConfig.DepositFee, 0)
	assertInt257(t, "withdraw_fee", poolData.PoolConfig.WithdrawFee, 0)
	assertInt257(t, "pool_fee", poolData.PoolConfig.PoolFee, 0)
	assertInt257(t, "receipt_price", poolData.PoolConfig.ReceiptPrice, 0)
	if poolData.Nominators.Exists {
		t.Fatalf("nominators should be empty")
	}
	if poolData.WithdrawRequests.Exists {
		t.Fatalf("withdraw_requests should be empty")
	}
	assertInt257(t, "stake_at", poolData.StakeAt, 0)
	assertInt257(t, "saved_validator_set_hash", poolData.SavedValidatorSetHash, 0)
	assertInt257(t, "validator_set_changes_count", poolData.ValidatorSetChangesCount, 2)
	assertInt257(t, "validator_set_change_time", poolData.ValidatorSetChangeTime, 0)
	assertInt257(t, "stake_held_for", poolData.StakeHeldFor, 0)
	if poolData.ConfigProposalVotings.Exists {
		t.Fatalf("config_proposal_votings should be empty")
	}
}

func assertAddress(t *testing.T, got tlb.InternalAddress, want string) {
	t.Helper()
	if got.String() != want {
		t.Errorf("got address %s, want %s", got.String(), want)
	}
}

func assertInt257(t *testing.T, name string, got tlb.Int257, want int64) {
	t.Helper()
	gotInt := big.Int(got)
	if gotInt.Int64() != want {
		t.Errorf("%s: got %s, want %d", name, gotInt.String(), want)
	}
}
