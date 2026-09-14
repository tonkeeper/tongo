package txemulator

import (
	_ "embed"
	"encoding/base64"
)

// mainnetConfigBoc holds result of client.GetConfigAll().Config
//
//go:embed mainnet_config.boc
var mainnetConfigBoc []byte

var DefaultConfig = base64.StdEncoding.EncodeToString(mainnetConfigBoc)
