package catalog_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/abi-tolk/catalog"
)

func TestReleaseURL(t *testing.T) {
	assert.Equal(t,
		"https://github.com/ton-blockchain/abis/releases/download/v0.1.1/abi-catalog.json",
		catalog.ReleaseURL("v0.1.1"))
}

func TestDownload(t *testing.T) {
	if testing.Short() {
		t.Skip("downloads a catalog release")
	}
	c, err := catalog.Download(context.Background(), catalog.ReleaseURL("v0.1.1"))
	require.NoError(t, err)
	assert.Len(t, c.Projects, 61)

	p, ok := c.Project("wallets")
	require.True(t, ok)
	contract, ok := p.Contract("WalletV5r1")
	require.True(t, ok)
	assert.NotNil(t, contract.ABI.Storage.StorageTyIdx)
}
