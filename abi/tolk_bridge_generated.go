// Code generated - DO NOT EDIT.

package abi

import (
	"context"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	abiElector "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/elector"
	abiFfVault "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/ffVault"
	abiPythOracle "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/pythOracle"
	abiSingleNominatorPool "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/singleNominatorPool"
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

	KnownMsgInTypes[abiCocoon.CocoonAddModelTypeMsgOp] = abiCocoon.AddModelType{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.AddModelType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixAddModelType), abiCocoon.CocoonAddModelTypeMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonAddProxyTypeMsgOp] = abiCocoon.AddProxyType{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.AddProxyType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixAddProxyType), abiCocoon.CocoonAddProxyTypeMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonAddWorkerTypeMsgOp] = abiCocoon.AddWorkerType{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.AddWorkerType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixAddWorkerType), abiCocoon.CocoonAddWorkerTypeMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonChangeFeesMsgOp] = abiCocoon.ChangeFees{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ChangeFees](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixChangeFees), abiCocoon.CocoonChangeFeesMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonChangeOwnerMsgOp] = abiCocoon.ChangeOwner{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ChangeOwner](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixChangeOwner), abiCocoon.CocoonChangeOwnerMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonChangeParamsMsgOp] = abiCocoon.ChangeParams{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ChangeParams](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixChangeParams), abiCocoon.CocoonChangeParamsMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonClientProxyRequestMsgOp] = abiCocoon.ClientProxyRequest{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ClientProxyRequest](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixClientProxyRequest), abiCocoon.CocoonClientProxyRequestMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonDelModelTypeMsgOp] = abiCocoon.DelModelType{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.DelModelType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixDelModelType), abiCocoon.CocoonDelModelTypeMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonDelProxyTypeMsgOp] = abiCocoon.DelProxyType{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.DelProxyType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixDelProxyType), abiCocoon.CocoonDelProxyTypeMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonDelWorkerTypeMsgOp] = abiCocoon.DelWorkerType{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.DelWorkerType](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixDelWorkerType), abiCocoon.CocoonDelWorkerTypeMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonExtClientChargeSignedMsgOp] = abiCocoon.ExtClientChargeSigned{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtClientChargeSigned](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtClientChargeSigned), abiCocoon.CocoonExtClientChargeSignedMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonExtClientGrantRefundSignedMsgOp] = abiCocoon.ExtClientGrantRefundSigned{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtClientGrantRefundSigned](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtClientGrantRefundSigned), abiCocoon.CocoonExtClientGrantRefundSignedMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonExtClientTopUpMsgOp] = abiCocoon.ExtClientTopUp{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtClientTopUp](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtClientTopUp), abiCocoon.CocoonExtClientTopUpMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonExtProxyCloseCompleteRequestSignedMsgOp] = abiCocoon.ExtProxyCloseCompleteRequestSigned{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtProxyCloseCompleteRequestSigned](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtProxyCloseCompleteRequestSigned), abiCocoon.CocoonExtProxyCloseCompleteRequestSignedMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonExtProxyCloseRequestSignedMsgOp] = abiCocoon.ExtProxyCloseRequestSigned{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtProxyCloseRequestSigned](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtProxyCloseRequestSigned), abiCocoon.CocoonExtProxyCloseRequestSignedMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonExtProxyIncreaseStakeMsgOp] = abiCocoon.ExtProxyIncreaseStake{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtProxyIncreaseStake](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtProxyIncreaseStake), abiCocoon.CocoonExtProxyIncreaseStakeMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonExtProxyPayoutRequestMsgOp] = abiCocoon.ExtProxyPayoutRequest{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ExtProxyPayoutRequest](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixExtProxyPayoutRequest), abiCocoon.CocoonExtProxyPayoutRequestMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonOwnerClientChangeSecretHashMsgOp] = abiCocoon.OwnerClientChangeSecretHash{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientChangeSecretHash](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientChangeSecretHash), abiCocoon.CocoonOwnerClientChangeSecretHashMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonOwnerClientChangeSecretHashAndTopUpMsgOp] = abiCocoon.OwnerClientChangeSecretHashAndTopUp{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientChangeSecretHashAndTopUp](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientChangeSecretHashAndTopUp), abiCocoon.CocoonOwnerClientChangeSecretHashAndTopUpMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonOwnerClientIncreaseStakeMsgOp] = abiCocoon.OwnerClientIncreaseStake{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientIncreaseStake](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientIncreaseStake), abiCocoon.CocoonOwnerClientIncreaseStakeMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonOwnerClientRegisterMsgOp] = abiCocoon.OwnerClientRegister{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientRegister](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientRegister), abiCocoon.CocoonOwnerClientRegisterMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonOwnerClientRequestRefundMsgOp] = abiCocoon.OwnerClientRequestRefund{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientRequestRefund](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientRequestRefund), abiCocoon.CocoonOwnerClientRequestRefundMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonOwnerClientWithdrawMsgOp] = abiCocoon.OwnerClientWithdraw{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerClientWithdraw](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerClientWithdraw), abiCocoon.CocoonOwnerClientWithdrawMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonOwnerProxyCloseMsgOp] = abiCocoon.OwnerProxyClose{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerProxyClose](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerProxyClose), abiCocoon.CocoonOwnerProxyCloseMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonOwnerWalletSendMessageMsgOp] = abiCocoon.OwnerWalletSendMessage{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.OwnerWalletSendMessage](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixOwnerWalletSendMessage), abiCocoon.CocoonOwnerWalletSendMessageMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonPayoutMsgOp] = abiCocoon.Payout{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.Payout](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixPayout), abiCocoon.CocoonPayoutMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonRegisterProxyMsgOp] = abiCocoon.RegisterProxy{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.RegisterProxy](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixRegisterProxy), abiCocoon.CocoonRegisterProxyMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonResetRootMsgOp] = abiCocoon.ResetRoot{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ResetRoot](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixResetRoot), abiCocoon.CocoonResetRootMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonReturnExcessesBackMsgOp] = abiCocoon.ReturnExcessesBack{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.ReturnExcessesBack](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixReturnExcessesBack), abiCocoon.CocoonReturnExcessesBackMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonTextCmdMsgOp] = abiCocoon.TextCmd{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.TextCmd](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixTextCmd), abiCocoon.CocoonTextCmdMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonTextCommandMsgOp] = abiCocoon.TextCommand{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.TextCommand](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixTextCommand), abiCocoon.CocoonTextCommandMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonUnregisterProxyMsgOp] = abiCocoon.UnregisterProxy{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.UnregisterProxy](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixUnregisterProxy), abiCocoon.CocoonUnregisterProxyMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonUpdateProxyMsgOp] = abiCocoon.UpdateProxy{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.UpdateProxy](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixUpdateProxy), abiCocoon.CocoonUpdateProxyMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonUpgradeCodeMsgOp] = abiCocoon.UpgradeCode{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.UpgradeCode](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixUpgradeCode), abiCocoon.CocoonUpgradeCodeMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonUpgradeContractsMsgOp] = abiCocoon.UpgradeContracts{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.UpgradeContracts](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixUpgradeContracts), abiCocoon.CocoonUpgradeContractsMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonUpgradeFullMsgOp] = abiCocoon.UpgradeFull{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.UpgradeFull](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixUpgradeFull), abiCocoon.CocoonUpgradeFullMsgOp)
	KnownMsgInTypes[abiCocoon.CocoonWorkerProxyRequestMsgOp] = abiCocoon.WorkerProxyRequest{}
	registerInMsgUnmarshalerForOpcode[*abiCocoon.WorkerProxyRequest](opcodedMsgInDecodeFunctions, uint32(abiCocoon.PrefixWorkerProxyRequest), abiCocoon.CocoonWorkerProxyRequestMsgOp)

}

func init() {
	tolkMethods = append(tolkMethods,
		MethodDescription{
			Name: "active_election_id",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiElector.GetActiveElectionId(ctx, executor, id)
				return "GetActiveElectionId_ElectorResult", r, err
			},
		},
		MethodDescription{
			Name: "participant_list_extended",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiElector.GetParticipantListExtended(ctx, executor, id)
				return "GetParticipantListExtended_ElectorResult", r, err
			},
		},
	)

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    Elector,
			Results: []string{"GetActiveElectionId_ElectorResult", "GetParticipantListExtended_ElectorResult"},
		},
	)

	KnownMsgInTypes[abiElector.ElectorComplaintResponseMsgOp] = abiElector.ComplaintResponse{}
	registerInMsgUnmarshalerForOpcode[*abiElector.ComplaintResponse](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixComplaintResponse), abiElector.ElectorComplaintResponseMsgOp)
	KnownMsgInTypes[abiElector.ElectorConfigAcceptedMsgOp] = abiElector.ConfigAccepted{}
	registerInMsgUnmarshalerForOpcode[*abiElector.ConfigAccepted](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixConfigAccepted), abiElector.ElectorConfigAcceptedMsgOp)
	KnownMsgInTypes[abiElector.ElectorConfigRejectedMsgOp] = abiElector.ConfigRejected{}
	registerInMsgUnmarshalerForOpcode[*abiElector.ConfigRejected](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixConfigRejected), abiElector.ElectorConfigRejectedMsgOp)
	KnownMsgInTypes[abiElector.ElectorErrorResponseMsgOp] = abiElector.ErrorResponse{}
	registerInMsgUnmarshalerForOpcode[*abiElector.ErrorResponse](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixErrorResponse), abiElector.ElectorErrorResponseMsgOp)
	KnownMsgInTypes[abiElector.ElectorNewStakeMsgOp] = abiElector.NewStake{}
	registerInMsgUnmarshalerForOpcode[*abiElector.NewStake](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixNewStake), abiElector.ElectorNewStakeMsgOp)
	KnownMsgInTypes[abiElector.ElectorNewStakeConfirmationMsgOp] = abiElector.NewStakeConfirmation{}
	registerInMsgUnmarshalerForOpcode[*abiElector.NewStakeConfirmation](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixNewStakeConfirmation), abiElector.ElectorNewStakeConfirmationMsgOp)
	KnownMsgInTypes[abiElector.ElectorRecoverStakeRequestMsgOp] = abiElector.RecoverStakeRequest{}
	registerInMsgUnmarshalerForOpcode[*abiElector.RecoverStakeRequest](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixRecoverStakeRequest), abiElector.ElectorRecoverStakeRequestMsgOp)
	KnownMsgInTypes[abiElector.ElectorRecoverStakeResponseMsgOp] = abiElector.RecoverStakeResponse{}
	registerInMsgUnmarshalerForOpcode[*abiElector.RecoverStakeResponse](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixRecoverStakeResponse), abiElector.ElectorRecoverStakeResponseMsgOp)
	KnownMsgInTypes[abiElector.ElectorRegisterComplaintMsgOp] = abiElector.RegisterComplaint{}
	registerInMsgUnmarshalerForOpcode[*abiElector.RegisterComplaint](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixRegisterComplaint), abiElector.ElectorRegisterComplaintMsgOp)
	KnownMsgInTypes[abiElector.ElectorReturnStakeMsgOp] = abiElector.ReturnStake{}
	registerInMsgUnmarshalerForOpcode[*abiElector.ReturnStake](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixReturnStake), abiElector.ElectorReturnStakeMsgOp)
	KnownMsgInTypes[abiElector.ElectorUpgradeCodeMsgOp] = abiElector.UpgradeCode{}
	registerInMsgUnmarshalerForOpcode[*abiElector.UpgradeCode](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixUpgradeCode), abiElector.ElectorUpgradeCodeMsgOp)
	KnownMsgInTypes[abiElector.ElectorUpgradeCodeResponseMsgOp] = abiElector.UpgradeCodeResponse{}
	registerInMsgUnmarshalerForOpcode[*abiElector.UpgradeCodeResponse](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixUpgradeCodeResponse), abiElector.ElectorUpgradeCodeResponseMsgOp)
	KnownMsgInTypes[abiElector.ElectorVoteComplaintMsgOp] = abiElector.VoteComplaint{}
	registerInMsgUnmarshalerForOpcode[*abiElector.VoteComplaint](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixVoteComplaint), abiElector.ElectorVoteComplaintMsgOp)
	KnownMsgInTypes[abiElector.ElectorVoteComplaintResponseMsgOp] = abiElector.VoteComplaintResponse{}
	registerInMsgUnmarshalerForOpcode[*abiElector.VoteComplaintResponse](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixVoteComplaintResponse), abiElector.ElectorVoteComplaintResponseMsgOp)

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

	KnownMsgInTypes[abiFfVault.FfVaultAssetDepositMsgOp] = abiFfVault.AssetDeposit{}
	registerInMsgUnmarshalerForOpcode[*abiFfVault.AssetDeposit](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixAssetDeposit), abiFfVault.FfVaultAssetDepositMsgOp)
	KnownMsgInTypes[abiFfVault.FfVaultStableDepositMsgOp] = abiFfVault.StableDeposit{}
	registerInMsgUnmarshalerForOpcode[*abiFfVault.StableDeposit](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixStableDeposit), abiFfVault.FfVaultStableDepositMsgOp)
	KnownMsgInTypes[abiFfVault.FfVaultStakeOperationMsgOp] = abiFfVault.StakeOperation{}
	registerInMsgUnmarshalerForOpcode[*abiFfVault.StakeOperation](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixStakeOperation), abiFfVault.FfVaultStakeOperationMsgOp)
	KnownMsgInTypes[abiFfVault.FfVaultUnstakeExecuteMsgOp] = abiFfVault.UnstakeExecute{}
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeExecute](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeExecute), abiFfVault.FfVaultUnstakeExecuteMsgOp)
	KnownMsgInTypes[abiFfVault.FfVaultUnstakeExecuteCancelMsgOp] = abiFfVault.UnstakeExecuteCancel{}
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeExecuteCancel](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeExecuteCancel), abiFfVault.FfVaultUnstakeExecuteCancelMsgOp)
	KnownMsgInTypes[abiFfVault.FfVaultUnstakeExecuteInternalMsgOp] = abiFfVault.UnstakeExecuteInternal{}
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeExecuteInternal](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeExecuteInternal), abiFfVault.FfVaultUnstakeExecuteInternalMsgOp)
	KnownMsgInTypes[abiFfVault.FfVaultUnstakeExecuteInternalCallbackMsgOp] = abiFfVault.UnstakeExecuteInternalCallback{}
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeExecuteInternalCallback](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeExecuteInternalCallback), abiFfVault.FfVaultUnstakeExecuteInternalCallbackMsgOp)
	KnownMsgInTypes[abiFfVault.FfVaultUnstakeOperationMsgOp] = abiFfVault.UnstakeOperation{}
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeOperation](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeOperation), abiFfVault.FfVaultUnstakeOperationMsgOp)
	KnownMsgInTypes[abiFfVault.FfVaultUnstakeRequestMsgOp] = abiFfVault.UnstakeRequest{}
	registerInMsgUnmarshalerForOpcode[*abiFfVault.UnstakeRequest](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixUnstakeRequest), abiFfVault.FfVaultUnstakeRequestMsgOp)
	KnownMsgInTypes[abiFfVault.FfVaultVaultStorageMsgOp] = abiFfVault.VaultStorage{}
	registerInMsgUnmarshalerForOpcode[*abiFfVault.VaultStorage](opcodedMsgInDecodeFunctions, uint32(abiFfVault.PrefixVaultStorage), abiFfVault.FfVaultVaultStorageMsgOp)
	KnownMsgInTypes[abiFfVault.FfVaultWithdrawJettonMsgOp] = abiFfVault.WithdrawJetton{}
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

	KnownMsgInTypes[abiPythOracle.PythOracleErrorResponseMsgOp] = abiPythOracle.ErrorResponse{}
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.ErrorResponse](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixErrorResponse), abiPythOracle.PythOracleErrorResponseMsgOp)
	KnownMsgInTypes[abiPythOracle.PythOracleOracleResponseSuccessMsgOp] = abiPythOracle.OracleResponseSuccess{}
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.OracleResponseSuccess](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixOracleResponseSuccess), abiPythOracle.PythOracleOracleResponseSuccessMsgOp)
	KnownMsgInTypes[abiPythOracle.PythOracleParsePriceFeedUpdatesMessageMsgOp] = abiPythOracle.ParsePriceFeedUpdatesMessage{}
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.ParsePriceFeedUpdatesMessage](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixParsePriceFeedUpdatesMessage), abiPythOracle.PythOracleParsePriceFeedUpdatesMessageMsgOp)
	KnownMsgInTypes[abiPythOracle.PythOracleParseUniquePriceFeedUpdatesMessageMsgOp] = abiPythOracle.ParseUniquePriceFeedUpdatesMessage{}
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.ParseUniquePriceFeedUpdatesMessage](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixParseUniquePriceFeedUpdatesMessage), abiPythOracle.PythOracleParseUniquePriceFeedUpdatesMessageMsgOp)

}

func init() {
	tolkMethods = append(tolkMethods,
		MethodDescription{
			Name: "get_roles",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiSingleNominatorPool.GetRoles(ctx, executor, id)
				return "GetRoles_SingleNominatorPoolResult", r, err
			},
		},
		MethodDescription{
			Name: "get_pool_data",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiSingleNominatorPool.GetPoolData(ctx, executor, id)
				return "GetPoolData_SingleNominatorPoolResult", r, err
			},
		},
	)

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    SingleNominatorPool,
			Results: []string{"GetRoles_SingleNominatorPoolResult", "GetPoolData_SingleNominatorPoolResult"},
		},
	)

	KnownMsgInTypes[abiSingleNominatorPool.SingleNominatorPoolChangeValidatorAddressMsgOp] = abiSingleNominatorPool.ChangeValidatorAddress{}
	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.ChangeValidatorAddress](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixChangeValidatorAddress), abiSingleNominatorPool.SingleNominatorPoolChangeValidatorAddressMsgOp)
	KnownMsgInTypes[abiSingleNominatorPool.SingleNominatorPoolSendRawMsgMsgOp] = abiSingleNominatorPool.SendRawMsg{}
	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.SendRawMsg](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixSendRawMsg), abiSingleNominatorPool.SingleNominatorPoolSendRawMsgMsgOp)
	KnownMsgInTypes[abiSingleNominatorPool.SingleNominatorPoolUpgradeMsgOp] = abiSingleNominatorPool.Upgrade{}
	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.Upgrade](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixUpgrade), abiSingleNominatorPool.SingleNominatorPoolUpgradeMsgOp)
	KnownMsgInTypes[abiSingleNominatorPool.SingleNominatorPoolWithdrawMsgOp] = abiSingleNominatorPool.Withdraw{}
	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.Withdraw](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixWithdraw), abiSingleNominatorPool.SingleNominatorPoolWithdrawMsgOp)

}
