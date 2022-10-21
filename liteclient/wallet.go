package liteclient

import (
	"context"
	"fmt"
	"github.com/startfellows/tongo"
)

func (c *Client) GetSeqno(ctx context.Context, account tongo.AccountID) (uint32, error) {
	stack, err := c.RunSmcMethod(ctx, 4, account, "seqno", tongo.VmStack{})
	if err != nil {
		return 0, err
	}
	if len(stack) != 1 || stack[0].SumType != "VmStkTinyInt" {
		return 0, fmt.Errorf("invalid stack")
	}
	return uint32(stack[0].VmStkTinyInt), nil
}
