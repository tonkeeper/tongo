package tvm

// #cgo linux LDFLAGS: -L ../lib/linux/ -Wl,-rpath,../lib/linux/ -l emulator
// #include "../lib/emulator-extern.h"
// #include <stdlib.h>
// #include <stdbool.h>
import "C"
import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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

type Emulator struct {
	emulator unsafe.Pointer
	config   string
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
	configBoc, err := config.ToBocBase64()
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
		config:   configBoc,
		balance:  uint64(balance),
	}
	runtime.SetFinalizer(&e, destroy)
	return &e, nil
}

// NewEmulatorFromBOCsBase64
// Verbosity level of VM log. 0 - log truncated to last 256 characters. 1 - unlimited length log.
// 2 - for each command prints its cell hash and offset. 3 - for each command log prints all stack values.
func NewEmulatorFromBOCsBase64(code, data, config string, balance int64, verbosityLevel txemulator.VerbosityLevel) (*Emulator, error) {
	cCodeStr := C.CString(code)
	defer C.free(unsafe.Pointer(cCodeStr))
	cDataStr := C.CString(data)
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

func (e *Emulator) setC7(address string, unixTime uint32, balance uint64, randSeed [32]byte, configBoc string) error {
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

func (e *Emulator) RunGetMethod(ctx context.Context, accountId tongo.AccountID, method string, params tlb.VmStack) (uint32, tlb.VmStack, error) {

	address := accountId.ToRaw()

	var seed [32]byte
	_, err := rand.Read(seed[:])

	err = e.setC7(address, uint32(time.Now().Unix()), e.balance, seed, e.config)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}

	paramsCell := boc.NewCell()
	err = tlb.Marshal(paramsCell, params)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	res, err := e.runGetMethod(method, paramsCell)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	if !res.Success {
		return 0, tlb.VmStack{}, fmt.Errorf("TVM emulation error: %v", res.Error)
	}

	b, err := base64.StdEncoding.DecodeString(res.Stack)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	c, err := boc.DeserializeBoc(b)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	var stack tlb.VmStack
	err = tlb.Unmarshal(c[0], &stack)
	if err != nil {
		return 0, tlb.VmStack{}, err
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

	methodID := utils.MethodIdFromName(methodName)

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
