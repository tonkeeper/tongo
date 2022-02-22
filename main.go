package main

import (
	"fmt"
	"tongo/boc"
	"tongo/tvm"
)

func main() {
	code, _ := boc.DeserializeBocBase64("te6cckEBBAEAGwABFP8A9KQT9LzyyAsBAgFiAwIAB6GX/0EAAtCnICBl")
	data, _ := boc.DeserializeBocBase64("te6cckEBAQEAAgAAAEysuc0=")

	args := []tvm.TvmStackEntry{
		tvm.NewIntStackEntry(1),
		tvm.NewIntStackEntry(2),
	}

	runTvm, _ := tvm.RunTvm(code[0], data[0], "sum", args, 0)

	fmt.Println(runTvm.Stack[0].Int64())
}
