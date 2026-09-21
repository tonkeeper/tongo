package catalog

import "github.com/tonkeeper/tongo/tolk/parser"

type Project struct {
	// Name is unique within the catalog.
	Name string
	// Contracts are sorted by name.
	Contracts []*Contract

	byName map[string]*Contract
}

func (p *Project) Contract(name string) (*Contract, bool) {
	c, ok := p.byName[name]
	return c, ok
}

type Contract struct {
	// Name is unique within the project.
	Name        string
	DisplayName string
	// Hashes are lowercase hex.
	Hashes []string
	// KnownAddresses are in the user-friendly form, and run one of Hashes.
	KnownAddresses []string
	Links          []Link
	ABI            parser.ContractABI

	// id is the name the bundle gives the contract, kept for diagnostics.
	id      string
	project *Project
}

func (c *Contract) Project() *Project { return c.project }
