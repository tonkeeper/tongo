package tolk

import (
	"github.com/tonkeeper/tongo/boc"
)

type TolkUnmarshaler interface {
	UnmarshalTolk(cell *boc.Cell, v *Value, d *Decoder) error
}

type abiRefs struct {
	structRefs  map[string]StructDeclaration
	aliasRefs   map[string]AliasDeclaration
	enumRefs    map[string]EnumDeclaration
	genericRefs map[string]Ty
}

type Decoder struct {
	abiRefs abiRefs
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (a *Decoder) WithABI(abi ABI) *Decoder {
	a.abiRefs = abiRefs{
		structRefs:  make(map[string]StructDeclaration),
		aliasRefs:   make(map[string]AliasDeclaration),
		enumRefs:    make(map[string]EnumDeclaration),
		genericRefs: make(map[string]Ty),
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

func UnmarshalTolk(cell *boc.Cell, ty Ty) (*Value, error) {
	a := NewDecoder()
	return a.UnmarshalTolk(cell, ty)
}

func (a *Decoder) UnmarshalTolk(cell *boc.Cell, ty Ty) (*Value, error) {
	res := &Value{}
	err := res.UnmarshalTolk(cell, ty, a)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func MarshalTolk(c *boc.Cell, v *Value) error {
	return nil
	//return v.valType.MarshalTolk(c, v)
}
