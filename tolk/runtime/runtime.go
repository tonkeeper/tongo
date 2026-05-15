package runtime

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

func Unmarshal(cell *boc.Cell, ty parser.Ty) (*Value, error) {
	d := &Decoder{}
	return d.Unmarshal(cell, ty)
}

func Marshal(v *Value, ty parser.Ty) (*boc.Cell, error) {
	e := &Encoder{}
	return e.Marshal(v, ty)
}

type customUnpackResolver = func(parser.AliasRef, *boc.Cell, *AliasValue) error

type Decoder struct {
	abiIndex             *parser.ABIIndex
	genericRefs          map[string]int
	opcodeRefs           map[uint64][]int
	customUnpackResolver customUnpackResolver
}

func NewDecoder(abi parser.ContractABI) *Decoder {
	index, opcodeRefs := newRuntimeIndex(abi)
	return &Decoder{abiIndex: index, genericRefs: make(map[string]int), opcodeRefs: opcodeRefs}
}

func newRuntimeIndex(abi parser.ContractABI) (*parser.ABIIndex, map[uint64][]int) {
	opcodeRefs := make(map[uint64][]int)
	for _, declr := range abi.Declarations {
		if declr.SumType == parser.DeclarationKindStruct && declr.StructDeclaration.Prefix != nil {
			prefix := uint64(declr.StructDeclaration.Prefix.PrefixNum)
			opcodeRefs[prefix] = append(opcodeRefs[prefix], declr.StructDeclaration.TyIdx)
		}
	}
	return parser.NewABIIndex(abi), opcodeRefs
}

func (d *Decoder) WithCustomUnpackResolver(customUnpackResolver customUnpackResolver) {
	d.customUnpackResolver = customUnpackResolver
}

func (d *Decoder) Unmarshal(cell *boc.Cell, ty parser.Ty) (*Value, error) {
	res := &Value{}
	err := res.Unmarshal(cell, ty, d)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tolk value: %w", err)
	}
	return res, nil
}

func (d *Decoder) UnmarshalTyIdx(cell *boc.Cell, tyIdx int) (*Value, error) {
	res := &Value{}
	err := res.UnmarshalTyIdx(cell, tyIdx, d)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tolk value: %w", err)
	}
	return res, nil
}

func (d *Decoder) UnmarshalMessage(cell *boc.Cell) (*Value, error) {
	res, isResolved, err := d.resolvePayload(cell)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tolk value: %w", err)
	}
	if isResolved {
		return &res, nil
	}
	res = Value{
		SumType: SumTypeRemaining,
		Remaining: &RemainingValue{
			IsRef: cell.BitsAvailableForRead() == 0 && cell.RefsAvailableForRead() > 0,
			Value: *cell,
		},
	}
	return &res, nil
}

func (d *Decoder) resolvePayload(payload *boc.Cell) (Value, bool, error) {
	payloadOpcode, err := payload.ReadUint(32) // payload always 32 bit length
	if err != nil {
		return Value{}, false, nil
	}
	payload.ResetCounters() // reset opcode

	guessedStructs := d.opcodeRefs[payloadOpcode]
	for _, tyIdx := range guessedStructs {
		v, err := d.UnmarshalTyIdx(payload, tyIdx)
		if err != nil {
			continue
		}
		return *v, true, nil
	}

	// todo: maybe try every known struct to unmarshal to?
	return Value{}, false, nil
}

type customPackResolver = func(parser.AliasRef, *boc.Cell, *AliasValue) error

type Encoder struct {
	abiIndex           *parser.ABIIndex
	genericRefs        map[string]int
	customPackResolver customPackResolver
}

func NewEncoder(abi parser.ContractABI) *Encoder {
	index, _ := newRuntimeIndex(abi)
	return &Encoder{abiIndex: index, genericRefs: make(map[string]int)}
}

func (e *Encoder) WithCustomPackResolver(customPackResolver customPackResolver) {
	e.customPackResolver = customPackResolver
}

func (e *Encoder) Marshal(v *Value, ty parser.Ty) (*boc.Cell, error) {
	cell := boc.NewCell()
	err := v.Marshal(cell, ty, e)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tolk value: %w", err)
	}
	return cell, nil
}

func (e *Encoder) MarshalTyIdx(v *Value, tyIdx int) (*boc.Cell, error) {
	cell := boc.NewCell()
	err := v.MarshalTyIdx(cell, tyIdx, e)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tolk value: %w", err)
	}
	return cell, nil
}
