package main

import (
	"encoding/hex"
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
)

// You can serialize and deserialize (to/from Cells) structures described by TL-B schemas.

func main() {

	// Deserialize
	// TL-B schema:
	// transfer#5fcc3d14 query_id:uint64 new_owner:MsgAddress
	// response_destination:MsgAddress custom_payload:(Maybe ^Cell)
	// forward_amount:(VarUInteger 16) forward_payload:(Either Cell ^Cell)  = InternalMsgBody;

	// Design new struct using ready-made TLB primitives
	type InternalMsgBody struct {
		tlb.SumType
		Transfer struct {
			QueryId             uint64
			NewOwner            tongo.MsgAddress
			ResponseDestination tongo.MsgAddress
			CustomPayload       tlb.Maybe[tlb.Ref[boc.Cell]]
			ForwardAmount       tlb.VarUInteger `tlb:"16bytes"`
			ForwardPayload      tlb.Either[boc.Cell, tlb.Ref[boc.Cell]]
		} `tlbSumType:"transfer#5fcc3d14"`
	}
	b, err := hex.DecodeString("b5ee9c72c10101010056000000a75fcc3d140000000000000000800c0674dd00e3a7231084788441cc873e60eb8681f44901cba3a9107c5c322dc4500034a37c6673343b360e10d4e438483b555805a20e5f056742b6a42ba35311994c802625a008a90c976e")
	if err != nil {
		panic(err)
	}
	cell, err := boc.DeserializeBoc(b) // deserialize to bag-of-cells with one root cell
	if err != nil {
		panic(err)
	}
	var res InternalMsgBody
	err = tlb.Unmarshal(cell[0], &res)
	if err != nil {
		panic(err)
	}
	if res.SumType == "Transfer" {
		newOwner, err := res.Transfer.NewOwner.AccountId() // convert tongo.MsgAddress to basic AccountID type
		if err != nil {
			panic(err)
		}
		fmt.Printf("Deserialized data:\n QueryId: %v\n NewOwner: %v\n",
			res.Transfer.QueryId, newOwner.String())
	} else {
		if err != nil {
			panic("invalid data")
		}
	}

	// Serialize
	newCell := boc.NewCell()
	var X struct {
		A uint32 `tlb:"8bits"`
		B tlb.Ref[struct {
			C uint32 `tlb:"8bits"`
		}]
	}
	X.A = 10
	X.B.Value.C = 11
	err = tlb.Marshal(newCell, X)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Serialized cell:\n %v", newCell.ToString())
}
