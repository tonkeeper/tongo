package parser

import (
	"encoding/json"

	"github.com/tonkeeper/tongo/tolk"
)

func ParseABI(s []byte) (tolk.ABI, error) {
	var abi tolk.ABI
	err := json.Unmarshal(s, &abi)
	return abi, err
}
