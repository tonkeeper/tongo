package liteclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/adnl"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/config"
	"github.com/startfellows/tongo/tl"
	"github.com/startfellows/tongo/tlb"
	"net/http"
)

type Client struct {
	adnlClient *adnl.Client
}

// NewClient
// Get options and create new lite client. If no options provided - download public config for mainnet from ton.org.
func NewClient(options *config.Options) (*Client, error) {
	// TODO: implement multiple server support
	if options == nil {
		var err error
		options, err = downloadConfig("https://ton.org/global-config.json")
		if err != nil {
			return nil, err
		}
	}
	if len(options.LiteServers) == 0 {
		return nil, fmt.Errorf("server list empty")
	}
	serverPubkey, err := base64.StdEncoding.DecodeString(options.LiteServers[0].Key)
	if err != nil {
		return nil, err
	}
	c, err := adnl.NewConnection(context.Background(), serverPubkey, options.LiteServers[0].Host)
	if err != nil {
		return nil, err
	}
	adnlClient := adnl.NewClient(c)
	return &Client{
		adnlClient: adnlClient,
	}, nil
}

func (c *Client) GetLastRawAccountState(ctx context.Context, accountId tongo.AccountID) (AccountState, error) {
	req := makeLiteServerQueryRequest(makeLiteServerGetMasterchainInfoRequest())
	resp, err := c.adnlClient.Request(ctx, req)
	if err != nil {
		return AccountState{}, err
	}
	parsedResp, err := parseLiteServerQueryResponse(resp)
	if err != nil {
		return AccountState{}, err
	}
	if parsedResp.Tag == LiteServerErrorTag {
		return AccountState{}, fmt.Errorf("lite server error: %v %v", parsedResp.LiteServerError.Code, parsedResp.LiteServerError.Message)
	}
	if parsedResp.Tag != LiteServerMasterchainInfoTag {
		return AccountState{}, fmt.Errorf("masterchain info not recieved")
	}
	asReq, err := makeLiteServerGetAccountStateRequest(parsedResp.LiteServerMasterchainInfo.Last, accountId)
	if err != nil {
		return AccountState{}, err
	}
	req = makeLiteServerQueryRequest(asReq)
	resp, err = c.adnlClient.Request(ctx, req)
	if err != nil {
		return AccountState{}, err
	}
	parsedResp, err = parseLiteServerQueryResponse(resp)
	if err != nil {
		return AccountState{}, err
	}
	if parsedResp.Tag == LiteServerErrorTag {
		return AccountState{}, fmt.Errorf("lite server error: %v %v", parsedResp.LiteServerError.Code, parsedResp.LiteServerError.Message)
	}
	if parsedResp.Tag != LiteServerAccountStateTag {
		return AccountState{}, fmt.Errorf("account state not recieved")
	}
	state, err := decodeRawAccountStateBoc(parsedResp.LiteServerAccountState.State)
	if err != nil {
		return AccountState{}, err
	}
	return state, nil
}

type AccountState struct {
	Status            tongo.AccountStatus
	Balance           uint64
	Data              []byte
	Code              []byte
	FrozenHash        tongo.Hash
	LastTransactionLt uint64
}

type LiteServerMasterchainInfo struct {
	Last          tongo.TonNodeBlockIdExt
	StateRootHash tongo.Hash
	// TODO: add init
}

type LiteServerAccountState struct {
	Id         tongo.TonNodeBlockIdExt
	ShardBlk   tongo.TonNodeBlockIdExt
	ShardProof []byte
	Proof      []byte
	State      []byte
}

type LiteServerError struct {
	Code    int32
	Message string
}

type ParsedLiteServerQueryResponse struct {
	Tag                       uint32
	LiteServerError           LiteServerError
	LiteServerMasterchainInfo LiteServerMasterchainInfo
	LiteServerAccountState    LiteServerAccountState
}

func ParseLiteServerMasterchainInfo(data []byte) (LiteServerMasterchainInfo, error) {
	if len(data) < 4+80+32+4+32+32 {
		return LiteServerMasterchainInfo{}, fmt.Errorf("invalid data length")
	}
	tag := binary.BigEndian.Uint32(data[:4])
	if tag != LiteServerMasterchainInfoTag {
		return LiteServerMasterchainInfo{}, fmt.Errorf("invalid tag")
	}
	var info LiteServerMasterchainInfo
	var last tongo.TonNodeBlockIdExt
	err := last.UnmarshalTL(data[4:84])
	if err != nil {
		return LiteServerMasterchainInfo{}, err
	}
	info.Last = last
	copy(info.StateRootHash[:], data[84:116])
	// TODO: fill init
	return info, nil
}

func ParseLiteServerAccountState(data []byte) (LiteServerAccountState, error) {
	if len(data) < 164 {
		return LiteServerAccountState{}, fmt.Errorf("invalid data length")
	}
	tag := binary.BigEndian.Uint32(data[:4])
	if tag != LiteServerAccountStateTag {
		return LiteServerAccountState{}, fmt.Errorf("invalid tag")
	}
	var state LiteServerAccountState
	var id, shardBlk tongo.TonNodeBlockIdExt
	err := id.UnmarshalTL(data[4:84])
	if err != nil {
		return LiteServerAccountState{}, err
	}
	state.Id = id
	err = shardBlk.UnmarshalTL(data[84:164])
	if err != nil {
		return LiteServerAccountState{}, err
	}
	state.ShardBlk = shardBlk
	data = data[164:]
	bytes, data, err := parseBytes(data)
	if err != nil {
		return LiteServerAccountState{}, err
	}
	state.ShardProof = append(state.ShardProof, bytes...)
	if len(data) == 0 {
		return LiteServerAccountState{}, fmt.Errorf("invalid length")
	}
	bytes, data, err = parseBytes(data)
	if err != nil {
		return LiteServerAccountState{}, err
	}
	state.Proof = append(state.Proof, bytes...)
	if len(data) == 0 {
		return LiteServerAccountState{}, fmt.Errorf("invalid length")
	}
	bytes, data, err = parseBytes(data)
	if err != nil {
		return LiteServerAccountState{}, err
	}
	state.State = append(state.State, bytes...)
	if len(data) != 0 {
		return LiteServerAccountState{}, fmt.Errorf("invalid length")
	}
	return state, nil
}

func parseLiteServerError(data []byte) (LiteServerError, error) {
	if len(data) < 8 {
		return LiteServerError{}, fmt.Errorf("invalid data length")
	}
	tag := binary.BigEndian.Uint32(data[:4])
	if tag != LiteServerErrorTag {
		return LiteServerError{}, fmt.Errorf("invalid tag")
	}
	code := binary.LittleEndian.Uint32(data[4:8])
	var bytes []byte
	if len(data) > 8 {
		var err error
		bytes, _, err = parseBytes(data[8:])
		if err != nil {
			return LiteServerError{}, err
		}
	}
	return LiteServerError{Code: int32(code), Message: string(bytes)}, nil
}

func makeLiteServerQueryRequest(payload []byte) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, LiteServerQueryTag)
	data = append(data, tl.EncodeLength(len(payload))...)
	data = append(data, payload...)
	left := len(data) % 4
	if left != 0 {
		data = append(data, make([]byte, 4-left)...)
	}
	return data
}

func decodeLength(b []byte) (int, []byte, error) {
	// TODO: import from ADNL
	if len(b) == 0 {
		return 0, nil, fmt.Errorf("size should contains at least one byte")
	}
	if b[0] == 255 {
		return 0, nil, fmt.Errorf("invalid first byte value %x", b[0])
	}
	if b[0] < 254 {
		return int(b[0]), b[1:], nil
	}
	if b[0] != 254 {
		panic("how it cat be possible? you are fucking wizard!")
	}
	if len(b) < 4 {
		return 0, nil, fmt.Errorf("not enought bytes for decoding size")
	}
	b[0] = 0
	i := binary.LittleEndian.Uint32(b[:4])
	b[0] = 254
	return int(i) >> 8, b[4:], nil
}

func parseBytes(source []byte) (read []byte, remaining []byte, err error) {
	ln, buffer, err := decodeLength(source)
	if err != nil {
		return nil, nil, err
	}
	if len(buffer) < ln {
		return nil, nil, fmt.Errorf("invalid length")
	}
	left := (len(source) - len(buffer) + ln) % 4
	index := ln
	if left != 0 {
		index = ln + 4 - left
	}
	return buffer[:ln], buffer[index:], nil
}

func makeLiteServerGetMasterchainInfoRequest() []byte {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, LiteServerGetMasterchainInfoTag)
	return payload
}

func makeLiteServerGetAccountStateRequest(blockIdExt tongo.TonNodeBlockIdExt, accountId tongo.AccountID) ([]byte, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, LiteServerGetAccountStateTag)
	block, err := blockIdExt.MarshalTL()
	if err != nil {
		return nil, err
	}
	payload = append(payload, block...)
	a, err := accountId.MarshalTL()
	if err != nil {
		return nil, err
	}
	payload = append(payload, a...)
	return payload, nil
}

func parseLiteServerQueryResponse(message adnl.Message) (ParsedLiteServerQueryResponse, error) {
	var response ParsedLiteServerQueryResponse
	if len(message) < 4 {
		return ParsedLiteServerQueryResponse{}, fmt.Errorf("invalid lenght")
	}
	tag := binary.BigEndian.Uint32(message[:4])
	switch tag {
	case LiteServerErrorTag:
		serverError, err := parseLiteServerError(message[:])
		if err != nil {
			return ParsedLiteServerQueryResponse{}, err
		}
		response.LiteServerError = serverError
		response.Tag = LiteServerErrorTag
	case LiteServerMasterchainInfoTag:
		info, err := ParseLiteServerMasterchainInfo(message[:])
		if err != nil {
			return ParsedLiteServerQueryResponse{}, err
		}
		response.LiteServerMasterchainInfo = info
		response.Tag = LiteServerMasterchainInfoTag
	case LiteServerAccountStateTag:
		res, err := ParseLiteServerAccountState(message[:])
		if err != nil {
			return ParsedLiteServerQueryResponse{}, err
		}
		response.LiteServerAccountState = res
		response.Tag = LiteServerAccountStateTag
	}
	return response, nil
}

func decodeRawAccountStateBoc(bocBytes []byte) (AccountState, error) {
	if bocBytes == nil {
		return AccountState{Status: tongo.AccountEmpty}, nil
	}
	cells, err := boc.DeserializeBoc(bocBytes)
	if err != nil {
		return AccountState{}, err
	}
	if len(cells) != 1 {
		return AccountState{}, fmt.Errorf("must be one root cell")
	}
	var acc tongo.Account
	err = tlb.Unmarshal(cells[0], &acc)
	if err != nil {
		return AccountState{}, err
	}
	return convertTlbAccountToAccountState(acc)
}

func convertTlbAccountToAccountState(acc tongo.Account) (AccountState, error) {
	if acc.SumType == "AccountNone" {
		return AccountState{Status: tongo.AccountNone}, nil
	}
	res := AccountState{
		Balance:           uint64(acc.Account.Storage.Balance.Grams),
		LastTransactionLt: acc.Account.Storage.LastTransLt,
	}
	if acc.Account.Storage.State.SumType == "AccountUninit" {
		res.Status = tongo.AccountUninit
		return res, nil
	}
	if acc.Account.Storage.State.SumType == "AccountFrozen" {
		res.FrozenHash = acc.Account.Storage.State.AccountFrozen.StateHash
		res.Status = tongo.AccountFrozen
		return res, nil
	}
	res.Status = tongo.AccountActive
	if !acc.Account.Storage.State.AccountActive.StateInit.Data.Null {
		data, err := acc.Account.Storage.State.AccountActive.StateInit.Data.Value.Value.ToBoc()
		if err != nil {
			return AccountState{}, err
		}
		res.Data = data
	}
	if !acc.Account.Storage.State.AccountActive.StateInit.Code.Null {
		code, err := acc.Account.Storage.State.AccountActive.StateInit.Code.Value.Value.ToBoc()
		if err != nil {
			return AccountState{}, err
		}
		res.Code = code
	}
	return res, nil
}

func (c *Client) GetTransactions(ctx context.Context, count uint32, accountId tongo.AccountID, lt uint64, hash tongo.Hash) ([]tongo.Transaction, error) {
	cells, err := c.GetRawTransactions(ctx, count, accountId, lt, hash)
	if err != nil {
		return nil, err
	}
	var res []tongo.Transaction
	for _, cell := range cells {
		var t tongo.Transaction
		cell.ResetCounters()
		err := tlb.Unmarshal(cell, &t)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

// GetRawTransactions
// Returns []boc.Cell of the transaction$0111 tlb constructor. Be careful when reading Cell. Some Cells are shared between slice elements. Use cell.ResetCounters()
func (c *Client) GetRawTransactions(ctx context.Context, count uint32, accountId tongo.AccountID, lt uint64, hash tongo.Hash) ([]*boc.Cell, error) {
	// TransactionList
	// liteServer.transactionList ids:(vector tonNode.blockIdExt) transactions:bytes = liteServer.TransactionList;
	type TransactionList struct {
		Ids          []tongo.TonNodeBlockIdExt
		Transactions []byte
	}
	type getTransactionsRequest struct {
		Count   uint32
		Account tongo.AccountID
		Lt      uint64
		Hash    tongo.Hash
	}
	r := struct {
		tl.SumType
		GetTransactionsRequest getTransactionsRequest `tlSumType:"a1e7401c"`
	}{
		SumType: "GetTransactionsRequest",
		GetTransactionsRequest: getTransactionsRequest{
			Count:   count,
			Account: accountId,
			Lt:      lt,
			Hash:    hash,
		},
	}
	rBytes, err := tl.Marshal(r)
	if err != nil {
		return nil, err
	}
	req := makeLiteServerQueryRequest(rBytes)
	resp, err := c.adnlClient.Request(ctx, req)
	if err != nil {
		return nil, err
	}
	var pResp struct {
		tlb.SumType
		TransactionList TransactionList `tlSumType:"0bc6266f"` // TODO: must be 9dd72eb9
		Error           LiteServerError `tlSumType:"48e1a9bb"`
	}
	reader := bytes.NewReader(resp)
	err = tl.Unmarshal(reader, &pResp)
	if err != nil {
		return nil, err
	}
	if pResp.SumType == "Error" {
		return nil, fmt.Errorf("error code: %v , message: %v", pResp.Error.Code, pResp.Error.Message)
	}
	cells, err := boc.DeserializeBoc(pResp.TransactionList.Transactions)
	if err != nil {
		return nil, err
	}
	if len(cells) != len(pResp.TransactionList.Ids) {
		return nil, fmt.Errorf("TonNodeBlockIdExt qty not equal transactions qty")
	}
	return cells, nil
}

// SendMessage
// Send binary payload to TON blockchain
// liteServer.sendMessage body:bytes = liteServer.SendMsgStatus;
// liteServer.sendMsgStatus status:int = liteServer.SendMsgStatus;
func (c *Client) SendMessage(ctx context.Context, payload []byte) error {
	request := struct {
		tl.SumType
		SendMessage struct {
			Body []byte
		} `tlSumType:"82d40a69"`
	}{
		SumType:     "SendMessage",
		SendMessage: struct{ Body []byte }{payload},
	}
	rBytes, err := tl.Marshal(request)
	if err != nil {
		return err
	}
	req := makeLiteServerQueryRequest(rBytes)
	resp, err := c.adnlClient.Request(ctx, req)
	if err != nil {
		return err
	}
	var response struct {
		tlb.SumType
		SendMsgStatus struct {
			Status int32
		} `tlSumType:"97e55039"`
		Error LiteServerError `tlSumType:"48e1a9bb"`
	}
	reader := bytes.NewReader(resp)
	err = tl.Unmarshal(reader, &response)
	if err != nil {
		return err
	}
	if response.SumType == "Error" {
		return fmt.Errorf("error code: %v , message: %v", response.Error.Code, response.Error.Message)
	}
	if response.SendMsgStatus.Status != 1 {
		return fmt.Errorf("message sending failed with status: %v", response.SendMsgStatus.Status)
	}
	return nil
}

func downloadConfig(path string) (*config.Options, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return config.ParseConfig(resp.Body)
}
