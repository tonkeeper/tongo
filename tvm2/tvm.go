package tvm2

import (
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

type TVM struct {
	// c0 — Contains the next continuation or return continuation (similar
	//to the subroutine return address in conventional designs). This value
	//must be a Continuation.
	c0 Continuation
	// c1 — Contains the alternative (return) continuation; this value must
	//be a Continuation. It is used in some (experimental) control flow
	//primitives, allowing TVM to define and call “subroutines with two exit
	//points”.
	с1 Continuation
	// c2 — Contains the exception handler. This value is a Continuation,
	//invoked whenever an exception is triggered.
	с2 Continuation
	// c3 — Contains the current dictionary, essentially a hashmap containing
	//the code of all functions used in the program. For reasons explained
	//later in 4.6, this value is also a Continuation, not a Cell as one might
	//expect.
	с3 Continuation
	// c4 — Contains the root of persistent data, or simply the data. This
	//value is a Cell. When the code of a smart contract is invoked, c4
	//points to the root cell of its persistent data kept in the blockchain
	//state. If the smart contract needs to modify this data, it changes c4
	//before returning.
	с4 boc.Cell
	// c5 — Contains the output actions. It is also a Cell initialized by a
	//reference to an empty cell, but its final value is considered one of the
	//smart contract outputs. For instance, the SENDMSG primitive, specific
	//for the TON Blockchain, simply
	c5 boc.Cell

	// c7 Contains the root of temporary data. It is a Tuple, initialized by
	//a reference to an empty Tuple before invoking the smart contract and
	//discarded after its termination.
	c7 Tuple

	stack tlb.VmStack

	// cc — Contains the current continuation (i.e., the
	//code that would be normally executed after the current primitive is
	//completed). This component is similar to the instruction pointer reg-
	//ister (ip) in other architectures.
	cc Continuation
	// cp — A special signed 16-bit integer value that selects
	//the way the next TVM opcode will be decoded. For example, future
	//versions of TVM might use different codepages to add new opcodes
	//while preserving backward compatibility.
	cp int16

	// gasLimits — Contains four signed 64-bit integers: the current gas
	//limit gl, the maximal gas limit gm, the remaining gas gr, and the gas
	//credit gc. Always 0 ≤ gl ≤ gm, gc ≥ 0, and gr ≤ gl + gc; gc is usually
	//initialized by zero, gr is initialized by gl + gc and gradually decreases
	//as the TVM runs. When gr becomes negative or if the final value of gr
	//is less than gc, an out of gas exception is triggered
	gasLimits struct {
		current   int64
		maximal   int64
		remaining int64
		credit    int64
	}
}

type Continuation struct {
}

type Tuple struct {
}
type StackValue struct {
}
