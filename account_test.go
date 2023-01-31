package tongo

import (
	"testing"

	"github.com/tonkeeper/tongo/boc"
)

func TestAccountStatus_MarshalTLB(t *testing.T) {
	tests := []struct {
		status AccountStatus
	}{
		{status: AccountActive},
		{status: AccountFrozen},
		{status: AccountUninit},
		{status: AccountNone},
	}
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			c := boc.NewCell()
			err := tt.status.MarshalTLB(c, "")
			if err != nil {
				t.Errorf("MarshalTLB() error = %v", err)
			}
			c.ResetCounters()
			var status AccountStatus
			err = (&status).UnmarshalTLB(c, "")
			if err != nil {
				t.Errorf("UnmarshalTLB() error = %v", err)
			}
			if tt.status != status {
				t.Fail()
			}
		})
	}
}
