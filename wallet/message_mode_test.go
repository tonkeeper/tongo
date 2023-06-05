package wallet

import (
	"testing"
)

func TestIsMessageModeSet(t *testing.T) {
	tests := []struct {
		name      string
		modeValue int
		mode      MessageMode
		want      bool
	}{
		{
			modeValue: DefaultMessageMode,
			mode:      AttachAllRemainingBalance,
			want:      false,
		},
		{
			modeValue: 128 + 32,
			mode:      AttachAllRemainingBalance,
			want:      true,
		},
		{
			modeValue: 128,
			mode:      AttachAllRemainingBalance,
			want:      true,
		},
		{
			modeValue: 128 + 2 + 1,
			mode:      AttachAllRemainingBalance,
			want:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := IsMessageModeSet(tt.modeValue, tt.mode)
			if set != tt.want {
				t.Fatalf("want set: %v, got: %v", tt.want, set)
			}
		})
	}
}
