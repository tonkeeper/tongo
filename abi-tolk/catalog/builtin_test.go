package catalog_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/abi-tolk/catalog"
)

func TestGetBuiltin(t *testing.T) {
	var c *catalog.Catalog
	require.NotPanics(t, func() {
		c = catalog.GetBuiltin()
	})
	require.NotEmpty(t, c.Projects)

	assert.Regexp(t, regexp.MustCompile(`^v\d+\.\d+\.\d+$`), catalog.BuiltinVersion())

	p, ok := c.Project("wallets")
	require.True(t, ok)
	contract, ok := p.Contract("WalletV5r1")
	require.True(t, ok)
	assert.NotNil(t, contract.ABI.Storage.StorageTyIdx, "the built-in ABIs carry a storage to decode")

	assert.Same(t, c, catalog.GetBuiltin(), "the built-in catalog is read once")
}

func TestLatestVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("asks github for the newest release")
	}
	tag, err := catalog.LatestVersion(context.Background())
	require.NoError(t, err)
	assert.Regexp(t, regexp.MustCompile(`^v\d+\.\d+\.\d+$`), tag)
}

func TestGetLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("may download a catalog release")
	}
	c, err := catalog.GetLatest(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, c.Projects)

	tag, err := catalog.LatestVersion(context.Background())
	require.NoError(t, err)
	if tag == catalog.BuiltinVersion() {
		assert.Same(t, catalog.GetBuiltin(), c, "the built-in catalog is up to date, so it should be reused")
	}
}
