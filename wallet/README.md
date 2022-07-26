## Tools to simplify the deployment and interaction with the wallet smart contract

This library supports:
1. Address generation from public or private keys for V1-V4 wallets
2. Message generation for deploying new wallets
3. Message generation for TON transfer (long comment supported) and custom payload 
(for V1-V4 supports up to 4 payload messages)

Wallets smart contracts description:
* [Wallets](https://github.com/toncenter/tonweb/blob/master/src/contract/wallet/WalletSources.md)

### Usage
[Example](../examples/wallet/main.go)