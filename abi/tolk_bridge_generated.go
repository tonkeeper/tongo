// Code generated - DO NOT EDIT.

package abi

import (
	"context"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	abiElector "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/elector"
	abiEvaa "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/evaa"
	abiFfVault "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/ffVault"
	abiPythOracle "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/pythOracle"
	abiSingleNominatorPool "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/singleNominatorPool"
	abiStonfi "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/stonfi"
	abiStonkspump "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/stonkspump"
	abiXtr "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/xtr"
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
			Name: "get_active",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiEvaa.GetActive(ctx, executor, id)
				return "GetActive_EvaaMasterResult", r, err
			},
		},
		MethodDescription{
			Name: "getTokensKeys",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiEvaa.GetTokensKeys(ctx, executor, id)
				return "GetTokensKeys_EvaaMasterResult", r, err
			},
		},
		MethodDescription{
			Name: "getLastUserScVersion",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiEvaa.GetLastUserScVersion(ctx, executor, id)
				return "GetLastUserScVersion_EvaaMasterResult", r, err
			},
		},
		MethodDescription{
			Name: "getUpgradeConfig",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiEvaa.GetUpgradeConfig(ctx, executor, id)
				return "GetUpgradeConfig_EvaaMasterResult", r, err
			},
		},
		MethodDescription{
			Name: "get_supervisor",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiEvaa.GetSupervisor(ctx, executor, id)
				return "GetSupervisor_EvaaMasterResult", r, err
			},
		},
		MethodDescription{
			Name: "getStore",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiEvaa.GetStore(ctx, executor, id)
				return "GetStore_EvaaMasterResult", r, err
			},
		},
		MethodDescription{
			Name: "codeVersion",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiEvaa.GetCodeVersion(ctx, executor, id)
				return "GetCodeVersion_EvaaUserResult", r, err
			},
		},
		MethodDescription{
			Name: "isUserSc",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiEvaa.GetIsUserSc(ctx, executor, id)
				return "GetIsUserSc_EvaaUserResult", r, err
			},
		},
		MethodDescription{
			Name: "getPrincipals",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiEvaa.GetPrincipals(ctx, executor, id)
				return "GetPrincipals_EvaaUserResult", r, err
			},
		},
		MethodDescription{
			Name: "getRewards",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiEvaa.GetRewards(ctx, executor, id)
				return "GetRewards_EvaaUserResult", r, err
			},
		},
		MethodDescription{
			Name: "getAllUserScData",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiEvaa.GetAllUserScData(ctx, executor, id)
				return "GetAllUserScData_EvaaUserResult", r, err
			},
		},
	)

	KnownGetMethodsDecoder["get_asset_sb_rate"] = append(KnownGetMethodsDecoder["get_asset_sb_rate"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAssetSbRate(&st)
		return "GetAssetSbRate_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["getAssetRates"] = append(KnownGetMethodsDecoder["getAssetRates"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAssetRates(&st)
		return "GetAssetRates_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["getAssetReserves"] = append(KnownGetMethodsDecoder["getAssetReserves"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAssetReserves(&st)
		return "GetAssetReserves_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["getAssetTotals"] = append(KnownGetMethodsDecoder["getAssetTotals"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAssetTotals(&st)
		return "GetAssetTotals_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["getUpdatedRates"] = append(KnownGetMethodsDecoder["getUpdatedRates"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetUpdatedRates(&st)
		return "GetUpdatedRates_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["getUpdatedRatesForAllAssets"] = append(KnownGetMethodsDecoder["getUpdatedRatesForAllAssets"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetUpdatedRatesForAllAssets(&st)
		return "GetUpdatedRatesForAllAssets_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["getCollateralQuote"] = append(KnownGetMethodsDecoder["getCollateralQuote"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetCollateralQuote(&st)
		return "GetCollateralQuote_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["get_user_address"] = append(KnownGetMethodsDecoder["get_user_address"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetUserAddress(&st)
		return "GetUserAddress_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["get_user_subaccount_address"] = append(KnownGetMethodsDecoder["get_user_subaccount_address"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetUserSubaccountAddress(&st)
		return "GetUserSubaccountAddress_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["get_active"] = append(KnownGetMethodsDecoder["get_active"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetActive(&st)
		return "GetActive_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["getTokensKeys"] = append(KnownGetMethodsDecoder["getTokensKeys"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetTokensKeys(&st)
		return "GetTokensKeys_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["getLastUserScVersion"] = append(KnownGetMethodsDecoder["getLastUserScVersion"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetLastUserScVersion(&st)
		return "GetLastUserScVersion_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["getUpgradeConfig"] = append(KnownGetMethodsDecoder["getUpgradeConfig"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetUpgradeConfig(&st)
		return "GetUpgradeConfig_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["get_asset_tracking_info"] = append(KnownGetMethodsDecoder["get_asset_tracking_info"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAssetTrackingInfo(&st)
		return "GetAssetTrackingInfo_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["get_supervisor"] = append(KnownGetMethodsDecoder["get_supervisor"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetSupervisor(&st)
		return "GetSupervisor_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["get_asset_total_principals"] = append(KnownGetMethodsDecoder["get_asset_total_principals"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAssetTotalPrincipals(&st)
		return "GetAssetTotalPrincipals_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["get_asset_balance"] = append(KnownGetMethodsDecoder["get_asset_balance"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAssetBalance(&st)
		return "GetAssetBalance_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["get_asset_liquidity_by_id"] = append(KnownGetMethodsDecoder["get_asset_liquidity_by_id"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAssetLiquidityById(&st)
		return "GetAssetLiquidityById_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["get_asset_liquidity_minus_reserves_by_id"] = append(KnownGetMethodsDecoder["get_asset_liquidity_minus_reserves_by_id"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAssetLiquidityMinusReservesById(&st)
		return "GetAssetLiquidityMinusReservesById_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["getStore"] = append(KnownGetMethodsDecoder["getStore"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetStore(&st)
		return "GetStore_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["claim_asset_reserves_min_attachment"] = append(KnownGetMethodsDecoder["claim_asset_reserves_min_attachment"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetClaimAssetReservesMinAttachment(&st)
		return "GetClaimAssetReservesMinAttachment_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["supply_min_attachment"] = append(KnownGetMethodsDecoder["supply_min_attachment"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetSupplyMinAttachment(&st)
		return "GetSupplyMinAttachment_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["withdraw_min_attachment"] = append(KnownGetMethodsDecoder["withdraw_min_attachment"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetWithdrawMinAttachment(&st)
		return "GetWithdrawMinAttachment_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["liquidate_min_attachment"] = append(KnownGetMethodsDecoder["liquidate_min_attachment"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetLiquidateMinAttachment(&st)
		return "GetLiquidateMinAttachment_EvaaMasterResult", r, err
	})
	KnownGetMethodsDecoder["codeVersion"] = append(KnownGetMethodsDecoder["codeVersion"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetCodeVersion(&st)
		return "GetCodeVersion_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["isUserSc"] = append(KnownGetMethodsDecoder["isUserSc"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetIsUserSc(&st)
		return "GetIsUserSc_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["getAccountAssetBalance"] = append(KnownGetMethodsDecoder["getAccountAssetBalance"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAccountAssetBalance(&st)
		return "GetAccountAssetBalance_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["getAccountBalances"] = append(KnownGetMethodsDecoder["getAccountBalances"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAccountBalances(&st)
		return "GetAccountBalances_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["getAccountHealth"] = append(KnownGetMethodsDecoder["getAccountHealth"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAccountHealth(&st)
		return "GetAccountHealth_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["getAvailableToBorrow"] = append(KnownGetMethodsDecoder["getAvailableToBorrow"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAvailableToBorrow(&st)
		return "GetAvailableToBorrow_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["getIsLiquidable"] = append(KnownGetMethodsDecoder["getIsLiquidable"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetIsLiquidable(&st)
		return "GetIsLiquidable_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["getAggregatedBalances"] = append(KnownGetMethodsDecoder["getAggregatedBalances"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAggregatedBalances(&st)
		return "GetAggregatedBalances_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["get_asset_principal"] = append(KnownGetMethodsDecoder["get_asset_principal"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAssetPrincipal(&st)
		return "GetAssetPrincipal_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["getPrincipals"] = append(KnownGetMethodsDecoder["getPrincipals"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetPrincipals(&st)
		return "GetPrincipals_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["getRewards"] = append(KnownGetMethodsDecoder["getRewards"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetRewards(&st)
		return "GetRewards_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["getAllUserScData"] = append(KnownGetMethodsDecoder["getAllUserScData"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetAllUserScData(&st)
		return "GetAllUserScData_EvaaUserResult", r, err
	})
	KnownGetMethodsDecoder["get_maximum_withdraw_amount"] = append(KnownGetMethodsDecoder["get_maximum_withdraw_amount"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiEvaa.DecodeGetMaximumWithdrawAmount(&st)
		return "GetMaximumWithdrawAmount_EvaaUserResult", r, err
	})

	KnownSimpleGetMethods[129327] = append(KnownSimpleGetMethods[129327], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiEvaa.GetActive(ctx, executor, id)
		return "GetActive_EvaaMasterResult", r, err
	})
	KnownSimpleGetMethods[98436] = append(KnownSimpleGetMethods[98436], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiEvaa.GetTokensKeys(ctx, executor, id)
		return "GetTokensKeys_EvaaMasterResult", r, err
	})
	KnownSimpleGetMethods[88592] = append(KnownSimpleGetMethods[88592], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiEvaa.GetLastUserScVersion(ctx, executor, id)
		return "GetLastUserScVersion_EvaaMasterResult", r, err
	})
	KnownSimpleGetMethods[73690] = append(KnownSimpleGetMethods[73690], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiEvaa.GetUpgradeConfig(ctx, executor, id)
		return "GetUpgradeConfig_EvaaMasterResult", r, err
	})
	KnownSimpleGetMethods[88756] = append(KnownSimpleGetMethods[88756], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiEvaa.GetSupervisor(ctx, executor, id)
		return "GetSupervisor_EvaaMasterResult", r, err
	})
	KnownSimpleGetMethods[87334] = append(KnownSimpleGetMethods[87334], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiEvaa.GetStore(ctx, executor, id)
		return "GetStore_EvaaMasterResult", r, err
	})
	KnownSimpleGetMethods[93886] = append(KnownSimpleGetMethods[93886], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiEvaa.GetCodeVersion(ctx, executor, id)
		return "GetCodeVersion_EvaaUserResult", r, err
	})
	KnownSimpleGetMethods[72773] = append(KnownSimpleGetMethods[72773], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiEvaa.GetIsUserSc(ctx, executor, id)
		return "GetIsUserSc_EvaaUserResult", r, err
	})
	KnownSimpleGetMethods[129778] = append(KnownSimpleGetMethods[129778], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiEvaa.GetPrincipals(ctx, executor, id)
		return "GetPrincipals_EvaaUserResult", r, err
	})
	KnownSimpleGetMethods[105294] = append(KnownSimpleGetMethods[105294], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiEvaa.GetRewards(ctx, executor, id)
		return "GetRewards_EvaaUserResult", r, err
	})
	KnownSimpleGetMethods[94250] = append(KnownSimpleGetMethods[94250], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiEvaa.GetAllUserScData(ctx, executor, id)
		return "GetAllUserScData_EvaaUserResult", r, err
	})

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    EvaaMaster,
			Results: []string{"GetActive_EvaaMasterResult", "GetTokensKeys_EvaaMasterResult", "GetLastUserScVersion_EvaaMasterResult", "GetUpgradeConfig_EvaaMasterResult", "GetSupervisor_EvaaMasterResult", "GetStore_EvaaMasterResult"},
		},
		InterfaceDescription{
			Name:    EvaaUser,
			Results: []string{"GetCodeVersion_EvaaUserResult", "GetIsUserSc_EvaaUserResult", "GetPrincipals_EvaaUserResult", "GetRewards_EvaaUserResult", "GetAllUserScData_EvaaUserResult"},
		},
	)

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

	KnownMsgInTypes[abiSingleNominatorPool.SingleNominatorPoolChangeValidatorAddressMsgOp] = abiSingleNominatorPool.ChangeValidatorAddress{}
	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.ChangeValidatorAddress](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixChangeValidatorAddress), abiSingleNominatorPool.SingleNominatorPoolChangeValidatorAddressMsgOp)
	KnownMsgInTypes[abiSingleNominatorPool.SingleNominatorPoolSendRawMsgMsgOp] = abiSingleNominatorPool.SendRawMsg{}
	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.SendRawMsg](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixSendRawMsg), abiSingleNominatorPool.SingleNominatorPoolSendRawMsgMsgOp)
	KnownMsgInTypes[abiSingleNominatorPool.SingleNominatorPoolUpgradeMsgOp] = abiSingleNominatorPool.Upgrade{}
	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.Upgrade](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixUpgrade), abiSingleNominatorPool.SingleNominatorPoolUpgradeMsgOp)
	KnownMsgInTypes[abiSingleNominatorPool.SingleNominatorPoolWithdrawMsgOp] = abiSingleNominatorPool.Withdraw{}
	registerInMsgUnmarshalerForOpcode[*abiSingleNominatorPool.Withdraw](opcodedMsgInDecodeFunctions, uint32(abiSingleNominatorPool.PrefixWithdraw), abiSingleNominatorPool.SingleNominatorPoolWithdrawMsgOp)

}

func init() {
	tolkMethods = append(tolkMethods,
		MethodDescription{
			Name: "getEscrowData",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiStonfi.GetEscrowData(ctx, executor, id)
				return "GetEscrowData_StonfiEscrowFactoryResult", r, err
			},
		},
		MethodDescription{
			Name: "getVersion",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiStonfi.GetVersion(ctx, executor, id)
				return "GetVersion_StonfiEscrowFactoryResult", r, err
			},
		},
		MethodDescription{
			Name: "getOrderData",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiStonfi.GetOrderData(ctx, executor, id)
				return "GetOrderData_StonfiEscrowPositionResult", r, err
			},
		},
		MethodDescription{
			Name: "get_cron_info",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiStonfi.GetCronInfo(ctx, executor, id)
				return "GetCronInfo_StonfiEscrowPositionResult", r, err
			},
		},
		MethodDescription{
			Name: "getVaultData",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiStonfi.GetVaultData(ctx, executor, id)
				return "GetVaultData_StonfiEscrowVaultResult", r, err
			},
		},
	)

	KnownGetMethodsDecoder["getEscrowData"] = append(KnownGetMethodsDecoder["getEscrowData"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiStonfi.DecodeGetEscrowData(&st)
		return "GetEscrowData_StonfiEscrowFactoryResult", r, err
	})
	KnownGetMethodsDecoder["getItemAddress"] = append(KnownGetMethodsDecoder["getItemAddress"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiStonfi.DecodeGetItemAddress(&st)
		return "GetItemAddress_StonfiEscrowFactoryResult", r, err
	})
	KnownGetMethodsDecoder["getVaultAddress"] = append(KnownGetMethodsDecoder["getVaultAddress"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiStonfi.DecodeGetVaultAddress(&st)
		return "GetVaultAddress_StonfiEscrowFactoryResult", r, err
	})
	KnownGetMethodsDecoder["getVersion"] = append(KnownGetMethodsDecoder["getVersion"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiStonfi.DecodeGetVersion(&st)
		return "GetVersion_StonfiEscrowFactoryResult", r, err
	})
	KnownGetMethodsDecoder["getOrderData"] = append(KnownGetMethodsDecoder["getOrderData"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiStonfi.DecodeGetOrderData(&st)
		return "GetOrderData_StonfiEscrowPositionResult", r, err
	})
	KnownGetMethodsDecoder["get_cron_info"] = append(KnownGetMethodsDecoder["get_cron_info"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiStonfi.DecodeGetCronInfo(&st)
		return "GetCronInfo_StonfiEscrowPositionResult", r, err
	})
	KnownGetMethodsDecoder["getVaultData"] = append(KnownGetMethodsDecoder["getVaultData"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiStonfi.DecodeGetVaultData(&st)
		return "GetVaultData_StonfiEscrowVaultResult", r, err
	})

	KnownSimpleGetMethods[77307] = append(KnownSimpleGetMethods[77307], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiStonfi.GetEscrowData(ctx, executor, id)
		return "GetEscrowData_StonfiEscrowFactoryResult", r, err
	})
	KnownSimpleGetMethods[74724] = append(KnownSimpleGetMethods[74724], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiStonfi.GetVersion(ctx, executor, id)
		return "GetVersion_StonfiEscrowFactoryResult", r, err
	})
	KnownSimpleGetMethods[119157] = append(KnownSimpleGetMethods[119157], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiStonfi.GetOrderData(ctx, executor, id)
		return "GetOrderData_StonfiEscrowPositionResult", r, err
	})
	KnownSimpleGetMethods[77915] = append(KnownSimpleGetMethods[77915], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiStonfi.GetCronInfo(ctx, executor, id)
		return "GetCronInfo_StonfiEscrowPositionResult", r, err
	})
	KnownSimpleGetMethods[114667] = append(KnownSimpleGetMethods[114667], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiStonfi.GetVaultData(ctx, executor, id)
		return "GetVaultData_StonfiEscrowVaultResult", r, err
	})

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    StonfiEscrowFactory,
			Results: []string{"GetEscrowData_StonfiEscrowFactoryResult", "GetVersion_StonfiEscrowFactoryResult"},
		},
		InterfaceDescription{
			Name:    StonfiEscrowPosition,
			Results: []string{"GetOrderData_StonfiEscrowPositionResult", "GetCronInfo_StonfiEscrowPositionResult"},
		},
		InterfaceDescription{
			Name:    StonfiEscrowVault,
			Results: []string{"GetVaultData_StonfiEscrowVaultResult"},
		},
	)

	KnownMsgInTypes[abiStonfi.StonfiEscrowWithdrawSignMessageMsgOp] = abiStonfi.EscrowWithdrawSignMessage{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.EscrowWithdrawSignMessage](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixEscrowWithdrawSignMessage), abiStonfi.StonfiEscrowWithdrawSignMessageMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiExternalCronTriggerMsgOp] = abiStonfi.ExternalCronTrigger{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.ExternalCronTrigger](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixExternalCronTrigger), abiStonfi.StonfiExternalCronTriggerMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiExternalItemWithdrawMsgOp] = abiStonfi.ExternalItemWithdraw{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.ExternalItemWithdraw](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixExternalItemWithdraw), abiStonfi.StonfiExternalItemWithdrawMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiItemInternalLockMsgOp] = abiStonfi.ItemInternalLock{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.ItemInternalLock](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixItemInternalLock), abiStonfi.StonfiItemInternalLockMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiItemInternalUnlockMsgOp] = abiStonfi.ItemInternalUnlock{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.ItemInternalUnlock](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixItemInternalUnlock), abiStonfi.StonfiItemInternalUnlockMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiItemLockSuccessMsgOp] = abiStonfi.ItemLockSuccess{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.ItemLockSuccess](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixItemLockSuccess), abiStonfi.StonfiItemLockSuccessMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiItemWithdrawMsgOp] = abiStonfi.ItemWithdraw{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.ItemWithdraw](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixItemWithdraw), abiStonfi.StonfiItemWithdrawMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiLockRejectNotificationMsgOp] = abiStonfi.LockRejectNotification{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.LockRejectNotification](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixLockRejectNotification), abiStonfi.StonfiLockRejectNotificationMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterCancelNextProtocolOwnerMsgOp] = abiStonfi.MinterCancelNextProtocolOwner{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterCancelNextProtocolOwner](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterCancelNextProtocolOwner), abiStonfi.StonfiMinterCancelNextProtocolOwnerMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterDepositVaultMsgOp] = abiStonfi.MinterDepositVault{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterDepositVault](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterDepositVault), abiStonfi.StonfiMinterDepositVaultMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterGiveProtocolOwnershipMsgOp] = abiStonfi.MinterGiveProtocolOwnership{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterGiveProtocolOwnership](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterGiveProtocolOwnership), abiStonfi.StonfiMinterGiveProtocolOwnershipMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterInitTransferMsgOp] = abiStonfi.MinterInitTransfer{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterInitTransfer](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterInitTransfer), abiStonfi.StonfiMinterInitTransferMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterInternalLockMsgOp] = abiStonfi.MinterInternalLock{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterInternalLock](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterInternalLock), abiStonfi.StonfiMinterInternalLockMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterInternalUnlockMsgOp] = abiStonfi.MinterInternalUnlock{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterInternalUnlock](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterInternalUnlock), abiStonfi.StonfiMinterInternalUnlockMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterInternalWithdrawTokensMsgOp] = abiStonfi.MinterInternalWithdrawTokens{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterInternalWithdrawTokens](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterInternalWithdrawTokens), abiStonfi.StonfiMinterInternalWithdrawTokensMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterLockPayloadMsgOp] = abiStonfi.MinterLockPayload{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterLockPayload](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterLockPayload), abiStonfi.StonfiMinterLockPayloadMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterRefundRequestMsgOp] = abiStonfi.MinterRefundRequest{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterRefundRequest](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterRefundRequest), abiStonfi.StonfiMinterRefundRequestMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterResetGasMsgOp] = abiStonfi.MinterResetGas{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterResetGas](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterResetGas), abiStonfi.StonfiMinterResetGasMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterTakeProtocolOwnershipMsgOp] = abiStonfi.MinterTakeProtocolOwnership{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterTakeProtocolOwnership](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterTakeProtocolOwnership), abiStonfi.StonfiMinterTakeProtocolOwnershipMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterUnlockPayloadMsgOp] = abiStonfi.MinterUnlockPayload{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterUnlockPayload](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterUnlockPayload), abiStonfi.StonfiMinterUnlockPayloadMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiMinterUpdateProtocolTierMsgOp] = abiStonfi.MinterUpdateProtocolTier{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.MinterUpdateProtocolTier](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixMinterUpdateProtocolTier), abiStonfi.StonfiMinterUpdateProtocolTierMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiVaultDepositNotificationMsgOp] = abiStonfi.VaultDepositNotification{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.VaultDepositNotification](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixVaultDepositNotification), abiStonfi.StonfiVaultDepositNotificationMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiVaultDepositTokensMsgOp] = abiStonfi.VaultDepositTokens{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.VaultDepositTokens](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixVaultDepositTokens), abiStonfi.StonfiVaultDepositTokensMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiVaultLockMsgOp] = abiStonfi.VaultLock{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.VaultLock](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixVaultLock), abiStonfi.StonfiVaultLockMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiVaultUnlockMsgOp] = abiStonfi.VaultUnlock{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.VaultUnlock](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixVaultUnlock), abiStonfi.StonfiVaultUnlockMsgOp)
	KnownMsgInTypes[abiStonfi.StonfiVaultWithdrawTokensMsgOp] = abiStonfi.VaultWithdrawTokens{}
	registerInMsgUnmarshalerForOpcode[*abiStonfi.VaultWithdrawTokens](opcodedMsgInDecodeFunctions, uint32(abiStonfi.PrefixVaultWithdrawTokens), abiStonfi.StonfiVaultWithdrawTokensMsgOp)

}

func init() {
	tolkMethods = append(tolkMethods,
		MethodDescription{
			Name: "get_bonding_data",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiStonkspump.GetBondingData(ctx, executor, id)
				return "GetBondingData_StonksPumpVirtualMinterResult", r, err
			},
		},
	)

	KnownGetMethodsDecoder["get_bonding_data"] = append(KnownGetMethodsDecoder["get_bonding_data"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiStonkspump.DecodeGetBondingData(&st)
		return "GetBondingData_StonksPumpVirtualMinterResult", r, err
	})

	KnownSimpleGetMethods[119688] = append(KnownSimpleGetMethods[119688], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiStonkspump.GetBondingData(ctx, executor, id)
		return "GetBondingData_StonksPumpVirtualMinterResult", r, err
	})

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    StonksPumpVirtualMinter,
			Results: []string{"GetBondingData_StonksPumpVirtualMinterResult"},
		},
	)

	KnownMsgInTypes[abiStonkspump.StonkspumpAskToPresaleSellMsgOp] = abiStonkspump.AskToPresaleSell{}
	registerInMsgUnmarshalerForOpcode[*abiStonkspump.AskToPresaleSell](opcodedMsgInDecodeFunctions, uint32(abiStonkspump.PrefixAskToPresaleSell), abiStonkspump.StonkspumpAskToPresaleSellMsgOp)
	KnownMsgInTypes[abiStonkspump.StonkspumpBuyFromPresaleMsgOp] = abiStonkspump.BuyFromPresale{}
	registerInMsgUnmarshalerForOpcode[*abiStonkspump.BuyFromPresale](opcodedMsgInDecodeFunctions, uint32(abiStonkspump.PrefixBuyFromPresale), abiStonkspump.StonkspumpBuyFromPresaleMsgOp)
	KnownMsgInTypes[abiStonkspump.StonkspumpCreateVirtualLiquidityJettonMsgOp] = abiStonkspump.CreateVirtualLiquidityJetton{}
	registerInMsgUnmarshalerForOpcode[*abiStonkspump.CreateVirtualLiquidityJetton](opcodedMsgInDecodeFunctions, uint32(abiStonkspump.PrefixCreateVirtualLiquidityJetton), abiStonkspump.StonkspumpCreateVirtualLiquidityJettonMsgOp)
	KnownMsgInTypes[abiStonkspump.StonkspumpPresaleSellNotificationForMinterMsgOp] = abiStonkspump.PresaleSellNotificationForMinter{}
	registerInMsgUnmarshalerForOpcode[*abiStonkspump.PresaleSellNotificationForMinter](opcodedMsgInDecodeFunctions, uint32(abiStonkspump.PrefixPresaleSellNotificationForMinter), abiStonkspump.StonkspumpPresaleSellNotificationForMinterMsgOp)
	KnownMsgInTypes[abiStonkspump.StonkspumpPresaleTradeEventMsgOp] = abiStonkspump.PresaleTradeEvent{}
	registerInMsgUnmarshalerForOpcode[*abiStonkspump.PresaleTradeEvent](opcodedMsgInDecodeFunctions, uint32(abiStonkspump.PrefixPresaleTradeEvent), abiStonkspump.StonkspumpPresaleTradeEventMsgOp)

}

func init() {
	tolkMethods = append(tolkMethods,
		MethodDescription{
			Name: "get_user_latest_version",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiXtr.GetUserLatestVersion(ctx, executor, id)
				return "GetUserLatestVersion_XtrMasterResult", r, err
			},
		},
		MethodDescription{
			Name: "get_payment_latest_version",
			InvokeFn: func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
				r, err := abiXtr.GetPaymentLatestVersion(ctx, executor, id)
				return "GetPaymentLatestVersion_XtrMasterResult", r, err
			},
		},
	)

	KnownGetMethodsDecoder["get_user_latest_version"] = append(KnownGetMethodsDecoder["get_user_latest_version"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiXtr.DecodeGetUserLatestVersion(&st)
		return "GetUserLatestVersion_XtrMasterResult", r, err
	})
	KnownGetMethodsDecoder["get_payment_latest_version"] = append(KnownGetMethodsDecoder["get_payment_latest_version"], func(stack tlb.VmStack) (string, any, error) {
		st := stack
		r, err := abiXtr.DecodeGetPaymentLatestVersion(&st)
		return "GetPaymentLatestVersion_XtrMasterResult", r, err
	})

	KnownSimpleGetMethods[98392] = append(KnownSimpleGetMethods[98392], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiXtr.GetUserLatestVersion(ctx, executor, id)
		return "GetUserLatestVersion_XtrMasterResult", r, err
	})
	KnownSimpleGetMethods[126127] = append(KnownSimpleGetMethods[126127], func(ctx context.Context, executor Executor, id ton.AccountID) (string, any, error) {
		r, err := abiXtr.GetPaymentLatestVersion(ctx, executor, id)
		return "GetPaymentLatestVersion_XtrMasterResult", r, err
	})

	tolkInterfaceOrder = append(tolkInterfaceOrder,
		InterfaceDescription{
			Name:    XtrMaster,
			Results: []string{"GetUserLatestVersion_XtrMasterResult", "GetPaymentLatestVersion_XtrMasterResult"},
		},
	)

	KnownMsgInTypes[abiXtr.XtrCommitXTRMsgOp] = abiXtr.CommitXTR{}
	registerInMsgUnmarshalerForOpcode[*abiXtr.CommitXTR](opcodedMsgInDecodeFunctions, uint32(abiXtr.PrefixCommitXTR), abiXtr.XtrCommitXTRMsgOp)
	KnownMsgInTypes[abiXtr.XtrPushXTRMsgOp] = abiXtr.PushXTR{}
	registerInMsgUnmarshalerForOpcode[*abiXtr.PushXTR](opcodedMsgInDecodeFunctions, uint32(abiXtr.PrefixPushXTR), abiXtr.XtrPushXTRMsgOp)
	KnownMsgInTypes[abiXtr.XtrUpdateContractAndProcessMessageMsgOp] = abiXtr.UpdateContractAndProcessMessage{}
	registerInMsgUnmarshalerForOpcode[*abiXtr.UpdateContractAndProcessMessage](opcodedMsgInDecodeFunctions, uint32(abiXtr.PrefixUpdateContractAndProcessMessage), abiXtr.XtrUpdateContractAndProcessMessageMsgOp)
	KnownMsgInTypes[abiXtr.XtrUpdatePaymentMsgOp] = abiXtr.UpdatePayment{}
	registerInMsgUnmarshalerForOpcode[*abiXtr.UpdatePayment](opcodedMsgInDecodeFunctions, uint32(abiXtr.PrefixUpdatePayment), abiXtr.XtrUpdatePaymentMsgOp)
	KnownMsgInTypes[abiXtr.XtrUpdateUserMsgOp] = abiXtr.UpdateUser{}
	registerInMsgUnmarshalerForOpcode[*abiXtr.UpdateUser](opcodedMsgInDecodeFunctions, uint32(abiXtr.PrefixUpdateUser), abiXtr.XtrUpdateUserMsgOp)

}
