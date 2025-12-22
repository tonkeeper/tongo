
# List of supported message opcodes

The first 4 bytes of a message's body identify the `operation` to be performed, or the `method` of the smart contract to be invoked.

The list below contains the supported message operations, their names and opcodes.

| Name        | Message operation code |
|-------------|------------------------|
| AboaLisa| 0xd53276db |
| AskToBurn| 0x595f07bc |
| AskToTransfer| 0x0f8a7ea5 |
| BurnNotificationForMinter| 0x7bdd97de |
| ChangeMinterAdmin| 0x00000003 |
| ChangeMinterContent| 0x00000004 |
| CocoonClientClientProxyRequest| 0x65448ff4 |
| CocoonClientExtClientChargeSigned| 0xbb63ff93 |
| CocoonClientExtClientGrantRefundSigned| 0xefd711e1 |
| CocoonClientExtClientTopUp| 0xf172e6c2 |
| CocoonClientOwnerClientChangeSecretHash| 0xa9357034 |
| CocoonClientOwnerClientChangeSecretHashAndTopUp| 0x8473b408 |
| CocoonClientOwnerClientIncreaseStake| 0x6a1f6a60 |
| CocoonClientOwnerClientRegister| 0xc45f9f3b |
| CocoonClientOwnerClientRequestRefund| 0xfafa6cc1 |
| CocoonClientOwnerClientWithdraw| 0xda068e78 |
| CocoonClientPayout| 0xc59a7cd3 |
| CocoonClientReturnExcessesBack| 0xd53276db |
| CocoonProxyClientProxyRequest| 0x65448ff4 |
| CocoonProxyExtProxyCloseCompleteRequestSigned| 0xe511abc7 |
| CocoonProxyExtProxyCloseRequestSigned| 0x636a4391 |
| CocoonProxyExtProxyIncreaseStake| 0x9713f187 |
| CocoonProxyExtProxyPayoutRequest| 0x7610e6eb |
| CocoonProxyOwnerProxyClose| 0xb51d5a01 |
| CocoonProxyPayout| 0xc59a7cd3 |
| CocoonProxyReturnExcessesBack| 0xd53276db |
| CocoonProxyTextCmd| 0x00000000 |
| CocoonProxyWorkerProxyRequest| 0x4d725d2c |
| CocoonRootAddModelType| 0xc146134d |
| CocoonRootAddProxyType| 0x71860e80 |
| CocoonRootAddWorkerType| 0xe34b1c60 |
| CocoonRootChangeFees| 0xc52ed8d4 |
| CocoonRootChangeOwner| 0xc4a1ae54 |
| CocoonRootChangeParams| 0x022fa189 |
| CocoonRootDelModelType| 0x92b11c18 |
| CocoonRootDelProxyType| 0x3c41d0b2 |
| CocoonRootDelWorkerType| 0x8d94a79a |
| CocoonRootPayout| 0xc59a7cd3 |
| CocoonRootRegisterProxy| 0x927c7cb5 |
| CocoonRootResetRoot| 0x563c1d96 |
| CocoonRootReturnExcessesBack| 0xd53276db |
| CocoonRootUnregisterProxy| 0x6d49eaf2 |
| CocoonRootUpdateProxy| 0x9c7924ba |
| CocoonRootUpgradeCode| 0x11aefd51 |
| CocoonRootUpgradeContracts| 0xa2370f61 |
| CocoonRootUpgradeFull| 0x4f7c5789 |
| CocoonWalletOwnerWalletSendMessage| 0x9c69f376 |
| CocoonWalletPayout| 0xc59a7cd3 |
| CocoonWalletReturnExcessesBack| 0xd53276db |
| CocoonWalletTextCommand| 0x00000000 |
| CocoonWorkerExtWorkerLastPayoutRequestSigned| 0xf5f26a36 |
| CocoonWorkerExtWorkerPayoutRequestSigned| 0xa040ad28 |
| CocoonWorkerOwnerWorkerRegister| 0x26ed7f65 |
| CocoonWorkerPayout| 0xc59a7cd3 |
| CocoonWorkerReturnExcessesBack| 0xd53276db |
| CocoonWorkerWorkerProxyRequest| 0x4d725d2c |
| InternalTransferStep| 0x178d4519 |
| MintNewJettons| 0x00000015 |
| RequestWalletAddress| 0x2c76b973 |
| ResponseWalletAddress| 0xd1735400 |
| ReturnExcessesBack| 0xd53276db |
| TransferNotificationForRecipient| 0x7362d09c |
