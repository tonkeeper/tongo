package func_testing

import (
	"context"
	"testing"

	"github.com/tonkeeper/tongo/code"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/tontest"
	"github.com/tonkeeper/tongo/txemulator"
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
	t.Skip("fixing")

	ctx := context.Background()
	compiler := code.NewFunCCompiler()

	_, codeBoc, err := compiler.Compile(map[string]string{"main.fc": SOURCE_CODE})
	if err != nil {
		t.Fatal(err)
	}

	payTo := ton.MustParseAccountID("0:ea27294687663e7c460be08e5bcef9341ae3ba6a1e2c4b0448b2f51e5eb45298")

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
	m := tontest.NewMessage().
		WithBody(tontest.MustAnyToCell(body)).
		WithInit(tontest.MustAnyToCell(codeBoc), tontest.MustAnyToCell(initData)).
		MustMessage()

	testAccount := tontest.Account().
		StateInit(tontest.MustAnyToCell(codeBoc), tontest.MustAnyToCell(initData)).
		Balance(ton.OneTON).
		MustShardAccount()

	tracer, err := txemulator.NewTraceBuilder(
		txemulator.WithAccounts(testAccount),
		txemulator.WithTestnet(),
	)
	if err != nil {
		t.Fatal(err)
	}
	result, err := tracer.Run(ctx, m)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Children) != 1 {
		t.Fatal("not enough transactions")
	}

}
