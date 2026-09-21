package abi

import (
	"encoding/hex"
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

// Every message below is a real one taken from mainnet, one per Hipo op declared in
// hipo_finance.xml that this file covers.
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

// request_loan has had three shapes. borrower_reward_share widened from 8 bits to 16 on
// 2026-09-05, and on 2026-09-21 it left the message altogether: the treasury sets one share
// for every loan in a round and refuses a request that names its own. Explorers reclassify
// history, so all three have to decode -- a request from an earlier era would otherwise turn
// back into raw hex. They are told apart by which reading consumes the body exactly, so none
// can be mistaken for another.
//
// The two older messages are real, one from each side of the widening; the older carries a
// share of 8 on the old scale out of 255, the same bid as 2056 on the new one. The third is
// the second one with its 16 share bits removed, built with this package rather than taken
// from a block: the release landed at 14:19:35 UTC on 2026-09-21, after the last request of
// that round, so the first request in this layout will be the next round's. Replace it with
// a real message once one exists.
func TestHipoFinanceRequestLoanEveryEra(t *testing.T) {
	for _, c := range []struct {
		name string
		body string
		want MsgOpName
	}{
		{"sent 2026-09-04, before the widening", "B5EE9C724101030100AE00013C36335DA9000000006A9ACF486A9B4F087038D7EA4C6800055D21DBA000080101903FB17DF20664C4E5A9C28B1DF3FE159CED3B01E67D7858CA2E9AA022DD2F94F16A9B4F080004800020791D63C50ED91D13E2B1F5D68589D5CC9A10B81AB775C3543BDFB9AC74F00C0200806DE2E5754D83CAFF1ED8654CA90D51E9D8224DFAF34E70B61BDEDD2D3FB255AD35333D0CBC4205D911B621782D22A47D4675C3EA31E3B07A3A7C371295E48605BA9A9BF8", HipoFinanceRequestLoanV1MsgOp},
		{"sent after the widening", "B5EE9C720101030100AF00013E36335DA9000000006AA5512B6AA64F087016BFFF2242C0052D7FFFD14007070101901781DB04A92CB9444B3043E9C8DB1BB6FC23B7DA38B9E7902A8ECCDEA2632F796AA64F08000300000F130A89779CC0D6E9CBA7571982050AF2AD2C59206879E5762DC21235D3BA66020080BFF1F56BB51083DB7B22CF827282C6D759FEBE75B5D011DEC36EF29933D7AC763757B03237DA9F88F85982223F4F91C9801C33631536342937184DD93D46FE06", HipoFinanceRequestLoanV2MsgOp},
		{"the layout in force since 2026-09-21, which names no share", "B5EE9C720101030100AD00013A36335DA9000000006AA5512B6AA64F087016BFFF2242C0052D7FFFD1400101901781DB04A92CB9444B3043E9C8DB1BB6FC23B7DA38B9E7902A8ECCDEA2632F796AA64F08000300000F130A89779CC0D6E9CBA7571982050AF2AD2C59206879E5762DC21235D3BA66020080BFF1F56BB51083DB7B22CF827282C6D759FEBE75B5D011DEC36EF29933D7AC763757B03237DA9F88F85982223F4F91C9801C33631536342937184DD93D46FE06", HipoFinanceRequestLoanMsgOp},
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
