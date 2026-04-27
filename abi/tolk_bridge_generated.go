// Code generated - DO NOT EDIT.

package abi

import (
	"context"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	abiFfVault "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/ffVault"
	abiPythOracle "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/pythOracle"
	"github.com/tonkeeper/tongo/ton"
)

func init() {
	tolkMethods = append(tolkMethods,
		MethodDescription{
			Name: "get_cocoon_client_data",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiCocoon.GetCocoonClientData(ctx, executor, id)
				return "GetCocoonClientData_CocoonClientResult", r, err
			},
		},
		MethodDescription{
			Name: "get_cocoon_proxy_data",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiCocoon.GetCocoonProxyData(ctx, executor, id)
				return "GetCocoonProxyData_CocoonProxyResult", r, err
			},
		},
		MethodDescription{
			Name: "last_proxy_seqno",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiCocoon.GetLastProxySeqno(ctx, executor, id)
				return "GetLastProxySeqno_CocoonRootResult", r, err
			},
		},
		MethodDescription{
			Name: "get_cocoon_data",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiCocoon.GetCocoonData(ctx, executor, id)
				return "GetCocoonData_CocoonRootResult", r, err
			},
		},
		MethodDescription{
			Name: "get_cur_params",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiCocoon.GetCurParams(ctx, executor, id)
				return "GetCurParams_CocoonRootResult", r, err
			},
		},
		MethodDescription{
			Name: "seqno",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiCocoon.GetSeqno(ctx, executor, id)
				return "GetSeqno_CocoonWalletResult", r, err
			},
		},
		MethodDescription{
			Name: "get_public_key",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiCocoon.GetPublicKey(ctx, executor, id)
				return "GetPublicKey_CocoonWalletResult", r, err
			},
		},
		MethodDescription{
			Name: "get_owner_address",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiCocoon.GetOwnerAddress(ctx, executor, id)
				return "GetOwnerAddress_CocoonWalletResult", r, err
			},
		},
		MethodDescription{
			Name: "get_cocoon_worker_data",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiCocoon.GetCocoonWorkerData(ctx, executor, id)
				return "GetCocoonWorkerData_CocoonWorkerResult", r, err
			},
		},
	)

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    CocoonClient,
			Results: []string{"GetCocoonClientData_CocoonClientResult"},
		},
		InterfaceDescription{
			Name:    CocoonProxy,
			Results: []string{"GetCocoonProxyData_CocoonProxyResult"},
		},
		InterfaceDescription{
			Name:    CocoonRoot,
			Results: []string{"GetLastProxySeqno_CocoonRootResult", "GetCocoonData_CocoonRootResult", "GetCurParams_CocoonRootResult"},
		},
		InterfaceDescription{
			Name:    CocoonWallet,
			Results: []string{"GetSeqno_CocoonWalletResult", "GetPublicKey_CocoonWalletResult", "GetOwnerAddress_CocoonWalletResult"},
		},
		InterfaceDescription{
			Name:    CocoonWorker,
			Results: []string{"GetCocoonWorkerData_CocoonWorkerResult"},
		},
	)

	registerInMsgUnmarshalerForOpcode[*abiCocoon.AddModelType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixAddModelType), abiCocoon.CocoonAddModelTypeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.AddProxyType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixAddProxyType), abiCocoon.CocoonAddProxyTypeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.AddWorkerType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixAddWorkerType), abiCocoon.CocoonAddWorkerTypeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ChangeFees](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixChangeFees), abiCocoon.CocoonChangeFeesMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ChangeOwner](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixChangeOwner), abiCocoon.CocoonChangeOwnerMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ChangeParams](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixChangeParams), abiCocoon.CocoonChangeParamsMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ClientProxyRequest](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixClientProxyRequest), abiCocoon.CocoonClientProxyRequestMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.DelModelType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixDelModelType), abiCocoon.CocoonDelModelTypeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.DelProxyType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixDelProxyType), abiCocoon.CocoonDelProxyTypeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.DelWorkerType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixDelWorkerType), abiCocoon.CocoonDelWorkerTypeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtClientChargeSigned](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtClientChargeSigned), abiCocoon.CocoonExtClientChargeSignedMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtClientGrantRefundSigned](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtClientGrantRefundSigned), abiCocoon.CocoonExtClientGrantRefundSignedMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtClientTopUp](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtClientTopUp), abiCocoon.CocoonExtClientTopUpMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtProxyCloseCompleteRequestSigned](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtProxyCloseCompleteRequestSigned), abiCocoon.CocoonExtProxyCloseCompleteRequestSignedMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtProxyCloseRequestSigned](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtProxyCloseRequestSigned), abiCocoon.CocoonExtProxyCloseRequestSignedMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtProxyIncreaseStake](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtProxyIncreaseStake), abiCocoon.CocoonExtProxyIncreaseStakeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtProxyPayoutRequest](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtProxyPayoutRequest), abiCocoon.CocoonExtProxyPayoutRequestMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientChangeSecretHash](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientChangeSecretHash), abiCocoon.CocoonOwnerClientChangeSecretHashMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientChangeSecretHashAndTopUp](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientChangeSecretHashAndTopUp), abiCocoon.CocoonOwnerClientChangeSecretHashAndTopUpMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientIncreaseStake](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientIncreaseStake), abiCocoon.CocoonOwnerClientIncreaseStakeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientRegister](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientRegister), abiCocoon.CocoonOwnerClientRegisterMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientRequestRefund](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientRequestRefund), abiCocoon.CocoonOwnerClientRequestRefundMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientWithdraw](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientWithdraw), abiCocoon.CocoonOwnerClientWithdrawMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerProxyClose](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerProxyClose), abiCocoon.CocoonOwnerProxyCloseMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerWalletSendMessage](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerWalletSendMessage), abiCocoon.CocoonOwnerWalletSendMessageMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.Payout](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixPayout), abiCocoon.CocoonPayoutMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.RegisterProxy](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixRegisterProxy), abiCocoon.CocoonRegisterProxyMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ResetRoot](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixResetRoot), abiCocoon.CocoonResetRootMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ReturnExcessesBack](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixReturnExcessesBack), abiCocoon.CocoonReturnExcessesBackMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.TextCmd](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixTextCmd), abiCocoon.CocoonTextCmdMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.TextCommand](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixTextCommand), abiCocoon.CocoonTextCommandMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.UnregisterProxy](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixUnregisterProxy), abiCocoon.CocoonUnregisterProxyMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.UpdateProxy](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixUpdateProxy), abiCocoon.CocoonUpdateProxyMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.UpgradeCode](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixUpgradeCode), abiCocoon.CocoonUpgradeCodeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.UpgradeContracts](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixUpgradeContracts), abiCocoon.CocoonUpgradeContractsMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.UpgradeFull](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixUpgradeFull), abiCocoon.CocoonUpgradeFullMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiCocoon.WorkerProxyRequest](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixWorkerProxyRequest), abiCocoon.CocoonWorkerProxyRequestMsgOp)

}

func init() {
	tolkMethods = append(tolkMethods,
		MethodDescription{
			Name: "get_nft_data",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiFfVault.GetNftData(ctx, executor, id)
				return "GetNftData_FfVaultPositionResult", r, err
			},
		},
		MethodDescription{
			Name: "get_stake_position_info",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiFfVault.GetStakePositionInfo(ctx, executor, id)
				return "GetStakePositionInfo_FfVaultPositionResult", r, err
			},
		},
		MethodDescription{
			Name: "get_collection_data",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiFfVault.GetCollectionData(ctx, executor, id)
				return "GetCollectionData_FfVaultResult", r, err
			},
		},
		MethodDescription{
			Name: "get_staking_data",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiFfVault.GetStakingData(ctx, executor, id)
				return "GetStakingData_FfVaultResult", r, err
			},
		},
		MethodDescription{
			Name: "get_balance",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiFfVault.GetBalance(ctx, executor, id)
				return "GetBalance_FfVaultResult", r, err
			},
		},
	)

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    FfVaultPosition,
			Results: []string{"GetNftData_FfVaultPositionResult", "GetStakePositionInfo_FfVaultPositionResult"},
		},
		InterfaceDescription{
			Name:    FfVault,
			Results: []string{"GetCollectionData_FfVaultResult", "GetStakingData_FfVaultResult", "GetBalance_FfVaultResult"},
		},
	)

	registerInMsgUnmarshalerForOpcode[*abiFfVault.AssetDeposit](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixAssetDeposit), abiFfVault.FfVaultAssetDepositMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiFfVault.StableDeposit](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixStableDeposit), abiFfVault.FfVaultStableDepositMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiFfVault.StakeOperation](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixStakeOperation), abiFfVault.FfVaultStakeOperationMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeExecute](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeExecute), abiFfVault.FfVaultUnstakeExecuteMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeExecuteCancel](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeExecuteCancel), abiFfVault.FfVaultUnstakeExecuteCancelMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeExecuteInternal](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeExecuteInternal), abiFfVault.FfVaultUnstakeExecuteInternalMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeExecuteInternalCallback](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeExecuteInternalCallback), abiFfVault.FfVaultUnstakeExecuteInternalCallbackMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeOperation](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeOperation), abiFfVault.FfVaultUnstakeOperationMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeRequest](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeRequest), abiFfVault.FfVaultUnstakeRequestMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiFfVault.VaultStorage](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixVaultStorage), abiFfVault.FfVaultVaultStorageMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiFfVault.WithdrawJetton](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixWithdrawJetton), abiFfVault.FfVaultWithdrawJettonMsgOp)

}

func init() {
	tolkMethods = append(tolkMethods,
		MethodDescription{
			Name: "get_single_update_fee",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiPythOracle.GetSingleUpdateFee(ctx, executor, id)
				return "GetSingleUpdateFee_PythOracleResult", r, err
			},
		},
		MethodDescription{
			Name: "get_governance_data_source_index",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiPythOracle.GetGovernanceDataSourceIndex(ctx, executor, id)
				return "GetGovernanceDataSourceIndex_PythOracleResult", r, err
			},
		},
		MethodDescription{
			Name: "get_governance_data_source",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiPythOracle.GetGovernanceDataSource(ctx, executor, id)
				return "GetGovernanceDataSource_PythOracleResult", r, err
			},
		},
		MethodDescription{
			Name: "get_last_executed_governance_sequence",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiPythOracle.GetLastExecutedGovernanceSequence(ctx, executor, id)
				return "GetLastExecutedGovernanceSequence_PythOracleResult", r, err
			},
		},
		MethodDescription{
			Name: "get_chain_id",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiPythOracle.GetChainId(ctx, executor, id)
				return "GetChainId_PythOracleResult", r, err
			},
		},
		MethodDescription{
			Name: "get_current_guardian_set_index",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiPythOracle.GetCurrentGuardianSetIndex(ctx, executor, id)
				return "GetCurrentGuardianSetIndex_PythOracleResult", r, err
			},
		},
		MethodDescription{
			Name: "get_governance_chain_id",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiPythOracle.GetGovernanceChainId(ctx, executor, id)
				return "GetGovernanceChainId_PythOracleResult", r, err
			},
		},
		MethodDescription{
			Name: "get_governance_contract",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiPythOracle.GetGovernanceContract(ctx, executor, id)
				return "GetGovernanceContract_PythOracleResult", r, err
			},
		},
	)

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    PythOracle,
			Results: []string{"GetSingleUpdateFee_PythOracleResult", "GetGovernanceDataSourceIndex_PythOracleResult", "GetGovernanceDataSource_PythOracleResult", "GetLastExecutedGovernanceSequence_PythOracleResult", "GetChainId_PythOracleResult", "GetCurrentGuardianSetIndex_PythOracleResult", "GetGovernanceChainId_PythOracleResult", "GetGovernanceContract_PythOracleResult"},
		},
	)

	registerInMsgUnmarshalerForOpcode[*abiPythOracle.ErrorResponse](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixErrorResponse), abiPythOracle.PythOracleErrorResponseMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.ExecuteGovernanceActionMessage](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixExecuteGovernanceActionMessage), abiPythOracle.PythOracleExecuteGovernanceActionMessageMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.OracleResponseSuccess](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixOracleResponseSuccess), abiPythOracle.PythOracleOracleResponseSuccessMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.ParsePriceFeedUpdatesMessage](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixParsePriceFeedUpdatesMessage), abiPythOracle.PythOracleParsePriceFeedUpdatesMessageMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.ParseUniquePriceFeedUpdatesMessage](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixParseUniquePriceFeedUpdatesMessage), abiPythOracle.PythOracleParseUniquePriceFeedUpdatesMessageMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.SuccessResponse](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixSuccessResponse), abiPythOracle.PythOracleSuccessResponseMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.UpdateGuardianSetMessage](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixUpdateGuardianSetMessage), abiPythOracle.PythOracleUpdateGuardianSetMessageMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.UpdatePriceFeedsMessage](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixUpdatePriceFeedsMessage), abiPythOracle.PythOracleUpdatePriceFeedsMessageMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.UpgradeContractMessage](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixUpgradeContractMessage), abiPythOracle.PythOracleUpgradeContractMessageMsgOp)

}
