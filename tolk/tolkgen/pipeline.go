package tolkgen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tonkeeper/tongo/tolk/parser"
	"github.com/tonkeeper/tongo/utils"
	"golang.org/x/exp/maps"
)

const DefaultModulePath = "github.com/tonkeeper/tongo"

type CodegenPipelineConfig struct {
	SchemasDir   string
	OutputDir    string
	ABIOutputDir string
	ModulePath   string
}

func DefaultCodegenPipelineConfig() CodegenPipelineConfig {
	return CodegenPipelineConfig{
		SchemasDir:   "abi-tolk/schemas",
		OutputDir:    "abi-tolk/abiGenerated",
		ABIOutputDir: "abi",
		ModulePath:   DefaultModulePath,
	}
}

type schemaEntry struct {
	outPath string
	abi     parser.ABI
}

func GenerateFromSchemas(cfg CodegenPipelineConfig) error {
	if cfg.ModulePath == "" {
		cfg.ModulePath = DefaultModulePath
	}

	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.MkdirAll(cfg.ABIOutputDir, 0755); err != nil {
		return fmt.Errorf("create abi output dir: %w", err)
	}

	// Phase 1: collect all schemas grouped by output directory.
	groups := make(map[string][]schemaEntry)
	var tolkJettonPayloads []JettonPayloadType

	if err := filepath.WalkDir(cfg.SchemasDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		var abi parser.ABI
		if err := json.Unmarshal(data, &abi); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}

		rel, _ := filepath.Rel(cfg.SchemasDir, path)
		outPath := filepath.Join(cfg.OutputDir, strings.TrimSuffix(rel, ".json")+".go")
		outDir := filepath.Dir(outPath)
		groupName := filepath.Base(outDir)

		filteredABI, payloads := StripJettonNotifyWrappers(abi)
		for _, p := range payloads {
			p.GroupName = groupName
			tolkJettonPayloads = append(tolkJettonPayloads, p)
		}

		groups[outDir] = append(groups[outDir], schemaEntry{outPath: outPath, abi: filteredABI})
		return nil
	}); err != nil {
		return err
	}

	groupKeys := maps.Keys(groups)
	sort.Strings(groupKeys)

	// Phase 2: deduplicate within each group and generate Go type files.
	for _, outDir := range groupKeys {
		if err := processGroup(outDir, groups[outDir], cfg.OutputDir); err != nil {
			return err
		}
	}

	// Phase 3: generate abi bridge files and constants file.
	if err := generateBridgeFiles(groupKeys, groups, cfg.OutputDir, cfg.ABIOutputDir, cfg.ModulePath); err != nil {
		return err
	}

	// Phase 4: generate Tolk-provided Jetton payload registrations.
	regs := BuildJettonRegistrations(cfg.ModulePath, cfg.OutputDir, tolkJettonPayloads)
	tolkJettonCode, err := GenerateTolkJettonMsgTypesFile(regs)
	if err != nil {
		return err
	}
	if err := utils.WriteFormattedGoCode(filepath.Join(cfg.ABIOutputDir, "tolk_jetton_msg_types.go"), tolkJettonCode); err != nil {
		return err
	}

	return nil
}

// generateBridgeFiles generates:
//   - <abiOutputDir>/tolk_consts_generated.go  — ContractInterface constants for all Tolk contracts
//   - <abiOutputDir>/tolk_bridge_generated.go  — combined bridge for all schema groups
func generateBridgeFiles(groupKeys []string, groups map[string][]schemaEntry, outputDir, abiOutputDir, modulePath string) error {
	var allDefs []TolkInterfaceDef
	var allGroups []BridgeGroup

	for _, outDir := range groupKeys {
		entries := groups[outDir]
		groupName := filepath.Base(outDir)

		var schemaEntries []SchemaEntry
		for _, e := range entries {
			baseName := strings.TrimSuffix(filepath.Base(e.outPath), ".go")
			allDefs = append(allDefs, TolkInterfaceDef{
				GoName:     InterfaceNameFromABI(e.abi, baseName),
				StringName: StringNameFromABI(e.abi, baseName),
			})
			schemaEntries = append(schemaEntries, SchemaEntry{BaseName: baseName, ABI: e.abi})
		}

		allGroups = append(allGroups, BridgeGroupFromEntries(groupName, modulePath, outputDir, schemaEntries))
	}

	bridgeContent, err := GenerateAllBridgesFile(allGroups)
	if err != nil {
		return fmt.Errorf("generate bridge: %w", err)
	}
	if err := utils.WriteFormattedGoCode(filepath.Join(abiOutputDir, "tolk_bridge_generated.go"), bridgeContent); err != nil {
		return err
	}

	constsContent, err := GenerateConstsFile(allDefs)
	if err != nil {
		return fmt.Errorf("generate consts: %w", err)
	}
	return utils.WriteFormattedGoCode(filepath.Join(abiOutputDir, "tolk_consts_generated.go"), constsContent)
}

func processGroup(outDir string, entries []schemaEntry, rootOutputDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", outDir, err)
	}

	rel, err := filepath.Rel(rootOutputDir, outDir)
	if err != nil || rel == "." || rel == "" {
		rel = "generated"
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	pkgName := "abi" + utils.ToCamelCase(parts[len(parts)-1])
	groupName := filepath.Base(outDir)
	camelGroupName := utils.ToCamelCase(groupName)

	abisByPath := make(map[string]parser.ABI, len(entries))
	groupABIs := make([]parser.ABI, 0, len(entries))
	for _, e := range entries {
		abisByPath[e.outPath] = e.abi
		groupABIs = append(groupABIs, e.abi)
	}

	deduped, sharedABI, err := ExtractShared(abisByPath)
	if err != nil {
		return err
	}

	if sharedABI != nil {
		if err := writeGoFile(filepath.Join(outDir, "shared.go"), pkgName, *sharedABI, ""); err != nil {
			return err
		}
	}

	msgOpsCode, err := GenerateMsgOpsFile(pkgName, CollectPrefixedStructs(groupABIs, camelGroupName))
	if err != nil {
		return fmt.Errorf("msg ops %s: %w", outDir, err)
	}
	if msgOpsCode != "" {
		if err := utils.WriteFormattedGoCode(filepath.Join(outDir, "msg_ops_generated.go"), msgOpsCode); err != nil {
			return err
		}
	}

	var execEntries []ExecutorEntry
	for _, e := range entries {
		if len(e.abi.GetMethods) == 0 {
			continue
		}
		baseName := strings.TrimSuffix(filepath.Base(e.outPath), ".go")
		gen, err := NewTolkGolangGenerator(e.abi)
		if err != nil {
			return fmt.Errorf("generator for %s: %w", e.outPath, err)
		}
		execEntries = append(execEntries, ExecutorEntry{
			IfaceName: utils.ToCamelCase(baseName),
			Gen:       gen,
		})
	}

	if len(execEntries) > 0 {
		code, err := GenerateExecutorFile(pkgName, execEntries)
		if err != nil {
			return fmt.Errorf("executor %s: %w", outDir, err)
		}
		if err := utils.WriteFormattedGoCode(filepath.Join(outDir, "executor.go"), code); err != nil {
			return err
		}
	}

	for _, e := range entries {
		abi := deduped[e.outPath]
		abi.GetMethods = e.abi.GetMethods
		abi.Storage = e.abi.Storage
		abi.IncomingExternal = e.abi.IncomingExternal
		abi.IncomingMessages = e.abi.IncomingMessages
		ifaceName := ""
		if len(abi.GetMethods) > 0 || abi.Storage.StorageTy != nil {
			baseName := strings.TrimSuffix(filepath.Base(e.outPath), ".go")
			ifaceName = utils.ToCamelCase(baseName)
		}
		if err := writeGoFile(e.outPath, pkgName, abi, ifaceName); err != nil {
			return err
		}
	}

	return nil
}

func writeGoFile(path, pkgName string, abi parser.ABI, ifaceName string) error {
	gen, err := NewTolkGolangGenerator(abi)
	if err != nil {
		return fmt.Errorf("generator for %s: %w", path, err)
	}

	code, err := gen.GenerateContractFile(pkgName, ifaceName)
	if err != nil {
		return fmt.Errorf("codegen %s: %w", path, err)
	}
	if code != "" {
		if err := utils.WriteFormattedGoCode(path, code); err != nil {
			return err
		}
	}

	marshalCode, err := gen.GenerateMarshalFile(pkgName)
	if err != nil {
		return fmt.Errorf("marshal codegen %s: %w", path, err)
	}
	if marshalCode != "" {
		marshalPath := strings.TrimSuffix(path, ".go") + "_marshal.go"
		if err := utils.WriteFormattedGoCode(marshalPath, marshalCode); err != nil {
			return err
		}
	}

	return nil
}
