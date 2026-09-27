package abi

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// Every message below is a real one taken from mainnet, one per Hipo op declared in
// hipo_finance.xml that this file covers. The one exception is send_unstake_all, which has
// never been sent and is built by its own test.
//
// Each case also asserts that decoding consumed the whole body. That is the point of the
// test rather than a detail: a field declared at the wrong width still decodes, it just
// reads the wrong bits and reports a plausible-looking number. finish_participation
// declared with a uint32 query_id parsed happily and reported the low half of the query id
// as the round, leaving 32 bits unread. Nothing but the leftover bits gives that away.
func TestHipoFinanceMessages(t *testing.T) {
	cases := []struct {
		name string
		ext  bool
		body string
		want MsgOpName
	}{
		{name: "participate_in_election", ext: true, body: "B5EE9C72010101010012000020574A297B000000006A9E2B966A9E4F08", want: HipoFinanceParticipateInElectionExtInMsgOp},
		{name: "vset_changed", ext: true, body: "B5EE9C720101010100120000202F0B5B3B000000006AA24F296AA14F08", want: HipoFinanceVsetChangedExtInMsgOp},
		{name: "finish_participation", ext: true, body: "B5EE9C7201010101001200002023274435000000006A95CF676A944F08", want: HipoFinanceFinishParticipationExtInMsgOp},
		{name: "decide_loan_requests", ext: false, body: "B5EE9C720101010100120000206A31D344000000006A9E2B966A9E4F08", want: HipoFinanceDecideLoanRequestsMsgOp},
		{name: "process_loan_requests", ext: false, body: "B5EE9C72010101010012000020071D07CC000000006A9E2B966A9E4F08", want: HipoFinanceProcessLoanRequestsMsgOp},
		{name: "recover_stakes", ext: false, body: "B5EE9C720101010100120000204F173D3E000000006A95CF676A944F08", want: HipoFinanceRecoverStakesMsgOp},
		{name: "proxy_new_stake", ext: false, body: "B5EE9C7201010301009C000118089CD4D0000000006A9E2B960101904539B3CE286825E1027B3558C95D016AA774DD24C5E6238EA57F967C0937595D6A9E4F0800030000A86F0B341C62D351E454B8ABBD6A6A041FC9A901F795F5AD307E6A2EF5AF029F020080C7F9BA5F103E7ABF9F9A328C7390739F50DFFD83F6530561945C489CFAD46438AE1D1D68B38021C686AA2BBD814004D32F8E71F76B3AD4BA1E320AB874E70F00", want: HipoFinanceProxyNewStakeMsgOp},
		{name: "proxy_recover_stake", ext: false, body: "B5EE9C7201010101000E000018407CB243000000006A95CF67", want: HipoFinanceProxyRecoverStakeMsgOp},
		{name: "request_rejected", ext: false, body: "B5EE9C7201010101000E000018CD0F2116000000006AA92B87", want: HipoFinanceRequestRejectedMsgOp},
		{name: "take_profit", ext: false, body: "B5EE9C7201010101000E0000188B556813000000006A95CF67", want: HipoFinanceTakeProfitMsgOp},
		{name: "take_borrower_fee", ext: false, body: "B5EE9C7201010101000E0000185E2D81F4000000006A9FCF60", want: HipoFinanceTakeBorrowerFeeMsgOp},
		{name: "withdraw_surplus (return_excess named)", ext: false, body: hipoWithdrawSurplusNamed, want: HipoFinanceWithdrawSurplusMsgOp},
		{name: "withdraw_surplus (return_excess addr_none)", ext: false, body: "B5EE9C7241010101000F00001923355FFB000000000000000020C0F0D6E5", want: HipoFinanceWithdrawSurplusMsgOp},
		{name: "proxy_unstake_all", ext: false, body: hipoProxyUnstakeAll, want: HipoFinanceProxyUnstakeAllMsgOp},
		{name: "unstake_all", ext: false, body: "B5EE9C7241010101000E0000185AE3014800000000000000001A99794C", want: HipoFinanceUnstakeAllMsgOp},
		{name: "request_loan (2026-09-05 era, names a share)", ext: false, body: "B5EE9C720101030100AF00013E36335DA9000000006AA5512B6AA64F087016BFFF2242C0052D7FFFD14007070101901781DB04A92CB9444B3043E9C8DB1BB6FC23B7DA38B9E7902A8ECCDEA2632F796AA64F08000300000F130A89779CC0D6E9CBA7571982050AF2AD2C59206879E5762DC21235D3BA66020080BFF1F56BB51083DB7B22CF827282C6D759FEBE75B5D011DEC36EF29933D7AC763757B03237DA9F88F85982223F4F91C9801C33631536342937184DD93D46FE06", want: HipoFinanceRequestLoanV2MsgOp},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			raw, err := hex.DecodeString(c.body)
			if err != nil {
				t.Fatal(err)
			}
			cells, err := boc.DeserializeBoc(raw)
			if err != nil {
				t.Fatal(err)
			}
			var op *MsgOpName
			if c.ext {
				_, op, _, err = ExtInMessageDecoder(cells[0], nil)
			} else {
				_, op, _, err = InternalMessageDecoder(cells[0], nil)
			}
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if op == nil {
				t.Fatalf("not decoded, so the op is still raw hex to every caller")
			}
			if *op != c.want {
				t.Fatalf("got %v, want %v", *op, c.want)
			}
			if left := cells[0].BitsAvailableForRead(); left != 0 {
				t.Errorf("%d bits left unread: a field is declared narrower than the contract writes it", left)
			}
		})
	}
}

// Two field widths in this schema are worth asserting rather than eyeballing.
func TestHipoFinanceFieldWidths(t *testing.T) {
	decode := func(t *testing.T, body string, ext bool) any {
		t.Helper()
		raw, err := hex.DecodeString(body)
		if err != nil {
			t.Fatal(err)
		}
		cells, err := boc.DeserializeBoc(raw)
		if err != nil {
			t.Fatal(err)
		}
		var v any
		if ext {
			_, _, v, err = ExtInMessageDecoder(cells[0], nil)
		} else {
			_, _, v, err = InternalMessageDecoder(cells[0], nil)
		}
		if err != nil {
			t.Fatal(err)
		}
		return v
	}

	// borrower_reward_share is 16 bits out of 65535, not 8 out of 255. Declared as uint8 it
	// still parses - it just yields the high byte, so this borrower's 1799 reads as 7.
	loan, ok := decode(t, "B5EE9C720101030100AF00013E36335DA9000000006AA5512B6AA64F087016BFFF2242C0052D7FFFD14007070101901781DB04A92CB9444B3043E9C8DB1BB6FC23B7DA38B9E7902A8ECCDEA2632F796AA64F08000300000F130A89779CC0D6E9CBA7571982050AF2AD2C59206879E5762DC21235D3BA66020080BFF1F56BB51083DB7B22CF827282C6D759FEBE75B5D011DEC36EF29933D7AC763757B03237DA9F88F85982223F4F91C9801C33631536342937184DD93D46FE06", false).(HipoFinanceRequestLoanV2MsgBody)
	if !ok {
		t.Fatalf("request_loan did not decode to the body type of the era that names a share")
	}
	if uint64(loan.BorrowerRewardShare) != 1799 {
		t.Errorf("borrower_reward_share = %d, want 1799", loan.BorrowerRewardShare)
	}

	// finish_participation's query_id is 64 bits like every other op here, whatever
	// contracts/schema.tlb says: treasury.fc reads load_uint(64) then load_uint(32) then
	// end_parse(). Declared as uint32 it reports the low half of the query id as the round.
	finish, ok := decode(t, "B5EE9C7201010101001200002023274435000000006A95CF676A944F08", true).(HipoFinanceFinishParticipationExtInMsgBody)
	if !ok {
		t.Fatalf("finish_participation did not decode to its body type")
	}
	if finish.RoundSince != 1788104456 {
		t.Errorf("round_since = %d, want 1788104456", finish.RoundSince)
	}
}

// request_loan has had four shapes. borrower_reward_share widened from 8 bits to 16 on
// 2026-09-05, and on 2026-09-21 it left the message altogether: the treasury sets one share
// for every loan in a round and refuses a request that names its own. Hipo's stake-cap
// release then made max_stake, a coins field, required after min_payment. Explorers
// reclassify history, so every shape has to decode: a request from an earlier era would
// otherwise turn back into raw hex. They are told apart by which reading consumes the body
// exactly, so none can be mistaken for another.
//
// All five are real messages; the oldest carries a share of 8 on the old scale out of 255, the
// same bid as 2056 on the new one. The last two carry max_stake, one of them 0 (no cap, the
// shortest the field can be).
func TestHipoFinanceRequestLoanEveryEra(t *testing.T) {
	for _, c := range []struct {
		name string
		body string
		want MsgOpName
	}{
		{"sent 2026-09-04, before the widening", "B5EE9C724101030100AE00013C36335DA9000000006A9ACF486A9B4F087038D7EA4C6800055D21DBA000080101903FB17DF20664C4E5A9C28B1DF3FE159CED3B01E67D7858CA2E9AA022DD2F94F16A9B4F080004800020791D63C50ED91D13E2B1F5D68589D5CC9A10B81AB775C3543BDFB9AC74F00C0200806DE2E5754D83CAFF1ED8654CA90D51E9D8224DFAF34E70B61BDEDD2D3FB255AD35333D0CBC4205D911B621782D22A47D4675C3EA31E3B07A3A7C371295E48605BA9A9BF8", HipoFinanceRequestLoanV1MsgOp},
		{"sent after the widening", "B5EE9C720101030100AF00013E36335DA9000000006AA5512B6AA64F087016BFFF2242C0052D7FFFD14007070101901781DB04A92CB9444B3043E9C8DB1BB6FC23B7DA38B9E7902A8ECCDEA2632F796AA64F08000300000F130A89779CC0D6E9CBA7571982050AF2AD2C59206879E5762DC21235D3BA66020080BFF1F56BB51083DB7B22CF827282C6D759FEBE75B5D011DEC36EF29933D7AC763757B03237DA9F88F85982223F4F91C9801C33631536342937184DD93D46FE06", HipoFinanceRequestLoanV2MsgOp},
		{"sent 2026-09-26, naming no share", "B5EE9C720101030100AD00013A36335DA9000000006AB72B706AB74F087031471B7D5E5F357F000000000101908BEFC2D391D594DB74129D7CF454338E8E2103D24A34A04F9D04F39E116C67726AB74F08000480000A23E5F061401A6CAF0B6F0319DE100CF7F7BB10398C0E57B2C8F3264818A97402008012DDE6DE3DFEB85DE761EC71BBACD92935E752AB2C1900B01B22C445B5F8A79F5196E06A469444948122E88647BDC2D5679863296206DCEBCADBCDE507D0C703", HipoFinanceRequestLoanV3MsgOp},
		{"sent 2026-09-26 14:37, max_stake 0", hipoRequestLoanNoCap, HipoFinanceRequestLoanMsgOp},
		{"sent 2026-09-26 09:39, max_stake 3067441.598 GRAM", hipoRequestLoanCapped, HipoFinanceRequestLoanMsgOp},
	} {
		t.Run(c.name, func(t *testing.T) {
			raw, err := hex.DecodeString(c.body)
			if err != nil {
				t.Fatal(err)
			}
			cells, err := boc.DeserializeBoc(raw)
			if err != nil {
				t.Fatal(err)
			}
			_, op, _, err := InternalMessageDecoder(cells[0], nil)
			if err != nil {
				t.Fatal(err)
			}
			if op == nil {
				t.Fatal("not decoded, so a request of this era is raw hex to every caller")
			}
			if *op != c.want {
				t.Fatalf("got %v, want %v", *op, c.want)
			}
			if left := cells[0].BitsAvailableForRead(); left != 0 {
				t.Errorf("%d bits left unread", left)
			}
		})
	}
}

// Two real requests sent 2026-09-26, after the stake-cap release: one at 14:37 UTC with no cap
// (max_stake 0), and one at 09:39 UTC carrying a cap.
const (
	hipoRequestLoanNoCap  = "B5EE9C720101030100AE00013B36335DA9000000006AB7D8906AB84F08704CBD15E7260005D78000000008010190E05B060871E5AF3693D58C9903E18B09A90063607E1B201331C2E6B7A3DAFBED6AB84F0800030000A86F0B341C62D351E454B8ABBD6A6A041FC9A901F795F5AD307E6A2EF5AF029F020080E1A5826D807E98351C0EE2F87912393D37834E72EC1B5D1183845CF329FEDCC10B83FF120630E50349BE81AAD8CC994CD7030D680BE7BA28E97BCCC1056C6706"
	hipoRequestLoanCapped = "B5EE9C720101030100B500014936335DA9000000006AB792B46AB84F087062B759BD8C3A85FEC000000070AE5D266D86D72801019090A730060B00B26CCC717450249B0FA9C0B66588678BB7111A211C205DC7FB496AB84F08000480000A23E5F061401A6CAF0B6F0319DE100CF7F7BB10398C0E57B2C8F3264818A974020080AC838AF88782854EEFA3718AD2FC4943ACFEB6E613A052CF80C8746488DBD7A0BDBD3F7E10EC4BFC5737313CAB5481B344164F14A12432F285E3AE12E492120B"
)

// max_stake is read as the value that was sent, and every field before it keeps its place.
func TestHipoFinanceRequestLoanMaxStake(t *testing.T) {
	for _, c := range []struct {
		body                       string
		round                      uint32
		loan, minPayment, maxStake string
	}{
		{hipoRequestLoanNoCap, 1790463752, "1350000000000000", "925565452288", "0"},
		{hipoRequestLoanCapped, 1790463752, "1736633986106280", "1094142918656", "3067441598459250"},
	} {
		raw, err := hex.DecodeString(c.body)
		if err != nil {
			t.Fatal(err)
		}
		cells, err := boc.DeserializeBoc(raw)
		if err != nil {
			t.Fatal(err)
		}
		_, _, body, err := InternalMessageDecoder(cells[0], nil)
		if err != nil {
			t.Fatal(err)
		}
		loan, ok := body.(HipoFinanceRequestLoanMsgBody)
		if !ok {
			t.Fatalf("request_loan with max_stake did not decode to its body type")
		}
		for name, got := range map[string]string{
			"max_stake":   bigString(loan.MaxStake),
			"loan_amount": bigString(loan.LoanAmount),
			"min_payment": bigString(loan.MinPayment),
		} {
			want := map[string]string{"max_stake": c.maxStake, "loan_amount": c.loan, "min_payment": c.minPayment}[name]
			if got != want {
				t.Errorf("%v = %v, want %v", name, got, want)
			}
		}
		if loan.RoundSince != c.round {
			t.Errorf("round_since = %v, want %v", loan.RoundSince, c.round)
		}
	}
}

// The unstake-all chain and withdraw_surplus, real messages from mainnet. The unstake-all pair
// was sent 2026-09-18 20:29 UTC for a holder who staked out by the comment "w": the treasury
// asks the parent for the owner's wallet (proxy_unstake_all), and the parent tells that wallet
// to burn everything (unstake_all). The named withdraw_surplus was sent to the treasury on
// 2026-04-28 and returns the excess to its own sender.
const (
	hipoProxyUnstakeAll      = "B5EE9C7241010101003000005B76BD27600000000000000000800AE312CA6033CE18AEA03884FB7B4C3EE66C9B77BFBB7A7ED5BE143646DB4FED10B57BFFDC"
	hipoWithdrawSurplusNamed = "B5EE9C7241010101003000005B23355FFB00000000000000008014C584C0176A7C32C8DEEE2841418B807482F9FB5EB86146284A4D05314C123DF0E01A8118"
)

// The addresses these carry are the ones sent: the owner whose wallet is unstaked, and where
// the surplus goes.
func TestHipoFinanceUnstakeAllAndSurplusAddresses(t *testing.T) {
	for _, c := range []struct {
		body, want string
		addr       func(any) (tlb.MsgAddress, bool)
	}{
		{hipoProxyUnstakeAll, "0:57189653019e70c57501c427dbda61f73364dbbdfddbd3f6adf0a1b236da7f68",
			func(b any) (tlb.MsgAddress, bool) { v, ok := b.(HipoFinanceProxyUnstakeAllMsgBody); return v.Owner, ok }},
		{hipoWithdrawSurplusNamed, "0:a62c2600bb53e19646f771420a0c5c03a417cfdaf5c30a31425268298a6091ef",
			func(b any) (tlb.MsgAddress, bool) {
				v, ok := b.(HipoFinanceWithdrawSurplusMsgBody)
				return v.ReturnExcess, ok
			}},
	} {
		raw, err := hex.DecodeString(c.body)
		if err != nil {
			t.Fatal(err)
		}
		cells, err := boc.DeserializeBoc(raw)
		if err != nil {
			t.Fatal(err)
		}
		_, _, body, err := InternalMessageDecoder(cells[0], nil)
		if err != nil {
			t.Fatal(err)
		}
		addr, ok := c.addr(body)
		if !ok {
			t.Fatalf("decoded to %T", body)
		}
		got, err := ton.AccountIDFromTlb(addr)
		if err != nil || got == nil {
			t.Fatalf("address did not decode: %v", err)
		}
		if got.ToRaw() != c.want {
			t.Errorf("address = %v, want %v", got.ToRaw(), c.want)
		}
	}
}

// send_unstake_all has never been sent on mainnet: the dapp and the comment "w" reach the same
// flow another way. So this one is built, from the layout treasury.fc reads -- a query_id and
// nothing else -- rather than taken from the chain.
func TestHipoFinanceSendUnstakeAll(t *testing.T) {
	c := boc.NewCell()
	if err := c.WriteUint(0x45baeda9, 32); err != nil {
		t.Fatal(err)
	}
	if err := c.WriteUint(1790463752, 64); err != nil {
		t.Fatal(err)
	}
	c.ResetCounters()
	_, op, body, err := InternalMessageDecoder(c, nil)
	if err != nil || op == nil || *op != HipoFinanceSendUnstakeAllMsgOp {
		t.Fatalf("got %v, %v", op, err)
	}
	if q := body.(HipoFinanceSendUnstakeAllMsgBody).QueryId; q != 1790463752 {
		t.Errorf("query_id = %v", q)
	}
	if left := c.BitsAvailableForRead(); left != 0 {
		t.Errorf("%d bits left unread", left)
	}
}

func bigString(v tlb.VarUInteger16) string {
	b := big.Int(v)
	return b.String()
}
