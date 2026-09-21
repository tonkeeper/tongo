package catalog

import "github.com/tonkeeper/tongo/tolk/parser"

// bundle Compiled ABI entries sorted by their stable contract IDs
type bundle struct {
	SchemaVersion int               `json:"schemaVersion"`
	Contracts     []BundledContract `json:"contracts"`
}

// BundledContract One compiled contract entry in the public ABI catalog
type BundledContract struct {
	// ID names the contract as "<project>.<contract>".
	ID             string             `json:"id"`
	DisplayName    string             `json:"displayName"`
	Hashes         []string           `json:"hashes"`
	KnownAddresses []string           `json:"knownAddresses"`
	Links          []Link             `json:"links"`
	ABI            parser.ContractABI `json:"compilerAbi"`
}

// Link External reference attached to a project or contract
type Link struct {
	Kind  string `json:"kind"`
	Title string `json:"title"`
	URL   string `json:"url"`
}
