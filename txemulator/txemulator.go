package txemulator

//#cgo linux LDFLAGS: -L ../lib/linux/ -Wl,-rpath,../lib/linux/ -l emulator
//#include "../lib/emulator-extern.h"
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
	"runtime"
	"unsafe"
)

type Emulator struct {
	emulator unsafe.Pointer
}

// { "success": false, "error": "Error description" }
// { "success": true, "transaction": "Base64 encoded Transaction boc", "shard_account": "Base64 encoded ShardAccount boc" }
type result struct {
	Success      bool   `json:"success"`
	Error        string `json:"error"`
	Transaction  string `json:"transaction"`
	ShardAccount string `json:"shard_account"`
}

type EmulatorError struct {
	text string
}

func (e EmulatorError) Error() string {
	return e.text
}

func NewEmulator(config *boc.Cell) (Emulator, error) {
	configBoc, err := config.ToBocBase64()
	if err != nil {
		return Emulator{}, err
	}
	var libs tlb.HashmapE[struct{}] // empty shard libs dict
	libsStr, err := tlbStructToBase64(libs)
	if err != nil {
		return Emulator{}, err
	}
	e := Emulator{emulator: C.transaction_emulator_create(C.CString(configBoc), C.CString(libsStr))}
	runtime.SetFinalizer(&e, destroy)
	return e, nil
}

func (e *Emulator) Emulate(shardAccount tongo.ShardAccount, message tongo.Message[tlb.Any]) (tongo.ShardAccount, tongo.Transaction, error) {

	msg, err := tlbStructToBase64(message)
	if err != nil {
		return tongo.ShardAccount{}, tongo.Transaction{}, err
	}

	acc, err := tlbStructToBase64(shardAccount)
	if err != nil {
		return tongo.ShardAccount{}, tongo.Transaction{}, err
	}

	r := C.transaction_emulator_emulate_transaction(e.emulator, C.CString(acc), C.CString(msg))
	rJSON := C.GoString(r)
	var (
		res     result
		account tongo.ShardAccount
		tx      tongo.Transaction
	)
	err = json.Unmarshal([]byte(rJSON), &res)
	if err != nil {
		return tongo.ShardAccount{}, tongo.Transaction{}, err
	}

	if res.Success == false {
		return tongo.ShardAccount{}, tongo.Transaction{}, EmulatorError{text: fmt.Sprintf("message rejected: %v", res.Error)}
	}

	accountCell, err := boc.DeserializeBocBase64(res.ShardAccount)
	if err != nil {
		return tongo.ShardAccount{}, tongo.Transaction{}, err
	}
	err = tlb.Unmarshal(accountCell[0], &account)
	if err != nil {
		return tongo.ShardAccount{}, tongo.Transaction{}, err
	}

	txCell, err := boc.DeserializeBocBase64(res.Transaction)
	if err != nil {
		return tongo.ShardAccount{}, tongo.Transaction{}, err
	}
	err = tlb.Unmarshal(txCell[0], &tx)
	if err != nil {
		return tongo.ShardAccount{}, tongo.Transaction{}, err
	}

	return account, tx, nil
}

func destroy(e *Emulator) {
	C.transaction_emulator_destroy(e.emulator)
	e.emulator = nil
}

func tlbStructToBase64(s any) (string, error) {
	cell := boc.NewCell()
	err := tlb.Marshal(cell, s)
	if err != nil {
		return "", err
	}
	return cell.ToBocBase64()
}
