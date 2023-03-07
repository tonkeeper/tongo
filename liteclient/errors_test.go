package liteclient

import (
	"errors"
	"testing"
)

func TestIsNotConnectedYet(t *testing.T) {
	c := Connection{}
	c.status = Connecting
	p, err := NewPacket([]byte("hello"))
	if err != nil {
		t.Fatalf("NewPacket() failed: %v", err)
	}
	err = c.Send(p)
	notConnected := IsNotConnectedYet(err)
	if notConnected != true {
		t.Fatal()
	}
}

func TestIsClientError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "client error",
			err:  newClientError("method() failed"),
			want: true,
		},
		{
			name: "not client error",
			err:  errors.New("some err"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsClientError(tt.err); got != tt.want {
				t.Errorf("IsClientError() = %v, want %v", got, tt.want)
			}
		})
	}
}
