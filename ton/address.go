package ton

import (
	"github.com/tonkeeper/tongo/tlb"
)

// Address identifies an account in the TON network.
// Comparing to AccountID which is a low-level identifier of a smart contract,
// Address is a high-level abstraction containing additional information besides AccountID,
// which is useful for building more advanced workflows.
type Address struct {
	ID          AccountID
	Bounce      bool
	StateInit   *tlb.StateInit
	TestnetOnly bool
}
