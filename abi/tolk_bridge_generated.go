// Code generated - DO NOT EDIT.

package abi

import (
	"context"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	abiElector "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/elector"
	abiFfVault "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/ffVault"
	abiPythOracle "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/pythOracle"
	abiSingleNominatorPool "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/singleNominatorPool"
	"github.com/tonkeeper/tongo/tlb"
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

	KnownGetMethodsDecoder["get_cocoon_client_data"] = append(KnownGetMethodsDecoder["get_cocoon_client_data"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetCocoonClientData(&st)
		return "GetCocoonClientData_CocoonClientResult", r, err
	})
	KnownGetMethodsDecoder["get_cocoon_proxy_data"] = append(KnownGetMethodsDecoder["get_cocoon_proxy_data"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetCocoonProxyData(&st)
		return "GetCocoonProxyData_CocoonProxyResult", r, err
	})
	KnownGetMethodsDecoder["last_proxy_seqno"] = append(KnownGetMethodsDecoder["last_proxy_seqno"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetLastProxySeqno(&st)
		return "GetLastProxySeqno_CocoonRootResult", r, err
	})
	KnownGetMethodsDecoder["get_cocoon_data"] = append(KnownGetMethodsDecoder["get_cocoon_data"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetCocoonData(&st)
		return "GetCocoonData_CocoonRootResult", r, err
	})
	KnownGetMethodsDecoder["get_cur_params"] = append(KnownGetMethodsDecoder["get_cur_params"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetCurParams(&st)
		return "GetCurParams_CocoonRootResult", r, err
	})
	KnownGetMethodsDecoder["proxy_hash_is_valid"] = append(KnownGetMethodsDecoder["proxy_hash_is_valid"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetProxyHashIsValid(&st)
		return "GetProxyHashIsValid_CocoonRootResult", r, err
	})
	KnownGetMethodsDecoder["worker_hash_is_valid"] = append(KnownGetMethodsDecoder["worker_hash_is_valid"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetWorkerHashIsValid(&st)
		return "GetWorkerHashIsValid_CocoonRootResult", r, err
	})
	KnownGetMethodsDecoder["model_hash_is_valid"] = append(KnownGetMethodsDecoder["model_hash_is_valid"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetModelHashIsValid(&st)
		return "GetModelHashIsValid_CocoonRootResult", r, err
	})
	KnownGetMethodsDecoder["seqno"] = append(KnownGetMethodsDecoder["seqno"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetSeqno(&st)
		return "GetSeqno_CocoonWalletResult", r, err
	})
	KnownGetMethodsDecoder["get_public_key"] = append(KnownGetMethodsDecoder["get_public_key"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetPublicKey(&st)
		return "GetPublicKey_CocoonWalletResult", r, err
	})
	KnownGetMethodsDecoder["get_owner_address"] = append(KnownGetMethodsDecoder["get_owner_address"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetOwnerAddress(&st)
		return "GetOwnerAddress_CocoonWalletResult", r, err
	})
	KnownGetMethodsDecoder["get_cocoon_worker_data"] = append(KnownGetMethodsDecoder["get_cocoon_worker_data"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiCocoon.DecodeGetCocoonWorkerData(&st)
		return "GetCocoonWorkerData_CocoonWorkerResult", r, err
	})

	KnownSimpleGetMethods[75156] = append(KnownSimpleGetMethods[75156], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiCocoon.GetCocoonClientData(ctx, executor, id)
		return "GetCocoonClientData_CocoonClientResult", r, err
	})
	KnownSimpleGetMethods[97687] = append(KnownSimpleGetMethods[97687], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiCocoon.GetCocoonProxyData(ctx, executor, id)
		return "GetCocoonProxyData_CocoonProxyResult", r, err
	})
	KnownSimpleGetMethods[65647] = append(KnownSimpleGetMethods[65647], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiCocoon.GetLastProxySeqno(ctx, executor, id)
		return "GetLastProxySeqno_CocoonRootResult", r, err
	})
	KnownSimpleGetMethods[96613] = append(KnownSimpleGetMethods[96613], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiCocoon.GetCocoonData(ctx, executor, id)
		return "GetCocoonData_CocoonRootResult", r, err
	})
	KnownSimpleGetMethods[89457] = append(KnownSimpleGetMethods[89457], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiCocoon.GetCurParams(ctx, executor, id)
		return "GetCurParams_CocoonRootResult", r, err
	})
	KnownSimpleGetMethods[85143] = append(KnownSimpleGetMethods[85143], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiCocoon.GetSeqno(ctx, executor, id)
		return "GetSeqno_CocoonWalletResult", r, err
	})
	KnownSimpleGetMethods[78748] = append(KnownSimpleGetMethods[78748], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiCocoon.GetPublicKey(ctx, executor, id)
		return "GetPublicKey_CocoonWalletResult", r, err
	})
	KnownSimpleGetMethods[114619] = append(KnownSimpleGetMethods[114619], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiCocoon.GetOwnerAddress(ctx, executor, id)
		return "GetOwnerAddress_CocoonWalletResult", r, err
	})
	KnownSimpleGetMethods[106427] = append(KnownSimpleGetMethods[106427], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiCocoon.GetCocoonWorkerData(ctx, executor, id)
		return "GetCocoonWorkerData_CocoonWorkerResult", r, err
	})

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

	KnownGetMethodsDecoder["active_election_id"] = append(KnownGetMethodsDecoder["active_election_id"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiElector.DecodeGetActiveElectionId(&st)
		return "GetActiveElectionId_ElectorResult", r, err
	})
	KnownGetMethodsDecoder["participates_in"] = append(KnownGetMethodsDecoder["participates_in"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiElector.DecodeGetParticipatesIn(&st)
		return "GetParticipatesIn_ElectorResult", r, err
	})
	KnownGetMethodsDecoder["compute_returned_stake"] = append(KnownGetMethodsDecoder["compute_returned_stake"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiElector.DecodeGetComputeReturnedStake(&st)
		return "GetComputeReturnedStake_ElectorResult", r, err
	})
	KnownGetMethodsDecoder["participant_list_extended"] = append(KnownGetMethodsDecoder["participant_list_extended"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiElector.DecodeGetParticipantListExtended(&st)
		return "GetParticipantListExtended_ElectorResult", r, err
	})

	KnownSimpleGetMethods[86535] = append(KnownSimpleGetMethods[86535], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiElector.GetActiveElectionId(ctx, executor, id)
		return "GetActiveElectionId_ElectorResult", r, err
	})
	KnownSimpleGetMethods[86698] = append(KnownSimpleGetMethods[86698], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiElector.GetParticipantListExtended(ctx, executor, id)
		return "GetParticipantListExtended_ElectorResult", r, err
	})

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    Elector,
			Results: []string{"GetActiveElectionId_ElectorResult", "GetParticipantListExtended_ElectorResult"},
		},
	)

	registerInMsgUnmarshalerForOpcode[*abiElector.ComplaintResponse](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixComplaintResponse), abiElector.ElectorComplaintResponseMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.ConfigAccepted](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixConfigAccepted), abiElector.ElectorConfigAcceptedMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.ConfigRejected](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixConfigRejected), abiElector.ElectorConfigRejectedMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.ErrorResponse](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixErrorResponse), abiElector.ElectorErrorResponseMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.NewStake](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixNewStake), abiElector.ElectorNewStakeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.NewStakeConfirmation](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixNewStakeConfirmation), abiElector.ElectorNewStakeConfirmationMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.RecoverStakeRequest](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixRecoverStakeRequest), abiElector.ElectorRecoverStakeRequestMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.RecoverStakeResponse](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixRecoverStakeResponse), abiElector.ElectorRecoverStakeResponseMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.RegisterComplaint](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixRegisterComplaint), abiElector.ElectorRegisterComplaintMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.ReturnStake](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixReturnStake), abiElector.ElectorReturnStakeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.UpgradeCode](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixUpgradeCode), abiElector.ElectorUpgradeCodeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.UpgradeCodeResponse](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixUpgradeCodeResponse), abiElector.ElectorUpgradeCodeResponseMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiElector.VoteComplaint](opcodedMsgInDecodeFunctions, uint32(abiElector.PrefixVoteComplaint), abiElector.ElectorVoteComplaintMsgOp)
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

	KnownGetMethodsDecoder["get_nft_data"] = append(KnownGetMethodsDecoder["get_nft_data"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiFfVault.DecodeGetNftData(&st)
		return "GetNftData_FfVaultPositionResult", r, err
	})
	KnownGetMethodsDecoder["get_stake_position_info"] = append(KnownGetMethodsDecoder["get_stake_position_info"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiFfVault.DecodeGetStakePositionInfo(&st)
		return "GetStakePositionInfo_FfVaultPositionResult", r, err
	})
	KnownGetMethodsDecoder["get_collection_data"] = append(KnownGetMethodsDecoder["get_collection_data"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiFfVault.DecodeGetCollectionData(&st)
		return "GetCollectionData_FfVaultResult", r, err
	})
	KnownGetMethodsDecoder["get_nft_address_by_index"] = append(KnownGetMethodsDecoder["get_nft_address_by_index"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiFfVault.DecodeGetNftAddressByIndex(&st)
		return "GetNftAddressByIndex_FfVaultResult", r, err
	})
	KnownGetMethodsDecoder["get_staking_data"] = append(KnownGetMethodsDecoder["get_staking_data"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiFfVault.DecodeGetStakingData(&st)
		return "GetStakingData_FfVaultResult", r, err
	})
	KnownGetMethodsDecoder["get_balance"] = append(KnownGetMethodsDecoder["get_balance"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiFfVault.DecodeGetBalance(&st)
		return "GetBalance_FfVaultResult", r, err
	})

	KnownSimpleGetMethods[102351] = append(KnownSimpleGetMethods[102351], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiFfVault.GetNftData(ctx, executor, id)
		return "GetNftData_FfVaultPositionResult", r, err
	})
	KnownSimpleGetMethods[102640] = append(KnownSimpleGetMethods[102640], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiFfVault.GetStakePositionInfo(ctx, executor, id)
		return "GetStakePositionInfo_FfVaultPositionResult", r, err
	})
	KnownSimpleGetMethods[102491] = append(KnownSimpleGetMethods[102491], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiFfVault.GetCollectionData(ctx, executor, id)
		return "GetCollectionData_FfVaultResult", r, err
	})
	KnownSimpleGetMethods[108033] = append(KnownSimpleGetMethods[108033], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiFfVault.GetStakingData(ctx, executor, id)
		return "GetStakingData_FfVaultResult", r, err
	})
	KnownSimpleGetMethods[130343] = append(KnownSimpleGetMethods[130343], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiFfVault.GetBalance(ctx, executor, id)
		return "GetBalance_FfVaultResult", r, err
	})

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

	KnownGetMethodsDecoder["get_update_fee"] = append(KnownGetMethodsDecoder["get_update_fee"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetUpdateFee(&st)
		return "GetUpdateFee_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_single_update_fee"] = append(KnownGetMethodsDecoder["get_single_update_fee"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetSingleUpdateFee(&st)
		return "GetSingleUpdateFee_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_governance_data_source_index"] = append(KnownGetMethodsDecoder["get_governance_data_source_index"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetGovernanceDataSourceIndex(&st)
		return "GetGovernanceDataSourceIndex_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_governance_data_source"] = append(KnownGetMethodsDecoder["get_governance_data_source"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetGovernanceDataSource(&st)
		return "GetGovernanceDataSource_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_last_executed_governance_sequence"] = append(KnownGetMethodsDecoder["get_last_executed_governance_sequence"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetLastExecutedGovernanceSequence(&st)
		return "GetLastExecutedGovernanceSequence_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_is_valid_data_source"] = append(KnownGetMethodsDecoder["get_is_valid_data_source"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetIsValidDataSource(&st)
		return "GetIsValidDataSource_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_price_unsafe"] = append(KnownGetMethodsDecoder["get_price_unsafe"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetPriceUnsafe(&st)
		return "GetPriceUnsafe_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_price_no_older_than"] = append(KnownGetMethodsDecoder["get_price_no_older_than"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetPriceNoOlderThan(&st)
		return "GetPriceNoOlderThan_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_ema_price_unsafe"] = append(KnownGetMethodsDecoder["get_ema_price_unsafe"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetEmaPriceUnsafe(&st)
		return "GetEmaPriceUnsafe_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_ema_price_no_older_than"] = append(KnownGetMethodsDecoder["get_ema_price_no_older_than"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetEmaPriceNoOlderThan(&st)
		return "GetEmaPriceNoOlderThan_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_chain_id"] = append(KnownGetMethodsDecoder["get_chain_id"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetChainId(&st)
		return "GetChainId_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_current_guardian_set_index"] = append(KnownGetMethodsDecoder["get_current_guardian_set_index"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetCurrentGuardianSetIndex(&st)
		return "GetCurrentGuardianSetIndex_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_guardian_set"] = append(KnownGetMethodsDecoder["get_guardian_set"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetGuardianSet(&st)
		return "GetGuardianSet_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_governance_chain_id"] = append(KnownGetMethodsDecoder["get_governance_chain_id"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetGovernanceChainId(&st)
		return "GetGovernanceChainId_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["get_governance_contract"] = append(KnownGetMethodsDecoder["get_governance_contract"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetGovernanceContract(&st)
		return "GetGovernanceContract_PythOracleResult", r, err
	})
	KnownGetMethodsDecoder["governance_action_is_consumed"] = append(KnownGetMethodsDecoder["governance_action_is_consumed"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiPythOracle.DecodeGetGovernanceActionIsConsumed(&st)
		return "GetGovernanceActionIsConsumed_PythOracleResult", r, err
	})

	KnownSimpleGetMethods[99955] = append(KnownSimpleGetMethods[99955], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiPythOracle.GetSingleUpdateFee(ctx, executor, id)
		return "GetSingleUpdateFee_PythOracleResult", r, err
	})
	KnownSimpleGetMethods[98238] = append(KnownSimpleGetMethods[98238], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiPythOracle.GetGovernanceDataSourceIndex(ctx, executor, id)
		return "GetGovernanceDataSourceIndex_PythOracleResult", r, err
	})
	KnownSimpleGetMethods[110935] = append(KnownSimpleGetMethods[110935], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiPythOracle.GetGovernanceDataSource(ctx, executor, id)
		return "GetGovernanceDataSource_PythOracleResult", r, err
	})
	KnownSimpleGetMethods[70196] = append(KnownSimpleGetMethods[70196], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiPythOracle.GetLastExecutedGovernanceSequence(ctx, executor, id)
		return "GetLastExecutedGovernanceSequence_PythOracleResult", r, err
	})
	KnownSimpleGetMethods[122952] = append(KnownSimpleGetMethods[122952], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiPythOracle.GetChainId(ctx, executor, id)
		return "GetChainId_PythOracleResult", r, err
	})
	KnownSimpleGetMethods[114628] = append(KnownSimpleGetMethods[114628], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiPythOracle.GetCurrentGuardianSetIndex(ctx, executor, id)
		return "GetCurrentGuardianSetIndex_PythOracleResult", r, err
	})
	KnownSimpleGetMethods[102302] = append(KnownSimpleGetMethods[102302], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiPythOracle.GetGovernanceChainId(ctx, executor, id)
		return "GetGovernanceChainId_PythOracleResult", r, err
	})
	KnownSimpleGetMethods[65842] = append(KnownSimpleGetMethods[65842], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiPythOracle.GetGovernanceContract(ctx, executor, id)
		return "GetGovernanceContract_PythOracleResult", r, err
	})

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    PythOracle,
			Results: []string{"GetSingleUpdateFee_PythOracleResult", "GetGovernanceDataSourceIndex_PythOracleResult", "GetGovernanceDataSource_PythOracleResult", "GetLastExecutedGovernanceSequence_PythOracleResult", "GetChainId_PythOracleResult", "GetCurrentGuardianSetIndex_PythOracleResult", "GetGovernanceChainId_PythOracleResult", "GetGovernanceContract_PythOracleResult"},
		},
	)

	registerInMsgUnmarshalerForOpcode[*abiPythOracle.ErrorResponse](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixErrorResponse), abiPythOracle.PythOracleErrorResponseMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.OracleResponseSuccess](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixOracleResponseSuccess), abiPythOracle.PythOracleOracleResponseSuccessMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiPythOracle.ParsePriceFeedUpdatesMessage](opcodedMsgInDecodeFunctions, uint32(abiPythOracle.PrefixParsePriceFeedUpdatesMessage), abiPythOracle.PythOracleParsePriceFeedUpdatesMessageMsgOp)
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

	KnownGetMethodsDecoder["get_roles"] = append(KnownGetMethodsDecoder["get_roles"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiSingleNominatorPool.DecodeGetRoles(&st)
		return "GetRoles_SingleNominatorPoolResult", r, err
	})
	KnownGetMethodsDecoder["get_pool_data"] = append(KnownGetMethodsDecoder["get_pool_data"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiSingleNominatorPool.DecodeGetPoolData(&st)
		return "GetPoolData_SingleNominatorPoolResult", r, err
	})

	KnownSimpleGetMethods[130208] = append(KnownSimpleGetMethods[130208], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiSingleNominatorPool.GetRoles(ctx, executor, id)
		return "GetRoles_SingleNominatorPoolResult", r, err
	})
	KnownSimpleGetMethods[81689] = append(KnownSimpleGetMethods[81689], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiSingleNominatorPool.GetPoolData(ctx, executor, id)
		return "GetPoolData_SingleNominatorPoolResult", r, err
	})

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    SingleNominatorPool,
			Results: []string{"GetRoles_SingleNominatorPoolResult", "GetPoolData_SingleNominatorPoolResult"},
		},
	)

	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.ChangeValidatorAddress](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixChangeValidatorAddress), abiSingleNominatorPool.SingleNominatorPoolChangeValidatorAddressMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.SendRawMsg](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixSendRawMsg), abiSingleNominatorPool.SingleNominatorPoolSendRawMsgMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.Upgrade](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixUpgrade), abiSingleNominatorPool.SingleNominatorPoolUpgradeMsgOp)
	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.Withdraw](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixWithdraw), abiSingleNominatorPool.SingleNominatorPoolWithdrawMsgOp)

}
