package tlb

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/boc"
)

const (
	// State update for account 0:1d7cf8c106bf7c1fcfd910679fde4162a3ae8f5191755783dc0d80b6e6602cc4 from block 40776306
	prunedAccountBoc = "te6ccgEBAwEAiQACZ8AB18+MEGv3wfz9kQZ5/eQWKjro9RkXVXg9wNgLbmYCzEIshYZDLC1neAAAnd6+eoEME0ABAihIAQH+tf9oIOL/DZSD5+DWLIF9hGeJ+0rlgMh4hm2VnavVwAAHAFEAAABqKamjF9Bck0MxvE0DsxXRadmJimTZeBNTMORgVZurR4LY7NW3QA=="
	walletV4R2Boc    = "te6cckECFAEAAtQAART/APSkE/S88sgLAQIBIAIDAgFIBAUE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8QERITAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNBgcCASAICQB4AfoA9AQw+CdvIjBQCqEhvvLgUIIQcGx1Z4MesXCAGFAEywUmzxZY+gIZ9ADLaRfLH1Jgyz8gyYBA+wAGAIpQBIEBCPRZMO1E0IEBQNcgyAHPFvQAye1UAXKwjiOCEGRzdHKDHrFwgBhQBcsFUAPPFiP6AhPLassfyz/JgED7AJJfA+ICASAKCwBZvSQrb2omhAgKBrkPoCGEcNQICEekk30pkQzmkD6f+YN4EoAbeBAUiYcVnzGEAgFYDA0AEbjJftRNDXCx+AA9sp37UTQgQFA1yH0BDACyMoHy//J0AGBAQj0Cm+hMYAIBIA4PABmtznaiaEAga5Drhf/AABmvHfaiaEAQa5DrhY/AAG7SB/oA1NQi+QAFyMoHFcv/ydB3dIAYyMsFywIizxZQBfoCFMtrEszMyXP7AMhAFIEBCPRR8qcCAHCBAQjXGPoA0z/IVCBHgQEI9FHyp4IQbm90ZXB0gBjIywXLAlAGzxZQBPoCFMtqEssfyz/Jc/sAAgBsgQEI1xj6ANM/MFIkgQEI9Fnyp4IQZHN0cnB0gBjIywXLAlAFzxZQA/oCE8tqyx8Syz/Jc/sAAAr0AMntVGliJeU="

	// what the block's HASH_UPDATE says this account hashes to after the block
	prunedAccountHash  = "a3a8380191d7e176c7c803ca531dcdcf16ddab48f94dba92a65ea722c8b3bc40"
	walletV4R2CodeHash = "feb5ff6820e2ff0d9483e7e0d62c817d846789fb4ae580c878866d959dabd5c0"
)

func TestUnmarshalKeepsPrunedBranch(t *testing.T) {
	cell, err := boc.DeserializeSinglRootBase64(prunedAccountBoc)
	require.Nil(t, err)

	var account Account
	require.Nil(t, Unmarshal(cell, &account))
	require.Equal(t, AccountActive, account.Status())

	stateInit := account.Account.Storage.State.AccountActive.StateInit

	// code is pruned (because not changed in transaction)
	require.True(t, stateInit.Code.Exists)
	code := stateInit.Code.Value.Value
	require.True(t, code.IsPruned())
	codeHash, err := code.PrunedHash()
	require.Nil(t, err)
	assert.Equal(t, walletV4R2CodeHash, hex.EncodeToString(codeHash[:]))

	// data is not pruned
	require.True(t, stateInit.Data.Exists)
	data := stateInit.Data.Value.Value
	assert.False(t, data.IsPruned())
	_, err = data.PrunedHash()
	require.NotNil(t, err, "data not pruned, should return error")

	t.Run("rebuild", func(t *testing.T) {
		rebuilt := boc.NewCell()
		require.Nil(t, Marshal(rebuilt, account))
		h, err := rebuilt.Hash256()
		require.Nil(t, err)
		require.Equal(t, prunedAccountHash, hex.EncodeToString(h[:]))
	})

	t.Run("restore pruned code", func(t *testing.T) {
		realCode, err := boc.DeserializeSinglRootBase64(walletV4R2Boc)
		require.Nil(t, err)
		realCodeHash, err := realCode.Hash256()
		require.Nil(t, err)
		require.Equal(t, walletV4R2CodeHash, hex.EncodeToString(realCodeHash[:]), "how is that?")

		restored := account
		restored.Account.Storage.State.AccountActive.StateInit.Code.Value.Value = *realCode
		require.False(t, restored.Account.Storage.State.AccountActive.StateInit.Code.Value.Value.IsPruned())

		rebuilt := boc.NewCell()
		require.Nil(t, Marshal(rebuilt, restored))
		h, err := rebuilt.Hash256()
		require.Nil(t, err)
		require.Equal(t, prunedAccountHash, hex.EncodeToString(h[:]))
	})
}
