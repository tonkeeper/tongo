package boc

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func DeserializeSingleRootBoc(boc []byte) (*Cell, error) {
	cells, err := DeserializeBoc(boc)
	if err != nil {
		return nil, err
	}
	if len(cells) != 1 {
		return nil, fmt.Errorf("invalid boc roots number %v", len(cells))
	}
	return cells[0], nil
}

func DeserializeBocBase64(boc string) ([]*Cell, error) {
	bocData, err := base64.StdEncoding.DecodeString(boc)
	if err != nil {
		return nil, err
	}
	return DeserializeBoc(bocData)
}

func DeserializeSinglRootBase64(boc string) (*Cell, error) {
	cells, err := DeserializeBocBase64(boc)
	if err != nil {
		return nil, err
	}
	if len(cells) != 1 {
		return nil, fmt.Errorf("invalid boc roots number %v", len(cells))
	}
	return cells[0], nil
}

func DeserializeBocHex(boc string) ([]*Cell, error) {
	bocData, err := hex.DecodeString(boc)
	if err != nil {
		return nil, err
	}
	return DeserializeBoc(bocData)
}

func DeserializeSinglRootHex(boc string) (*Cell, error) {
	cells, err := DeserializeBocHex(boc)
	if err != nil {
		return nil, err
	}
	if len(cells) != 1 {
		return nil, fmt.Errorf("invalid boc roots number %v", len(cells))
	}
	return cells[0], nil
}

func MustDeserializeSinglRootHex(boc string) *Cell {
	c, err := DeserializeSinglRootHex(boc)
	if err != nil {
		panic(err)
	}
	return c
}

func MustDeserializeSinglRootBase64(boc string) *Cell {
	c, err := DeserializeSinglRootBase64(boc)
	if err != nil {
		panic(err)
	}
	return c
}

func MustBitStringFromFiftHex(hexRepr string) *BitString {
	bs, err := BitStringFromFiftHex(hexRepr)
	if err != nil {
		panic(err)
	}
	return bs
}
