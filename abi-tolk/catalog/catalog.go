// Package catalog reads the contract catalog published at
// https://github.com/ton-blockchain/abis.
//
// That repository describes contracts in info.toml files and compiles the Tolk
// types of each of them into an ABI; it publishes both as a single JSON bundle,
// which is what this package reads:
//
//	c, err := catalog.LoadFile("abi-catalog.json")
//	for _, p := range c.Projects {
//		for _, contract := range p.Contracts {
//			fmt.Println(p.Name, contract.Name, contract.Hashes)
//		}
//	}
//
// The ABI of a contract is ready for tolk/runtime to decode the storage and the
// messages of accounts running it.
package catalog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

const SchemaVersion = 1

type Catalog struct {
	// Projects are sorted by name.
	Projects []*Project

	byName map[string]*Project
}

func Load(r io.Reader) (*Catalog, error) {
	var b bundle
	if err := json.NewDecoder(r).Decode(&b); err != nil {
		return nil, err
	}
	if b.SchemaVersion != SchemaVersion {
		return nil, fmt.Errorf("catalog schema version %d, want %d", b.SchemaVersion, SchemaVersion)
	}
	c := &Catalog{byName: make(map[string]*Project)}
	for _, entry := range b.Contracts {
		if err := c.add(entry); err != nil {
			return nil, err
		}
	}
	slices.SortFunc(c.Projects, func(a, b *Project) int { return strings.Compare(a.Name, b.Name) })
	for _, p := range c.Projects {
		slices.SortFunc(p.Contracts, func(a, b *Contract) int { return strings.Compare(a.Name, b.Name) })
	}
	return c, nil
}

func LoadFile(name string) (*Catalog, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	c, err := Load(f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	return c, nil
}

func (c *Catalog) Project(name string) (*Project, bool) {
	p, ok := c.byName[name]
	return p, ok
}

func (c *Catalog) add(entry BundledContract) error {
	projectName, name, ok := strings.Cut(entry.ID, ".")
	if !ok {
		return fmt.Errorf("contract id %q is not \"<project>.<contract>\"", entry.ID)
	}
	p := c.project(projectName)
	if declared, ok := p.byName[name]; ok {
		return fmt.Errorf("contract %q is declared twice, as %q and %q", name, declared.id, entry.ID)
	}
	contract := &Contract{
		Name:           name,
		DisplayName:    entry.DisplayName,
		Hashes:         entry.Hashes,
		KnownAddresses: entry.KnownAddresses,
		Links:          entry.Links,
		ABI:            entry.ABI,
		id:             entry.ID,
		project:        p,
	}
	p.Contracts = append(p.Contracts, contract)
	p.byName[name] = contract
	return nil
}

// project finds or creates the project a bundle name belongs to.
//
// The bundle names a project after the directory describing it, and some of
// them sit one directory deeper: the wallets are spread over "wallets/w4r2",
// "wallets/lockup_vesting" and the like. Taking the first segment of that path
// puts every wallet in the same "wallets" project.
func (c *Catalog) project(name string) *Project {
	name, _, _ = strings.Cut(name, "/")
	p, ok := c.byName[name]
	if !ok {
		p = &Project{Name: name, byName: make(map[string]*Contract)}
		c.byName[name] = p
		c.Projects = append(c.Projects, p)
	}
	return p
}
