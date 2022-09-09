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
	"github.com/startfellows/tongo/utils"
	"net/http"
	"sync"
)

type Client struct {
	adnlClient *adnl.Client
}

func NewClientWithDefaultMainnet() (*Client, error) {
	options, err := downloadConfig("https://ton-blockchain.github.io/global.config.json")
	if err != nil {
		return nil, err
	}
	return NewClient(*options)
}
func NewClientWithDefaultTestnet() (*Client, error) {
	options, err := downloadConfig("https://ton-blockchain.github.io/testnet-global.config.json")
	if err != nil {
		return nil, err
	}
	return NewClient(*options)
}

// NewClient
// Get options and create new lite client. If no options provided - download public config for mainnet from ton.org.
func NewClient(options config.Options) (*Client, error) {
	// TODO: implement multiple server support

	if len(options.LiteServers) == 0 {
		return nil, fmt.Errorf("server list empty")
	}
	for _, ls := range options.LiteServers {
		serverPubkey, err := base64.StdEncoding.DecodeString(ls.Key)
		if err != nil {
			continue
		}
		c, err := adnl.NewConnection(context.Background(), serverPubkey, ls.Host)
		if err != nil {
			continue
		}
		adnlClient := adnl.NewClient(c)
		return &Client{
			adnlClient: adnlClient,
		}, nil
	}
	return nil, fmt.Errorf("all liteservers are unavailable")
}

func (c *Client) GetAccountState(ctx context.Context, accountId tongo.AccountID) (tongo.AccountInfo, error) {
	b, err := c.getLastRawAccount(ctx, accountId)
	if err != nil {
		return tongo.AccountInfo{}, err
	}
	if b == nil {
		return tongo.AccountInfo{Status: tongo.AccountEmpty}, nil
	}
	account, err := decodeRawAccountBoc(b)
	if err != nil {
		return tongo.AccountInfo{}, err
	}
	return convertTlbAccountToAccountState(account)
}

func (c *Client) GetLastRawAccount(ctx context.Context, accountId tongo.AccountID) (tongo.Account, error) {
	b, err := c.getLastRawAccount(ctx, accountId)
	if err != nil {
		return tongo.Account{}, err
	}
	if b == nil {
		acc := tongo.Account{SumType: "AccountNone"}
		return acc, nil
	}
	return decodeRawAccountBoc(b)
}

func (c *Client) getLastRawAccount(ctx context.Context, accountId tongo.AccountID) ([]byte, error) {
	masterchainInfo, err := c.GetMasterchainInfo(ctx)
	if err != nil {
		return nil, err
	}
	asReq, err := makeLiteServerGetAccountStateRequest(masterchainInfo, accountId)
	if err != nil {
		return nil, err
	}
	req := makeLiteServerQueryRequest(asReq)
	resp, err := c.adnlClient.Request(ctx, req)
	if err != nil {
		return nil, err
	}
	var parsedResp struct {
		tl.SumType
		LiteServerError        LiteServerError        `tlSumType:"48e1a9bb"`
		LiteServerAccountState LiteServerAccountState `tlSumType:"51c77970"`
	}
	err = tl.Unmarshal(bytes.NewReader(resp), &parsedResp)
	if err != nil {
		return nil, err
	}
	switch parsedResp.SumType {
	case "LiteServerError":
		return nil, fmt.Errorf("lite server error: %v %v", parsedResp.LiteServerError.Code, parsedResp.LiteServerError.Message)
	case "LiteServerAccountState":
	default:
		return nil, fmt.Errorf("account state not recieved")
	}
	return parsedResp.LiteServerAccountState.State, nil
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

type LiteServerAllShardsInfo struct {
	Id    tongo.TonNodeBlockIdExt
	Proof []byte
	Data  []byte
}

type LiteServerError struct {
	Code    int32
	Message string
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

func makeLiteServerGetMasterchainInfoRequest() []byte {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, LiteServerGetMasterchainInfoTag)
	return payload
}

func makeLiteServerAllShardsInfoRequest(blockIdExt tongo.TonNodeBlockIdExt) ([]byte, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, LiteServerGetAllShardsInfoTag)
	block, err := blockIdExt.MarshalTL()
	if err != nil {
		return nil, err
	}
	payload = append(payload, block...)
	return payload, nil
}

func makeLiteServerGetAccountStateRequest(blockIdExt tongo.TonNodeBlockIdExt, accountId tongo.AccountID) ([]byte, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, LiteServerGetAccountStateTag)
	block, err := tl.Marshal(blockIdExt)
	if err != nil {
		return nil, err
	}
	payload = append(payload, block...)
	a, err := tl.Marshal(accountId)
	if err != nil {
		return nil, err
	}
	payload = append(payload, a...)
	return payload, nil
}

func decodeRawAccountBoc(bocBytes []byte) (tongo.Account, error) {
	cells, err := boc.DeserializeBoc(bocBytes)
	if err != nil {
		return tongo.Account{}, err
	}
	if len(cells) != 1 {
		return tongo.Account{}, fmt.Errorf("must be one root cell")
	}
	var acc tongo.Account
	err = tlb.Unmarshal(cells[0], &acc)
	if err != nil {
		return tongo.Account{}, err
	}
	return acc, nil
}

func convertTlbAccountToAccountState(acc tongo.Account) (tongo.AccountInfo, error) {
	if acc.SumType == "AccountNone" {
		return tongo.AccountInfo{Status: tongo.AccountNone}, nil
	}
	res := tongo.AccountInfo{
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
			return tongo.AccountInfo{}, err
		}
		res.Data = data
	}
	if !acc.Account.Storage.State.AccountActive.StateInit.Code.Null {
		code, err := acc.Account.Storage.State.AccountActive.StateInit.Code.Value.Value.ToBoc()
		if err != nil {
			return tongo.AccountInfo{}, err
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
		tl.SumType
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

// SendRawMessage
// Send binary payload to TON blockchain
// liteServer.sendMessage body:bytes = liteServer.SendMsgStatus;
// liteServer.sendMsgStatus status:int = liteServer.SendMsgStatus;
func (c *Client) SendRawMessage(ctx context.Context, payload []byte) error {
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
		tl.SumType
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
	if response.SumType != "SendMsgStatus" {
		return fmt.Errorf("not SendMsgStatus response")
	}
	if response.SendMsgStatus.Status != 1 {
		return fmt.Errorf("message sending failed with status: %v", response.SendMsgStatus.Status)
	}
	return nil
}

var configCache = make(map[string]*config.Options)
var configCacheMutex sync.RWMutex

func downloadConfig(path string) (*config.Options, error) {
	configCacheMutex.RLock()
	o, prs := configCache[path]
	configCacheMutex.RUnlock()
	if prs {
		return o, nil
	}
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	o, err = config.ParseConfig(resp.Body)
	if err == nil {
		configCacheMutex.Lock()
		configCache[path] = o
		configCacheMutex.Unlock()
	}
	return o, err
}

// GetLastConfigAll
// liteServer.getConfigAll mode:# id:tonNode.blockIdExt = liteServer.ConfigInfo;
// liteServer.configInfo mode:# id:tonNode.blockIdExt state_proof:bytes config_proof:bytes = liteServer.ConfigInfo;
// Returns config: (Hashmap 32 ^Cell) as a Cell from config_proof
func (c *Client) GetLastConfigAll(ctx context.Context) (*boc.Cell, error) {
	type getConfigAllRequest struct {
		Mode uint32
		ID   tongo.TonNodeBlockIdExt
	}
	type configInfo struct {
		Mode        uint32
		ID          tongo.TonNodeBlockIdExt
		StateProof  []byte
		ConfigProof []byte
	}

	lastBlock, err := c.GetMasterchainInfo(ctx)
	if err != nil {
		return nil, err
	}

	r := struct {
		tl.SumType
		GetConfigAllRequest getConfigAllRequest `tlSumType:"b7261b91"`
	}{
		SumType: "GetConfigAllRequest",
		GetConfigAllRequest: getConfigAllRequest{
			Mode: 0,
			ID:   lastBlock,
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
		tl.SumType
		ConfigInfo configInfo      `tlSumType:"2f277bae"`
		Error      LiteServerError `tlSumType:"48e1a9bb"`
	}
	reader := bytes.NewReader(resp)
	err = tl.Unmarshal(reader, &pResp)
	if err != nil {
		return nil, err
	}
	if pResp.SumType == "Error" {
		return nil, fmt.Errorf("error code: %v , message: %v", pResp.Error.Code, pResp.Error.Message)
	}

	cell, err := boc.DeserializeBoc(pResp.ConfigInfo.ConfigProof)
	if err != nil {
		return nil, err
	}
	var proof struct {
		Proof tongo.MerkleProof[tongo.ShardStateUnsplit]
	}
	err = tlb.Unmarshal(cell[0], &proof)
	if err != nil {
		return nil, err
	}
	if proof.Proof.SumType != "MerkleProof" ||
		proof.Proof.MerkleProof.VirtualRoot.Value.SumType != "ShardStateUnsplit" ||
		proof.Proof.MerkleProof.VirtualRoot.Value.ShardStateUnsplit.Custom.Null ||
		proof.Proof.MerkleProof.VirtualRoot.Value.ShardStateUnsplit.Custom.Value.Value.SumType != "McStateExtra" {
		return nil, fmt.Errorf("can not extract config")
	}
	conf := boc.Cell(proof.Proof.MerkleProof.VirtualRoot.Value.ShardStateUnsplit.Custom.Value.Value.McStateExtra.Config.Config.Value)
	return &conf, nil
}

// GetMasterchainInfo
// liteServer.getMasterchainInfo = liteServer.MasterchainInfo;
// liteServer.masterchainInfo last:tonNode.blockIdExt state_root_hash:int256
// init:tonNode.zeroStateIdExt = liteServer.MasterchainInfo;
func (c *Client) GetMasterchainInfo(ctx context.Context) (tongo.TonNodeBlockIdExt, error) {
	req := makeLiteServerQueryRequest(makeLiteServerGetMasterchainInfoRequest())
	resp, err := c.adnlClient.Request(ctx, req)
	if err != nil {
		return tongo.TonNodeBlockIdExt{}, err
	}
	var parsedResp struct {
		tl.SumType
		LiteServerError           LiteServerError           `tlSumType:"48e1a9bb"`
		LiteServerMasterchainInfo LiteServerMasterchainInfo `tlSumType:"81288385"`
	}
	err = tl.Unmarshal(bytes.NewReader(resp), &parsedResp)
	if err != nil {
		return tongo.TonNodeBlockIdExt{}, err
	}
	switch parsedResp.SumType {
	case "LiteServerError":
		return tongo.TonNodeBlockIdExt{}, fmt.Errorf("lite server error: %v %v", parsedResp.LiteServerError.Code, parsedResp.LiteServerError.Message)
	case "LiteServerMasterchainInfo":
		return parsedResp.LiteServerMasterchainInfo.Last, nil
	default:
		return tongo.TonNodeBlockIdExt{}, fmt.Errorf("masterchain info not recieved")
	}

}

// RunSmcMethod
// Run smart contract method by name and parameters
// liteServer.runSmcMethod mode:# id:tonNode.blockIdExt account:liteServer.accountId method_id:long params:bytes = liteServer.RunMethodResult;
// liteServer.runMethodResult mode:# id:tonNode.blockIdExt shardblk:tonNode.blockIdExt shard_proof:mode.0?bytes
// proof:mode.0?bytes state_proof:mode.1?bytes init_c7:mode.3?bytes lib_extras:mode.4?bytes exit_code:int result:mode.2?bytes = liteServer.RunMethodResult;
func (c *Client) RunSmcMethod(ctx context.Context, mode uint32, accountId tongo.AccountID, method string, params tongo.VmStack) (tongo.VmStack, error) {
	type runSmcRequest struct {
		Mode     uint32
		Id       tongo.TonNodeBlockIdExt
		Account  tongo.AccountID
		MethodId uint64
		Params   tongo.VmStack
	}
	info, err := c.GetMasterchainInfo(ctx)
	if err != nil {
		return tongo.VmStack{}, err
	}
	r := struct {
		tl.SumType
		RunSmcRequest runSmcRequest `tlSumType:"d25dc65c"`
	}{
		SumType: "RunSmcRequest",
		RunSmcRequest: runSmcRequest{
			Mode:     mode,
			Id:       info,
			Account:  accountId,
			MethodId: uint64(utils.Crc16String(method)&0xffff) | 0x10000,
			Params:   params,
		},
	}
	payload, err := tl.Marshal(r)
	if err != nil {
		return tongo.VmStack{}, err
	}
	fmt.Printf("%x\n", payload)
	req := makeLiteServerQueryRequest(payload)
	resp, err := c.adnlClient.Request(ctx, req)
	if err != nil {
		return tongo.VmStack{}, err
	}
	var response struct {
		tl.SumType
		RunMethodResult struct {
			Mode     uint32
			Id       tongo.TonNodeBlockIdExt
			ShardBlk tongo.TonNodeBlockIdExt
			// TODO: add proofs support
			ExitCode uint32
			Result   tongo.VmStack
		} `tlSumType:"6b619aa3"`
		Error LiteServerError `tlSumType:"48e1a9bb"`
	}
	reader := bytes.NewReader(resp)
	err = tl.Unmarshal(reader, &response)
	if err != nil {
		return tongo.VmStack{}, err
	}
	if response.SumType == "Error" {
		return tongo.VmStack{}, fmt.Errorf("error code: %v , message: %v", response.Error.Code, response.Error.Message)
	}
	if response.SumType != "RunMethodResult" {
		return tongo.VmStack{}, fmt.Errorf("not RunMethodResult response")
	}
	if response.RunMethodResult.ExitCode != 0 && response.RunMethodResult.ExitCode != 1 {
		return tongo.VmStack{}, fmt.Errorf("method execution failed with code: %v", response.RunMethodResult.ExitCode)
	}
	return response.RunMethodResult.Result, nil
}

func (c *Client) BlocksGetShards(ctx context.Context, last tongo.TonNodeBlockIdExt) ([]tongo.TonNodeBlockIdExt, error) {
	asReq, err := makeLiteServerAllShardsInfoRequest(last)
	if err != nil {
		return nil, err
	}
	req := makeLiteServerQueryRequest(asReq)
	resp, err := c.adnlClient.Request(ctx, req)
	if err != nil {
		return nil, err
	}
	parsedResp, err := parseLiteServerQueryResponse(resp)
	if err != nil {
		return nil, err
	}
	if parsedResp.Tag == LiteServerErrorTag {
		return nil, fmt.Errorf("lite server error: %v %v", parsedResp.LiteServerError.Code, parsedResp.LiteServerError.Message)
	}
	if parsedResp.Tag != LiteServerAllShardsInfoTag {
		return nil, fmt.Errorf("all shard info not recieved")
	}

	cells, err := boc.DeserializeBoc(parsedResp.LiteServerAllShardsInfo.Data)
	if err != nil {
		return nil, err
	}
	var inf tongo.AllShardsInfo
	err = tlb.Unmarshal(cells[0], &inf)
	if err != nil {
		return nil, err
	}
	var shards []tongo.TonNodeBlockIdExt
	for i, v := range inf.ShardHashes.Values() {
		wc, err := inf.ShardHashes.Keys()[i].ReadUint(32)
		if err != nil {
			return nil, err
		}
		for _, vv := range v.Value.BinTree.Values {
			shards = append(shards, tongo.TonNodeBlockIdExt{
				Workchain: int32(wc),
				Shard:     vv.New.NextValidatorShard,
				Seqno:     int32(vv.New.SeqNo),
				RootHash:  vv.New.RootHash,
				FileHash:  vv.New.FileHash,
			})
		}

	}
	return shards, nil
}
