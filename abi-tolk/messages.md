
# List of supported message opcodes

The first 4 bytes of a message's body identify the `operation` to be performed, or the `method` of the smart contract to be invoked.

The list below contains the supported message operations, their names and opcodes.

| Name        | Message operation code |
|-------------|------------------------|
| TonCocoonAddModelType| 0xc146134d |
| TonCocoonAddProxyType| 0x71860e80 |
| TonCocoonAddWorkerType| 0xe34b1c60 |
| TonCocoonChangeFees| 0xc52ed8d4 |
| TonCocoonChangeOwner| 0xc4a1ae54 |
| TonCocoonChangeParams| 0x022fa189 |
| TonCocoonClientProxyRequest| 0x65448ff4 |
| TonCocoonDelModelType| 0x92b11c18 |
| TonCocoonDelProxyType| 0x3c41d0b2 |
| TonCocoonDelWorkerType| 0x8d94a79a |
| TonCocoonExtClientChargeSigned| 0xbb63ff93 |
| TonCocoonExtClientGrantRefundSigned| 0xefd711e1 |
| TonCocoonExtClientTopUp| 0xf172e6c2 |
| TonCocoonExtProxyCloseCompleteRequestSigned| 0xe511abc7 |
| TonCocoonExtProxyCloseRequestSigned| 0x636a4391 |
| TonCocoonExtProxyIncreaseStake| 0x9713f187 |
| TonCocoonExtProxyPayoutRequest| 0x7610e6eb |
| TonCocoonExtWorkerLastPayoutRequestSigned| 0xf5f26a36 |
| TonCocoonExtWorkerPayoutRequestSigned| 0xa040ad28 |
| TonCocoonOwnerClientChangeSecretHash| 0xa9357034 |
| TonCocoonOwnerClientChangeSecretHashAndTopUp| 0x8473b408 |
| TonCocoonOwnerClientIncreaseStake| 0x6a1f6a60 |
| TonCocoonOwnerClientRegister| 0xc45f9f3b |
| TonCocoonOwnerClientRequestRefund| 0xfafa6cc1 |
| TonCocoonOwnerClientWithdraw| 0xda068e78 |
| TonCocoonOwnerProxyClose| 0xb51d5a01 |
| TonCocoonOwnerWalletSendMessage| 0x9c69f376 |
| TonCocoonOwnerWorkerRegister| 0x26ed7f65 |
| TonCocoonPayout| 0xc59a7cd3 |
| TonCocoonRegisterProxy| 0x927c7cb5 |
| TonCocoonResetRoot| 0x563c1d96 |
| TonCocoonReturnExcessesBack| 0xd53276db |
| TonCocoonTextCmd| 0x00000000 |
| TonCocoonTextCommand| 0x00000000 |
| TonCocoonUnregisterProxy| 0x6d49eaf2 |
| TonCocoonUpdateProxy| 0x9c7924ba |
| TonCocoonUpgradeCode| 0x11aefd51 |
| TonCocoonUpgradeContracts| 0xa2370f61 |
| TonCocoonUpgradeFull| 0x4f7c5789 |
| TonCocoonWorkerProxyRequest| 0x4d725d2c |
| TonTep74AboaLisa| 0xd53276db |
| TonTep74AskToBurn| 0x595f07bc |
| TonTep74AskToTransfer| 0x0f8a7ea5 |
| TonTep74BurnNotificationForMinter| 0x7bdd97de |
| TonTep74ChangeMinterAdmin| 0x00000003 |
| TonTep74ChangeMinterContent| 0x00000004 |
| TonTep74InternalTransferStep| 0x178d4519 |
| TonTep74MintNewJettons| 0x00000015 |
| TonTep74RequestWalletAddress| 0x2c76b973 |
| TonTep74ResponseWalletAddress| 0xd1735400 |
| TonTep74ReturnExcessesBack| 0xd53276db |
| TonTep74TransferNotificationForRecipient| 0x7362d09c |
