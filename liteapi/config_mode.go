package liteapi

type ConfigMode uint32

// constants below are different flags for "mode" param of GetConfigAll.
// one needs to OR them to get different aspects of TON's configuration.
const (
	NeedStateRoot      ConfigMode = 1
	NeedLibraries      ConfigMode = 2
	NeedStateExtraRoot ConfigMode = 4
	NeedShardHashes    ConfigMode = 8
	NeedValidatorSet   ConfigMode = 16
	NeedSpecialSmc     ConfigMode = 32
	NeedAccountsRoot   ConfigMode = 64
	NeedPrevBlocks     ConfigMode = 128
	NeedWorkchainInfo  ConfigMode = 256
	NeedCapabilities   ConfigMode = 512
)
