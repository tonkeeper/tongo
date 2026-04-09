package abi

import "fmt"

var (
	KnownJettonTypes          = mergeNoConflicts[JettonOpName, any](xmlKnownJettonTypes, tolkKnownJettonTypes)
	JettonOpCodes             = mergeNoConflicts[JettonOpName, JettonOpCode](xmlJettonOpCodes, tolkJettonOpCodes)
	funcJettonDecodersMapping = mergeNoConflicts[JettonOpCode, jettonDecoder](xmlJettonDecodersMapping, tolkJettonDecodersMapping)
)

func mergeNoConflicts[K comparable, V any](parts ...map[K]V) map[K]V {
	total := 0
	for _, p := range parts {
		total += len(p)
	}
	out := make(map[K]V, total)
	for _, p := range parts {
		for k, v := range p {
			if _, exists := out[k]; exists {
				panic(fmt.Sprintf("duplicate jetton registry key: %v", k))
			}
			out[k] = v
		}
	}
	return out
}
