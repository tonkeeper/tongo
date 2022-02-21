package main

import "C"
import (
	"encoding/json"
	"fmt"
	"tongo/tvm"
)

// #cgo darwin CFLAGS: -I ./lib/
// #cgo LDFLAGS: -L ./lib/ -Wl,-rpath,./lib -l vm-exec-lib
// #include "./lib/libvm-exec-lib.h"
import "C"

type TVMExecutionResult struct {
	ExitCode       int                 `json:"exit_code"`
	GasConsumed    int                 `json:"gas_consumed"`
	DataCell       string              `json:"data_cell"`
	ActionListCell string              `json:"action_list_cell"`
	Logs           string              `json:"logs"`
	Stack          []tvm.TvmStackEntry `json:"stack"`
}

func main() {

	//config := `{"function_selector":117759,"init_stack":[{"type":"int","value":"123"},{"type":"int","value":"123"}],"code":"te6cckEBBAEAGwABFP8A9KQT9LzyyAsBAgFiAwIAB6GX/0EAAtCnICBl","data":"te6cckEBAQEAAgAAAEysuc0=","time":1645034711694}`
	//config := `{"function_selector":117759,"init_stack":[{"type":"int","value":"123"},{"type":"int","value":"123"}],"code":"te6cckEBBAEAJAABFP8A9KQT9LzyyAsBAgFiAwIAGaGX/wIGE5GWP5KzQAMAAtA1XPLA","data":"te6cckEBAQEAAgAAAEysuc0=","time":1645206666065}
	config := `{"function_selector":117759,"init_stack":[{"type":"int","value":"123"},{"type":"int","value":"123"}],"code":"te6cckEBBAEAKAABFP8A9KQT9LzyyAsBAgFiAwIAIaGX/wIGE5GWP5KzQOLk3gQlAALQeuKz9Q==","data":"te6cckEBAQEAAgAAAEysuc0=","time":1645207136290}`

	res := C.vm_exec(C.int(len(config)), C.CString(config))
	resJson := C.GoString(res)

	var executeResult TVMExecutionResult
	err := json.Unmarshal([]byte(resJson), &executeResult)

	fmt.Println(err)
	fmt.Println(executeResult.Stack[2].Tuple()[0].Int())
	fmt.Println(executeResult.Stack[2].Tuple()[1].Int())
	//parse := executeResult.Stack[1].cellVal.BeginParse()
	//fmt.Println(parse.ReadBigUint(32))
	//fmt.Print("Hello world!")
}
