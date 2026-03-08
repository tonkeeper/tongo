package tolk

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

func Unmarshal(cell *boc.Cell, ty parser.Ty) (*Value, error) {
	d := NewDecoder()
	return d.Unmarshal(cell, ty)
}

func Marshal(v *Value, ty parser.Ty) (*boc.Cell, error) {
	e := NewEncoder()
	return e.Marshal(v, ty)
}

type abiRefs struct {
	structRefs  map[string]parser.StructDeclaration
	aliasRefs   map[string]parser.AliasDeclaration
	enumRefs    map[string]parser.EnumDeclaration
	genericRefs map[string]parser.Ty
	opcodeRefs  map[uint64][]parser.StructDeclaration
}

type customUnpackResolver = func(parser.AliasRef, *boc.Cell, *AliasValue) error

type Decoder struct {
	abiRefs              abiRefs
	customUnpackResolver customUnpackResolver
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) WithABIs(abis ...parser.ABI) error {
	d.abiRefs = abiRefs{
		structRefs:  make(map[string]parser.StructDeclaration),
		aliasRefs:   make(map[string]parser.AliasDeclaration),
		enumRefs:    make(map[string]parser.EnumDeclaration),
		genericRefs: make(map[string]parser.Ty),
		opcodeRefs:  make(map[uint64][]parser.StructDeclaration),
	}
	for _, abi := range abis {
		for _, declr := range abi.Declarations {
			switch declr.SumType {
			case "Struct":
				d.abiRefs.structRefs[declr.StructDeclaration.Name] = declr.StructDeclaration
				if declr.StructDeclaration.Prefix != nil {
					prefix, err := binHexToUint64(declr.StructDeclaration.Prefix.PrefixStr)
					if err != nil {
						return fmt.Errorf("failed to parse prefix struct %v prefix: %w", declr.StructDeclaration.Name, err)
					}
					d.abiRefs.opcodeRefs[prefix] = append(d.abiRefs.opcodeRefs[prefix], declr.StructDeclaration)
				}
			case "Alias":
				d.abiRefs.aliasRefs[declr.AliasDeclaration.Name] = declr.AliasDeclaration
			case "Enum":
				d.abiRefs.enumRefs[declr.EnumDeclaration.Name] = declr.EnumDeclaration
			}
		}
	}
	return nil
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

func (d *Decoder) UnmarshalMessage(cell *boc.Cell) (*Value, error) {
	res, isResolved, err := d.resolvePayload(cell)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tolk value: %w", err)
	}
	if isResolved {
		return &res, nil
	}
	res = Value{
		SumType: "Remaining",
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

	guessedStructs := d.abiRefs.opcodeRefs[payloadOpcode]
	for _, strct := range guessedStructs {
		v, err := d.Unmarshal(payload, parser.NewStructType(strct.Name))
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
	abiRefs            abiRefs
	customPackResolver customPackResolver
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) WithABIs(abis ...parser.ABI) error {
	e.abiRefs = abiRefs{
		structRefs:  make(map[string]parser.StructDeclaration),
		aliasRefs:   make(map[string]parser.AliasDeclaration),
		enumRefs:    make(map[string]parser.EnumDeclaration),
		genericRefs: make(map[string]parser.Ty),
		opcodeRefs:  make(map[uint64][]parser.StructDeclaration),
	}
	for _, abi := range abis {
		for _, declr := range abi.Declarations {
			switch declr.SumType {
			case "Struct":
				e.abiRefs.structRefs[declr.StructDeclaration.Name] = declr.StructDeclaration
				if declr.StructDeclaration.Prefix != nil {
					prefix, err := binHexToUint64(declr.StructDeclaration.Prefix.PrefixStr)
					if err != nil {
						return fmt.Errorf("failed to parse prefix struct %v prefix: %w", declr.StructDeclaration.Name, err)
					}
					e.abiRefs.opcodeRefs[prefix] = append(e.abiRefs.opcodeRefs[prefix], declr.StructDeclaration)
				}
			case "Alias":
				e.abiRefs.aliasRefs[declr.AliasDeclaration.Name] = declr.AliasDeclaration
			case "Enum":
				e.abiRefs.enumRefs[declr.EnumDeclaration.Name] = declr.EnumDeclaration
			}
		}
	}
	return nil
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
