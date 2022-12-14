package liteclient

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/startfellows/tongo/tl"
	"io"
)

type LiteServerTransactionList struct {
	Ids          []TonNodeBlockIdExt
	Transactions []byte
}

func (t LiteServerTransactionList) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Ids)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Transactions)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerTransactionList) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Ids)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Transactions)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerTransactionId3 struct {
	Account tl.Int256
	Lt      int64
}

func (t LiteServerTransactionId3) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Account)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Lt)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerTransactionId3) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Account)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Lt)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerBlockData struct {
	Id   TonNodeBlockIdExt
	Data []byte
}

func (t LiteServerBlockData) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Data)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerBlockData) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Data)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerSendMsgStatus struct {
	Status int32
}

func (t LiteServerSendMsgStatus) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Status)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerSendMsgStatus) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Status)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerMasterchainInfoExt struct {
	Mode          int32
	Version       int32
	Capabilities  int64
	Last          TonNodeBlockIdExt
	LastUtime     int32
	Now           int32
	StateRootHash tl.Int256
	Init          TonNodeZeroStateIdExt
}

func (t LiteServerMasterchainInfoExt) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Version)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Capabilities)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Last)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.LastUtime)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Now)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.StateRootHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Init)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerMasterchainInfoExt) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Mode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Version)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Capabilities)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Last)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.LastUtime)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Now)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.StateRootHash)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Init)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerTransactionInfo struct {
	Id          TonNodeBlockIdExt
	Proof       []byte
	Transaction []byte
}

func (t LiteServerTransactionInfo) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Proof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Transaction)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerTransactionInfo) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Proof)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Transaction)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerBlockTransactions struct {
	Id         TonNodeBlockIdExt
	ReqCount   int32
	Incomplete bool
	Ids        []LiteServerTransactionId
	Proof      []byte
}

func (t LiteServerBlockTransactions) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ReqCount)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Incomplete)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Ids)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Proof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerBlockTransactions) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ReqCount)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Incomplete)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Ids)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Proof)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerSignatureSet struct {
	ValidatorSetHash int32
	CatchainSeqno    int32
	Signatures       []LiteServerSignature
}

func (t LiteServerSignatureSet) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.ValidatorSetHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.CatchainSeqno)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Signatures)
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
	var err error
	err = tl.Unmarshal(r, &t.ValidatorSetHash)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.CatchainSeqno)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Signatures)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerValidatorStats struct {
	Mode       int32
	Id         TonNodeBlockIdExt
	Count      int32
	Complete   bool
	StateProof []byte
	DataProof  []byte
}

func (t LiteServerValidatorStats) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Count)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Complete)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.StateProof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.DataProof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerValidatorStats) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Mode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Count)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Complete)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.StateProof)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.DataProof)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerError struct {
	Code    int32
	Message string
}

func (t LiteServerError) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Code)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Message)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerError) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Code)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Message)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerAccountId struct {
	Workchain int32
	Id        tl.Int256
}

func (t LiteServerAccountId) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Workchain)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerAccountId) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Workchain)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	return nil
}

type TonNodeBlockIdExt struct {
	Workchain int32
	Shard     int64
	Seqno     int32
	RootHash  tl.Int256
	FileHash  tl.Int256
}

func (t TonNodeBlockIdExt) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Workchain)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Shard)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Seqno)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RootHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.FileHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *TonNodeBlockIdExt) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Workchain)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Shard)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Seqno)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RootHash)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.FileHash)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerVersion struct {
	Mode         int32
	Version      int32
	Capabilities int64
	Now          int32
}

func (t LiteServerVersion) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Version)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Capabilities)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Now)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerVersion) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Mode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Version)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Capabilities)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Now)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerShardInfo struct {
	Id         TonNodeBlockIdExt
	Shardblk   TonNodeBlockIdExt
	ShardProof []byte
	ShardDescr []byte
}

func (t LiteServerShardInfo) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Shardblk)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ShardProof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ShardDescr)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerShardInfo) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Shardblk)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ShardProof)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ShardDescr)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerAllShardsInfo struct {
	Id    TonNodeBlockIdExt
	Proof []byte
	Data  []byte
}

func (t LiteServerAllShardsInfo) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Proof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Data)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerAllShardsInfo) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Proof)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Data)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerSignature struct {
	NodeIdShort tl.Int256
	Signature   []byte
}

func (t LiteServerSignature) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.NodeIdShort)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Signature)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerSignature) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.NodeIdShort)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Signature)
	if err != nil {
		return err
	}
	return nil
}

type AdnlMessage struct {
	tl.SumType
	AdnlMessageQuery struct {
		QueryId tl.Int256
		Query   []byte
	}
	AdnlMessageAnswer struct {
		QueryId tl.Int256
		Answer  []byte
	}
}

func (t AdnlMessage) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	switch t.SumType {
	case "AdnlMessageQuery":
		b, err = tl.Marshal(uint32(0x7af98bb4))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.AdnlMessageQuery.QueryId)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.AdnlMessageQuery.Query)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	case "AdnlMessageAnswer":
		b, err = tl.Marshal(uint32(0x1684ac0f))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.AdnlMessageAnswer.QueryId)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.AdnlMessageAnswer.Answer)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid sum type")
	}
	return buf.Bytes(), nil
}

func (t *AdnlMessage) UnmarshalTL(r io.Reader) error {
	var err error
	var b [4]byte
	_, err = io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	tag := int(binary.LittleEndian.Uint32(b[:]))
	switch tag {
	case 0x7af98bb4:
		t.SumType = "AdnlMessageQuery"
		err = tl.Unmarshal(r, &t.AdnlMessageQuery.QueryId)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.AdnlMessageQuery.Query)
		if err != nil {
			return err
		}
	case 0x1684ac0f:
		t.SumType = "AdnlMessageAnswer"
		err = tl.Unmarshal(r, &t.AdnlMessageAnswer.QueryId)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.AdnlMessageAnswer.Answer)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid tag")
	}
	return nil
}

type LiteServerCurrentTime struct {
	Now int32
}

func (t LiteServerCurrentTime) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Now)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerCurrentTime) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Now)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerBlockHeader struct {
	Id          TonNodeBlockIdExt
	Mode        int32
	HeaderProof []byte
}

func (t LiteServerBlockHeader) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.HeaderProof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerBlockHeader) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Mode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.HeaderProof)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerBlockLink struct {
	tl.SumType
	LiteServerBlockLinkBack struct {
		ToKeyBlock bool
		From       TonNodeBlockIdExt
		To         TonNodeBlockIdExt
		DestProof  []byte
		Proof      []byte
		StateProof []byte
	}
	LiteServerBlockLinkForward struct {
		ToKeyBlock  bool
		From        TonNodeBlockIdExt
		To          TonNodeBlockIdExt
		DestProof   []byte
		ConfigProof []byte
		Signatures  LiteServerSignatureSet
	}
}

func (t LiteServerBlockLink) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	switch t.SumType {
	case "LiteServerBlockLinkBack":
		b, err = tl.Marshal(uint32(0xef1b7eef))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.LiteServerBlockLinkBack.ToKeyBlock)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.LiteServerBlockLinkBack.From)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.LiteServerBlockLinkBack.To)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.LiteServerBlockLinkBack.DestProof)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.LiteServerBlockLinkBack.Proof)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.LiteServerBlockLinkBack.StateProof)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	case "LiteServerBlockLinkForward":
		b, err = tl.Marshal(uint32(0x1cce0f52))
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		b, err = tl.Marshal(t.LiteServerBlockLinkForward.ToKeyBlock)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.LiteServerBlockLinkForward.From)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.LiteServerBlockLinkForward.To)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.LiteServerBlockLinkForward.DestProof)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.LiteServerBlockLinkForward.ConfigProof)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
		b, err = tl.Marshal(t.LiteServerBlockLinkForward.Signatures)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid sum type")
	}
	return buf.Bytes(), nil
}

func (t *LiteServerBlockLink) UnmarshalTL(r io.Reader) error {
	var err error
	var b [4]byte
	_, err = io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	tag := int(binary.LittleEndian.Uint32(b[:]))
	switch tag {
	case 0xef1b7eef:
		t.SumType = "LiteServerBlockLinkBack"
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkBack.ToKeyBlock)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkBack.From)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkBack.To)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkBack.DestProof)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkBack.Proof)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkBack.StateProof)
		if err != nil {
			return err
		}
	case 0x1cce0f52:
		t.SumType = "LiteServerBlockLinkForward"
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkForward.ToKeyBlock)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkForward.From)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkForward.To)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkForward.DestProof)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkForward.ConfigProof)
		if err != nil {
			return err
		}
		err = tl.Unmarshal(r, &t.LiteServerBlockLinkForward.Signatures)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid tag")
	}
	return nil
}

type LiteServerConfigInfo struct {
	Mode        int32
	Id          TonNodeBlockIdExt
	StateProof  []byte
	ConfigProof []byte
}

func (t LiteServerConfigInfo) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.StateProof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ConfigProof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerConfigInfo) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Mode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.StateProof)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ConfigProof)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerLibraryEntry struct {
	Hash tl.Int256
	Data []byte
}

func (t LiteServerLibraryEntry) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Hash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Data)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerLibraryEntry) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Hash)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Data)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerBlockState struct {
	Id       TonNodeBlockIdExt
	RootHash tl.Int256
	FileHash tl.Int256
	Data     []byte
}

func (t LiteServerBlockState) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RootHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.FileHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Data)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerBlockState) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RootHash)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.FileHash)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Data)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerRunMethodResult struct {
	Mode       int32
	Id         TonNodeBlockIdExt
	Shardblk   TonNodeBlockIdExt
	ShardProof []byte
	Proof      []byte
	StateProof []byte
	InitC7     []byte
	LibExtras  []byte
	ExitCode   int32
	Result     []byte
}

func (t LiteServerRunMethodResult) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Shardblk)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (t.Mode>>0)&1 == 1 {
		b, err = tl.Marshal(t.ShardProof)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Mode>>0)&1 == 1 {
		b, err = tl.Marshal(t.Proof)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Mode>>1)&1 == 1 {
		b, err = tl.Marshal(t.StateProof)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Mode>>3)&1 == 1 {
		b, err = tl.Marshal(t.InitC7)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Mode>>4)&1 == 1 {
		b, err = tl.Marshal(t.LibExtras)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	b, err = tl.Marshal(t.ExitCode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (t.Mode>>2)&1 == 1 {
		b, err = tl.Marshal(t.Result)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (t *LiteServerRunMethodResult) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Mode)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Shardblk)
	if err != nil {
		return err
	}
	if (t.Mode>>0)&1 == 1 {
		var tempShardProof []byte
		err = tl.Unmarshal(r, &tempShardProof)
		if err != nil {
			return err
		}
		t.ShardProof = tempShardProof
	}
	if (t.Mode>>0)&1 == 1 {
		var tempProof []byte
		err = tl.Unmarshal(r, &tempProof)
		if err != nil {
			return err
		}
		t.Proof = tempProof
	}
	if (t.Mode>>1)&1 == 1 {
		var tempStateProof []byte
		err = tl.Unmarshal(r, &tempStateProof)
		if err != nil {
			return err
		}
		t.StateProof = tempStateProof
	}
	if (t.Mode>>3)&1 == 1 {
		var tempInitC7 []byte
		err = tl.Unmarshal(r, &tempInitC7)
		if err != nil {
			return err
		}
		t.InitC7 = tempInitC7
	}
	if (t.Mode>>4)&1 == 1 {
		var tempLibExtras []byte
		err = tl.Unmarshal(r, &tempLibExtras)
		if err != nil {
			return err
		}
		t.LibExtras = tempLibExtras
	}
	err = tl.Unmarshal(r, &t.ExitCode)
	if err != nil {
		return err
	}
	if (t.Mode>>2)&1 == 1 {
		var tempResult []byte
		err = tl.Unmarshal(r, &tempResult)
		if err != nil {
			return err
		}
		t.Result = tempResult
	}
	return nil
}

type LiteServerPartialBlockProof struct {
	Complete bool
	From     TonNodeBlockIdExt
	To       TonNodeBlockIdExt
	Steps    []LiteServerBlockLink
}

func (t LiteServerPartialBlockProof) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Complete)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.From)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.To)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Steps)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerPartialBlockProof) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Complete)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.From)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.To)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Steps)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerLibraryResult struct {
	Result []LiteServerLibraryEntry
}

func (t LiteServerLibraryResult) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Result)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerLibraryResult) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Result)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerDebugVerbosity struct {
	Value int32
}

func (t LiteServerDebugVerbosity) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Value)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerDebugVerbosity) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Value)
	if err != nil {
		return err
	}
	return nil
}

type TonNodeBlockId struct {
	Workchain int32
	Shard     int64
	Seqno     int32
}

func (t TonNodeBlockId) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Workchain)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Shard)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Seqno)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *TonNodeBlockId) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Workchain)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Shard)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Seqno)
	if err != nil {
		return err
	}
	return nil
}

type TonNodeZeroStateIdExt struct {
	Workchain int32
	RootHash  tl.Int256
	FileHash  tl.Int256
}

func (t TonNodeZeroStateIdExt) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Workchain)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.RootHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.FileHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *TonNodeZeroStateIdExt) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Workchain)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.RootHash)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.FileHash)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerTransactionId struct {
	Mode    int32
	Account *tl.Int256
	Lt      *int64
	Hash    *tl.Int256
}

func (t LiteServerTransactionId) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (t.Mode>>0)&1 == 1 {
		b, err = tl.Marshal(t.Account)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Mode>>1)&1 == 1 {
		b, err = tl.Marshal(t.Lt)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Mode>>2)&1 == 1 {
		b, err = tl.Marshal(t.Hash)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (t *LiteServerTransactionId) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Mode)
	if err != nil {
		return err
	}
	if (t.Mode>>0)&1 == 1 {
		var tempAccount tl.Int256
		err = tl.Unmarshal(r, &tempAccount)
		if err != nil {
			return err
		}
		t.Account = &tempAccount
	}
	if (t.Mode>>1)&1 == 1 {
		var tempLt int64
		err = tl.Unmarshal(r, &tempLt)
		if err != nil {
			return err
		}
		t.Lt = &tempLt
	}
	if (t.Mode>>2)&1 == 1 {
		var tempHash tl.Int256
		err = tl.Unmarshal(r, &tempHash)
		if err != nil {
			return err
		}
		t.Hash = &tempHash
	}
	return nil
}

type LiteServerShardBlockLink struct {
	Id    TonNodeBlockIdExt
	Proof []byte
}

func (t LiteServerShardBlockLink) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Proof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerShardBlockLink) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Proof)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerShardBlockProof struct {
	MasterchainId TonNodeBlockIdExt
	Links         []LiteServerShardBlockLink
}

func (t LiteServerShardBlockProof) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.MasterchainId)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Links)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerShardBlockProof) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.MasterchainId)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Links)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerMasterchainInfo struct {
	Last          TonNodeBlockIdExt
	StateRootHash tl.Int256
	Init          TonNodeZeroStateIdExt
}

func (t LiteServerMasterchainInfo) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Last)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.StateRootHash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Init)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerMasterchainInfo) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Last)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.StateRootHash)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Init)
	if err != nil {
		return err
	}
	return nil
}

type LiteServerAccountState struct {
	Id         TonNodeBlockIdExt
	Shardblk   TonNodeBlockIdExt
	ShardProof []byte
	Proof      []byte
	State      []byte
}

func (t LiteServerAccountState) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Shardblk)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ShardProof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Proof)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.State)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *LiteServerAccountState) UnmarshalTL(r io.Reader) error {
	var err error
	err = tl.Unmarshal(r, &t.Id)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Shardblk)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.ShardProof)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.Proof)
	if err != nil {
		return err
	}
	err = tl.Unmarshal(r, &t.State)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) LiteServerGetMasterchainInfo(ctx context.Context) (LiteServerMasterchainInfo, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, 0x2ee6b589)
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerMasterchainInfo{}, err
	}
	if len(resp) < 4 {
		return LiteServerMasterchainInfo{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerMasterchainInfo{}, err
		}
		return LiteServerMasterchainInfo{}, errRes
	}
	if tag == 0x81288385 {
		var res LiteServerMasterchainInfo
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerMasterchainInfo{}, fmt.Errorf("invalid tag")
}

type LiteServerGetMasterchainInfoExtRequest struct {
	Mode int32
}

func (t LiteServerGetMasterchainInfoExtRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetMasterchainInfoExt(ctx context.Context, request LiteServerGetMasterchainInfoExtRequest) (LiteServerMasterchainInfoExt, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetMasterchainInfoExtRequest `tlSumType:"df71a670"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerMasterchainInfoExt{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerMasterchainInfoExt{}, err
	}
	if len(resp) < 4 {
		return LiteServerMasterchainInfoExt{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerMasterchainInfoExt{}, err
		}
		return LiteServerMasterchainInfoExt{}, errRes
	}
	if tag == 0xf5e0cca8 {
		var res LiteServerMasterchainInfoExt
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerMasterchainInfoExt{}, fmt.Errorf("invalid tag")
}

func (c *Client) LiteServerGetTime(ctx context.Context) (LiteServerCurrentTime, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, 0x345aad16)
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerCurrentTime{}, err
	}
	if len(resp) < 4 {
		return LiteServerCurrentTime{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerCurrentTime{}, err
		}
		return LiteServerCurrentTime{}, errRes
	}
	if tag == 0xd0053e9 {
		var res LiteServerCurrentTime
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerCurrentTime{}, fmt.Errorf("invalid tag")
}

func (c *Client) LiteServerGetVersion(ctx context.Context) (LiteServerVersion, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, 0xb942b23)
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerVersion{}, err
	}
	if len(resp) < 4 {
		return LiteServerVersion{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerVersion{}, err
		}
		return LiteServerVersion{}, errRes
	}
	if tag == 0xe591045a {
		var res LiteServerVersion
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerVersion{}, fmt.Errorf("invalid tag")
}

type LiteServerGetBlockRequest struct {
	Id TonNodeBlockIdExt
}

func (t LiteServerGetBlockRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetBlock(ctx context.Context, request LiteServerGetBlockRequest) (LiteServerBlockData, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetBlockRequest `tlSumType:"0dcf7763"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerBlockData{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerBlockData{}, err
	}
	if len(resp) < 4 {
		return LiteServerBlockData{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerBlockData{}, err
		}
		return LiteServerBlockData{}, errRes
	}
	if tag == 0x6ced74a5 {
		var res LiteServerBlockData
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerBlockData{}, fmt.Errorf("invalid tag")
}

type LiteServerGetStateRequest struct {
	Id TonNodeBlockIdExt
}

func (t LiteServerGetStateRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetState(ctx context.Context, request LiteServerGetStateRequest) (LiteServerBlockState, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetStateRequest `tlSumType:"b62e6eba"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerBlockState{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerBlockState{}, err
	}
	if len(resp) < 4 {
		return LiteServerBlockState{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerBlockState{}, err
		}
		return LiteServerBlockState{}, errRes
	}
	if tag == 0xcdcadab {
		var res LiteServerBlockState
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerBlockState{}, fmt.Errorf("invalid tag")
}

type LiteServerGetBlockHeaderRequest struct {
	Id   TonNodeBlockIdExt
	Mode int32
}

func (t LiteServerGetBlockHeaderRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetBlockHeader(ctx context.Context, request LiteServerGetBlockHeaderRequest) (LiteServerBlockHeader, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetBlockHeaderRequest `tlSumType:"9e06ec21"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerBlockHeader{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerBlockHeader{}, err
	}
	if len(resp) < 4 {
		return LiteServerBlockHeader{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerBlockHeader{}, err
		}
		return LiteServerBlockHeader{}, errRes
	}
	if tag == 0x19822d75 {
		var res LiteServerBlockHeader
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerBlockHeader{}, fmt.Errorf("invalid tag")
}

type LiteServerSendMessageRequest struct {
	Body []byte
}

func (t LiteServerSendMessageRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Body)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerSendMessage(ctx context.Context, request LiteServerSendMessageRequest) (LiteServerSendMsgStatus, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerSendMessageRequest `tlSumType:"82d40a69"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerSendMsgStatus{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerSendMsgStatus{}, err
	}
	if len(resp) < 4 {
		return LiteServerSendMsgStatus{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerSendMsgStatus{}, err
		}
		return LiteServerSendMsgStatus{}, errRes
	}
	if tag == 0x97e55039 {
		var res LiteServerSendMsgStatus
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerSendMsgStatus{}, fmt.Errorf("invalid tag")
}

type LiteServerGetAccountStateRequest struct {
	Id      TonNodeBlockIdExt
	Account LiteServerAccountId
}

func (t LiteServerGetAccountStateRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Account)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetAccountState(ctx context.Context, request LiteServerGetAccountStateRequest) (LiteServerAccountState, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetAccountStateRequest `tlSumType:"250e896b"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerAccountState{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerAccountState{}, err
	}
	if len(resp) < 4 {
		return LiteServerAccountState{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerAccountState{}, err
		}
		return LiteServerAccountState{}, errRes
	}
	if tag == 0x51c77970 {
		var res LiteServerAccountState
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerAccountState{}, fmt.Errorf("invalid tag")
}

type LiteServerRunSmcMethodRequest struct {
	Mode     int32
	Id       TonNodeBlockIdExt
	Account  LiteServerAccountId
	MethodId int64
	Params   []byte
}

func (t LiteServerRunSmcMethodRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Account)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.MethodId)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Params)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerRunSmcMethod(ctx context.Context, request LiteServerRunSmcMethodRequest) (LiteServerRunMethodResult, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerRunSmcMethodRequest `tlSumType:"d25dc65c"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerRunMethodResult{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerRunMethodResult{}, err
	}
	if len(resp) < 4 {
		return LiteServerRunMethodResult{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerRunMethodResult{}, err
		}
		return LiteServerRunMethodResult{}, errRes
	}
	if tag == 0x6b619aa3 {
		var res LiteServerRunMethodResult
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerRunMethodResult{}, fmt.Errorf("invalid tag")
}

type LiteServerGetShardInfoRequest struct {
	Id        TonNodeBlockIdExt
	Workchain int32
	Shard     int64
	Exact     bool
}

func (t LiteServerGetShardInfoRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Workchain)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Shard)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Exact)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetShardInfo(ctx context.Context, request LiteServerGetShardInfoRequest) (LiteServerShardInfo, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetShardInfoRequest `tlSumType:"25f4a246"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerShardInfo{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerShardInfo{}, err
	}
	if len(resp) < 4 {
		return LiteServerShardInfo{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerShardInfo{}, err
		}
		return LiteServerShardInfo{}, errRes
	}
	if tag == 0x84cde69f {
		var res LiteServerShardInfo
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerShardInfo{}, fmt.Errorf("invalid tag")
}

type LiteServerGetAllShardsInfoRequest struct {
	Id TonNodeBlockIdExt
}

func (t LiteServerGetAllShardsInfoRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetAllShardsInfo(ctx context.Context, request LiteServerGetAllShardsInfoRequest) (LiteServerAllShardsInfo, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetAllShardsInfoRequest `tlSumType:"6bfdd374"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerAllShardsInfo{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerAllShardsInfo{}, err
	}
	if len(resp) < 4 {
		return LiteServerAllShardsInfo{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerAllShardsInfo{}, err
		}
		return LiteServerAllShardsInfo{}, errRes
	}
	if tag == 0x2de78f09 {
		var res LiteServerAllShardsInfo
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerAllShardsInfo{}, fmt.Errorf("invalid tag")
}

type LiteServerGetOneTransactionRequest struct {
	Id      TonNodeBlockIdExt
	Account LiteServerAccountId
	Lt      int64
}

func (t LiteServerGetOneTransactionRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Account)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Lt)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetOneTransaction(ctx context.Context, request LiteServerGetOneTransactionRequest) (LiteServerTransactionInfo, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetOneTransactionRequest `tlSumType:"ea240fd4"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerTransactionInfo{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerTransactionInfo{}, err
	}
	if len(resp) < 4 {
		return LiteServerTransactionInfo{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerTransactionInfo{}, err
		}
		return LiteServerTransactionInfo{}, errRes
	}
	if tag == 0x47edde0e {
		var res LiteServerTransactionInfo
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerTransactionInfo{}, fmt.Errorf("invalid tag")
}

type LiteServerGetTransactionsRequest struct {
	Count   int32
	Account LiteServerAccountId
	Lt      int64
	Hash    tl.Int256
}

func (t LiteServerGetTransactionsRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Count)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Account)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Lt)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Hash)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetTransactions(ctx context.Context, request LiteServerGetTransactionsRequest) (LiteServerTransactionList, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetTransactionsRequest `tlSumType:"a1e7401c"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerTransactionList{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerTransactionList{}, err
	}
	if len(resp) < 4 {
		return LiteServerTransactionList{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerTransactionList{}, err
		}
		return LiteServerTransactionList{}, errRes
	}
	if tag == 0xbc6266f {
		var res LiteServerTransactionList
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerTransactionList{}, fmt.Errorf("invalid tag")
}

type LiteServerLookupBlockRequest struct {
	Mode  int32
	Id    TonNodeBlockId
	Lt    *int64
	Utime *int32
}

func (t LiteServerLookupBlockRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (t.Mode>>1)&1 == 1 {
		b, err = tl.Marshal(t.Lt)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Mode>>2)&1 == 1 {
		b, err = tl.Marshal(t.Utime)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerLookupBlock(ctx context.Context, request LiteServerLookupBlockRequest) (LiteServerBlockHeader, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerLookupBlockRequest `tlSumType:"1ef7c8fa"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerBlockHeader{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerBlockHeader{}, err
	}
	if len(resp) < 4 {
		return LiteServerBlockHeader{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerBlockHeader{}, err
		}
		return LiteServerBlockHeader{}, errRes
	}
	if tag == 0x19822d75 {
		var res LiteServerBlockHeader
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerBlockHeader{}, fmt.Errorf("invalid tag")
}

type LiteServerListBlockTransactionsRequest struct {
	Id    TonNodeBlockIdExt
	Mode  int32
	Count int32
	After *LiteServerTransactionId3
}

func (t LiteServerListBlockTransactionsRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Count)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (t.Mode>>7)&1 == 1 {
		b, err = tl.Marshal(t.After)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerListBlockTransactions(ctx context.Context, request LiteServerListBlockTransactionsRequest) (LiteServerBlockTransactions, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerListBlockTransactionsRequest `tlSumType:"dac7fcad"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerBlockTransactions{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerBlockTransactions{}, err
	}
	if len(resp) < 4 {
		return LiteServerBlockTransactions{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerBlockTransactions{}, err
		}
		return LiteServerBlockTransactions{}, errRes
	}
	if tag == 0x5c6c542f {
		var res LiteServerBlockTransactions
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerBlockTransactions{}, fmt.Errorf("invalid tag")
}

type LiteServerGetBlockProofRequest struct {
	Mode        int32
	KnownBlock  TonNodeBlockIdExt
	TargetBlock *TonNodeBlockIdExt
}

func (t LiteServerGetBlockProofRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.KnownBlock)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (t.Mode>>0)&1 == 1 {
		b, err = tl.Marshal(t.TargetBlock)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetBlockProof(ctx context.Context, request LiteServerGetBlockProofRequest) (LiteServerPartialBlockProof, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetBlockProofRequest `tlSumType:"449cea8a"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerPartialBlockProof{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerPartialBlockProof{}, err
	}
	if len(resp) < 4 {
		return LiteServerPartialBlockProof{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerPartialBlockProof{}, err
		}
		return LiteServerPartialBlockProof{}, errRes
	}
	if tag == 0xc1d2d08e {
		var res LiteServerPartialBlockProof
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerPartialBlockProof{}, fmt.Errorf("invalid tag")
}

type LiteServerGetConfigAllRequest struct {
	Mode int32
	Id   TonNodeBlockIdExt
}

func (t LiteServerGetConfigAllRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetConfigAll(ctx context.Context, request LiteServerGetConfigAllRequest) (LiteServerConfigInfo, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetConfigAllRequest `tlSumType:"b7261b91"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerConfigInfo{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerConfigInfo{}, err
	}
	if len(resp) < 4 {
		return LiteServerConfigInfo{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerConfigInfo{}, err
		}
		return LiteServerConfigInfo{}, errRes
	}
	if tag == 0x2f277bae {
		var res LiteServerConfigInfo
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerConfigInfo{}, fmt.Errorf("invalid tag")
}

type LiteServerGetConfigParamsRequest struct {
	Mode      int32
	Id        TonNodeBlockIdExt
	ParamList []int32
}

func (t LiteServerGetConfigParamsRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.ParamList)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetConfigParams(ctx context.Context, request LiteServerGetConfigParamsRequest) (LiteServerConfigInfo, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetConfigParamsRequest `tlSumType:"638df89e"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerConfigInfo{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerConfigInfo{}, err
	}
	if len(resp) < 4 {
		return LiteServerConfigInfo{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerConfigInfo{}, err
		}
		return LiteServerConfigInfo{}, errRes
	}
	if tag == 0x2f277bae {
		var res LiteServerConfigInfo
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerConfigInfo{}, fmt.Errorf("invalid tag")
}

type LiteServerGetValidatorStatsRequest struct {
	Mode          int32
	Id            TonNodeBlockIdExt
	Limit         int32
	StartAfter    *tl.Int256
	ModifiedAfter *int32
}

func (t LiteServerGetValidatorStatsRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Mode)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	b, err = tl.Marshal(t.Limit)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	if (t.Mode>>0)&1 == 1 {
		b, err = tl.Marshal(t.StartAfter)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	if (t.Mode>>2)&1 == 1 {
		b, err = tl.Marshal(t.ModifiedAfter)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetValidatorStats(ctx context.Context, request LiteServerGetValidatorStatsRequest) (LiteServerValidatorStats, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetValidatorStatsRequest `tlSumType:"091a58bc"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerValidatorStats{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerValidatorStats{}, err
	}
	if len(resp) < 4 {
		return LiteServerValidatorStats{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerValidatorStats{}, err
		}
		return LiteServerValidatorStats{}, errRes
	}
	if tag == 0xd896f7b9 {
		var res LiteServerValidatorStats
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerValidatorStats{}, fmt.Errorf("invalid tag")
}

type LiteServerGetLibrariesRequest struct {
	LibraryList []tl.Int256
}

func (t LiteServerGetLibrariesRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.LibraryList)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetLibraries(ctx context.Context, request LiteServerGetLibrariesRequest) (LiteServerLibraryResult, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetLibrariesRequest `tlSumType:"99181e7e"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerLibraryResult{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerLibraryResult{}, err
	}
	if len(resp) < 4 {
		return LiteServerLibraryResult{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerLibraryResult{}, err
		}
		return LiteServerLibraryResult{}, errRes
	}
	if tag == 0xc43848b {
		var res LiteServerLibraryResult
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerLibraryResult{}, fmt.Errorf("invalid tag")
}

type LiteServerGetShardBlockProofRequest struct {
	Id TonNodeBlockIdExt
}

func (t LiteServerGetShardBlockProofRequest) MarshalTL() ([]byte, error) {
	var (
		err error
		b   []byte
	)
	buf := new(bytes.Buffer)
	b, err = tl.Marshal(t.Id)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) LiteServerGetShardBlockProof(ctx context.Context, request LiteServerGetShardBlockProofRequest) (LiteServerShardBlockProof, error) {
	payload, err := tl.Marshal(struct {
		tl.SumType
		Req LiteServerGetShardBlockProofRequest `tlSumType:"5003a64c"`
	}{SumType: "Req", Req: request})
	if err != nil {
		return LiteServerShardBlockProof{}, err
	}
	resp, err := c.liteServerRequest(ctx, payload)
	if err != nil {
		return LiteServerShardBlockProof{}, err
	}
	if len(resp) < 4 {
		return LiteServerShardBlockProof{}, fmt.Errorf("not enought bytes for tag")
	}
	tag := binary.BigEndian.Uint32(resp[:4])
	if tag == 0x48e1a9bb {
		var errRes LiteServerError
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &errRes)
		if err != nil {
			return LiteServerShardBlockProof{}, err
		}
		return LiteServerShardBlockProof{}, errRes
	}
	if tag == 0x70347608 {
		var res LiteServerShardBlockProof
		err = tl.Unmarshal(bytes.NewReader(resp[4:]), &res)
		return res, err
	}
	return LiteServerShardBlockProof{}, fmt.Errorf("invalid tag")
}
