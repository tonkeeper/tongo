package tolkParser

func NewIntNType(n int) Ty {
	return Ty{
		SumType: "IntN",
		IntN: &IntN{
			N: n,
		},
	}
}

func NewUIntNType(n int) Ty {
	return Ty{
		SumType: "UintN",
		UintN: &UintN{
			N: n,
		},
	}
}

func NewVarInt16Type() Ty {
	return Ty{
		SumType: "VarIntN",
		VarIntN: &VarIntN{
			N: 16,
		},
	}
}

func NewVarInt32Type() Ty {
	return Ty{
		SumType: "VarIntN",
		VarIntN: &VarIntN{
			N: 32,
		},
	}
}

func NewVarIntType(n int) Ty {
	return Ty{
		SumType: "VarIntN",
		VarIntN: &VarIntN{
			N: n,
		},
	}
}

func NewVarUInt16Type() Ty {
	return Ty{
		SumType: "VarUintN",
		VarUintN: &VarUintN{
			N: 16,
		},
	}
}

func NewVarUInt32Type() Ty {
	return Ty{
		SumType: "VarUintN",
		VarUintN: &VarUintN{
			N: 32,
		},
	}
}

func NewVarUIntType(n int) Ty {
	return Ty{
		SumType: "VarUintN",
		VarUintN: &VarUintN{
			N: n,
		},
	}
}

func NewBitsNType(n int) Ty {
	return Ty{
		SumType: "BitsN",
		BitsN: &BitsN{
			N: n,
		},
	}
}

func NewCoinsType() Ty {
	return Ty{
		SumType: "Coins",
		Coins:   &Coins{},
	}
}

func NewBoolType() Ty {
	return Ty{
		SumType: "Bool",
		Bool:    &Bool{},
	}
}

func NewCellType() Ty {
	return Ty{
		SumType: "Cell",
		Cell:    &Cell{},
	}
}

func NewRemainingType() Ty {
	return Ty{
		SumType:   "Remaining",
		Remaining: &Remaining{},
	}
}

func NewAddressType() Ty {
	return Ty{
		SumType: "Address",
		Address: &Address{},
	}
}

func NewAddressOptType() Ty {
	return Ty{
		SumType:    "AddressOpt",
		AddressOpt: &AddressOpt{},
	}
}

func NewAddressExtType() Ty {
	return Ty{
		SumType:    "AddressExt",
		AddressExt: &AddressExt{},
	}
}

func NewAddressAnyType() Ty {
	return Ty{
		SumType:    "AddressAny",
		AddressAny: &AddressAny{},
	}
}

func NewNullableType(of Ty) Ty {
	return Ty{
		SumType: "Nullable",
		Nullable: &Nullable{
			Inner: of,
		},
	}
}

func NewCellOfType(of Ty) Ty {
	return Ty{
		SumType: "CellOf",
		CellOf: &CellOf{
			Inner: of,
		},
	}
}

func NewTensorType(of ...Ty) Ty {
	return Ty{
		SumType: "Tensor",
		Tensor: &Tensor{
			Items: of,
		},
	}
}

func NewTupleWithType(of ...Ty) Ty {
	return Ty{
		SumType: "TupleWith",
		TupleWith: &TupleWith{
			Items: of,
		},
	}
}

func NewMapType(key, value Ty) Ty {
	return Ty{
		SumType: "Map",
		Map: &Map{
			K: key,
			V: value,
		},
	}
}

func NewEnumType(name string) Ty {
	return Ty{
		SumType: "EnumRef",
		EnumRef: &EnumRef{
			EnumName: name,
		},
	}
}

func NewAliasType(name string, typeArgs ...Ty) Ty {
	return Ty{
		SumType: "AliasRef",
		AliasRef: &AliasRef{
			AliasName: name,
			TypeArgs:  typeArgs,
		},
	}
}

func NewStructType(name string, typeArgs ...Ty) Ty {
	return Ty{
		SumType: "StructRef",
		StructRef: &StructRef{
			StructName: name,
			TypeArgs:   typeArgs,
		},
	}
}

func NewGenericType(nameT string) Ty {
	return Ty{
		SumType: "Generic",
		Generic: &Generic{
			NameT: nameT,
		},
	}
}

func NewUnionVariant(variant Ty, prefixValue string) UnionVariant {
	return UnionVariant{
		PrefixStr: prefixValue,
		VariantTy: variant,
	}
}

func NewUnionType(prefixLen int, prefixEatIntPlace bool, variants ...UnionVariant) Ty {
	for i := range variants {
		variants[i].PrefixLen = prefixLen
		variants[i].PrefixEatInPlace = prefixEatIntPlace
	}

	return Ty{
		SumType: "Union",
		Union: &Union{
			Variants: variants,
		},
	}
}
