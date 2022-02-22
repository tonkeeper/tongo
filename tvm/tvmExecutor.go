package tvm

// #cgo darwin CFLAGS: -I ../lib/
// #cgo LDFLAGS: -L ../lib/ -Wl,-rpath,./lib -l vm-exec-lib
// #include "../lib/libvm-exec-lib.h"
import "C"
import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"tongo/boc"
)

type tvmExecutionResultInternal struct {
	ExitCode       int             `json:"exit_code"`
	GasConsumed    int             `json:"gas_consumed"`
	DataCell       string          `json:"data_cell"`
	ActionListCell string          `json:"action_list_cell"`
	Logs           string          `json:"logs"`
	Stack          []TvmStackEntry `json:"stack"`
}

type tvmExecConfig struct {
	FunctionSelector int             `json:"function_selector"`
	InitStack        []TvmStackEntry `json:"init_stack"`
	Code             string          `json:"code"`
	Data             string          `json:"data"`
	Time             string          `json:"time"`
}

type TvmExecutionResult struct {
	ExitCode       int
	GasConsumed    int
	DataCell       *boc.Cell
	ActionListCell *boc.Cell
	Logs           string
	Stack          []TvmStackEntry
}

func getVmFunctionSelector(name string) int {
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

func RunTvm(code *boc.Cell, data *boc.Cell, funcName string, args []TvmStackEntry, time int64) (TvmExecutionResult, error) {
	codeBoc, err := code.ToBocBase64Custom(false, true, false, 0)
	if err != nil {
		return TvmExecutionResult{}, err
	}

	dataBoc, err := data.ToBocBase64Custom(false, true, false, 0)
	if err != nil {
		return TvmExecutionResult{}, err
	}

	config := tvmExecConfig{
		FunctionSelector: getVmFunctionSelector(funcName),
		InitStack:        args,
		Code:             codeBoc,
		Data:             dataBoc,
		Time:             strconv.FormatInt(time, 10),
	}

	configStr, err := json.Marshal(config)
	if err != nil {
		return TvmExecutionResult{}, err
	}

	res := C.vm_exec(C.int(len(string(configStr))), C.CString(string(configStr)))
	resJson := C.GoString(res)

	var executeResult tvmExecutionResultInternal
	err = json.Unmarshal([]byte(resJson), &executeResult)
	if err != nil {
		return TvmExecutionResult{}, err
	}

	dataCell, err := boc.DeserializeBocBase64(executeResult.DataCell)
	if err != nil {
		return TvmExecutionResult{}, err
	}
	actionListCell, err := boc.DeserializeBocBase64(executeResult.ActionListCell)
	if err != nil {
		return TvmExecutionResult{}, err
	}
	logs, err := base64.StdEncoding.DecodeString(executeResult.Logs)
	if err != nil {
		return TvmExecutionResult{}, err
	}

	result := TvmExecutionResult{
		ExitCode:       executeResult.ExitCode,
		GasConsumed:    executeResult.GasConsumed,
		DataCell:       dataCell[0],
		ActionListCell: actionListCell[0],
		Logs:           string(logs),
		Stack:          executeResult.Stack,
	}

	return result, nil
}
