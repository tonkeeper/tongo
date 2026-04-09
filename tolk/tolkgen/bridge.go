package tolkgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/utils"
)

// BridgeMethod describes a 0-parameter get method that can be wrapped as abi.InvokeFn.
type BridgeMethod struct {
	TolkName   string // e.g. "get_cocoon_data"
	GoFuncName string // e.g. "GetCocoonData"
	ResultType string // e.g. "GetCocoonData_CocoonRootResult"
}

// PrefixedStruct describes a TLB struct with a fixed opcode prefix.
type PrefixedStruct struct {
	StructName  string // e.g. "Payout"
	PrefixConst string // e.g. "PrefixPayout"
	MsgOpName   string // e.g. "CocoonPayout"
	MsgOpConst  string // e.g. "CocoonPayoutMsgOp"
}

// BridgeContract holds the abi-integration info for one contract JSON file.
type BridgeContract struct {
	InterfaceName string         // abi.ContractInterface const, e.g. "CocoonRoot"
	StringName    string         // snake_case identifier, e.g. "cocoon_root"
	Methods       []BridgeMethod // 0-param get methods only
	Results       []string       // method result type strings for interface detection
}

// BridgeGroup holds everything needed to emit abi/<group>_tolk_generated.go.
type BridgeGroup struct {
	GroupName       string
	PackageAlias    string // e.g. "abiCocoon"
	ImportPath      string // full Go import path of the generated package
	Contracts       []BridgeContract
	PrefixedStructs []PrefixedStruct // all prefixed message structs in the group
}

// TolkInterfaceDef is one ContractInterface constant entry for the constants file.
type TolkInterfaceDef struct {
	GoName     string // e.g. "CocoonRoot"
	StringName string // e.g. "cocoon_root"
}

// SchemaEntry is a single ABI file within a group, identified by its base name.
type SchemaEntry struct {
	BaseName string // file basename without extension, e.g. "cocoon_root"
	ABI      parser.ABI
}

// BridgeGroupFromEntries builds a BridgeGroup from a slice of schema entries.
// modulePath is the Go module path (e.g. "github.com/tonkeeper/tongo").
// outputDir is the root output directory for generated packages (e.g. "abi-tolk/abiGenerated").
func BridgeGroupFromEntries(groupName, modulePath, outputDir string, entries []SchemaEntry) BridgeGroup {
	camelGroup := utils.ToCamelCase(groupName) // e.g. "Cocoon", "PythOracle"

	var contracts []BridgeContract
	var allABIs []parser.ABI

	for _, e := range entries {
		allABIs = append(allABIs, e.ABI)
		ifaceName := InterfaceNameFromABI(e.ABI, e.BaseName)

		var methods []BridgeMethod
		var results []string
		for _, m := range e.ABI.GetMethods {
			if len(m.Parameters) > 0 {
				continue // only 0-param methods can be wrapped as InvokeFn
			}
			goFunc := MethodGoName(m.Name)
			resultType := goFunc + "_" + ifaceName + "Result"
			methods = append(methods, BridgeMethod{
				TolkName:   m.Name,
				GoFuncName: goFunc,
				ResultType: resultType,
			})
			results = append(results, resultType)
		}

		contracts = append(contracts, BridgeContract{
			InterfaceName: ifaceName,
			StringName:    StringNameFromABI(e.ABI, e.BaseName),
			Methods:       methods,
			Results:       results,
		})
	}

	return BridgeGroup{
		GroupName:       groupName,
		PackageAlias:    "abi" + camelGroup,
		ImportPath:      modulePath + "/" + strings.ReplaceAll(outputDir, "\\", "/") + "/" + groupName,
		Contracts:       contracts,
		PrefixedStructs: CollectPrefixedStructs(allABIs, camelGroup),
	}
}

// InterfaceNameFromABI derives the abi.ContractInterface Go constant name.
// Uses contractName from the ABI JSON if set, otherwise derives from the file basename.
func InterfaceNameFromABI(abi parser.ABI, baseName string) string {
	if abi.ContractName != "" {
		return utils.ToCamelCase(abi.ContractName)
	}
	return utils.ToCamelCase(baseName)
}

// StringNameFromABI returns the snake_case name used in String() / ContractInterfaceFromString().
func StringNameFromABI(abi parser.ABI, baseName string) string {
	if abi.ContractName != "" {
		return abi.ContractName
	}
	return baseName
}

// CollectPrefixedStructs returns all structs with a TLB prefix across the given ABIs,
// deduplicated by struct name and sorted for deterministic output.
func CollectPrefixedStructs(abis []parser.ABI, camelGroupName string) []PrefixedStruct {
	seen := map[string]struct{}{}
	var result []PrefixedStruct
	for _, a := range abis {
		for _, decl := range a.Declarations {
			if decl.SumType != parser.DeclarationKindStruct {
				continue
			}
			sd := decl.StructDeclaration
			if sd.Prefix == nil {
				continue
			}
			if _, ok := seen[sd.Name]; ok {
				continue
			}
			seen[sd.Name] = struct{}{}
			goName := safeGoIdent(sd.Name)
			result = append(result, PrefixedStruct{
				StructName:  goName,
				PrefixConst: "Prefix" + goName,
				MsgOpName:   camelGroupName + goName,
				MsgOpConst:  camelGroupName + goName + "MsgOp",
			})
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].StructName < result[j].StructName })
	return result
}

// GenerateAllBridgesFile returns the content of abi/tolk_bridge_generated.go,
// combining the registrations for all groups into a single file.
func GenerateAllBridgesFile(groups []BridgeGroup) (string, error) {
	body := &strings.Builder{}
	for _, g := range groups {
		if err := writeBridgeBody(body, g); err != nil {
			return "", err
		}
	}

	extra := make(map[string]string, len(groups))
	for _, g := range groups {
		extra[g.PackageAlias] = g.ImportPath
	}
	imports := deriveImportBlock(body.String(), extra)
	return fmt.Sprintf("// Code generated - DO NOT EDIT.\n\npackage abi\n\nimport (\n%s\n)\n\n%s",
		imports, body.String()), nil
}

// writeBridgeBody writes one init() block for the given group into out.
func writeBridgeBody(out *strings.Builder, g BridgeGroup) error {
	fmt.Fprintf(out, "func init() {\n")

	if hasAnyMethods(g) {
		fmt.Fprintf(out, "\ttolkMethods = append(tolkMethods,\n")
		for _, c := range g.Contracts {
			for _, m := range c.Methods {
				fmt.Fprintf(out,
					`		MethodDescription{
			Name: %q,
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := %s.%s(ctx, executor, id)
				return %q, r, err
			},
		},
`, m.TolkName, g.PackageAlias, m.GoFuncName, m.ResultType)
			}
		}
		fmt.Fprintf(out, "\t)\n\n")
	}

	if hasDetectableIfaces(g) {
		fmt.Fprintf(out, "\ttolkInterfaceOrder = append(tolkInterfaceOrder,\n")
		for _, c := range g.Contracts {
			if len(c.Results) == 0 {
				continue
			}
			fmt.Fprintf(out,
				`		InterfaceDescription{
			Name: %s,
			Results: []string{%s},
		},
`, c.InterfaceName, quoteList(c.Results))
		}
		fmt.Fprintf(out, "\t)\n\n")
	}

	if len(g.PrefixedStructs) > 0 {
		for _, s := range g.PrefixedStructs {
			fmt.Fprintf(out, "\tregisterInMsgUnmarshalerForOpcode[*%s.%s](opcodedMsgInDecodeFunctions, uint32(%s.%s), %s.%s)\n",
				g.PackageAlias, s.StructName,
				g.PackageAlias, s.PrefixConst,
				g.PackageAlias, s.MsgOpConst)
		}
		fmt.Fprintf(out, "\n")
	}

	fmt.Fprintf(out, "}\n\n")
	return nil
}

func hasDetectableIfaces(g BridgeGroup) bool {
	for _, c := range g.Contracts {
		if len(c.Results) > 0 {
			return true
		}
	}
	return false
}

func hasAnyMethods(g BridgeGroup) bool {
	for _, c := range g.Contracts {
		if len(c.Methods) > 0 {
			return true
		}
	}
	return false
}

// GenerateConstsFile returns the content of abi/tolk_consts_generated.go.
// defs must be in the order that determines their iota offsets.
func GenerateConstsFile(defs []TolkInterfaceDef) (string, error) {
	out := &strings.Builder{}
	fmt.Fprintf(out, "// Code generated - DO NOT EDIT.\n\npackage abi\n\n")

	fmt.Fprintf(out, "const (\n")
	for i, d := range defs {
		// we skip _xmlContractInterfaceEnd + 0,
		// because in a debugger it will be shown as "_xmlContractInterfaceEnd"
		// instead of the first declaration here
		fmt.Fprintf(out, "\t%s ContractInterface = _xmlContractInterfaceEnd + %d\n", d.GoName, i+1)
	}
	fmt.Fprintf(out, ")\n\n")

	fmt.Fprintf(out, "func init() {\n")
	for _, d := range defs {
		fmt.Fprintf(out, "\ttolkContractInterfaceStrings[%s] = %q\n", d.GoName, d.StringName)
	}
	fmt.Fprintf(out, "\n")
	for _, d := range defs {
		fmt.Fprintf(out, "\ttolkContractInterfaceFromString[%q] = %s\n", d.StringName, d.GoName)
	}
	fmt.Fprintf(out, "}\n")

	return out.String(), nil
}

// GenerateMsgOpsFile returns a group-level constants file content for MsgOp names.
// Constants are named with MsgOp suffix, e.g. "CocoonPayoutMsgOp".
func GenerateMsgOpsFile(pkgName string, prefixed []PrefixedStruct) (string, error) {
	if len(prefixed) == 0 {
		return "", nil
	}
	out := &strings.Builder{}
	fmt.Fprintf(out, "// Code generated - DO NOT EDIT.\n\npackage %s\n\n", pkgName)
	fmt.Fprintf(out, "const (\n")
	for _, p := range prefixed {
		fmt.Fprintf(out, "\t%s = %q\n", p.MsgOpConst, p.MsgOpName)
	}
	fmt.Fprintf(out, ")\n")
	return out.String(), nil
}
