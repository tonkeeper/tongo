package tolkgen

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/tonkeeper/tongo/tolk/parser"
)

// ExtractShared analyses a group of ABIs that target the same package and extracts any identical declaration
// If the same identifier is used for non-identical decl., an error will be thrown
func ExtractShared(abis map[string]parser.ABI) (map[string]parser.ABI, *parser.ABI, error) {
	type declRecord struct {
		decl        parser.Declaration
		fingerprint string
		sourceKey   string
	}
	declRegistry := make(map[string]declRecord)
	sharedDecls := make(map[string]struct{})

	type errRecord struct {
		thrownErr parser.ThrownError
		sourceKey string
	}
	errRegistry := make(map[string]errRecord)
	sharedErrs := make(map[string]struct{})

	for _, key := range sortedKeys(abis) {
		abi := abis[key]
		for _, decl := range abi.Declarations {
			name := DeclName(decl)
			fp := declFingerprint(decl)
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

	result := make(map[string]parser.ABI, len(abis)+1)

	var sharedABI *parser.ABI
	if len(sharedDecls) > 0 || len(sharedErrs) > 0 {
		sDecls := make([]parser.Declaration, 0, len(sharedDecls))
		for name := range sharedDecls {
			sDecls = append(sDecls, declRegistry[name].decl)
		}
		sort.Slice(sDecls, func(i, j int) bool { return DeclName(sDecls[i]) < DeclName(sDecls[j]) })

		sErrs := make([]parser.ThrownError, 0, len(sharedErrs))
		for name := range sharedErrs {
			sErrs = append(sErrs, errRegistry[name].thrownErr)
		}
		sort.Slice(sErrs, func(i, j int) bool { return sErrs[i].Name < sErrs[j].Name })

		sharedABI = &parser.ABI{Declarations: sDecls, ThrownErrors: sErrs}
	}

	// Rebuild each individual ABI without the shared items.
	for _, key := range sortedKeys(abis) {
		abi := abis[key]
		uniqueDecls := make([]parser.Declaration, 0, len(abi.Declarations))
		for _, d := range abi.Declarations {
			if _, isShared := sharedDecls[DeclName(d)]; !isShared {
				uniqueDecls = append(uniqueDecls, d)
			}
		}
		uniqueErrs := make([]parser.ThrownError, 0, len(abi.ThrownErrors))
		for _, te := range abi.ThrownErrors {
			if _, isShared := sharedErrs[te.Name]; !isShared {
				uniqueErrs = append(uniqueErrs, te)
			}
		}
		result[key] = parser.ABI{Declarations: uniqueDecls, ThrownErrors: uniqueErrs}
	}

	return result, sharedABI, nil
}

func DeclName(d parser.Declaration) string {
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

func declFingerprint(d parser.Declaration) string {
	b, err := json.Marshal(d)
	if err != nil {
		panic(fmt.Sprintf("declFingerprint invariant violated: %v", err))
	}
	return string(b)
}

func sortedKeys(m map[string]parser.ABI) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
