package tonconnect

import (
	"context"
	"fmt"

	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type stubExecutor struct{}

func (stubExecutor) RunSmcMethodByID(_ context.Context, _ ton.AccountID, _ int, _ tlb.VmStack) (uint32, tlb.VmStack, error) {
	return 0, tlb.VmStack{}, fmt.Errorf("offline stub executor: account is not deployed")
}
