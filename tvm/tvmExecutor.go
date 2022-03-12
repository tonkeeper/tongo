package tvm

// #cgo darwin CFLAGS: -I ../lib/
// #cgo darwin,arm64 LDFLAGS: -L ../lib/darwin/arm64 -Wl,-rpath,../lib/darwin/arm64 -l vm-exec-lib
// #cgo darwin,x86_64 LDFLAGS: -L ../lib/darwin/arm64 -Wl,-rpath,../lib/darwin/x86 -l vm-exec-lib
// #cgo linux LDFLAGS: -L ../lib/linux/ -Wl,-rpath,../lib/linux/ -l vm-exec-lib
// #include "../lib/libvm-exec-lib.h"
import "C"
import (
	"encoding/base64"
	"encoding/json"
	"github.com/startfellows/tongo/boc"
	"time"
)

type tvmExecutionResultInternal struct {
	ExitCode       int          `json:"exit_code"`
	GasConsumed    int          `json:"gas_consumed"`
	DataCell       string       `json:"data_cell"`
	ActionListCell string       `json:"action_list_cell"`
	Logs           string       `json:"logs"`
	Stack          []StackEntry `json:"stack"`
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
		return int(Crc16String(name)&0xffff) | 0x10000
	}
}

func buildDefaultC7Register() StackEntry {
	now := int(time.Now().Unix())

	balance := NewTupleStackEntry([]StackEntry{
		NewIntStackEntry(1000),
		NewNullStackEntry(),
	})

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
			NewNullStackEntry(),          // myself:MsgAddressInt
			NewNullStackEntry(),          // global_config:(Maybe Cell) ] = SmartContractInfo;
		}),
	})
}

func RunTvm(code *boc.Cell, data *boc.Cell, funcName string, args []StackEntry) (ExecutionResult, error) {
	codeBoc, err := code.ToBocBase64Custom(false, true, false, 0)
	if err != nil {
		return ExecutionResult{}, err
	}

	dataBoc, err := data.ToBocBase64Custom(false, true, false, 0)
	if err != nil {
		return ExecutionResult{}, err
	}

	config := tvmExecConfig{
		FunctionSelector: getVMFunctionSelector(funcName),
		InitStack:        args,
		Code:             codeBoc,
		Data:             dataBoc,
		C7Register:       buildDefaultC7Register(),
	}

	configStr, err := json.Marshal(config)
	if err != nil {
		return ExecutionResult{}, err
	}

	res := C.vm_exec(C.int(len(string(configStr))), C.CString(string(configStr)))
	resJSON := C.GoString(res)

	var executeResult tvmExecutionResultInternal
	err = json.Unmarshal([]byte(resJSON), &executeResult)
	if err != nil {
		return ExecutionResult{}, err
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
