package liteclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	mrand "math/rand"
	"net/http"
	"sync"

	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/adnl"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/config"
	"github.com/startfellows/tongo/tl"
	"github.com/startfellows/tongo/tlb"
	"github.com/startfellows/tongo/utils"
)

var (
	ErrBlockNotApplied = fmt.Errorf("block is not applied")
)

type liteserverConnection struct {
	workchain   int32
	shardPrefix tongo.ShardID
	client      *adnl.Client
}

type Client struct {
	adnlClient []liteserverConnection
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
		c, err := adnl.NewConnection(context.Background(), serverPubkey, ls.Host)
		if err != nil {
			continue
		}
		client.adnlClient = append(client.adnlClient, liteserverConnection{
			workchain:   0,
			shardPrefix: tongo.MustParseShardID(-0x8000000000000000),
			client:      adnl.NewClient(c),
		})
		return &client, nil
	}
	return nil, fmt.Errorf("all liteservers are unavailable")
}

func (c *Client) getMasterchainServer() *adnl.Client {
	return c.adnlClient[mrand.Intn(len(c.adnlClient))].client
}

func (c *Client) getServerByAccountID(a tongo.AccountID) (*adnl.Client, error) {
	if a.Workchain == -1 {
		return c.getMasterchainServer(), nil
	}
	for _, server := range c.adnlClient {
		if server.workchain != a.Workchain {
			continue
		}
		if server.shardPrefix.MatchAccountID(a) {
			return server.client, nil
		}
	}
	return nil, fmt.Errorf("can't find server for account %v", a.ToRaw())
}

func (c *Client) getServerByBlockID(block tongo.TonNodeBlockId) (*adnl.Client, error) {
	if block.Workchain == -1 {
		return c.getMasterchainServer(), nil
	}
	for _, server := range c.adnlClient {
		if server.shardPrefix.MatchBlockID(block) {
			return server.client, nil
		}
	}
	return nil, fmt.Errorf("can't find server for block %v", block.String())
}

func (c *Client) GetAccountState(ctx context.Context, accountId tongo.AccountID) (tongo.AccountInfo, error) {
	a, err := c.getLastRawAccountState(ctx, accountId)
	if err != nil {
		return tongo.AccountInfo{}, err
	}
	if len(a.State) == 0 {
		return tongo.AccountInfo{Status: tongo.AccountEmpty}, nil
	}
	account, err := decodeRawAccountBoc(a.State)
	if err != nil {
		return tongo.AccountInfo{}, err
	}
	return convertTlbAccountToAccountState(account)
}

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

func (c *Client) GetRawAccountById(ctx context.Context, accountId tongo.AccountID, blockId tongo.TonNodeBlockIdExt) (tongo.Account, error) {
	a, err := c.getRawAccountStateById(ctx, accountId, blockId)
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

func (c *Client) getLastRawAccountState(ctx context.Context, accountId tongo.AccountID) (LiteServerAccountState, error) {
	mcInfo, err := c.GetMasterchainInfo(ctx)
	if err != nil {
		return LiteServerAccountState{}, err
	}
	st, err := c.getRawAccountState(ctx, mcInfo, accountId)
	if err != nil && errors.Is(err, ErrBlockNotApplied) {
		prevMcInfo, _, err := c.LookupBlock(ctx, 1, tongo.TonNodeBlockId{Workchain: mcInfo.Workchain, Shard: mcInfo.Shard, Seqno: mcInfo.Seqno - 1}, 0, 0)
		if err != nil {
			return LiteServerAccountState{}, err
		}
		return c.getRawAccountState(ctx, prevMcInfo, accountId)
	}
	return st, err
}

func (c *Client) getRawAccountStateById(ctx context.Context, accountId tongo.AccountID, blockId tongo.TonNodeBlockIdExt) (LiteServerAccountState, error) {
	st, err := c.getRawAccountState(ctx, blockId, accountId)
	if err != nil && errors.Is(err, ErrBlockNotApplied) {
		prevMcInfo, _, err := c.LookupBlock(ctx, 1, tongo.TonNodeBlockId{Workchain: blockId.Workchain, Shard: blockId.Shard, Seqno: blockId.Seqno - 1}, 0, 0)
		if err != nil {
			return LiteServerAccountState{}, err
		}
		return c.getRawAccountState(ctx, prevMcInfo, accountId)
	}
	return st, err
}

func (c *Client) getRawAccountState(ctx context.Context, masterchainInfo tongo.TonNodeBlockIdExt, accountId tongo.AccountID) (LiteServerAccountState, error) {
	asReq, err := makeLiteServerGetAccountStateRequest(masterchainInfo, accountId)
	if err != nil {
		return LiteServerAccountState{}, err
	}
	req := makeLiteServerQueryRequest(asReq)
	server, err := c.getServerByAccountID(accountId)
	if err != nil {
		return LiteServerAccountState{}, err
	}
	resp, err := server.Request(ctx, req)
	if err != nil {
		return LiteServerAccountState{}, err
	}
	var parsedResp struct {
		tl.SumType
		LiteServerError        LiteServerError        `tlSumType:"48e1a9bb"`
		LiteServerAccountState LiteServerAccountState `tlSumType:"51c77970"`
	}
	err = tl.Unmarshal(bytes.NewReader(resp), &parsedResp)
	if err != nil {
		return LiteServerAccountState{}, err
	}
	switch parsedResp.SumType {
	case "LiteServerError":
		if parsedResp.LiteServerError.Message == "block is not applied" {
			return LiteServerAccountState{}, ErrBlockNotApplied
		}
		return LiteServerAccountState{}, fmt.Errorf("lite server error: %v %v", parsedResp.LiteServerError.Code, parsedResp.LiteServerError.Message)
	case "LiteServerAccountState":
	default:
		return LiteServerAccountState{}, fmt.Errorf("account state not recieved")
	}
	return parsedResp.LiteServerAccountState, nil
}

type LiteServerMasterchainInfo struct {
	Last          tongo.TonNodeBlockIdExt
	StateRootHash tongo.Hash
	// TODO: add init
}

type LiteServerMasterchainInfoExt struct {
	Mode          uint32
	Version       uint32
	Capabilities  uint64
	Last          tongo.TonNodeBlockIdExt
	LastUTime     uint32
	Now           uint32
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

func makeLiteServerGetMasterchainInfoExtRequest(mode uint32) ([]byte, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, LiteServerGetMasterchainInfoExtTag)
	m, err := tl.Marshal(mode)
	if err != nil {
		return nil, err
	}
	payload = append(payload, m...)
	return payload, nil
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

// liteServer.getValidatorStats#091a58bc mode:# id:tonNode.blockIdExt
// limit:int start_after:mode.0?int256 modified_after:mode.2?int
// = liteServer.ValidatorStats;
func makeLiteServerGetValidatorStatsRequest(mode uint32,
	blockIdExt tongo.TonNodeBlockIdExt,
	limit uint32,
	startAfter *tongo.Hash,
	modifiedAfter *uint32,
) ([]byte, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, LiteServerGetValidatorStatsTag)
	m, err := tl.Marshal(mode)
	if err != nil {
		return nil, err
	}
	payload = append(payload, m...)
	block, err := blockIdExt.MarshalTL()
	if err != nil {
		return nil, err
	}
	payload = append(payload, block...)
	l, err := tl.Marshal(limit)
	if err != nil {
		return nil, err
	}
	payload = append(payload, l...)
	if (mode & 0x01) == 0x1 { // start_after:mode.0?int256
		if startAfter == nil {
			return nil, fmt.Errorf("startAfter is null, but mode.0 is true. Please, see lite_api.tl")
		}
		sA, err := startAfter.MarshalTL()
		if err != nil {
			return nil, err
		}
		payload = append(payload, sA...)
	}
	if (mode & 0x04) == 0x4 { // modified_after:mode.2?int
		if modifiedAfter == nil {
			return nil, fmt.Errorf("modifiedAfter is null, but mode.2 is true. Please, see lite_api.tl")
		}
		mA, err := tl.Marshal(modifiedAfter)
		if err != nil {
			return nil, err
		}
		payload = append(payload, mA...)
	}
	return payload, nil
}

// liteServer.lookupBlock mode:# id:tonNode.blockId lt:mode.1?long
// utime:mode.2?int = liteServer.BlockHeader;
func makeLiteServerLookupBlockRequest(mode uint32,
	blockId tongo.TonNodeBlockId,
	lt uint64, utime uint32,
) ([]byte, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, LiteServerLookupBlockTag)
	m, err := tl.Marshal(mode)
	if err != nil {
		return nil, err
	}
	payload = append(payload, m...)
	block, err := tl.Marshal(blockId)
	if err != nil {
		return nil, err
	}
	payload = append(payload, block...)
	if (mode & 0x02) == 0x2 { // lt:mode.1?long
		l, err := tl.Marshal(lt)
		if err != nil {
			return nil, err
		}
		payload = append(payload, l...)
	}
	if (mode & 0x04) == 0x4 { // utime:mode.2?int
		u, err := tl.Marshal(utime)
		if err != nil {
			return nil, err
		}
		payload = append(payload, u...)
	}
	return payload, nil
}

// liteServer.getBlockProof mode:# known_block:tonNode.blockIdExt
// target_block:mode.0?tonNode.blockIdExt = liteServer.PartialBlockProof;
func makeLiteServerGetBlockProofRequest(mode uint32,
	knownBlock tongo.TonNodeBlockIdExt,
	targetBlock *tongo.TonNodeBlockIdExt,
) ([]byte, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, LiteServerGetBlockProofTag)
	m, err := tl.Marshal(mode)
	if err != nil {
		return nil, err
	}
	payload = append(payload, m...)
	block, err := tl.Marshal(knownBlock)
	if err != nil {
		return nil, err
	}
	payload = append(payload, block...)
	if (mode & 0x01) == 0x1 { // target_block:mode.0?tonNode.blockIdExt
		if targetBlock == nil {
			return nil, fmt.Errorf("target_block is nil")
		}
		block, err = tl.Marshal(*targetBlock)
		if err != nil {
			return nil, err
		}
		payload = append(payload, block...)
	}
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
	resp, err := c.getMasterchainServer().Request(ctx, req)
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

// GetMasterchainInfo
// liteServer.getMasterchainInfo = liteServer.MasterchainInfo;
// liteServer.masterchainInfo last:tonNode.blockIdExt state_root_hash:int256
// init:tonNode.zeroStateIdExt = liteServer.MasterchainInfo;
func (c *Client) GetMasterchainInfo(ctx context.Context) (tongo.TonNodeBlockIdExt, error) {
	req := makeLiteServerQueryRequest(makeLiteServerGetMasterchainInfoRequest())
	resp, err := c.getMasterchainServer().Request(ctx, req)
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

// GetMasterchainInfoExt
// liteServer.getMasterchainInfoExt mode:# = liteServer.MasterchainInfoExt;
// liteServer.masterchainInfoExt mode:# version:int capabilities:long
// last:tonNode.blockIdExt last_utime:int now:int state_root_hash:int256
// init:tonNode.zeroStateIdExt = liteServer.MasterchainInfoExt;
func (c *Client) GetMasterchainInfoExt(ctx context.Context, mode uint32) (LiteServerMasterchainInfoExt, error) {
	asReq, err := makeLiteServerGetMasterchainInfoExtRequest(mode)
	if err != nil {
		return LiteServerMasterchainInfoExt{}, err
	}
	req := makeLiteServerQueryRequest(asReq)
	resp, err := c.getMasterchainServer().Request(ctx, req)
	if err != nil {
		return LiteServerMasterchainInfoExt{}, err
	}
	var parsedResp struct {
		tl.SumType
		LiteServerError              LiteServerError              `tlSumType:"48e1a9bb"`
		LiteServerMasterchainInfoExt LiteServerMasterchainInfoExt `tlSumType:"f5e0cca8"`
	}
	err = tl.Unmarshal(bytes.NewReader(resp), &parsedResp)
	if err != nil {
		return LiteServerMasterchainInfoExt{}, err
	}
	switch parsedResp.SumType {
	case "LiteServerError":
		return LiteServerMasterchainInfoExt{}, fmt.Errorf("lite server error: %v %v", parsedResp.LiteServerError.Code, parsedResp.LiteServerError.Message)
	case "LiteServerMasterchainInfoExt":
		return parsedResp.LiteServerMasterchainInfoExt, nil
	default:
		return LiteServerMasterchainInfoExt{}, fmt.Errorf("masterchain info not recieved")
	}

}

// RunSmcMethod
// Run smart contract method by name and parameters
// liteServer.runSmcMethod mode:# id:tonNode.blockIdExt account:liteServer.accountId method_id:long params:bytes = liteServer.RunMethodResult;
// liteServer.runMethodResult mode:# id:tonNode.blockIdExt shardblk:tonNode.blockIdExt shard_proof:mode.0?bytes
// proof:mode.0?bytes state_proof:mode.1?bytes init_c7:mode.3?bytes lib_extras:mode.4?bytes exit_code:int result:mode.2?bytes = liteServer.RunMethodResult;
func (c *Client) RunSmcMethod(ctx context.Context, mode uint32, accountId tongo.AccountID, method string, params tongo.VmStack) (uint32, tongo.VmStack, error) {
	type runSmcRequest struct {
		Mode     uint32
		Id       tongo.TonNodeBlockIdExt
		Account  tongo.AccountID
		MethodId uint64
		Params   tongo.VmStack
	}
	info, err := c.GetMasterchainInfo(ctx)
	if err != nil {
		return 0, tongo.VmStack{}, err
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
		return 0, tongo.VmStack{}, err
	}
	req := makeLiteServerQueryRequest(payload)
	server, err := c.getServerByAccountID(accountId)
	if err != nil {
		return 0, nil, err
	}
	resp, err := server.Request(ctx, req)
	if err != nil {
		return 0, tongo.VmStack{}, err
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
		return 0, tongo.VmStack{}, err
	}
	if response.SumType == "Error" {
		return 0, tongo.VmStack{}, fmt.Errorf("error code: %v , message: %v", response.Error.Code, response.Error.Message)
	}
	if response.SumType != "RunMethodResult" {
		return 0, tongo.VmStack{}, fmt.Errorf("not RunMethodResult response")
	}
	return response.RunMethodResult.ExitCode, response.RunMethodResult.Result, nil
}

// RunSmcMethod
// Run smart contract method by name and parameters
// liteServer.runSmcMethod mode:# id:tonNode.blockIdExt account:liteServer.accountId method_id:long params:bytes = liteServer.RunMethodResult;
// liteServer.runMethodResult mode:# id:tonNode.blockIdExt shardblk:tonNode.blockIdExt shard_proof:mode.0?bytes
// proof:mode.0?bytes state_proof:mode.1?bytes init_c7:mode.3?bytes lib_extras:mode.4?bytes exit_code:int result:mode.2?bytes = liteServer.RunMethodResult;
func (c *Client) RunSmcMethodByExtBlockId(ctx context.Context, mode uint32, id tongo.TonNodeBlockIdExt, accountId tongo.AccountID, method string, params tongo.VmStack) (tongo.VmStack, error) {
	type runSmcRequest struct {
		Mode     uint32
		Id       tongo.TonNodeBlockIdExt
		Account  tongo.AccountID
		MethodId uint64
		Params   tongo.VmStack
	}
	r := struct {
		tl.SumType
		RunSmcRequest runSmcRequest `tlSumType:"d25dc65c"`
	}{
		SumType: "RunSmcRequest",
		RunSmcRequest: runSmcRequest{
			Mode:     mode,
			Id:       id,
			Account:  accountId,
			MethodId: uint64(utils.Crc16String(method)&0xffff) | 0x10000,
			Params:   params,
		},
	}
	payload, err := tl.Marshal(r)
	if err != nil {
		return tongo.VmStack{}, err
	}
	req := makeLiteServerQueryRequest(payload)
	server, err := c.getServerByAccountID(accountId)
	if err != nil {
		return nil, err
	}
	resp, err := server.Request(ctx, req)
	if err != nil {
		return tongo.VmStack{}, err
	}
	var response struct {
		tl.SumType
		RunMethodResult struct {
			Mode     uint32
			Id       tongo.TonNodeBlockIdExt
			ShardBlk tongo.TonNodeBlockIdExt
			// ShardProof []byte
			// Proof      []byte
			// StateProof []byte
			// InitC7     []byte
			// LibExtras  []byte
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
	var server *adnl.Client
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

// GetBlock
// liteServer.getBlock id:tonNode.blockIdExt = liteServer.BlockData;
// liteServer.blockData id:tonNode.blockIdExt data:bytes = liteServer.BlockData;
func (c *Client) GetBlock(ctx context.Context, blockID tongo.TonNodeBlockIdExt) (*tongo.TonNodeBlockIdExt, tongo.Block, error) {
	r := struct {
		tl.SumType
		GetBlockRequest tongo.TonNodeBlockIdExt `tlSumType:"0dcf7763"`
	}{
		SumType:         "GetBlockRequest",
		GetBlockRequest: blockID,
	}

	rBytes, err := tl.Marshal(r)
	if err != nil {
		return nil, tongo.Block{}, err
	}
	req := makeLiteServerQueryRequest(rBytes)
	server, err := c.getServerByBlockID(blockID.TonNodeBlockId)
	if err != nil {
		return nil, tongo.Block{}, err
	}
	resp, err := server.Request(ctx, req)
	if err != nil {
		return nil, tongo.Block{}, err
	}
	var pResp struct {
		tl.SumType
		BlockData struct {
			ID   tongo.TonNodeBlockIdExt
			Data []byte
		} `tlSumType:"6ced74a5"`
		Error LiteServerError `tlSumType:"48e1a9bb"`
	}
	reader := bytes.NewReader(resp)
	err = tl.Unmarshal(reader, &pResp)
	if err != nil {
		return nil, tongo.Block{}, err
	}
	if pResp.SumType == "Error" {
		if pResp.Error.Message == "block is not applied" {
			return nil, tongo.Block{}, ErrBlockNotApplied
		}
		return nil, tongo.Block{}, fmt.Errorf("error code: %v , message: %v", pResp.Error.Code, pResp.Error.Message)
	}
	cell, err := boc.DeserializeBoc(pResp.BlockData.Data)
	if err != nil {
		return nil, tongo.Block{}, err
	}
	var data tongo.Block
	err = tlb.Unmarshal(cell[0], &data)
	if err != nil {
		return nil, tongo.Block{}, err
	}
	return &pResp.BlockData.ID, data, nil
}

// LookupBlock
// liteServer.lookupBlock mode:# id:tonNode.blockId lt:mode.1?long utime:mode.2?int = liteServer.BlockHeader;
// liteServer.blockHeader id:tonNode.blockIdExt mode:# header_proof:bytes = liteServer.BlockHeader;
func (c *Client) LookupBlock(ctx context.Context, mode uint32, blockID tongo.TonNodeBlockId, lt uint64, utime uint32) (tongo.TonNodeBlockIdExt, tongo.BlockInfo, error) {
	asReq, err := makeLiteServerLookupBlockRequest(mode, blockID, lt, utime)
	if err != nil {
		return tongo.TonNodeBlockIdExt{}, tongo.BlockInfo{}, err
	}
	server, err := c.getServerByBlockID(blockID)
	if err != nil {
		return tongo.TonNodeBlockIdExt{}, tongo.BlockInfo{}, err
	}
	resp, err := server.Request(ctx, makeLiteServerQueryRequest(asReq))
	if err != nil {
		return tongo.TonNodeBlockIdExt{}, tongo.BlockInfo{}, err
	}
	var pResp struct {
		tl.SumType
		BlockHeader struct {
			ID          tongo.TonNodeBlockIdExt
			Mode        uint32
			HeaderProof []byte
		} `tlSumType:"19822d75"`
		Error LiteServerError `tlSumType:"48e1a9bb"`
	}
	reader := bytes.NewReader(resp)
	err = tl.Unmarshal(reader, &pResp)
	if err != nil {
		return tongo.TonNodeBlockIdExt{}, tongo.BlockInfo{}, err
	}
	if pResp.SumType == "Error" {
		return tongo.TonNodeBlockIdExt{}, tongo.BlockInfo{}, fmt.Errorf("error code: %v , message: %v", pResp.Error.Code, pResp.Error.Message)
	}
	cells, err := boc.DeserializeBoc(pResp.BlockHeader.HeaderProof)
	if err != nil {
		return tongo.TonNodeBlockIdExt{}, tongo.BlockInfo{}, err
	}

	var proof struct {
		Proof tongo.MerkleProof[tongo.BlockHeader]
	}
	err = tlb.Unmarshal(cells[0], &proof)
	if err != nil {
		return tongo.TonNodeBlockIdExt{}, tongo.BlockInfo{}, err
	}
	return pResp.BlockHeader.ID, proof.Proof.VirtualRoot.Info, nil
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
