package liteapi

import (
	"encoding/base64"
	"testing"
)

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
