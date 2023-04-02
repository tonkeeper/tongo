package tvm

// #cgo darwin LDFLAGS: -L ../lib/darwin/ -Wl,-rpath,../lib/darwin/ -l emulator
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
	"math/rand"
	"runtime"
	"time"
	"unsafe"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/txemulator"
	"github.com/tonkeeper/tongo/utils"
)

type Emulator struct {
	emulator unsafe.Pointer
	config   string
	balance  uint64
	lazyC7   bool
}

type Options struct {
	verbosityLevel txemulator.VerbosityLevel
	balance        int64
	lazyC7         bool
}

type Option func(o *Options)

func WithVerbosityLevel(level txemulator.VerbosityLevel) Option {
	return func(o *Options) {
		o.verbosityLevel = level
	}
}
func WithBalance(balance int64) Option {
	return func(o *Options) {
		o.balance = balance
	}
}

// WithLazyC7Optimization allows to make two attempts to execute a get method.
// At the first attempt an emulator invokes a get method without C7.
// This works for most get methods and significantly decreases the execution time.
// If the first attempt fails,
// an emulator invokes the same get method again but with configured C7.
func WithLazyC7Optimization() Option {
	return func(o *Options) {
		o.lazyC7 = true
	}
}

func defaultOptions() Options {
	return Options{
		lazyC7:         false,
		balance:        1_000_000_000,
		verbosityLevel: txemulator.LogTruncated,
	}
}

// NewEmulator
// Verbosity level of VM log. 0 - log truncated to last 256 characters. 1 - unlimited length log.
// 2 - for each command prints its cell hash and offset. 3 - for each command log prints all stack values.
func NewEmulator(code, data, config *boc.Cell, opts ...Option) (*Emulator, error) {
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
	return NewEmulatorFromBOCsBase64(codeBoc, dataBoc, configBoc, opts...)
}

// NewEmulatorFromBOCsBase64
// Verbosity level of VM log. 0 - log truncated to last 256 characters. 1 - unlimited length log.
// 2 - for each command prints its cell hash and offset. 3 - for each command log prints all stack values.
func NewEmulatorFromBOCsBase64(code, data, config string, opts ...Option) (*Emulator, error) {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}
	cCodeStr := C.CString(code)
	defer C.free(unsafe.Pointer(cCodeStr))
	cDataStr := C.CString(data)
	defer C.free(unsafe.Pointer(cDataStr))
	level := C.int(options.verbosityLevel)
	e := Emulator{
		emulator: C.tvm_emulator_create(cCodeStr, cDataStr, level),
		config:   config,
		lazyC7:   options.lazyC7,
		balance:  uint64(options.balance),
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

func (e *Emulator) setC7(address string, unixTime uint32) error {
	var seed [32]byte
	_, err := rand.Read(seed[:])
	if err != nil {
		return err
	}
	cConfigStr := C.CString(e.config)
	defer C.free(unsafe.Pointer(cConfigStr))
	cAddressStr := C.CString(address)
	defer C.free(unsafe.Pointer(cAddressStr))
	cSeedStr := C.CString(hex.EncodeToString(seed[:]))
	defer C.free(unsafe.Pointer(cSeedStr))
	ok := C.tvm_emulator_set_c7(e.emulator, cAddressStr, C.uint32_t(unixTime), C.uint64_t(e.balance), cSeedStr, cConfigStr)
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

func (e *Emulator) RunSmcMethod(ctx context.Context, accountId tongo.AccountID, method string, params tlb.VmStack) (uint32, tlb.VmStack, error) {
	methodID := utils.MethodIdFromName(method)
	return e.RunSmcMethodByID(ctx, accountId, methodID, params)
}

func (e *Emulator) RunSmcMethodByID(ctx context.Context, accountId tongo.AccountID, methodID int, params tlb.VmStack) (uint32, tlb.VmStack, error) {
	if !e.lazyC7 {
		err := e.setC7(accountId.ToRaw(), uint32(time.Now().Unix()))
		if err != nil {
			return 0, tlb.VmStack{}, err
		}
	}
	res, err := e.runGetMethod(methodID, params)
	if err != nil {
		return 0, tlb.VmStack{}, err
	}
	if res.Success && res.VmExitCode != 0 && res.VmExitCode != 1 && e.lazyC7 {
		err = e.setC7(accountId.ToRaw(), uint32(time.Now().Unix()))
		if err != nil {
			return 0, tlb.VmStack{}, err
		}
		res, err = e.runGetMethod(methodID, params)
		if err != nil {
			return 0, tlb.VmStack{}, err
		}
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

func (e *Emulator) runGetMethod(methodID int, params tlb.VmStack) (result, error) {
	stack := boc.NewCell()
	err := tlb.Marshal(stack, params)
	if err != nil {
		return result{}, err
	}
	stackBoc, err := stack.ToBocBase64()
	if err != nil {
		return result{}, err
	}
	cStackStr := C.CString(stackBoc)
	defer C.free(unsafe.Pointer(cStackStr))

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
