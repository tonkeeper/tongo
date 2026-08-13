package wallet

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/tontest"
	"github.com/tonkeeper/tongo/txemulator"
)

func TestGetW5ExtensionsList(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("hangs in CI")
	}

	tests := []struct {
		name           string
		accountID      string
		wantErr        bool
		wantExtensions map[ton.AccountID]struct{}
	}{
		{
			name:      "wallet v5 with extensions",
			accountID: "0:8a189456e840670f5c09c5db75c829454cc6b1e5dd81cc20065fdb1999a9cbad",
			wantExtensions: map[ton.AccountID]struct{}{
				ton.MustParseAccountID("0:6ccd325a858c379693fae2bcaab1c2906831a4e10a6c3bb44ee8b615bca1d220"): {},
				ton.MustParseAccountID("0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb"): {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountID := ton.MustParseAccountID(tt.accountID)
			cli, err := liteapi.NewClient(liteapi.Testnet())
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			state, err := cli.GetAccountState(context.Background(), accountID)
			if err != nil {
				t.Fatalf("GetAccountState() error = %v", err)
			}
			extensions, err := GetW5BetaExtensionsList(state)
			if err != nil {
				t.Fatalf("GetW5BetaExtensionsList() error = %v", err)
			}
			if !reflect.DeepEqual(extensions, tt.wantExtensions) {
				t.Errorf("GetW5BetaExtensionsList() = %v, want %v", extensions, tt.wantExtensions)
			}
		})
	}
}

func TestNewWalletV5R1(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		want *walletV5R1
	}{
		{
			name: "workchain 0, testnet",
			opts: []Option{
				WithWorkchain(0),
				WithNetworkGlobalID(TestnetGlobalID),
			},
			want: &walletV5R1{
				workchain: 0,
				walletID:  2147483645,
			},
		},
		{
			name: "workchain 0, mainnet",
			opts: []Option{
				WithWorkchain(0),
				WithNetworkGlobalID(MainnetGlobalID),
			},
			want: &walletV5R1{
				workchain: 0,
				walletID:  2147483409,
			},
		},
		{
			name: "workchain 0, mainnet, subwallet 1",
			opts: []Option{
				WithWorkchain(0),
				WithNetworkGlobalID(MainnetGlobalID),
				WithSubWalletID(1),
			},
			want: &walletV5R1{
				workchain: 0,
				walletID:  2147483408,
			},
		},
		{
			name: "workchain -1, mainnet",
			opts: []Option{
				WithWorkchain(-1),
				WithNetworkGlobalID(MainnetGlobalID),
			},
			want: &walletV5R1{
				workchain: -1,
				walletID:  8388369,
			},
		},
		{
			name: "workchain -1, testnet",
			opts: []Option{
				WithWorkchain(-1),
				WithNetworkGlobalID(TestnetGlobalID),
			},
			want: &walletV5R1{
				workchain: -1,
				walletID:  8388605,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			wallet := NewWalletV5R1(nil, applyOptions(tt.opts...))
			if wallet.walletID != tt.want.walletID {
				t.Errorf("NewWalletV5R1() = %v, want %v", wallet.walletID, tt.want.walletID)
			}
			if wallet.workchain != tt.want.workchain {
				t.Errorf("NewWalletV5R1() = %v, want %v", wallet.workchain, tt.want.workchain)
			}
		})
	}
}

func Test_walletV5R1_generateAddress(t *testing.T) {
	tests := []struct {
		name       string
		privateKey string
		opts       []Option
		want       ton.AccountID
	}{
		{
			name:       "workchain 0, testnet",
			privateKey: "7c94066ee822c97aa6992fa1c506bfd56d0d8fed2f1027070af7e0a683d46fb671ced1c4c69e53eb7ede24658375f56c142d22cdb21d0728138cb53b817e454e",
			opts: []Option{
				WithWorkchain(0),
				WithNetworkGlobalID(TestnetGlobalID),
			},
			want: ton.MustParseAccountID("0:aa7bd5aa1614bc01f5460cfdb14224cb8db7a89612c6b2f21b6e043a0b75d3e6"),
		},
		{
			name:       "workchain 0, mainnet",
			privateKey: "7c94066ee822c97aa6992fa1c506bfd56d0d8fed2f1027070af7e0a683d46fb671ced1c4c69e53eb7ede24658375f56c142d22cdb21d0728138cb53b817e454e",
			opts: []Option{
				WithWorkchain(0),
				WithNetworkGlobalID(MainnetGlobalID),
			},
			want: ton.MustParseAccountID("0:827137ba7a1ad871a8a8605e8dba799666abb952dfe9eff6e9dfa96700ae16f4"),
		},
		{
			name:       "workchain 0, mainnet, subwallet 1",
			privateKey: "7c94066ee822c97aa6992fa1c506bfd56d0d8fed2f1027070af7e0a683d46fb671ced1c4c69e53eb7ede24658375f56c142d22cdb21d0728138cb53b817e454e",
			opts: []Option{
				WithWorkchain(0),
				WithNetworkGlobalID(MainnetGlobalID),
				WithSubWalletID(1),
			},
			want: ton.MustParseAccountID("0:f6b01bbe21f27c1bf9a1bf43dac3b1c3e42489b9d2fdd7776a1e68908ed699ee"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, err := hex.DecodeString(tt.privateKey)
			if err != nil {
				t.Fatalf("hex.DecodeString() error = %v", err)
			}

			publicKey := ed25519.PrivateKey(privateKey).Public().(ed25519.PublicKey)
			w := NewWalletV5R1(publicKey, applyOptions(tt.opts...))
			address, err := w.generateAddress()
			if err != nil {
				t.Fatalf("generateAddress() error = %v", err)
			}
			if address.ToRaw() != tt.want.ToRaw() {
				t.Errorf("generateAddress() got = %v, want %v", address, tt.want)
			}
		})
	}
}

func TestGetW5R1ExtensionsList(t *testing.T) {
	tests := []struct {
		name           string
		accountID      string
		wantExtensions map[ton.AccountID]struct{}
	}{
		{
			name:      "wallet v5 with extensions",
			accountID: "0:d7391407f03695b3af341b29ab3137a8c269ab091e0043b4301a15a034828e1c",
			wantExtensions: map[ton.AccountID]struct{}{
				ton.MustParseAccountID("0:6ccd325a858c379693fae2bcaab1c2906831a4e10a6c3bb44ee8b615bca1d221"): {},
				ton.MustParseAccountID("0:6ccd325a858c379693fae2bcaab1c2906831a4e10a6c3bb44ee8b615bca1d225"): {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountID := ton.MustParseAccountID(tt.accountID)
			cli, err := liteapi.NewClient(liteapi.Testnet())
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			state, err := cli.GetAccountState(context.Background(), accountID)
			if err != nil {
				t.Fatalf("GetAccountState() error = %v", err)
			}
			extensions, err := GetW5R1ExtensionsList(state, 0)
			if err != nil {
				t.Fatalf("GetW5BetaExtensionsList() error = %v", err)
			}
			if !reflect.DeepEqual(extensions, tt.wantExtensions) {
				t.Errorf("GetW5BetaExtensionsList() \ngot  = %v\nwant = %v\n", extensions, tt.wantExtensions)
			}
		})
	}
}

func TestW5ActionsEmulatedExecutionOrder(t *testing.T) {
	privateKey, err := SeedToPrivateKey(RandomSeed())
	require.NoError(t, err)
	w, err := New(privateKey, V5R1, nil, WithNetworkGlobalID(MainnetGlobalID))
	require.NoError(t, err)
	stateInit, err := w.StateInit()
	require.NoError(t, err)
	walletState := tontest.Account().
		Address(w.GetAddress()).
		State(tlb.AccountActive).
		Balance(220_000_000).
		StateInit(&stateInit.Code.Value.Value, &stateInit.Data.Value.Value).
		MustShardAccount()

	recipients := []ton.AccountID{
		ton.MustParseAccountID("0:00000000000000000000000000000000000000000000000000000000000000a1"),
		ton.MustParseAccountID("0:00000000000000000000000000000000000000000000000000000000000000b2"),
		ton.MustParseAccountID("0:00000000000000000000000000000000000000000000000000000000000000c3"),
	}
	transfers := []Sendable{
		SimpleTransfer{Amount: 70_000_000, Address: recipients[0], Comment: "first"},
		SimpleTransfer{Amount: 80_000_000, Address: recipients[1], Comment: "second"},
		SimpleTransfer{Amount: 500_000_000, Address: recipients[2], Comment: "too much"},
	}

	body, err := w.CreateMessageBody(MessageConfig{
		Seqno:      0,
		ValidUntil: time.Now().Add(time.Minute),
		V5MsgType:  V5MsgTypeSignedExternal,
	}, transfers...)
	require.NoError(t, err)
	msg, err := ton.CreateExternalMessage(w.GetAddress(), body, nil, tlb.VarUInteger16{})
	require.NoError(t, err)
	accountSource := walletV5OrderTestAccountSource{
		accounts: map[ton.AccountID]tlb.ShardAccount{
			w.GetAddress(): walletState,
		},
	}
	tracer, err := txemulator.NewTraceBuilder(
		txemulator.WithAccountsSource(accountSource),
		txemulator.WithLimit(10),
	)
	require.NoError(t, err)
	tree, err := tracer.Run(context.Background(), msg)
	require.NoError(t, err)

	action := tree.TX.Description.TransOrd.Action
	require.True(t, action.Exists)
	actionPhase := action.Value.Value
	assert.Equal(t, uint16(len(transfers)), actionPhase.TotActions)
	assert.Equal(t, uint16(1), actionPhase.SkippedActions)
	assert.Equal(t, uint16(2), actionPhase.MsgsCreated)
	require.Len(t, tree.Children, 2)
	for i, child := range tree.Children {
		msg := child.TX.Msgs.InMsg.Value.Value
		gotDest, err := ton.AccountIDFromTlb(msg.Info.IntMsgInfo.Dest)
		if assert.NoError(t, err) {
			assert.NotNil(t, gotDest)
			assert.Equal(t, recipients[i], *gotDest)
			wantAmount := transfers[i].(SimpleTransfer).Amount
			assert.Equal(t, wantAmount, msg.Info.IntMsgInfo.Value.Grams)
		}
	}
}

func TestW5ActionsMarshalUnmarshalLogicalOrder(t *testing.T) {
	actions := W5ActionList{
		{Mode: 1, Msg: mustW5ActionMsgCell(t, 1)},
		{Mode: 2, Msg: mustW5ActionMsgCell(t, 2)},
		{Mode: 3, Msg: mustW5ActionMsgCell(t, 3)},
	}

	cell := boc.NewCell()
	require.NoError(t, tlb.Marshal(cell, actions))

	var got W5ActionList
	require.NoError(t, tlb.Unmarshal(cell.CopyCell(), &got))
	requireW5ActionsEqual(t, actions, got)
}

func mustW5ActionMsgCell(t *testing.T, value uint64) *boc.Cell {
	t.Helper()
	cell := boc.NewCell()
	require.NoError(t, cell.WriteUint(value, 32))
	return cell
}

func requireW5ActionsEqual(t *testing.T, want, got W5ActionList) {
	t.Helper()
	require.Len(t, got, len(want))
	for i := range want {
		require.Equal(t, want[i].Mode, got[i].Mode)
		wantHash, err := want[i].Msg.HashString()
		require.NoError(t, err)
		gotHash, err := got[i].Msg.HashString()
		require.NoError(t, err)
		require.Equal(t, wantHash, gotHash)
	}
}

type walletV5OrderTestAccountSource struct {
	accounts map[ton.AccountID]tlb.ShardAccount
}

func (s walletV5OrderTestAccountSource) GetAccountState(_ context.Context, account ton.AccountID) (tlb.ShardAccount, error) {
	if state, ok := s.accounts[account]; ok {
		return state, nil
	}
	return tontest.Account().Address(account).MustShardAccount(), nil
}

func (s walletV5OrderTestAccountSource) GetLibraries(_ context.Context, _ []ton.Bits256) (map[ton.Bits256]*boc.Cell, error) {
	return nil, nil
}

func (s walletV5OrderTestAccountSource) GetAllShardsInfo(_ context.Context, _ ton.BlockIDExt) ([]ton.BlockIDExt, error) {
	return []ton.BlockIDExt{{
		BlockID: ton.BlockID{
			Workchain: 0,
			Shard:     0x8000000000000000,
			Seqno:     1,
		},
	}}, nil
}

func (s walletV5OrderTestAccountSource) GetMasterchainInfo(_ context.Context) (liteclient.LiteServerMasterchainInfoC, error) {
	return liteclient.LiteServerMasterchainInfoC{
		Last: liteclient.TonNodeBlockIdExtC{Seqno: 1},
	}, nil
}
