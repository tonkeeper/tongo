package tvm

// #cgo linux LDFLAGS: -L ../lib/linux/ -Wl,-rpath,../lib/linux/ -l emulator
// #cgo darwin CFLAGS: -I ../lib/
// #cgo darwin,arm64 LDFLAGS: -L ../lib/darwin/arm64 -Wl,-rpath,../lib/darwin/arm64 -l vm-exec-lib
// #cgo darwin,x86_64 LDFLAGS: -L ../lib/darwin/arm64 -Wl,-rpath,../lib/darwin/x86 -l vm-exec-lib
// #cgo linux LDFLAGS: -L ../lib/linux/ -Wl,-rpath,../lib/linux/ -l vm-exec-lib
// #include "../lib/libvm-exec-lib.h"
// #include "../lib/emulator-extern.h"
// #include <stdlib.h>
// #include <stdbool.h>
import "C"
import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
	"github.com/startfellows/tongo/txemulator"
	"github.com/startfellows/tongo/utils"
	"math/rand"
	"runtime"
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

type Emulator struct {
	emulator unsafe.Pointer
	config   *boc.Cell
	balance  uint64
}

// NewEmulator
// Verbosity level of VM log. 0 - log truncated to last 256 characters. 1 - unlimited length log.
// 2 - for each command prints its cell hash and offset. 3 - for each command log prints all stack values.
func NewEmulator(code, data, config *boc.Cell, balance int64, verbosityLevel txemulator.VerbosityLevel) (*Emulator, error) {
	codeBoc, err := code.ToBocBase64()
	if err != nil {
		return nil, err
	}
	dataBoc, err := data.ToBocBase64()
	if err != nil {
		return nil, err
	}
	cCodeStr := C.CString(codeBoc)
	defer C.free(unsafe.Pointer(cCodeStr))
	cDataStr := C.CString(dataBoc)
	defer C.free(unsafe.Pointer(cDataStr))
	level := C.int(verbosityLevel)
	e := Emulator{
		emulator: C.tvm_emulator_create(cCodeStr, cDataStr, level),
		config:   config,
		balance:  uint64(balance),
	}
	runtime.SetFinalizer(&e, destroy)
	return &e, nil
}

func destroy(e *Emulator) {
	C.tvm_emulator_destroy(e.emulator)
}

// SetVerbosityLevel
// verbosity level (0 - never, 1 - error, 2 - warning, 3 - info, 4 - debug)
func (e *Emulator) SetVerbosityLevel(level int) error {
	ok := C.emulator_set_verbosity_level(C.int(level))
	if !ok {
		return fmt.Errorf("set VerbosityLevel error")
	}
	return nil
}

func (e *Emulator) SetBalance(balance int64) {
	e.balance = uint64(balance)
}

func (e *Emulator) SetLibs(libs *boc.Cell) error {
	libsBoc, err := libs.ToBocBase64()
	if err != nil {
		return err
	}
	cLibsStr := C.CString(libsBoc)
	defer C.free(unsafe.Pointer(cLibsStr))
	ok := C.tvm_emulator_set_libraries(e.emulator, cLibsStr)
	if !ok {
		return fmt.Errorf("set libs error")
	}
	return nil
}

func (e *Emulator) SetGasLimit(gasLimit int64) error {
	ok := C.tvm_emulator_set_gas_limit(e.emulator, C.int64_t(gasLimit))
	if !ok {
		return fmt.Errorf("set gas limit error")
	}
	return nil
}

func (e *Emulator) setC7(address string, unixTime uint32, balance uint64, randSeed [32]byte, config *boc.Cell) error {
	configBoc, err := config.ToBocBase64()
	if err != nil {
		return err
	}
	cConfigStr := C.CString(configBoc)
	defer C.free(unsafe.Pointer(cConfigStr))
	cAddressStr := C.CString(address)
	defer C.free(unsafe.Pointer(cAddressStr))
	cSeedStr := C.CString(hex.EncodeToString(randSeed[:]))
	defer C.free(unsafe.Pointer(cSeedStr))
	ok := C.tvm_emulator_set_c7(e.emulator, cAddressStr, C.uint32_t(unixTime), C.uint64_t(balance), cSeedStr, cConfigStr)
	if !ok {
		return fmt.Errorf("set C7 error")
	}
	return nil
}

/**
 * @brief Run get method
 * @param tvm_emulator Pointer to TVM emulator
 * @param method_id Integer method id
 * @param stack_boc Base64 encoded BoC serialized stack (VmStack)
 * @return Json object with error:
 * {
 *   "success": false,
 *   "error": "Error description"
 * }
 * Or success:
 * {
 *   "success": true
 *   "vm_log": "...",
 *   "vm_exit_code": 0,
 *   "stack": "Base64 encoded BoC serialized stack (VmStack)",
 *   "missing_library": null,
 *   "gas_used": 1212
 * }
 */
type result struct {
	Success        bool   `json:"success"`
	Error          string `json:"error"`
	VmLog          string `json:"vm_log"`
	VmExitCode     int    `json:"vm_exit_code"`
	Stack          string `json:"stack"`
	MissingLibrary string `json:"missing_library"`
	GasUsed        string `json:"gas_used"`
}

func (e *Emulator) RunGetMethod(ctx context.Context, accountId tongo.AccountID, method string, params tongo.VmStack) (uint32, tongo.VmStack, error) {

	address := accountId.ToRaw()

	var seed [32]byte
	_, err := rand.Read(seed[:])

	err = e.setC7(address, uint32(time.Now().Unix()), e.balance, seed, e.config)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}

	paramsCell := boc.NewCell()
	err = tlb.Marshal(paramsCell, params)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	res, err := e.runGetMethod(method, paramsCell)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	if !res.Success {
		return 0, tongo.VmStack{}, fmt.Errorf("TVM emulation error: %v", res.Error)
	}

	b, err := base64.StdEncoding.DecodeString(res.Stack)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	c, err := boc.DeserializeBoc(b)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	var stack tongo.VmStack
	err = tlb.Unmarshal(c[0], &stack)
	if err != nil {
		return 0, tongo.VmStack{}, err
	}
	return uint32(res.VmExitCode), stack, nil
}

func (e *Emulator) runGetMethod(methodName string, stack *boc.Cell) (result, error) {
	stackBoc, err := stack.ToBocBase64()
	if err != nil {
		return result{}, err
	}
	cStackStr := C.CString(stackBoc)
	defer C.free(unsafe.Pointer(cStackStr))

	methodID := int(utils.Crc16String(methodName)&0xffff) | 0x10000

	var res result
	r := C.tvm_emulator_run_get_method(e.emulator, C.int(methodID), cStackStr)
	rJSON := C.GoString(r)
	defer C.free(unsafe.Pointer(r))

	err = json.Unmarshal([]byte(rJSON), &res)
	if err != nil {
		return result{}, err
	}

	return res, nil
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
