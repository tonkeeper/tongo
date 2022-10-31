# tongo
***
Go implementation of libraries for TON blockchain.

## Library structure
1. [ADNL](adnl/README.md) - low level adnl protocol implementation
2. [Lite client](liteclient/README.md) - interaction with TON node as lite client
3. [BOC](boc/README.md) - cells and bag-of-cells methods and primitives
4. [TL](tl/README.md) - interaction with binary data described by TL (Type Language) schemas
5. [TLB](tlb/README.md) - interaction with binary data (in Cells) described by TL-B (Typed Language - Binary) schemas
6. [TVM](tvm/README.md) - interaction with TVM (TON Virtual Machine)
7. [Wallet](wallet/README.md) - tools to simplify the deployment and interaction with the wallet smart contract
8. [Contract](contract/README.md) - tools to simplify the interaction with the smart contracts like Jettons and NFT
9. [Examples](examples)

## Dependencies
### Libraries
For TVM executor you need a libraries from `lib/darwin` (MAC) or `lib/linux`
### Connection to TON node
For connect to TON node you need to know public key and `ip:port`. In most cases you can use public config files. 
Download `global-config.json` (mainnet) or `testnet-global.config.json` (testnet) file from [ton.org](https://ton.org/docs/#/)
Lite client supports auto-download mainnet config from ton.org.
## Package installation

```shell
go get github.com/startfellows/tongo
```

## Basic types
Tongo operates with TON blockchain structures described in [block.tlb](https://github.com/ton-blockchain/ton/blob/master/crypto/block/block.tlb)
and some types described in [lite_api.tl](https://github.com/ton-blockchain/ton/blob/master/tl/generate/scheme/lite_api.tl).
Go definitions of this types you can find in files: `account.go`, `transactions.go`, `models.go` ... 

