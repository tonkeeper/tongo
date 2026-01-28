package tolk

import (
	"github.com/tonkeeper/tongo/boc"
)

type TolkUnmarshaler interface {
	UnmarshalTolk(cell *boc.Cell, v *Value, d *Decoder) error
}

type abiContext struct {
	structRefs  map[string]StructDeclaration
	aliasRefs   map[string]AliasDeclaration
	enumRefs    map[string]EnumDeclaration
	genericRefs map[string]Ty
}

type Decoder struct {
	abiCtx *abiContext
}

func NewDecoder() *Decoder {
	return &Decoder{
		abiCtx: &abiContext{
			genericRefs: make(map[string]Ty),
		},
	}
}

func (d *Decoder) WithABI(abi ABI) *Decoder {
	d.abiCtx = &abiContext{
		structRefs: make(map[string]StructDeclaration),
		aliasRefs:  make(map[string]AliasDeclaration),
		enumRefs:   make(map[string]EnumDeclaration),
	}
	for _, declr := range abi.Declarations {
		switch declr.SumType {
		case "Struct":
			d.abiCtx.structRefs[declr.StructDeclaration.Name] = declr.StructDeclaration
		case "Alias":
			d.abiCtx.aliasRefs[declr.AliasDeclaration.Name] = declr.AliasDeclaration
		case "Enum":
			d.abiCtx.enumRefs[declr.EnumDeclaration.Name] = declr.EnumDeclaration
		}
	}
	return d
}

func UnmarshalTolk(cell *boc.Cell, ty Ty) (*Value, error) {
	d := NewDecoder()
	return d.UnmarshalTolk(cell, ty)
}

// todo: maybe use only abi (guess best struct to unmarshal)
func (d *Decoder) UnmarshalTolk(cell *boc.Cell, ty Ty) (*Value, error) {
	res := &Value{}
	err := ty.UnmarshalTolk(cell, res, d)
	if err != nil {
		return nil, err
	}
	return res, nil
}
