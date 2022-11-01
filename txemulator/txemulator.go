package txemulator

// #cgo linux LDFLAGS: -L ../lib/linux/ -Wl,-rpath,../lib/linux/ -l emulator
// #include "../lib/emulator-extern.h"
// #include <stdlib.h>
import "C"
import (
	"encoding/json"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
	"runtime"
	"unsafe"
)

type VerbosityLevel int

const (
	LogTruncated VerbosityLevel = iota
	LogUnlimited
	CellHashAndOffsetForCommand
	PrintsAllStackValuesForCommand
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
	VmLog        string `json:"vm_log"`
	VmExitCode   int    `json:"vm_exit_code"`
}

type EmulationResult struct {
	Success   bool
	Emulation *struct {
		ShardAccount tongo.ShardAccount
		Transaction  tongo.Transaction
	}
	Logs  string
	Error *struct {
		ExitCode int
		Text     string
	}
}

// NewEmulator
// Verbosity level of VM log. 0 - log truncated to last 256 characters. 1 - unlimited length log.
// 2 - for each command prints its cell hash and offset. 3 - for each command log prints all stack values.
func NewEmulator(config *boc.Cell, verbosityLevel VerbosityLevel) (*Emulator, error) {
	configBoc, err := config.ToBocBase64()
	if err != nil {
		return nil, err
	}
	var libs tlb.HashmapE[struct{}] // empty shard libs dict
	libsStr, err := tlbStructToBase64(libs)
	if err != nil {
		return nil, err
	}
	cConfigStr := C.CString(configBoc)
	defer C.free(unsafe.Pointer(cConfigStr))
	cLibStr := C.CString(libsStr)
	defer C.free(unsafe.Pointer(cLibStr))
	level := C.int(verbosityLevel)
	e := Emulator{emulator: C.transaction_emulator_create(cConfigStr, cLibStr, level)}
	runtime.SetFinalizer(&e, destroy)
	return &e, nil
}

func (e *Emulator) Emulate(shardAccount tongo.ShardAccount, message tongo.Message) (EmulationResult, error) {
	msg, err := tlbStructToBase64(message)
	if err != nil {
		return EmulationResult{}, err
	}

	acc, err := tlbStructToBase64(shardAccount)
	if err != nil {
		return EmulationResult{}, err
	}

	cAccStr := C.CString(acc)
	defer C.free(unsafe.Pointer(cAccStr))
	cMsgStr := C.CString(msg)
	defer C.free(unsafe.Pointer(cMsgStr))

	r := C.transaction_emulator_emulate_transaction(e.emulator, cAccStr, cMsgStr)
	rJSON := C.GoString(r)
	defer C.free(unsafe.Pointer(r))
	var (
		res     result
		account tongo.ShardAccount
		tx      tongo.Transaction
	)
	err = json.Unmarshal([]byte(rJSON), &res)
	if err != nil {
		return EmulationResult{}, err
	}

	if res.Success == false {
		err1 := struct {
			ExitCode int
			Text     string
		}{
			ExitCode: res.VmExitCode,
			Text:     res.Error,
		}
		return EmulationResult{Success: false, Logs: res.VmLog, Error: &err1}, nil
	}

	accountCell, err := boc.DeserializeBocBase64(res.ShardAccount)
	if err != nil {
		return EmulationResult{}, err
	}
	err = tlb.Unmarshal(accountCell[0], &account)
	if err != nil {
		return EmulationResult{}, err
	}

	txCell, err := boc.DeserializeBocBase64(res.Transaction)
	if err != nil {
		return EmulationResult{}, err
	}
	err = tlb.Unmarshal(txCell[0], &tx)
	if err != nil {
		return EmulationResult{}, err
	}
	em := struct {
		ShardAccount tongo.ShardAccount
		Transaction  tongo.Transaction
	}{
		ShardAccount: account,
		Transaction:  tx,
	}
	return EmulationResult{Success: true, Logs: res.VmLog, Emulation: &em}, nil
}

func destroy(e *Emulator) {
	C.transaction_emulator_destroy(e.emulator)
}

// SetVerbosityLevel
// verbosity level (0 - never, 1 - error, 2 - warning, 3 - info, 4 - debug)
func (e *Emulator) SetVerbosityLevel(level int) {
	C.transaction_emulator_set_verbosity_level(C.int(level))
}

func tlbStructToBase64(s any) (string, error) {
	cell := boc.NewCell()
	err := tlb.Marshal(cell, s)
	if err != nil {
		return "", err
	}
	return cell.ToBocBase64()
}
