package runtime

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/tonkeeper/tongo/tolk/parser"
)

func loadTestABI(t testing.TB, filename string) parser.ContractABI {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	var abi parser.ContractABI
	if err = json.Unmarshal(data, &abi); err != nil {
		t.Fatal(err)
	}
	return abi
}

func testTyIdx(t testing.TB, abi parser.ContractABI, typeName string) int {
	t.Helper()
	idx := parser.NewABIIndex(abi)
	if alias, ok := idx.Aliases[typeName]; ok {
		return alias.TyIdx
	}
	if strct, ok := idx.Structs[typeName]; ok {
		return strct.TyIdx
	}
	for i := range abi.UniqueTypes {
		rendered, err := idx.RenderTy(i)
		if err != nil {
			continue
		}
		if rendered == typeName {
			return i
		}
	}
	t.Fatalf("type %q not found in ABI", typeName)
	return 0
}
