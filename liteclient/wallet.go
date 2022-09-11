package liteclient

import (
	"context"
	"github.com/startfellows/tongo"
)

func (c *Client) GetWalletSeqno(account tongo.AccountID) (uint32, error) {
	stack, err := c.RunSmcMethod(context.Background(), 4, account, "seqno", tongo.VmStack{})
	if err != nil {
		return 0, err
	}
	return uint32(stack[0].VmStkTinyInt), nil
}
