package wallet

import (
	"crypto/ed25519"
	"fmt"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func mustFromHex(msg string) *boc.Cell {
	c, err := boc.DeserializeSinglRootBase64(msg)
	if err != nil {
		panic(err)
	}
	return c
}

func TestExtractRawMessages(t *testing.T) {
	tests := []struct {
		name    string
		ver     Version
		boc     string
		want    []tlb.RawMessage
		wantErr bool
	}{
		{
			name: "v4",
			boc:  "te6ccgECAwEAAQUAAeGIANmaZLULGG8tJ/XFeVVjhSDQY0nCFNh3aJ3RbCt5Q6RABSMjS4x6Gq0Zqdbt/8u9KDhBmpjeDE1mJwmaGkKpoKmNpuFpsf2j6g/KVbw9kWLcEdc/rCcX6euh2ksWAyZx6AFNTRi7I89J2AAAASAAHAEBaGIAS1ZNypaCh7zgPRcvBcpDlS3gxPwxnEFWGfVBevyzhRwhMS0AAAAAAAAAAAAAAAAAAAECALAPin6lAAAAAAAAAAAxtgM4AKZ+YbyuRCr3COPqoHc/iwAZGwcvzy6H7y1iPME1tc0/ABszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSIAgIAAAAA",
			ver:  V4R1,
			want: []tlb.RawMessage{
				{
					Message: mustFromHex("te6ccgEBAgEAkQABaGIAS1ZNypaCh7zgPRcvBcpDlS3gxPwxnEFWGfVBevyzhRwhMS0AAAAAAAAAAAAAAAAAAAEBALAPin6lAAAAAAAAAAAxtgM4AKZ+YbyuRCr3COPqoHc/iwAZGwcvzy6H7y1iPME1tc0/ABszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSIAgIAAAAA"),
					Mode:    3,
				},
			},
		},
		{
			name: "v4",
			boc:  "te6ccgEBAgEAqgAB4YgA2ZpktQsYby0n9cV5VWOFINBjScIU2HdondFsK3lDpEAAQ+B903cV6YIMdtd4QtdyekehadSk+QjIgoIiRgjZD9v81PVGEXBKHPgPUknVvxvr/LGcKkLNhY+I1Wuwi/7ACU1NGLsi5dhQAAAA8AAcAQBoQgApn5hvK5EKvcI4+qgdz+LABkbBy/PLofvLWI8wTW1zT6WWgvAAAAAAAAAAAAAAAAAAAA==",
			ver:  V4R1,
			want: []tlb.RawMessage{
				{
					Message: mustFromHex("te6ccgEBAQEANgAAaEIAKZ+YbyuRCr3COPqoHc/iwAZGwcvzy6H7y1iPME1tc0+lloLwAAAAAAAAAAAAAAAAAAA="),
					Mode:    3,
				},
			},
		},
		{
			name: "v4",
			boc:  "te6ccgECAwEAAQAAAeGIANmaZLULGG8tJ/XFeVVjhSDQY0nCFNh3aJ3RbCt5Q6RAAR/y7WiDk/zi6/QObgK7qDZRawFY0k5TaspQuK98GHfLWcVcMgc/kdpXj+nNrmpWHO2mJ6nyxhuxwzzphZVmuBlNTRi7I88W+AAAARAAHAEBaGIAYeITnAruocV3ZaCBjfbcIK27S8GFMv5jOh6XPwNuAUkgFykzCAAAAAAAAAAAAAAAAAECAKVfzD0UAAAAAAAAAACACmfmG8rkQq9wjj6qB3P4sAGRsHL88uh+8tYjzBNbXNPwAbM0yWoWMN5aT+uK8qrHCkGgxpOEKbDu0Tui2Fbyh0iAcxLQCA==",
			ver:  V4R1,
			want: []tlb.RawMessage{
				{
					Message: mustFromHex("te6ccgEBAgEAjAABaGIAYeITnAruocV3ZaCBjfbcIK27S8GFMv5jOh6XPwNuAUkgFykzCAAAAAAAAAAAAAAAAAEBAKVfzD0UAAAAAAAAAACACmfmG8rkQq9wjj6qB3P4sAGRsHL88uh+8tYjzBNbXNPwAbM0yWoWMN5aT+uK8qrHCkGgxpOEKbDu0Tui2Fbyh0iAcxLQCA=="),
					Mode:    3,
				},
			},
		},
		{
			name: "highload",
			boc:  "te6ccgECCQEAAUMAAUWIAbeTPaOhIeFpX00pVBankGP2F/kaObq5EAdGLvI+omE+DAEBmXzKceTPz+weyz8nYZbOkpsBYbvy6gN7h38ZVL6RTqln7XbUzHkQqxRp1B1ZYkBgMW1NtE7r8Jwg26HcS3qPiwYAAYiUZMJyTpfTrVXAAgIFngACAwQBAwDgBQEDAOAHAWJCADZmmS1CxhvLSf1xXlVY4Ug0GNJwhTYd2id0WwreUOkQCKAAAAAAAAAAAAAAAAABBgBQAAAAADcwMzBhYzQ2LWI5NWMtNDRjNy04ZDdiLTYxMjMyNmU2ZTUxMgFiQgA2ZpktQsYby0n9cV5VWOFINBjScIU2HdondFsK3lDpEAlAAAAAAAAAAAAAAAAAAQgAUAAAAAAzYjA2OTU1YS03YjRjLTQ1YWEtOTVlNy0wNTI4ZWZhYjAyM2E=",
			ver:  HighLoadV2R2,
			want: []tlb.RawMessage{
				{
					Message: mustFromHex("te6ccgEBAgEAXgABYkIANmaZLULGG8tJ/XFeVVjhSDQY0nCFNh3aJ3RbCt5Q6RAIoAAAAAAAAAAAAAAAAAEBAFAAAAAANzAzMGFjNDYtYjk1Yy00NGM3LThkN2ItNjEyMzI2ZTZlNTEy"),
					Mode:    3,
				},
				{
					Message: mustFromHex("te6ccgEBAgEAXgABYkIANmaZLULGG8tJ/XFeVVjhSDQY0nCFNh3aJ3RbCt5Q6RAJQAAAAAAAAAAAAAAAAAEBAFAAAAAAM2IwNjk1NWEtN2I0Yy00NWFhLTk1ZTctMDUyOGVmYWIwMjNh"),
					Mode:    3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := boc.DeserializeBocBase64(tt.boc)
			if err != nil {
				t.Fatal(err)
			}
			rawMessages, err := ExtractRawMessages(tt.ver, c[0])
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(rawMessages, tt.want) {
				for i, msg := range rawMessages {
					bocBase64, err := msg.Message.ToBocBase64()
					if err != nil {
						t.Fatal(err)
					}
					wantBase64, err := tt.want[i].Message.ToBocBase64()
					if err != nil {
						t.Fatal(err)
					}
					fmt.Printf(" got message: %v\n", bocBase64)
					fmt.Printf("want message: %v\n", wantBase64)
					fmt.Printf(" got mode: %v\n", msg.Mode)
					fmt.Printf("want mode: %v\n", tt.want[i].Mode)
				}
				t.Fatalf("wrong raw messages")
			}
		})
	}
}

func TestSignedMsgBody_Verify(t *testing.T) {
	seed1 := RandomSeed()
	privateKey1, _ := SeedToPrivateKey(seed1)

	seed2 := RandomSeed()
	privateKey2, _ := SeedToPrivateKey(seed2)

	tests := []struct {
		name              string
		privateKey        ed25519.PrivateKey
		invalidPublicKeys []ed25519.PublicKey
	}{
		{
			name:       "signed by privateKey1",
			privateKey: privateKey1,
			invalidPublicKeys: []ed25519.PublicKey{
				privateKey2.Public().(ed25519.PublicKey),
			},
		},
		{
			name:       "signed by privateKey2",
			privateKey: privateKey2,
			invalidPublicKeys: []ed25519.PublicKey{
				privateKey1.Public().(ed25519.PublicKey),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyCell := boc.NewCell()
			if err := bodyCell.WriteBytes([]byte("hello")); err != nil {
				t.Fatalf("WriteBytes() failed: %v", err)
			}
			signBytes, err := bodyCell.Sign(tt.privateKey)
			if err != nil {
				t.Fatalf("Sign() failed: %v", err)
			}
			bits512 := tlb.Bits512{}
			copy(bits512[:], signBytes[:])
			signedBody := SignedMsgBody{
				Sign:    bits512,
				Message: tlb.Any(*bodyCell),
			}
			publicKey := tt.privateKey.Public().(ed25519.PublicKey)
			if err = signedBody.Verify(publicKey); err != nil {
				t.Fatalf("Verify() failed: %v", err)
			}
			for _, invalidKey := range tt.invalidPublicKeys {
				if err = signedBody.Verify(invalidKey); err == nil {
					t.Fatalf("Verify() had to fail but it didn't")
				}
			}
		})
	}
}
