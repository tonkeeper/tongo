package tvm

// #cgo darwin CFLAGS: -I ../lib/
// #cgo darwin,arm64 LDFLAGS: -L ../lib/darwin/arm64 -Wl,-rpath,../lib/darwin/arm64 -l vm-exec-lib
// #cgo darwin,x86_64 LDFLAGS: -L ../lib/darwin/arm64 -Wl,-rpath,../lib/darwin/x86 -l vm-exec-lib
// #cgo linux LDFLAGS: -L ../lib/linux/ -Wl,-rpath,../lib/linux/ -l vm-exec-lib
// #include "../lib/libvm-exec-lib.h"
// #include <stdlib.h>
import "C"
import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
	"github.com/startfellows/tongo/utils"
	"time"
	"unsafe"
)

type tvmExecutionResultInternal struct {
	ExitCode       int          `json:"exit_code"`
	GasConsumed    int          `json:"gas_consumed"`
	DataCell       string       `json:"data_cell"`
	ActionListCell string       `json:"action_list_cell"`
	Logs           string       `json:"logs"`
	Stack          []StackEntry `json:"stack"`
	Ok             bool         `json:"ok"`
	Error          string       `json:"error"`
}

type tvmExecConfig struct {
	FunctionSelector int          `json:"function_selector"`
	InitStack        []StackEntry `json:"init_stack"`
	Code             string       `json:"code"`
	Data             string       `json:"data"`
	C7Register       StackEntry   `json:"c7_register"`
}

type ExecutionResult struct {
	ExitCode       int
	GasConsumed    int
	DataCell       *boc.Cell
	ActionListCell *boc.Cell
	Logs           string
	Stack          []StackEntry
}

func getVMFunctionSelector(name string) int {
	if name == "main" {
		return 0
	} else if name == "recv_internal" {
		return 0
	} else if name == "recv_external" {
		return -1
	} else {
		return int(utils.Crc16String(name)&0xffff) | 0x10000
	}
}

func buildDefaultC7Register(address *tongo.AccountID) (StackEntry, error) {
	now := int(time.Now().Unix())

	balance := NewTupleStackEntry([]StackEntry{
		NewIntStackEntry(1000),
		NewNullStackEntry(),
	})
	var addrStack StackEntry
	if address != nil {
		addrCell := boc.NewCell()
		err := tlb.Marshal(addrCell, address)
		if err != nil {
			return StackEntry{}, err
		}
		addrStack = NewCellSliceStackEntry(addrCell)
	} else {
		addrStack = NewNullStackEntry()
	}

	return NewTupleStackEntry([]StackEntry{
		NewTupleStackEntry([]StackEntry{
			NewIntStackEntry(0x076ef1ea), // [ magic:0x076ef1ea
			NewIntStackEntry(0),          // actions:Integer
			NewIntStackEntry(0),          // msgs_sent:Integer
			NewIntStackEntry(now),        // unixtime:Integer
			NewIntStackEntry(now),        // block_lt:Integer
			NewIntStackEntry(now),        // trans_lt:Integer
			NewIntStackEntry(now),        // rand_seed:Integer
			balance,                      // balance_remaining:[Integer (Maybe Cell)]
			addrStack,                    // myself:MsgAddressInt
			NewNullStackEntry(),          // global_config:(Maybe Cell) ] = SmartContractInfo;
		}),
	}), nil
}

func RunTvm(code *boc.Cell, data *boc.Cell, funcName string, args []StackEntry, destAccount *tongo.AccountID) (ExecutionResult, error) {
	codeBoc, err := code.ToBocBase64Custom(false, true, false, 0)
	if err != nil {
		return ExecutionResult{}, err
	}

	dataBoc, err := data.ToBocBase64Custom(false, true, false, 0)
	if err != nil {
		return ExecutionResult{}, err
	}

	register, err := buildDefaultC7Register(destAccount)
	if err != nil {
		return ExecutionResult{}, err
	}
	config := tvmExecConfig{
		FunctionSelector: getVMFunctionSelector(funcName),
		InitStack:        args,
		Code:             codeBoc,
		Data:             dataBoc,
		C7Register:       register,
	}

	configStr, err := json.Marshal(config)
	if err != nil {
		return ExecutionResult{}, err
	}
	CconfigStr := C.CString(string(configStr))
	defer C.free(unsafe.Pointer(CconfigStr))
	res := C.vm_exec(C.int(len(string(configStr))), CconfigStr)
	resJSON := C.GoString(res)
	defer C.free(unsafe.Pointer(res))
	var executeResult tvmExecutionResultInternal
	err = json.Unmarshal([]byte(resJSON), &executeResult)
	if err != nil {
		return ExecutionResult{}, err
	}

	if !executeResult.Ok {
		return ExecutionResult{}, errors.New(executeResult.Error)
	}

	if executeResult.ExitCode != 0 {
		logs, err := base64.StdEncoding.DecodeString(executeResult.Logs)
		if err != nil {
			return ExecutionResult{}, err
		}

		result := ExecutionResult{
			ExitCode:    executeResult.ExitCode,
			GasConsumed: executeResult.GasConsumed,
			Logs:        string(logs),
			Stack:       []StackEntry{},
		}
		return result, nil
	}

	dataCell, err := boc.DeserializeBocBase64(executeResult.DataCell)
	if err != nil {
		return ExecutionResult{}, err
	}
	actionListCell, err := boc.DeserializeBocBase64(executeResult.ActionListCell)
	if err != nil {
		return ExecutionResult{}, err
	}
	logs, err := base64.StdEncoding.DecodeString(executeResult.Logs)
	if err != nil {
		return ExecutionResult{}, err
	}

	result := ExecutionResult{
		ExitCode:       executeResult.ExitCode,
		GasConsumed:    executeResult.GasConsumed,
		DataCell:       dataCell[0],
		ActionListCell: actionListCell[0],
		Logs:           string(logs),
		Stack:          executeResult.Stack,
	}

	return result, nil
}
