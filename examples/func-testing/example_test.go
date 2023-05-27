package func_testing

import (
	"context"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/code"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/txemulator"
	"testing"
)

const SOURCE_CODE = `
(int) load_data() inline {
    var ds = get_data().begin_parse();
    return (ds~load_uint(64));
}

() save_data(int counter) impure inline {
    set_data(begin_cell()
            .store_uint(counter, 64)
            .end_cell());
}

() recv_internal(int msg_value, cell in_msg, slice in_msg_body) impure {
    int op = in_msg_body~load_uint(32);
    var addr = in_msg_body~load_msg_addr();
    var (counter) = load_data();
    if (op != 1) {
        throw(32);
    }
	if (counter > 4) {
		throw(33);
    }
    var msg = begin_cell()
            .store_uint(0x10, 6)
            .store_slice(addr)
            .store_coins(100000000)
            .store_uint(0, 1 + 4 + 4 + 64 + 32 + 1 + 1);
    send_raw_message(msg.end_cell(), 3);
    var msg = begin_cell()
            .store_uint(0x10, 6)
            .store_slice(my_address())
            .store_coins(1)
            .store_uint(0, 1 + 4 + 4 + 64 + 32 + 1 + 1)
            .store_uint(1,32)
            .store_slice(addr);
    send_raw_message(msg.end_cell(), 128);

    save_data(counter + 1);
}

int counter() method_id {
    var (counter) = load_data();
    return counter;
}

() recv_external(slice in_msg) impure {
	accept_message();
    int op = in_msg~load_uint(32);
    var addr = in_msg~load_msg_addr();

    var (counter) = load_data();
    if (counter > 10) {
        throw(30);
    }
    if (op != 1) {
        throw(31);
    }
    var msg = begin_cell()
            .store_uint(0x18, 6)
            .store_slice(my_address())
            .store_coins(1)
            .store_uint(0, 1 + 4 + 4 + 64 + 32 + 1 + 1)
            .store_uint(1,32)
            .store_slice(addr);
    send_raw_message(msg.end_cell(), 128);
    save_data(counter + 1);
}
`

func TestExample(t *testing.T) {
	ctx := context.Background()
	compiler := code.NewFunCCompiler()
	client, err := liteapi.NewClientWithDefaultTestnet()
	if err != nil {
		t.Fatal(err)
	}

	_, codeBoc, err := compiler.Compile(map[string]string{"main.fc": SOURCE_CODE})
	if err != nil {
		t.Fatal(err)
	}

	payTo := tongo.MustParseAccountID("0:ea27294687663e7c460be08e5bcef9341ae3ba6a1e2c4b0448b2f51e5eb45298")
	oldDestinationState, err := client.GetAccountState(ctx, payTo)

	initData := struct {
		Counter uint64
	}{1}
	body := struct {
		OP            uint32
		PayoutAddress tlb.MsgAddress
	}{
		OP:            1,
		PayoutAddress: payTo.ToMsgAddress(),
	}
	external, a, err := externalMessageToNewAccount(codeBoc, initData, body)
	if err != nil {
		t.Fatal(err)
	}
	tracer, err := txemulator.NewTraceBuilder(
		txemulator.WithAccountsSource(txemulator.NewAccountGetterMixin(client, map[tongo.AccountID]tlb.ShardAccount{
			a: fakeUninitAccount(a, tongo.OneTON),
		})),
	)
	if err != nil {
		t.Fatal(err)
	}
	result, err := tracer.Run(ctx, external)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Children) != 1 {
		t.Fatal("not enough transactions")
	}
	if tracer.FinalStates()[payTo].Account.Account.Storage.Balance.Grams < oldDestinationState.Account.Account.Storage.Balance.Grams+200000000 {
		t.Fatal("invalid result balance")
	}
}
