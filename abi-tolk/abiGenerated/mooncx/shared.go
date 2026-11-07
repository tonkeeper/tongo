// Code generated - DO NOT EDIT.

package abiMooncx

import (
	"github.com/tonkeeper/tongo/tlb"
)

type ClmmDepositBin struct {
	X tlb.Coins // coins
	Y tlb.Coins // coins
}

type ClmmBin struct {
	L    tlb.Uint256 // uint256
	XFee tlb.Uint256 // uint256
	YFee tlb.Uint256 // uint256
}

type LpSplit struct {
	ARecipient tlb.InternalAddress // address
	ABps       tlb.Uint16          // uint16
	BRecipient tlb.InternalAddress // address
}
