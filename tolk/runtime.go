package tolk

import (
	"fmt"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tolk/parser"
)

func Unmarshal(cell *boc.Cell, ty tolkParser.Ty) (*Value, error) {
	d := NewDecoder()
	return d.Unmarshal(cell, ty)
}

func Marshal(v *Value, ty tolkParser.Ty) (*boc.Cell, error) {
	e := NewEncoder()
	return e.Marshal(v, ty)
}

type abiRefs struct {
	structRefs  map[string]tolkParser.StructDeclaration
	aliasRefs   map[string]tolkParser.AliasDeclaration
	enumRefs    map[string]tolkParser.EnumDeclaration
	genericRefs map[string]tolkParser.Ty
	opcodeRefs  map[uint64][]tolkParser.StructDeclaration
}

type customUnpackResolver = func(tolkParser.AliasRef, *boc.Cell, *AliasValue) error

type Decoder struct {
	abiRefs              abiRefs
	customUnpackResolver customUnpackResolver
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) WithABIs(abis ...tolkParser.ABI) error {
	d.abiRefs = abiRefs{
		structRefs:  make(map[string]tolkParser.StructDeclaration),
		aliasRefs:   make(map[string]tolkParser.AliasDeclaration),
		enumRefs:    make(map[string]tolkParser.EnumDeclaration),
		genericRefs: make(map[string]tolkParser.Ty),
		opcodeRefs:  make(map[uint64][]tolkParser.StructDeclaration),
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

func (d *Decoder) Unmarshal(cell *boc.Cell, ty tolkParser.Ty) (*Value, error) {
	res := &Value{}
	err := res.Unmarshal(cell, ty, d)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tolk value: %w", err)
	}
	return res, nil
}

func (d *Decoder) UnmarshalMessage(cell *boc.Cell) (*Value, error) {
	res, isResolved, err := resolvePayload(cell, d)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tolk value: %w", err)
	}
	if isResolved {
		return &res, nil
	}
	res = Value{
		SumType:   "Remaining",
		Remaining: (*RemainingValue)(cell),
	}
	return &res, nil
}

type customPackResolver = func(tolkParser.AliasRef, *boc.Cell, *AliasValue) error

type Encoder struct {
	abiRefs            abiRefs
	customPackResolver customPackResolver
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (a *Encoder) WithABIs(abis ...tolkParser.ABI) error {
	a.abiRefs = abiRefs{
		structRefs:  make(map[string]tolkParser.StructDeclaration),
		aliasRefs:   make(map[string]tolkParser.AliasDeclaration),
		enumRefs:    make(map[string]tolkParser.EnumDeclaration),
		genericRefs: make(map[string]tolkParser.Ty),
		opcodeRefs:  make(map[uint64][]tolkParser.StructDeclaration),
	}
	for _, abi := range abis {
		for _, declr := range abi.Declarations {
			switch declr.SumType {
			case "Struct":
				a.abiRefs.structRefs[declr.StructDeclaration.Name] = declr.StructDeclaration
				if declr.StructDeclaration.Prefix != nil {
					prefix, err := binHexToUint64(declr.StructDeclaration.Prefix.PrefixStr)
					if err != nil {
						return fmt.Errorf("failed to parse prefix struct %v prefix: %w", declr.StructDeclaration.Name, err)
					}
					a.abiRefs.opcodeRefs[prefix] = append(a.abiRefs.opcodeRefs[prefix], declr.StructDeclaration)
				}
			case "Alias":
				a.abiRefs.aliasRefs[declr.AliasDeclaration.Name] = declr.AliasDeclaration
			case "Enum":
				a.abiRefs.enumRefs[declr.EnumDeclaration.Name] = declr.EnumDeclaration
			}
		}
	}
	return nil
}

func (a *Encoder) WithCustomPackResolver(customPackResolver customPackResolver) {
	a.customPackResolver = customPackResolver
}

func (a *Encoder) Marshal(v *Value, ty tolkParser.Ty) (*boc.Cell, error) {
	cell := boc.NewCell()
	err := v.Marshal(cell, ty, a)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tolk value: %w", err)
	}
	return cell, nil
}
