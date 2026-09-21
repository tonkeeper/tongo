package catalog_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/abi-tolk/catalog"
)

func TestLoad(t *testing.T) {
	c := catalog.GetBuiltin()
	assert.Greater(t, len(c.Projects), 50)
	assert.True(t, slices.IsSortedFunc(c.Projects, func(a, b *catalog.Project) int {
		return strings.Compare(a.Name, b.Name)
	}), "projects are sorted by name")
	for _, p := range c.Projects {
		assert.True(t, slices.IsSortedFunc(p.Contracts, func(a, b *catalog.Contract) int {
			return strings.Compare(a.Name, b.Name)
		}), "%s: contracts are sorted by name", p.Name)
	}

	p, ok := c.Project("airdrop")
	require.True(t, ok)
	contract, ok := p.Contract("AirdropInterlockerV2")
	require.True(t, ok)
	assert.Same(t, p, contract.Project())
	assert.Equal(t, "Airdrop Interlocker v2", contract.DisplayName)
	assert.Equal(t, []string{"0ff44a96fc2481111236d8d9920fd6e29b6108215e3ec5ba761ead9602adada4"}, contract.Hashes)
	assert.NotEmpty(t, contract.KnownAddresses)

	_, ok = c.Project("nothing")
	assert.False(t, ok)
	_, ok = p.Contract("nothing")
	assert.False(t, ok)
}

// TestMergedProject covers the wallets, which the bundle spreads over a
// "wallets" project and one project per wallet directory below it.
func TestMergedProject(t *testing.T) {
	c := catalog.GetBuiltin()
	p, ok := c.Project("wallets")
	require.True(t, ok)

	names := make([]string, 0, len(p.Contracts))
	for _, contract := range p.Contracts {
		names = append(names, contract.Name)
	}
	assert.Contains(t, names, "WalletV5r1", "declared by the wallets project itself")
	assert.Contains(t, names, "WalletV1r1", "declared by wallets/w1r1")

	_, ok = c.Project("wallets/w1r1")
	assert.False(t, ok, "a wallet does not make a project of its own")
	for _, project := range c.Projects {
		assert.NotContains(t, project.Name, "/", "a project name is one path segment")
	}
}

func TestContractABI(t *testing.T) {
	p, ok := catalog.GetBuiltin().Project("wallets")
	require.True(t, ok)
	contract, ok := p.Contract("WalletV5r1")
	require.True(t, ok)

	abi := contract.ABI
	assert.Equal(t, contract.Name, abi.ContractName)
	require.NotNil(t, abi.Storage.StorageTyIdx, "the ABI describes a storage to decode")
	assert.NotEmpty(t, abi.UniqueTypes)
	assert.NotEmpty(t, abi.IncomingMessages)

	methods := make([]string, 0, len(abi.GetMethods))
	for _, method := range abi.GetMethods {
		methods = append(methods, method.Name)
	}
	assert.Contains(t, methods, "get_subwallet_id")

	kinds := make([]string, 0, len(contract.Links))
	for _, link := range contract.Links {
		assert.True(t, strings.HasPrefix(link.URL, "https://"), "link %q", link.Title)
		kinds = append(kinds, link.Kind)
	}
	assert.Equal(t, []string{"docs", "source", "spec"}, kinds, "a contract inherits the links of its project")
}

func TestLoadFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "abi-catalog.json")
	require.NoError(t, os.WriteFile(path, []byte(`{
		"schemaVersion": 1,
		"contracts": [{"id": "wallets/w1r1.WalletV1r1", "compilerAbi": {"contract_name": "WalletV1r1"}}]
	}`), 0o644))

	c, err := catalog.LoadFile(path)
	require.NoError(t, err)
	p, ok := c.Project("wallets")
	require.True(t, ok)
	contract, ok := p.Contract("WalletV1r1")
	require.True(t, ok)
	assert.Equal(t, "WalletV1r1", contract.ABI.ContractName)
}

func TestLoadErrors(t *testing.T) {
	_, err := catalog.Load(strings.NewReader(`{"schemaVersion": 0, "contracts": []}`))
	assert.ErrorContains(t, err, "want 1")

	_, err = catalog.Load(strings.NewReader(`{"schemaVersion": 1, "contracts": [{"id": "Wallet"}]}`))
	assert.ErrorContains(t, err, `"<project>.<contract>"`)

	_, err = catalog.Load(strings.NewReader(`{"schemaVersion": 1, "contracts": [
		{"id": "wallets.W"}, {"id": "wallets/w1r1.W"}
	]}`))
	assert.ErrorContains(t, err, "declared twice")

	_, err = catalog.LoadFile(filepath.Join(t.TempDir(), "nothing.json"))
	assert.ErrorIs(t, err, os.ErrNotExist)
}
