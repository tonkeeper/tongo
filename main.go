package main

import (
	"fmt"
	"tongo/boc"
	"tongo/tvm"
)

func main() {
	//  () main() {
	//		;; noop
	//	}
	//
	//	(int) sum(int a, int b) method_id {
	//		return (a + b);
	//	}
	code, _ := boc.DeserializeBocBase64("te6cckEBBAEAGwABFP8A9KQT9LzyyAsBAgFiAwIAB6GX/0EAAtCnICBl")
	// Empty data
	data, _ := boc.DeserializeBocBase64("te6cckEBAQEAAgAAAEysuc0=")

	args := []tvm.TvmStackEntry{
		tvm.NewIntStackEntry(1),
		tvm.NewIntStackEntry(2),
	}

	runTvm, _ := tvm.RunTvm(code[0], data[0], "sum", args, 0)

	// Prints 3
	fmt.Println(runTvm.Stack[0].Int64())
}
