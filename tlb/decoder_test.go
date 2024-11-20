package tlb

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/tonkeeper/tongo/boc"
)

func Test_UnmarshalMessage(t *testing.T) {
	rawBody := "b5ee9c7201010101000e000018d53276db546de4efe9d175f0"
	ex, err := newParseTokenExcesses(rawBody)
	assert.NoError(t, err)
	t.Logf("%+v", ex)
	ex1, err := parseTokenExcesses(rawBody)
	assert.Equal(t, ex1.QueryId, ex.QueryId)
}

type TokenExcesses struct {
	QueryId uint64
}

func newParseTokenExcesses(rawBody string) (*TokenExcesses, error) {
	cells, err := boc.DeserializeBocHex(rawBody)
	if err != nil {
		return nil, err
	}
	ex := TokenExcesses{}
	opCode, err := UnmarshalMessage(cells[0], &ex)
	if err != nil {
		return nil, err
	}
	if opCode != uint32(0xd53276db) {
		return nil, fmt.Errorf("unexpected opCode: %x", opCode)
	}
	//return struct without opCode
	return &ex, nil
}

type TokenExcessesWithOpCode struct {
	OpCode  uint32
	QueryId uint64
}

func parseTokenExcesses(rawBody string) (*TokenExcesses, error) {
	cells, err := boc.DeserializeBocHex(rawBody)
	if err != nil {
		return nil, err
	}
	//define new struct with opCode
	ex := TokenExcessesWithOpCode{}
	err = Unmarshal(cells[0], &ex)
	if err != nil {
		return nil, err
	}
	//return struct without opCode
	return &TokenExcesses{QueryId: ex.QueryId}, nil
}
func Test_Unmarshal(t *testing.T) {
	rawBody := "b5ee9c7201010101000e000018d53276db546de4efe9d175f0"
	ex, err := parseTokenExcesses(rawBody)
	assert.NoError(t, err)
	t.Logf("%+v", ex)
}
