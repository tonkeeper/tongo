package tolkgen

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/utils"
)

// JettonPayloadType describes one concrete Jetton payload type discovered in Tolk ABIs.
type JettonPayloadType struct {
	GroupName string
	TypeName  string
	OpCode    uint64
}

// JettonRegistration is a render-ready descriptor for abi/tolk_jetton_msg_types.go.
type JettonRegistration struct {
	ImportAlias string
	ImportPath  string
	TypeName    string
	OpConst     string
	OpCodeConst string
	OpName      string
	OpCode      uint64
}

const jettonNotifyMsgOpCode uint64 = 0x7362d09c

// StripJettonNotifyWrappers removes generic notification wrappers from declarations and incoming messages,
// and returns concrete Jetton payload types referenced in wrapper specializations.
//
// Wrapper detection is structural:
// - struct has one type parameter
// - struct has exactly one field
// - field type is genericT with the same type parameter name
// - incoming message references this struct with one type argument
func StripJettonNotifyWrappers(abi parser.ABI) (parser.ABI, []JettonPayloadType) {
	structs := map[string]parser.StructDeclaration{}
	aliases := map[string]parser.AliasDeclaration{}
	for _, decl := range abi.Declarations {
		switch decl.SumType {
		case parser.DeclarationKindStruct:
			structs[decl.StructDeclaration.Name] = decl.StructDeclaration
		case parser.DeclarationKindAlias:
			aliases[decl.AliasDeclaration.Name] = decl.AliasDeclaration
		}
	}

	wrappers := map[string]struct{}{}
	payloadByType := map[string]JettonPayloadType{}

	for _, msg := range abi.IncomingMessages {
		body := msg.BodyTy
		if body.SumType != parser.TyKindStructRef {
			continue
		}
		sName := body.StructRef.StructName
		sDecl, ok := structs[sName]
		if !ok || !isJettonPayloadWrapper(sDecl) || len(body.StructRef.TypeArgs) != 1 {
			continue
		}
		wrappers[sName] = struct{}{}

		for _, p := range collectPayloadTypes(body.StructRef.TypeArgs[0], structs, aliases, nil) {
			prev, exists := payloadByType[p.TypeName]
			if !exists || prev.OpCode == 0 {
				payloadByType[p.TypeName] = p
			}
		}
	}

	if len(wrappers) == 0 {
		return abi, nil
	}

	filtered := abi
	filtered.Declarations = make([]parser.Declaration, 0, len(abi.Declarations))
	for _, decl := range abi.Declarations {
		if decl.SumType == parser.DeclarationKindStruct {
			if _, skip := wrappers[decl.StructDeclaration.Name]; skip {
				continue
			}
		}
		filtered.Declarations = append(filtered.Declarations, decl)
	}

	filtered.IncomingMessages = make([]parser.IncomingMessage, 0, len(abi.IncomingMessages))
	for _, msg := range abi.IncomingMessages {
		keep := true
		if msg.BodyTy.SumType == parser.TyKindStructRef {
			_, keep = wrappers[msg.BodyTy.StructRef.StructName]
			keep = !keep
		}
		if keep {
			filtered.IncomingMessages = append(filtered.IncomingMessages, msg)
		}
	}

	payloads := make([]JettonPayloadType, 0, len(payloadByType))
	for _, p := range payloadByType {
		payloads = append(payloads, p)
	}
	sort.Slice(payloads, func(i, j int) bool {
		if payloads[i].TypeName == payloads[j].TypeName {
			return payloads[i].OpCode < payloads[j].OpCode
		}
		return payloads[i].TypeName < payloads[j].TypeName
	})

	return filtered, payloads
}

func isJettonPayloadWrapper(decl parser.StructDeclaration) bool {
	if len(decl.TypeParams) != 1 || len(decl.Fields) != 1 {
		return false
	}
	if op, ok := parseOpcode32(decl.Prefix); !ok || op != jettonNotifyMsgOpCode {
		return false
	}
	fieldTy := decl.Fields[0].Ty
	if fieldTy.SumType != parser.TyKindGenericT {
		return false
	}
	return fieldTy.Generic.NameT == decl.TypeParams[0]
}

func collectPayloadTypes(
	ty parser.Ty,
	structs map[string]parser.StructDeclaration,
	aliases map[string]parser.AliasDeclaration,
	aliasStack map[string]struct{},
) []JettonPayloadType {
	switch ty.SumType {
	case parser.TyKindStructRef:
		if len(ty.StructRef.TypeArgs) > 0 {
			var out []JettonPayloadType
			for _, a := range ty.StructRef.TypeArgs {
				out = append(out, collectPayloadTypes(a, structs, aliases, aliasStack)...)
			}
			return out
		}
		if s, ok := structs[ty.StructRef.StructName]; ok {
			if op, ok := parseOpcode32(s.Prefix); ok {
				return []JettonPayloadType{{TypeName: s.Name, OpCode: op}}
			}
		}
		return nil
	case parser.TyKindUnion:
		out := make([]JettonPayloadType, 0, len(ty.Union.Variants))
		for _, v := range ty.Union.Variants {
			if v.VariantTy.SumType == parser.TyKindStructRef {
				if op, ok := parseOpcode32(&parser.Prefix{PrefixStr: v.PrefixStr, PrefixLen: v.PrefixLen}); ok {
					out = append(out, JettonPayloadType{
						TypeName: v.VariantTy.StructRef.StructName,
						OpCode:   op,
					})
					continue
				}
			}
			out = append(out, collectPayloadTypes(v.VariantTy, structs, aliases, aliasStack)...)
		}
		return out
	case parser.TyKindAliasRef:
		alias, ok := aliases[ty.AliasRef.AliasName]
		if !ok {
			return nil
		}
		if aliasStack == nil {
			aliasStack = map[string]struct{}{}
		}
		if _, loop := aliasStack[alias.Name]; loop {
			return nil
		}
		aliasStack[alias.Name] = struct{}{}
		target := alias.TargetTy
		if len(alias.TypeParams) > 0 {
			target = target.InstantiateGenerics(alias.TypeParams, ty.AliasRef.TypeArgs)
		}
		out := collectPayloadTypes(target, structs, aliases, aliasStack)
		delete(aliasStack, alias.Name)
		return out
	case parser.TyKindNullable:
		return collectPayloadTypes(ty.Nullable.Inner, structs, aliases, aliasStack)
	case parser.TyKindCellOf:
		return collectPayloadTypes(ty.CellOf.Inner, structs, aliases, aliasStack)
	case parser.TyKindArrayOf:
		return collectPayloadTypes(ty.ArrayOf.Inner, structs, aliases, aliasStack)
	default:
		return nil
	}
}

func parseOpcode32(p *parser.Prefix) (uint64, bool) {
	if p == nil || p.PrefixLen != 32 || p.PrefixStr == "" {
		return 0, false
	}
	v, err := strconv.ParseUint(p.PrefixStr, 0, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// BuildJettonRegistrations converts collected payload types into render-ready registration records.
func BuildJettonRegistrations(modulePath, outputDir string, payloads []JettonPayloadType) []JettonRegistration {
	seen := map[string]struct{}{}
	result := make([]JettonRegistration, 0, len(payloads))

	for _, p := range payloads {
		groupCamel := utils.ToCamelCase(p.GroupName)
		opBase := groupCamel + safeGoIdent(p.TypeName)
		key := opBase + fmt.Sprintf("#%08x", p.OpCode)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		result = append(result, JettonRegistration{
			ImportAlias: "abi" + groupCamel,
			ImportPath:  modulePath + "/" + strings.ReplaceAll(outputDir, "\\", "/") + "/" + p.GroupName,
			TypeName:    safeGoIdent(p.TypeName),
			OpConst:     opBase + "JettonOp",
			OpCodeConst: opBase + "JettonOpCode",
			OpName:      opBase,
			OpCode:      p.OpCode,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].OpCode == result[j].OpCode {
			return result[i].OpConst < result[j].OpConst
		}
		return result[i].OpCode < result[j].OpCode
	})
	return result
}

// GenerateTolkJettonMsgTypesFile renders abi/tolk_jetton_msg_types.go.
func GenerateTolkJettonMsgTypesFile(regs []JettonRegistration) (string, error) {
	out := &strings.Builder{}
	out.WriteString("// Code generated - DO NOT EDIT.\n\npackage abi\n\n")

	if len(regs) > 0 {
		imports := map[string]string{}
		for _, r := range regs {
			imports[r.ImportAlias] = r.ImportPath
		}
		aliases := make([]string, 0, len(imports))
		for a := range imports {
			aliases = append(aliases, a)
		}
		sort.Strings(aliases)
		out.WriteString("import (\n")
		for _, alias := range aliases {
			fmt.Fprintf(out, "\t%s %q\n", alias, imports[alias])
		}
		out.WriteString(")\n\n")
	}

	out.WriteString("const (\n")
	for _, r := range regs {
		fmt.Fprintf(out, "\t%s JettonOpName = %q\n", r.OpConst, r.OpName)
	}
	for _, r := range regs {
		fmt.Fprintf(out, "\t%s JettonOpCode = 0x%08x\n", r.OpCodeConst, r.OpCode)
	}
	out.WriteString(")\n\n")

	out.WriteString("var tolkKnownJettonTypes = map[JettonOpName]any{\n")
	for _, r := range regs {
		fmt.Fprintf(out, "\t%s: %s.%s{},\n", r.OpConst, r.ImportAlias, r.TypeName)
	}
	out.WriteString("}\n\n")

	out.WriteString("var tolkJettonOpCodes = map[JettonOpName]JettonOpCode{\n")
	for _, r := range regs {
		fmt.Fprintf(out, "\t%s: %s,\n", r.OpConst, r.OpCodeConst)
	}
	out.WriteString("}\n\n")

	out.WriteString("var tolkJettonDecodersMapping = map[JettonOpCode]jettonDecoder{\n")
	for _, r := range regs {
		fmt.Fprintf(out, "\t%s: decodeJettonPayload[%s.%s](%s, %s, false, false),\n",
			r.OpCodeConst, r.ImportAlias, r.TypeName, r.OpConst, r.OpCodeConst)
	}
	out.WriteString("}\n")
	return out.String(), nil
}
