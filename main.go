package main

import "C"
import (
	"encoding/json"
	"fmt"
)

// #cgo darwin CFLAGS: -I ./lib/
// #cgo LDFLAGS: -L ./lib/ -Wl,-rpath,./lib -l vm-exec-lib
// #include "./lib/libvm-exec-lib.h"
import "C"

type TvmStackEntry interface{}

type TvmStackEntryNull struct {
	Type string `json:"type"`
}

type TVMStackEntryInt struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type TVMExecutionResult struct {
	ExitCode       int             `json:"exit_code"`
	GasConsumed    int             `json:"gas_consumed"`
	DataCell       string          `json:"data_cell"`
	ActionListCell string          `json:"action_list_cell"`
	Logs           string          `json:"logs"`
	Stack          []TvmStackEntry `json:"stack"`
}

func main() {

	config := `{"function_selector":117759,"init_stack":[{"type":"int","value":"123"},{"type":"int","value":"123"}],"code":"te6cckEBBAEAGwABFP8A9KQT9LzyyAsBAgFiAwIAB6GX/0EAAtCnICBl","data":"te6cckEBAQEAAgAAAEysuc0=","time":1645034711694}`

	res := C.vm_exec(C.int(len(config)), C.CString(config))
	resJson := C.GoString(res)

	var executeResult TVMExecutionResult
	err := json.Unmarshal([]byte(resJson), &executeResult)

	fmt.Println(err)
	fmt.Println(executeResult.Stack[0])
	//fmt.Print("Hello world!")
}
