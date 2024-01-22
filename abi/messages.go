package abi

import (
	"fmt"
	"reflect"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func pointer[T any](t T) *T {
	return &t
}

func decodeMsg(tag tlb.Tag, name MsgOpName, bodyType any) msgDecoderFunc {
	return func(cell *boc.Cell) (*uint32, *MsgOpName, any, error) {
		readTag, err := cell.ReadUint(tag.Len)
		if err != nil {
			return nil, nil, nil, err
		}
		var uintTag *uint32
		if readTag != tag.Val {
			cell.ResetCounters()
			return nil, nil, nil, fmt.Errorf("invalid tag")
		}
		if tag.Len == 32 {
			uintTag = pointer(uint32(readTag))
		}
		body := reflect.New(reflect.TypeOf(bodyType))
		err = tlb.Unmarshal(cell, body.Interface())
		if err == nil {
			return uintTag, pointer(name), body.Elem().Interface(), nil
		}
		cell.ResetCounters()
		return nil, nil, nil, err
	}
}

func decodeMultipleMsgs(funcs []msgDecoderFunc, tag string) msgDecoderFunc {
	return func(cell *boc.Cell) (*uint32, *MsgOpName, any, error) {
		for _, f := range funcs {
			tag, opName, object, err := f(cell)
			if err == nil && cell.BitsAvailableForRead() == 0 && cell.RefsAvailableForRead() == 0 {
				return tag, opName, object, err
			}
			cell.ResetCounters()
		}
		return nil, nil, nil, fmt.Errorf("no one message can be unmarshled for %v", tag)
	}
}

type msgDecoderFunc func(cell *boc.Cell) (*uint32, *MsgOpName, any, error)

// MessageDecoder takes in a message body as a cell and tries to decode it based on the contract type or the first 4 bytes.
// It returns an opcode, an operation name and a decoded body.
func InternalMessageDecoder(cell *boc.Cell, interfaces []ContractInterface) (*uint32, *MsgOpName, any, error) {
	if cell.BitsAvailableForRead() < 32 {
		return nil, nil, nil, nil
	}
	if cell.BitsAvailableForRead() != cell.BitSize() {
		cell = cell.CopyRemaining() //because body can be part of the message cell
	}
	for _, i := range interfaces {
		for _, f := range i.IntMsgs() {
			t, m, v, err := f(cell)
			if err == nil {
				return t, m, v, nil
			}
		}
	}
	tag64, err := cell.PickUint(32)
	if err != nil {
		return nil, nil, nil, err
	}
	tag := uint32(tag64)
	f := opcodedMsgInDecodeFunctions[tag]
	if f != nil {
		return f(cell)
	}
	return &tag, nil, nil, nil
}

func completedRead(cell *boc.Cell) bool {
	return cell.RefsAvailableForRead() == 0 && cell.BitsAvailableForRead() == 0
}

// MsgOpName is a human-friendly name for a message's operation which is identified by the first 4 bytes of the message's body.
type MsgOpName = string

// MsgOpCode is the first 4 bytes of a message body identifying an operation to be performed.
type MsgOpCode = uint32
