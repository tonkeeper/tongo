package liteapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	mrand "math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/config"
	"github.com/startfellows/tongo/liteclient"
	"github.com/startfellows/tongo/tl"
	"github.com/startfellows/tongo/tlb"
	"github.com/startfellows/tongo/utils"
)

var (
	ErrBlockNotApplied = fmt.Errorf("block is not applied")
)

type connection struct {
	workchain   int32
	shardPrefix tongo.ShardID
	client      *liteclient.Client
}

type Client struct {
	connectionPool                 []connection
	targetBlockID                  *tongo.TonNodeBlockIdExt
	masterchainLastBlockCache      *tongo.TonNodeBlockIdExt
	masterchainLastBlockUpdateTime time.Time
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
	client := Client{}
	for _, ls := range options.LiteServers {
		serverPubkey, err := base64.StdEncoding.DecodeString(ls.Key)
		if err != nil {
			continue
		}
		c, err := liteclient.NewConnection(context.Background(), serverPubkey, ls.Host)
		if err != nil {
			continue
		}
		client.connectionPool = append(client.connectionPool, connection{
			workchain:   0,
			shardPrefix: tongo.MustParseShardID(-0x8000000000000000),
			client:      liteclient.NewClient(c),
		})
		return &client, nil
	}
	return nil, fmt.Errorf("all liteservers are unavailable")
}

func (c *Client) getMasterchainServer() *liteclient.Client {
	return c.connectionPool[mrand.Intn(len(c.connectionPool))].client
}

func (c *Client) getServerByAccountID(a tongo.AccountID) (*liteclient.Client, error) {
	if a.Workchain == -1 {
		return c.getMasterchainServer(), nil
	}
	for _, server := range c.connectionPool {
		if server.workchain != a.Workchain {
			continue
		}
		if server.shardPrefix.MatchAccountID(a) {
			return server.client, nil
		}
	}
	return nil, fmt.Errorf("can't find server for account %v", a.ToRaw())
}

func (c *Client) getServerByBlockID(block tongo.TonNodeBlockId) (*liteclient.Client, error) {
	if block.Workchain == -1 {
		return c.getMasterchainServer(), nil
	}
	for _, server := range c.connectionPool {
		if server.shardPrefix.MatchBlockID(block) {
			return server.client, nil
		}
	}
	return nil, fmt.Errorf("can't find server for block %v", block.String())
}

func (c Client) targetBlock(ctx context.Context) (tongo.TonNodeBlockIdExt, error) {
	if c.targetBlockID != nil {
		return *c.targetBlockID, nil
	}
	if time.Since(c.masterchainLastBlockUpdateTime) < 20*time.Second {
		return c.masterchainLastBlockCache, nil
	}
	return c.getlastBlock(ctx)
}

func (c *Client) refreshMasterchainTask() {
	// TODO: implement
	for {
		c.masterchainLastBlockCache = blahblah
	}
}

func (c Client) WithBlock(block tongo.TonNodeBlockIdExt) Client {
	c.targetBlockID = &block
	return c
}

func (c *Client) GetMasterchainInfo(ctx context.Context) (liteclient.LiteServerMasterchainInfo, error) {
	return c.getMasterchainServer().LiteServerGetMasterchainInfo(ctx)
}

func (c *Client) GetMasterchainInfoExt(ctx context.Context, mode uint32) (liteclient.LiteServerMasterchainInfoExt, error) {
	return c.getMasterchainServer().LiteServerGetMasterchainInfoExt(ctx, liteclient.LiteServerGetMasterchainInfoExtRequest(mode))
}

func (c *Client) GetTime(ctx context.Context) (uint32, error) {
	res, err := c.getMasterchainServer().LiteServerGetTime(ctx)
	return uint32(res), err
}

func (c *Client) GetVersion(ctx context.Context) (liteclient.LiteServerVersion, error) {
	return c.getMasterchainServer().LiteServerGetVersion(ctx)
}

func (c *Client) GetBlock(ctx context.Context, blockID tongo.TonNodeBlockIdExt) (tongo.Block, error) {
	server, err := c.getServerByBlockID(blockID.TonNodeBlockId)
	if err != nil {
		return tongo.Block{}, err
	}
	res, err := server.LiteServerGetBlock(ctx, liteclient.LiteServerGetBlockRequest(blockID))
	if err != nil {
		return tongo.Block{}, err
	}
	cells, err := boc.DeserializeBoc(res.Data)
	if err != nil {
		return tongo.Block{}, err
	}
	if len(cells) != 1 {
		return tongo.Block{}, boc.ErrNotSingleRoot
	}
	var data tongo.Block
	err = tlb.Unmarshal(cells[0], &data)
	if err != nil {
		return tongo.Block{}, err
	}
	return data, nil
}

func (c *Client) GetState(ctx context.Context, blockID tongo.TonNodeBlockIdExt) (tongo.State, tongo.Hash, tongo.Hash, error) {
	server, err := c.getServerByBlockID(blockID.TonNodeBlockId)
	if err != nil {
		return tongo.Block{}, tongo.Hash{}, tongo.Hash{}, err
	}
	res, err := server.LiteServerGetState(ctx, liteclient.LiteServerGetStateRequest(blockID))
	if err != nil {
		return tongo.Block{}, tongo.Hash{}, tongo.Hash{}, err
	}
	cells, err := boc.DeserializeBoc(res.Data)
	if err != nil {
		return tongo.Block{}, tongo.Hash{}, tongo.Hash{}, err
	}
	if len(cells) != 1 {
		return tongo.Block{}, tongo.Hash{}, tongo.Hash{}, boc.ErrNotSingleRoot
	}
	var state tongo.State // TODO: add State tlb type
	err = tlb.Unmarshal(cells[0], &state)
	if err != nil {
		return tongo.Block{}, tongo.Hash{}, tongo.Hash{}, err
	}
	return state, tongo.Hash(res.RootHash), tongo.Hash(res.FileHash), nil
}

func (c *Client) GetBlockHeader(ctx context.Context, blockID tongo.TonNodeBlockIdExt, mode uint32) (tongo.BlockInfo, error) {
	server, err := c.getServerByBlockID(blockID.TonNodeBlockId)
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	res, err := server.LiteServerGetBlockHeader(ctx, liteclient.LiteServerGetBlockHeaderRequest{
		Id:   blockID,
		Mode: mode,
	})
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	return decodeBlockHeader(res)
}

func (c *Client) LookupBlock(ctx context.Context, blockID tongo.TonNodeBlockId, mode uint32, lt *uint64, utime *uint32) (tongo.BlockInfo, error) {
	server, err := c.getServerByBlockID(blockID)
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	res, err := server.LiteServerLookupBlock(ctx, liteclient.LiteServerLookupBlockRequest{
		Mode: mode,
		Id: liteclient.TonNodeBlockId{
			Workchain: uint32(blockID.Workchain),
			Shard:     blockID.Shard,
			Seqno:     blockID.Seqno,
		},
		Lt:    lt,
		Utime: utime,
	})
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	return decodeBlockHeader(res)
}

func decodeBlockHeader(header liteclient.LiteServerBlockHeader) (tongo.BlockInfo, error) {
	cells, err := boc.DeserializeBoc(header.HeaderProof)
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	if len(cells) != 1 {
		return tongo.BlockInfo{}, boc.ErrNotSingleRoot
	}
	var proof struct {
		Proof tongo.MerkleProof[tongo.BlockHeader]
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return tongo.BlockInfo{}, err
	}
	return proof.Proof.VirtualRoot.Info, nil // TODO: maybe decode more
}

func (c *Client) SendMessage(ctx context.Context, payload []byte) (int32, error) {
	res, err := c.getMasterchainServer().LiteServerSendMessage(ctx, payload)
	return int32(res), err
}

func (c *Client) RunSmcMethod(
	ctx context.Context,
	accountID tongo.AccountID,
	method string,
	params tongo.VmStack,
) (uint32, tongo.VmStack, error) {
	stack, err := tl.Marshal(params)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	id, err := c.targetBlock(ctx)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	req := liteclient.LiteServerRunSmcMethodRequest{
		Mode:     0, // TODO: default mode
		Id:       id,
		Account:  accountID,
		MethodId: uint64(utils.Crc16String(method)&0xffff) | 0x10000,
		Params:   stack,
	}
	server, err := c.getServerByAccountID(accountID)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	res, err := server.LiteServerRunSmcMethod(ctx, req)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	cells, err := boc.DeserializeBoc(res.Result)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	if len(cells) != 1 {
		return 0, tongo.VmStack{}, boc.ErrNotSingleRoot
	}
	var result tongo.VmStack
	err = tlb.Unmarshal(cells[0], &result)
	return res.ExitCode, result, err
}

func (c *Client) GetAccountState(ctx context.Context, accountID tongo.AccountID) (tongo.AccountInfo, error) {
	id, err := c.targetBlock(ctx)
	if err != nil {
		return tongo.AccountInfo{}, err
	}
	server, err := c.getServerByAccountID(accountID)
	if err != nil {
		return tongo.AccountInfo{}, err
	}
	res, err := server.LiteServerGetAccountState(ctx, liteclient.LiteServerGetAccountStateRequest{
		Account: accountID,
		Id:      id,
	})
	if err != nil {
		return tongo.AccountInfo{}, err
	}
	if checkForNotApplied(err.(liteclient.LiteServerError)) { // TODO: add to other methods
		return tongo.AccountInfo{}, ErrBlockNotApplied
	}
	if len(res.State) == 0 {
		return tongo.AccountInfo{Status: tongo.AccountEmpty}, nil
	}
	acc, err := decodeRawAccountBoc(res.State)
	if err != nil {
		return tongo.AccountInfo{}, err
	}
	// TODO: proof check?
	// TODO: save raw account into account info?
	return convertTlbAccountToAccountState(acc)
}

func checkForNotApplied(e liteclient.LiteServerError) bool {
	return e.Message == "block is not applied"
}

func decodeRawAccountBoc(bocBytes []byte) (tongo.Account, error) {
	cells, err := boc.DeserializeBoc(bocBytes)
	if err != nil {
		return tongo.Account{}, err
	}
	if len(cells) != 1 {
		return tongo.Account{}, boc.ErrNotSingleRoot
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

//--------------------------------------------------------------------------------------------------------------------//

func (c *Client) GetLastRawAccount(ctx context.Context, accountId tongo.AccountID) (tongo.Account, error) {
	a, err := c.getLastRawAccountState(ctx, accountId)
	if err != nil {
		return tongo.Account{}, err
	}
	if len(a.State) == 0 {
		acc := tongo.Account{SumType: "AccountNone"}
		return acc, nil
	}
	return decodeRawAccountBoc(a.State)
}

func (c *Client) GetLastShardAccount(ctx context.Context, accountId tongo.AccountID) (tongo.ShardAccount, error) {
	a, err := c.getLastRawAccountState(ctx, accountId)
	if err != nil {
		return tongo.ShardAccount{}, err
	}
	var sa tongo.ShardAccount
	if len(a.State) == 0 {
		sa.Account.SumType = "AccountNone"
		return sa, nil
	}
	account, err := decodeRawAccountBoc(a.State)
	if err != nil {
		return tongo.ShardAccount{}, err
	}
	if len(a.Proof) == 0 {
		return tongo.ShardAccount{}, fmt.Errorf("empty proof")
	}
	lt, hash, err := decodeAccountDataFromProof(a.Proof, accountId)
	if err != nil {
		return tongo.ShardAccount{}, err
	}
	sa.LastTransHash = hash
	sa.LastTransLt = lt
	sa.Account = account
	return sa, nil
}

func decodeAccountDataFromProof(bocBytes []byte, account tongo.AccountID) (uint64, tongo.Hash, error) {
	cells, err := boc.DeserializeBoc(bocBytes)
	if err != nil {
		return 0, tongo.Hash{}, err
	}
	if len(cells) < 1 {
		return 0, tongo.Hash{}, fmt.Errorf("must be at least one root cell")
	}
	var proof struct {
		Proof tongo.MerkleProof[tongo.ShardStateUnsplit]
	}
	err = tlb.Unmarshal(cells[1], &proof) // cells order must be strictly defined
	if err != nil {
		return 0, tongo.Hash{}, err
	}
	values := proof.Proof.VirtualRoot.ShardStateUnsplit.Accounts.Accounts.Values()
	keys := proof.Proof.VirtualRoot.ShardStateUnsplit.Accounts.Accounts.Keys()
	for i, k := range keys {
		keyVal, err := k.ReadBytes(32)
		if err != nil {
			return 0, tongo.Hash{}, err
		}
		if bytes.Equal(keyVal, account.Address[:]) {
			return values[i].LastTransLt, values[i].LastTransHash, nil
		}
	}
	return 0, tongo.Hash{}, fmt.Errorf("account not found in ShardAccounts")
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
	server, err := c.getServerByAccountID(accountId)
	if err != nil {
		return nil, err
	}
	resp, err := server.Request(ctx, req)
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
	resp, err := c.getMasterchainServer().Request(ctx, req)
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

	conf := boc.NewCell()
	tlb.Marshal(conf, proof.Proof.VirtualRoot.ShardStateUnsplit.Custom.Value.Value.Config.Config)
	return conf, nil
}

// GetConfigAll
// liteServer.getConfigAll mode:# id:tonNode.blockIdExt = liteServer.ConfigInfo;
// liteServer.configInfo mode:# id:tonNode.blockIdExt state_proof:bytes config_proof:bytes = liteServer.ConfigInfo;
// Returns config: (Hashmap 32 ^Cell) as a Cell from config_proof
func (c *Client) GetConfigAll(ctx context.Context) (*tongo.McStateExtra, error) {
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

	lastBlock, err := c.GetMasterchainInfoExt(ctx, 0)
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
			ID:   lastBlock.Last,
		},
	}

	rBytes, err := tl.Marshal(r)
	if err != nil {
		return nil, err
	}
	req := makeLiteServerQueryRequest(rBytes)
	resp, err := c.getMasterchainServer().Request(ctx, req)
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
	return &proof.Proof.VirtualRoot.ShardStateUnsplit.Custom.Value.Value, nil
}

func (c *Client) GetConfigAllById(ctx context.Context, last tongo.TonNodeBlockIdExt) (*tongo.ShardState, error) {
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

	r := struct {
		tl.SumType
		GetConfigAllRequest getConfigAllRequest `tlSumType:"b7261b91"`
	}{
		SumType: "GetConfigAllRequest",
		GetConfigAllRequest: getConfigAllRequest{
			Mode: 0,
			ID:   last,
		},
	}

	rBytes, err := tl.Marshal(r)
	if err != nil {
		return nil, err
	}
	req := makeLiteServerQueryRequest(rBytes)
	resp, err := c.getMasterchainServer().Request(ctx, req)
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
		Proof tongo.MerkleProof[tongo.ShardState]
	}
	err = tlb.Unmarshal(cell[0], &proof)
	if err != nil {
		return nil, err
	}
	return &proof.Proof.VirtualRoot, nil
}

func (c *Client) BlocksGetShards(ctx context.Context, last tongo.TonNodeBlockIdExt) ([]tongo.TonNodeBlockIdExt, error) {
	asReq, err := makeLiteServerAllShardsInfoRequest(last)
	if err != nil {
		return nil, err
	}
	req := makeLiteServerQueryRequest(asReq)
	resp, err := c.getMasterchainServer().Request(ctx, req)
	if err != nil {
		return nil, err
	}
	var response struct {
		tl.SumType
		LiteServerAllShardsInfo struct {
			Id    tongo.TonNodeBlockIdExt
			Proof []byte
			Data  []byte
		} `tlSumType:"2de78f09"`
		Error LiteServerError `tlSumType:"48e1a9bb"`
	}
	err = tl.Unmarshal(bytes.NewReader(resp), &response)
	if err != nil {
		return nil, err
	}
	cells, err := boc.DeserializeBoc(response.LiteServerAllShardsInfo.Data)
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
			shards = append(shards, vv.ToBlockId(int32(wc)))
		}

	}
	return shards, nil
}

func (c *Client) ValidatorStats(ctx context.Context, mode uint32, last tongo.TonNodeBlockIdExt, limit uint32, startAfter *tongo.Hash, modifiedAfter *uint32) (*tongo.ShardStateUnsplit, error) {
	asReq, err := makeLiteServerGetValidatorStatsRequest(mode, last, limit, startAfter, modifiedAfter)
	if err != nil {
		return nil, err
	}
	req := makeLiteServerQueryRequest(asReq)
	resp, err := c.getMasterchainServer().Request(ctx, req)
	if err != nil {
		return nil, err
	}
	var response struct {
		tl.SumType
		LiteServerValidatorStats struct {
			Mode       uint32
			Id         tongo.TonNodeBlockIdExt
			Count      uint32
			Complete   uint32
			StateProof []byte
			DataProof  []byte
		} `tlSumType:"d896f7b9"`
		Error LiteServerError `tlSumType:"48e1a9bb"`
	}

	err = tl.Unmarshal(bytes.NewReader(resp), &response)
	if err != nil {
		return nil, err
	}
	if response.SumType == "Error" {
		return nil, fmt.Errorf(response.Error.Message)
	}
	cells, err := boc.DeserializeBoc(response.LiteServerValidatorStats.DataProof)
	if err != nil {
		return nil, err
	}

	var proof struct {
		Proof tongo.MerkleProof[tongo.ShardStateUnsplit]
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return nil, err
	}

	// }
	return &proof.Proof.VirtualRoot, nil //shards, nil
}
func (c *Client) GetBlockProof(ctx context.Context, mode uint32, knownBlock tongo.TonNodeBlockIdExt, targetBlock *tongo.TonNodeBlockIdExt) ([]tongo.BlockProof, error) {
	asReq, err := makeLiteServerGetBlockProofRequest(mode, knownBlock, targetBlock)
	if err != nil {
		return nil, err
	}
	req := makeLiteServerQueryRequest(asReq)
	var server *liteclient.Client
	if targetBlock != nil {
		server, err = c.getServerByBlockID(targetBlock.TonNodeBlockId)
	} else {
		server, err = c.getServerByBlockID(knownBlock.TonNodeBlockId)
	}
	if err != nil {
		return nil, err
	}
	resp, err := server.Request(ctx, req)
	if err != nil {
		return nil, err
	}

	// liteServer.signature node_id_short:int256 signature:bytes = liteServer.Signature;
	type Signature struct {
		NodeIdShort tongo.Hash
		Signature   []byte
	}

	// liteServer.signatureSet validator_set_hash:int catchain_seqno:int
	// signatures:(vector liteServer.signature) = liteServer.SignatureSet;
	type SignatureSet struct {
		ValidatorSetHash int32
		CatchainSeqNo    int32
		Signatures       Signature
	}

	type blockLink struct {
		tl.SumType
		// liteServer.blockLinkForward to_key_block:Bool from:tonNode.blockIdExt
		// to:tonNode.blockIdExt dest_proof:bytes config_proof:bytes
		// signatures:liteServer.SignatureSet
		// = liteServer.BlockLink;
		BlockLinkForward struct {
			ToKeyBlock  uint32
			From        tongo.TonNodeBlockIdExt
			To          tongo.TonNodeBlockIdExt
			DestProof   []byte
			ConfigProof []byte
			Signatures  SignatureSet
		} `tlSumType:"1cce0f52"`
		// liteServer.blockLinkBack to_key_block:Bool from:tonNode.blockIdExt
		// to:tonNode.blockIdExt dest_proof:bytes proof:bytes state_proof:bytes
		// = liteServer.BlockLink;
		BlockLinkBack struct {
			ToKeyBlock uint32
			From       tongo.TonNodeBlockIdExt
			To         tongo.TonNodeBlockIdExt
			DestProof  []byte
			Proof      []byte
			StateProof []byte
		} `tlSumType:"ef1b7eef"`
	}

	// liteServer.partialBlockProof complete:Bool from:tonNode.blockIdExt
	// to:tonNode.blockIdExt steps:(vector liteServer.BlockLink) =
	// liteServer.PartialBlockProof
	var response struct {
		tl.SumType
		LiteServerGetProofBlock struct {
			Complete uint32
			From     tongo.TonNodeBlockIdExt
			To       tongo.TonNodeBlockIdExt
			Steps    []blockLink
		} `tlSumType:"c1d2d08e"`
		Error LiteServerError `tlSumType:"48e1a9bb"`
	}

	err = tl.Unmarshal(bytes.NewReader(resp), &response)
	if err != nil {
		return nil, err
	}
	if response.SumType == "Error" {
		return nil, fmt.Errorf(response.Error.Message)
	}

	for i := range response.LiteServerGetProofBlock.Steps {
		if response.LiteServerGetProofBlock.Steps[i].SumType == "BlockLinkBack" {
			cells, err := boc.DeserializeBoc(response.LiteServerGetProofBlock.Steps[i].BlockLinkBack.StateProof)
			if err != nil {
				return nil, err
			}
			cells, err = boc.DeserializeBoc(response.LiteServerGetProofBlock.Steps[i].BlockLinkBack.Proof)
			if err != nil {
				return nil, err
			}
			cells, err = boc.DeserializeBoc(response.LiteServerGetProofBlock.Steps[i].BlockLinkBack.DestProof)
			if err != nil {
				return nil, err
			}

			var proof struct {
				Proof tongo.MerkleProof[tongo.BlockHeader]
			}
			err = tlb.Unmarshal(cells[0], &proof)
			if err != nil {
				return nil, err
			}

		}
		if response.LiteServerGetProofBlock.Steps[i].SumType == "BlockLinkForward" {
			cells, err := boc.DeserializeBoc(response.LiteServerGetProofBlock.Steps[i].BlockLinkForward.ConfigProof)
			if err != nil {
				return nil, err
			}
			cells, err = boc.DeserializeBoc(response.LiteServerGetProofBlock.Steps[i].BlockLinkForward.DestProof)
			if err != nil {
				return nil, err
			}

			var proof struct {
				Proof tongo.MerkleProof[tongo.BlockHeader]
			}
			err = tlb.Unmarshal(cells[0], &proof)
			if err != nil {
				return nil, err
			}

		}

	}
	return nil, nil

}

// GetOneRawTransaction
// liteServer.getOneTransaction id:tonNode.blockIdExt account:liteServer.accountId lt:long = liteServer.TransactionInfo;
// liteServer.transactionInfo id:tonNode.blockIdExt proof:bytes transaction:bytes = liteServer.TransactionInfo;
func (c *Client) GetOneRawTransaction(ctx context.Context, id tongo.TonNodeBlockIdExt, accountId tongo.AccountID, lt uint64) ([]*boc.Cell, []byte, error) {
	type getOneTransactionRequest struct {
		ID      tongo.TonNodeBlockIdExt
		Account tongo.AccountID
		Lt      uint64
	}
	type transactionInfo struct {
		Id          tongo.TonNodeBlockIdExt
		Proof       []byte
		Transaction []byte
	}
	r := struct {
		tl.SumType
		GetOneTransactionRequest getOneTransactionRequest `tlSumType:"ea240fd4"`
	}{
		SumType: "GetOneTransactionRequest",
		GetOneTransactionRequest: getOneTransactionRequest{
			ID:      id,
			Account: accountId,
			Lt:      lt,
		},
	}
	rBytes, err := tl.Marshal(r)
	if err != nil {
		return nil, nil, err
	}
	req := makeLiteServerQueryRequest(rBytes)
	server, err := c.getServerByAccountID(accountId)
	if err != nil {
		return nil, nil, err
	}
	resp, err := server.Request(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	var pResp struct {
		tl.SumType
		TransactionInfo transactionInfo `tlSumType:"47edde0e"`
		Error           LiteServerError `tlSumType:"48e1a9bb"`
	}
	reader := bytes.NewReader(resp)
	err = tl.Unmarshal(reader, &pResp)
	if err != nil {
		return nil, nil, err
	}
	if pResp.SumType == "Error" {
		return nil, nil, fmt.Errorf("error code: %v , message: %v", pResp.Error.Code, pResp.Error.Message)
	}
	cells, err := boc.DeserializeBoc(pResp.TransactionInfo.Transaction)
	if err != nil {
		return nil, nil, err
	}
	if len(cells) != 1 {
		return nil, nil, fmt.Errorf("must be one root cell")
	}
	return cells, pResp.TransactionInfo.Proof, nil
}
