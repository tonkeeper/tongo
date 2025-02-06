package liteclient

//go:generate go run generator.go

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tonkeeper/tongo/tl"
	"github.com/tonkeeper/tongo/ton"
)

var (
	ErrBlockNotApplied = fmt.Errorf("block is not applied")
)

func (t LiteServerErrorC) Error() string {
	return fmt.Sprintf("error code: %d message: %s", t.Code, t.Message)
}

func (t LiteServerErrorC) IsNotApplied() bool {
	return t.Message == "block is not applied"
}

type LiteServerSignatureSet LiteServerSignatureSetC

func (t LiteServerSignatureSet) MarshalTL() ([]byte, error) {
	var tag uint32 = 0xf644a6e6
	buf := new(bytes.Buffer)
	b, err := tl.Marshal(tag)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(LiteServerSignatureSetC(t))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerSignatureSet) UnmarshalTL(r io.Reader) error {
	var (
		res LiteServerSignatureSetC
		tag uint32
	)
	err := tl.Unmarshal(r, &tag)
	if err != nil {
		return err
	}
	if tag != 0xf644a6e6 {
		return fmt.Errorf("invalid tag")
	}
	err = tl.Unmarshal(r, &res)
	if err != nil {
		return err
	}
	*t = LiteServerSignatureSet(res)
	return nil
}

func (t TonNodeBlockIdExtC) ToBlockIdExt() ton.BlockIDExt {
	res := ton.BlockIDExt{
		RootHash: ton.Bits256(t.RootHash),
		FileHash: ton.Bits256(t.FileHash),
	}
	res.Seqno = t.Seqno
	res.Shard = t.Shard
	res.Workchain = t.Workchain
	return res
}

func AccountID(id ton.AccountID) LiteServerAccountIdC {
	return LiteServerAccountIdC{
		Workchain: id.Workchain,
		Id:        id.Address,
	}
}

func BlockIDExt(id ton.BlockIDExt) TonNodeBlockIdExtC {
	return TonNodeBlockIdExtC{
		Workchain: id.Workchain,
		Shard:     id.Shard,
		Seqno:     id.Seqno,
		RootHash:  tl.Int256(id.RootHash),
		FileHash:  tl.Int256(id.FileHash),
	}
}
