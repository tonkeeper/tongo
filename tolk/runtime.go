package tolk

import (
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
}

type Decoder struct {
	abiRefs abiRefs
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (a *Decoder) WithABI(abi tolkParser.ABI) *Decoder {
	a.abiRefs = abiRefs{
		structRefs:  make(map[string]tolkParser.StructDeclaration),
		aliasRefs:   make(map[string]tolkParser.AliasDeclaration),
		enumRefs:    make(map[string]tolkParser.EnumDeclaration),
		genericRefs: make(map[string]tolkParser.Ty),
	}
	for _, declr := range abi.Declarations {
		switch declr.SumType {
		case "Struct":
			a.abiRefs.structRefs[declr.StructDeclaration.Name] = declr.StructDeclaration
		case "Alias":
			a.abiRefs.aliasRefs[declr.AliasDeclaration.Name] = declr.AliasDeclaration
		case "Enum":
			a.abiRefs.enumRefs[declr.EnumDeclaration.Name] = declr.EnumDeclaration
		}
	}
	return a
}

func (a *Decoder) Unmarshal(cell *boc.Cell, ty tolkParser.Ty) (*Value, error) {
	res := &Value{}
	err := res.Unmarshal(cell, ty, a)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type Encoder struct {
	abiRefs abiRefs
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (a *Encoder) WithABI(abi tolkParser.ABI) *Encoder {
	a.abiRefs = abiRefs{
		structRefs:  make(map[string]tolkParser.StructDeclaration),
		aliasRefs:   make(map[string]tolkParser.AliasDeclaration),
		enumRefs:    make(map[string]tolkParser.EnumDeclaration),
		genericRefs: make(map[string]tolkParser.Ty),
	}
	for _, declr := range abi.Declarations {
		switch declr.SumType {
		case "Struct":
			a.abiRefs.structRefs[declr.StructDeclaration.Name] = declr.StructDeclaration
		case "Alias":
			a.abiRefs.aliasRefs[declr.AliasDeclaration.Name] = declr.AliasDeclaration
		case "Enum":
			a.abiRefs.enumRefs[declr.EnumDeclaration.Name] = declr.EnumDeclaration
		}
	}
	return a
}

func (a *Encoder) Marshal(v *Value, ty tolkParser.Ty) (*boc.Cell, error) {
	cell := boc.NewCell()
	err := v.Marshal(cell, ty, a)
	if err != nil {
		return nil, err
	}
	return cell, nil
}
