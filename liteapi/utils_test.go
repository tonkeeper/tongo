package liteapi

import (
	"encoding/base64"
	"errors"
	"fmt"
	"testing"

	"github.com/tonkeeper/tongo/liteclient"
)

func TestBlockNotInDB(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "code 651",
			err:  liteclient.LiteServerErrorC{Code: 651, Message: "cannot load block"},
			want: true,
		},
		{
			name: "message contains not in db",
			err:  liteclient.LiteServerErrorC{Code: 500, Message: "block seqno not in db"},
			want: true,
		},
		{
			name: "wrapped code 651",
			err:  fmt.Errorf("lookup failed: %w", liteclient.LiteServerErrorC{Code: 651, Message: "cannot load block"}),
			want: true,
		},
		{
			name: "plain error containing not in db",
			err:  errors.New("something not in db"),
			want: true,
		},
		{
			name: "unrelated lite server error",
			err:  liteclient.LiteServerErrorC{Code: 500, Message: "cannot compute block"},
			want: false,
		},
		{
			name: "unrelated plain error",
			err:  errors.New("connection refused"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := blockNotInDB(tt.err); got != tt.want {
				t.Fatalf("blockNotInDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifySendMessagePayload(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr string
	}{
		{
			name:    "all good",
			payload: "te6ccgEBBAEAtwABRYgBvVXMoxQj+kmDtTinWnFdumvpTNo33p48YQKOWyTtUkAMAQGcEai7Dc89wZxdeCdFylIZpkyPHeryzz3UVi/Hz2KuK/vIwuubT1KsFMdJJVGwVNEh4CUvlpMzSjZjDZzoUTADASmpoxdkR96PAAAAcAADAgFkQgAoPvU+sDeRbPQrPGn3bxzd8JnUNGlQcfA/qoFluFxSiRE4gAAAAAAAAAAAAAAAAAEDABIAAAAAaGVsbG8=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := base64.StdEncoding.DecodeString(tt.payload)
			if err != nil {
				t.Fatalf("base64.StdEncoding.DecodeString() failed: %v", err)
			}
			err = VerifySendMessagePayload(payload)
			if len(tt.wantErr) > 0 {
				if err == nil {
					t.Fatalf("expected to get an error")
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("want err: %v, got err: %v", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("VerifySendMessagePayload() failed: %v", err)
			}
		})
	}
}
