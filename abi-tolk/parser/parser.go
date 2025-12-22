package parser

import (
	"encoding/json"

	"github.com/tonkeeper/tongo/tolk/abi"
)

func ParseABI(s []byte) (tolkAbi.ABI, error) {
	var abi tolkAbi.ABI
	err := json.Unmarshal(s, &abi)
	return abi, err
}
