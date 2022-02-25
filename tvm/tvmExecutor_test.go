package tvm

import (
	"testing"
	"tongo/boc"
)

func TestExec(t *testing.T) {
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

	args := []TvmStackEntry{
		NewIntStackEntry(1),
		NewIntStackEntry(2),
	}

	runTvm, _ := RunTvm(code[0], data[0], "sum", args, 0)

	if runTvm.Stack[0].Int64() != 3 {
		t.Fail()
	}
}
