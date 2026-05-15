package tolkgen

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/tonkeeper/tongo/tolk/parser"
)

// ExtractShared analyses a group of ABIs that target the same package and extracts any identical declaration
// If the same identifier is used for non-identical decl., an error will be thrown
func ExtractShared(abis map[string]parser.ContractABI) (map[string]parser.ContractABI, *parser.ContractABI, error) {
	type declRecord struct {
		decl        parser.ABIDeclaration
		fingerprint string
		sourceKey   string
	}
	declRegistry := make(map[string]declRecord)
	sharedDecls := make(map[string]struct{})

	type errRecord struct {
		thrownErr parser.ABIThrownError
		sourceKey string
	}
	errRegistry := make(map[string]errRecord)
	sharedErrs := make(map[string]struct{})

	for _, key := range sortedKeys(abis) {
		abi := abis[key]
		symbols := symTable{
			ABIIndex: parser.NewABIIndex(abi),
		}
		for _, decl := range abi.Declarations {
			name := DeclName(decl)
			fp, err := declFingerprint(decl, &symbols)
			if err != nil {
				return nil, nil, fmt.Errorf("fingerprint declaration %q in %s: %w", name, key, err)
			}
			if prev, exists := declRegistry[name]; exists {
				if prev.fingerprint != fp {
					return nil, nil, fmt.Errorf("declaration %q is defined differently in %s and %s", name, prev.sourceKey, key)
				}
				sharedDecls[name] = struct{}{}
			} else {
				declRegistry[name] = declRecord{decl: decl, fingerprint: fp, sourceKey: key}
			}
		}
		for _, te := range abi.ThrownErrors {
			if te.Name == "" {
				continue
			}
			if prev, exists := errRegistry[te.Name]; exists {
				if prev.thrownErr.ErrCode != te.ErrCode {
					return nil, nil, fmt.Errorf("error constant %q has code %d in %s but %d in %s",
						te.Name, prev.thrownErr.ErrCode, prev.sourceKey, te.ErrCode, key)
				}
				sharedErrs[te.Name] = struct{}{}
			} else {
				errRegistry[te.Name] = errRecord{thrownErr: te, sourceKey: key}
			}
		}
	}

	result := make(map[string]parser.ContractABI, len(abis)+1)

	var sharedABI *parser.ContractABI
	if len(sharedDecls) > 0 || len(sharedErrs) > 0 {
		sDecls := make([]parser.ABIDeclaration, 0, len(sharedDecls))
		sharedUniqueTypes := make([]parser.Ty, 0)
		remappers := make(map[string]*typeRemapper)
		for _, name := range sortedStringSet(sharedDecls) {
			record := declRegistry[name]
			remapper := remappers[record.sourceKey]
			if remapper == nil {
				abi := abis[record.sourceKey]
				remapper = &typeRemapper{
					src:  abi.UniqueTypes,
					dst:  &sharedUniqueTypes,
					memo: map[int]int{},
				}
				remappers[record.sourceKey] = remapper
			}
			sDecls = append(sDecls, remapper.remapDecl(record.decl))
		}

		sErrs := make([]parser.ABIThrownError, 0, len(sharedErrs))
		for _, name := range sortedStringSet(sharedErrs) {
			sErrs = append(sErrs, errRegistry[name].thrownErr)
		}

		sharedABI = &parser.ContractABI{UniqueTypes: sharedUniqueTypes, Declarations: sDecls, ThrownErrors: sErrs}
	}

	// Rebuild each individual ABI without the shared items.
	for _, key := range sortedKeys(abis) {
		abi := abis[key]
		uniqueDecls := make([]parser.ABIDeclaration, 0, len(abi.Declarations))
		for _, d := range abi.Declarations {
			if _, isShared := sharedDecls[DeclName(d)]; !isShared {
				uniqueDecls = append(uniqueDecls, d)
			}
		}
		uniqueErrs := make([]parser.ABIThrownError, 0, len(abi.ThrownErrors))
		for _, te := range abi.ThrownErrors {
			if _, isShared := sharedErrs[te.Name]; !isShared {
				uniqueErrs = append(uniqueErrs, te)
			}
		}
		abi.Declarations = uniqueDecls
		abi.ThrownErrors = uniqueErrs
		result[key] = abi
	}

	return result, sharedABI, nil
}

func sortedStringSet(set map[string]struct{}) []string {
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

type typeRemapper struct {
	src  []parser.Ty
	dst  *[]parser.Ty
	memo map[int]int
}

func (r *typeRemapper) remapIdx(idx int) int {
	if remapped, ok := r.memo[idx]; ok {
		return remapped
	}
	remapped := len(*r.dst)
	r.memo[idx] = remapped
	*r.dst = append(*r.dst, parser.Ty{})
	(*r.dst)[remapped] = r.remapTy(r.src[idx])
	return remapped
}

func (r *typeRemapper) remapTy(ty parser.Ty) parser.Ty {
	switch ty.SumType {
	case parser.TyKindNullable:
		ty.Nullable.InnerTyIdx = r.remapIdx(ty.Nullable.InnerTyIdx)
	case parser.TyKindCellOf:
		ty.CellOf.InnerTyIdx = r.remapIdx(ty.CellOf.InnerTyIdx)
	case parser.TyKindArrayOf:
		ty.ArrayOf.InnerTyIdx = r.remapIdx(ty.ArrayOf.InnerTyIdx)
	case parser.TyKindLispListOf:
		ty.LispListOf.InnerTyIdx = r.remapIdx(ty.LispListOf.InnerTyIdx)
	case parser.TyKindTensor:
		for i, idx := range ty.Tensor.ItemsTyIdx {
			ty.Tensor.ItemsTyIdx[i] = r.remapIdx(idx)
		}
	case parser.TyKindShapedTuple:
		for i, idx := range ty.ShapedTuple.ItemsTyIdx {
			ty.ShapedTuple.ItemsTyIdx[i] = r.remapIdx(idx)
		}
	case parser.TyKindMapKV:
		ty.MapKV.KeyTyIdx = r.remapIdx(ty.MapKV.KeyTyIdx)
		ty.MapKV.ValueTyIdx = r.remapIdx(ty.MapKV.ValueTyIdx)
	case parser.TyKindStructRef:
		for i, idx := range ty.StructRef.TypeArgsTyIdx {
			ty.StructRef.TypeArgsTyIdx[i] = r.remapIdx(idx)
		}
	case parser.TyKindAliasRef:
		for i, idx := range ty.AliasRef.TypeArgsTyIdx {
			ty.AliasRef.TypeArgsTyIdx[i] = r.remapIdx(idx)
		}
	case parser.TyKindUnion:
		for i := range ty.Union.Variants {
			ty.Union.Variants[i].VariantTyIdx = r.remapIdx(ty.Union.Variants[i].VariantTyIdx)
		}
	}
	return ty
}

func (r *typeRemapper) remapDecl(decl parser.ABIDeclaration) parser.ABIDeclaration {
	switch decl.SumType {
	case parser.DeclarationKindStruct:
		decl.StructDeclaration.TyIdx = r.remapIdx(decl.StructDeclaration.TyIdx)
		for i := range decl.StructDeclaration.Fields {
			decl.StructDeclaration.Fields[i].TyIdx = r.remapIdx(decl.StructDeclaration.Fields[i].TyIdx)
			if decl.StructDeclaration.Fields[i].ClientTyIdx != nil {
				clientTyIdx := r.remapIdx(*decl.StructDeclaration.Fields[i].ClientTyIdx)
				decl.StructDeclaration.Fields[i].ClientTyIdx = &clientTyIdx
			}
		}
	case parser.DeclarationKindAlias:
		decl.AliasDeclaration.TyIdx = r.remapIdx(decl.AliasDeclaration.TyIdx)
		decl.AliasDeclaration.TargetTyIdx = r.remapIdx(decl.AliasDeclaration.TargetTyIdx)
	case parser.DeclarationKindEnum:
		decl.EnumDeclaration.TyIdx = r.remapIdx(decl.EnumDeclaration.TyIdx)
		decl.EnumDeclaration.EncodedAsTyIdx = r.remapIdx(decl.EnumDeclaration.EncodedAsTyIdx)
	}
	return decl
}

func DeclName(d parser.ABIDeclaration) string {
	switch d.SumType {
	case parser.DeclarationKindStruct:
		return d.StructDeclaration.Name
	case parser.DeclarationKindAlias:
		return d.AliasDeclaration.Name
	case parser.DeclarationKindEnum:
		return d.EnumDeclaration.Name
	default:
		return ""
	}
}

func declFingerprint(d parser.ABIDeclaration, st *symTable) (string, error) {
	type fieldFP struct {
		Name        string `json:"name"`
		Ty          string `json:"ty"`
		ClientTy    string `json:"client_ty,omitempty"`
		Description string `json:"description,omitempty"`
	}
	var fp any
	switch d.SumType {
	case parser.DeclarationKindStruct:
		fields := make([]fieldFP, 0, len(d.StructDeclaration.Fields))
		for _, f := range d.StructDeclaration.Fields {
			ty, err := st.RenderTy(f.TyIdx)
			if err != nil {
				return "", fmt.Errorf("field %q type: %w", f.Name, err)
			}
			var clientTy string
			if f.ClientTyIdx != nil {
				clientTy, err = st.RenderTy(*f.ClientTyIdx)
				if err != nil {
					return "", fmt.Errorf("field %q client type: %w", f.Name, err)
				}
			}
			fields = append(fields, fieldFP{Name: f.Name, Ty: ty, ClientTy: clientTy, Description: f.Description})
		}
		fp = struct {
			Kind             parser.ABIDeclarationKind   `json:"kind"`
			Name             string                      `json:"name"`
			TypeParams       []string                    `json:"type_params,omitempty"`
			Prefix           *parser.Prefix              `json:"prefix,omitempty"`
			Fields           []fieldFP                   `json:"fields"`
			CustomPackUnpack parser.ABICustomSerializers `json:"custom_pack_unpack,omitempty"`
			Description      string                      `json:"description,omitempty"`
		}{d.SumType, d.StructDeclaration.Name, d.StructDeclaration.TypeParams, d.StructDeclaration.Prefix, fields, d.StructDeclaration.CustomPackUnpack, d.StructDeclaration.Description}
	case parser.DeclarationKindAlias:
		targetTy, err := st.RenderTy(d.AliasDeclaration.TargetTyIdx)
		if err != nil {
			return "", fmt.Errorf("alias target: %w", err)
		}
		fp = struct {
			Kind             parser.ABIDeclarationKind   `json:"kind"`
			Name             string                      `json:"name"`
			TargetTy         string                      `json:"target_ty"`
			TypeParams       []string                    `json:"type_params,omitempty"`
			CustomPackUnpack parser.ABICustomSerializers `json:"custom_pack_unpack,omitempty"`
			Description      string                      `json:"description,omitempty"`
		}{d.SumType, d.AliasDeclaration.Name, targetTy, d.AliasDeclaration.TypeParams, d.AliasDeclaration.CustomPackUnpack, d.AliasDeclaration.Description}
	case parser.DeclarationKindEnum:
		encodedAs, err := st.RenderTy(d.EnumDeclaration.EncodedAsTyIdx)
		if err != nil {
			return "", fmt.Errorf("enum encoded type: %w", err)
		}
		fp = struct {
			Kind             parser.ABIDeclarationKind   `json:"kind"`
			Name             string                      `json:"name"`
			EncodedAs        string                      `json:"encoded_as"`
			Members          []parser.ABIEnumMember      `json:"members"`
			CustomPackUnpack parser.ABICustomSerializers `json:"custom_pack_unpack,omitempty"`
			Description      string                      `json:"description,omitempty"`
		}{d.SumType, d.EnumDeclaration.Name, encodedAs, d.EnumDeclaration.Members, d.EnumDeclaration.CustomPackUnpack, d.EnumDeclaration.Description}
	default:
		return "", fmt.Errorf("unexpected kind %q", d.SumType)
	}
	b, err := json.Marshal(fp)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func sortedKeys(m map[string]parser.ContractABI) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
