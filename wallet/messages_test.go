package wallet

import (
	"crypto/ed25519"
	"encoding/hex"
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

func mustPubkeyFromHex(hexPubkey string) ed25519.PublicKey {
	bytes, err := hex.DecodeString(hexPubkey)
	if err != nil {
		panic(err)
	}
	return ed25519.PublicKey(bytes)
}

func TestExtractRawMessages(t *testing.T) {
	tests := []struct {
		name    string
		ver     Version
		boc     string
		want    []RawMessage
		wantErr bool
	}{
		{
			name: "v4",
			boc:  "te6ccgECAwEAAQUAAeGIANmaZLULGG8tJ/XFeVVjhSDQY0nCFNh3aJ3RbCt5Q6RABSMjS4x6Gq0Zqdbt/8u9KDhBmpjeDE1mJwmaGkKpoKmNpuFpsf2j6g/KVbw9kWLcEdc/rCcX6euh2ksWAyZx6AFNTRi7I89J2AAAASAAHAEBaGIAS1ZNypaCh7zgPRcvBcpDlS3gxPwxnEFWGfVBevyzhRwhMS0AAAAAAAAAAAAAAAAAAAECALAPin6lAAAAAAAAAAAxtgM4AKZ+YbyuRCr3COPqoHc/iwAZGwcvzy6H7y1iPME1tc0/ABszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSIAgIAAAAA",
			ver:  V4R1,
			want: []RawMessage{
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
			want: []RawMessage{
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
			want: []RawMessage{
				{
					Message: mustFromHex("te6ccgEBAgEAjAABaGIAYeITnAruocV3ZaCBjfbcIK27S8GFMv5jOh6XPwNuAUkgFykzCAAAAAAAAAAAAAAAAAEBAKVfzD0UAAAAAAAAAACACmfmG8rkQq9wjj6qB3P4sAGRsHL88uh+8tYjzBNbXNPwAbM0yWoWMN5aT+uK8qrHCkGgxpOEKbDu0Tui2Fbyh0iAcxLQCA=="),
					Mode:    3,
				},
			},
		},
		{
			name: "v5",
			boc:  "te6ccgECCAEAAZ4AAfGIAehvqHPiQ2Ru+zkowjJx/7oJbqEYRnlCOuPe5+2gm24WA5tLO3f////oAAAAAAADMYd8kAAAAAEHzN670eqqNU3yWGkX1dOynyAbT7DN4cFDpE0r+nInTomGrifjPTaZvG3YxYzTHpLoNesGc9s5Q0tHlLNcFNQeAQIKDsPIbQMCAwIKDsPIbQMEBQCpaAHob6hz4kNkbvs5KMIycf+6CW6hGEZ5Qjrj3uftoJtuFwAbM0yWoWMN5aT+uK8qrHCkGgxpOEKbDu0Tui2Fbyh0iAy3GwAAAAAAAAAAAAAAAAAAQAIKDsPIbQMGBwCpaAHob6hz4kNkbvs5KMIycf+6CW6hGEZ5Qjrj3uftoJtuFwAbM0yWoWMN5aT+uK8qrHCkGgxpOEKbDu0Tui2Fbyh0iAx6EgAAAAAAAAAAAAAAAAAAQAAAAKloAehvqHPiQ2Ru+zkowjJx/7oJbqEYRnlCOuPe5+2gm24XABszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSIDD0JAAAAAAAAAAAAAAAAAABA",
			ver:  V5R1,
			want: []RawMessage{
				{
					Message: mustFromHex("te6ccgEBAQEAVwAAqWgB6G+oc+JDZG77OSjCMnH/ugluoRhGeUI6497n7aCbbhcAGzNMlqFjDeWk/rivKqxwpBoMaThCmw7tE7othW8odIgMtxsAAAAAAAAAAAAAAAAAAEA="),
					Mode:    3,
				},
				{
					Message: mustFromHex("te6ccgEBAQEAVwAAqWgB6G+oc+JDZG77OSjCMnH/ugluoRhGeUI6497n7aCbbhcAGzNMlqFjDeWk/rivKqxwpBoMaThCmw7tE7othW8odIgMehIAAAAAAAAAAAAAAAAAAEA="),
					Mode:    3,
				},
				{
					Message: mustFromHex("te6ccgEBAQEAVwAAqWgB6G+oc+JDZG77OSjCMnH/ugluoRhGeUI6497n7aCbbhcAGzNMlqFjDeWk/rivKqxwpBoMaThCmw7tE7othW8odIgMPQkAAAAAAAAAAAAAAAAAAEA="),
					Mode:    3,
				},
			},
		},
		{
			name: "highload",
			boc:  "te6ccgECCQEAAUMAAUWIAbeTPaOhIeFpX00pVBankGP2F/kaObq5EAdGLvI+omE+DAEBmXzKceTPz+weyz8nYZbOkpsBYbvy6gN7h38ZVL6RTqln7XbUzHkQqxRp1B1ZYkBgMW1NtE7r8Jwg26HcS3qPiwYAAYiUZMJyTpfTrVXAAgIFngACAwQBAwDgBQEDAOAHAWJCADZmmS1CxhvLSf1xXlVY4Ug0GNJwhTYd2id0WwreUOkQCKAAAAAAAAAAAAAAAAABBgBQAAAAADcwMzBhYzQ2LWI5NWMtNDRjNy04ZDdiLTYxMjMyNmU2ZTUxMgFiQgA2ZpktQsYby0n9cV5VWOFINBjScIU2HdondFsK3lDpEAlAAAAAAAAAAAAAAAAAAQgAUAAAAAAzYjA2OTU1YS03YjRjLTQ1YWEtOTVlNy0wNTI4ZWZhYjAyM2E=",
			ver:  HighLoadV2R2,
			want: []RawMessage{
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

func TestMessageV5VerifySignature(t *testing.T) {
	tests := []struct {
		name              string
		boc               string
		publicKey         ed25519.PublicKey
		invalidPublicKeys []ed25519.PublicKey
		wantErr           bool
	}{
		{
			name:      "wallet v5",
			boc:       "te6ccgECCAEAAZ4AAfGIAehvqHPiQ2Ru+zkowjJx/7oJbqEYRnlCOuPe5+2gm24WA5tLO3f////oAAAAAAADMY8YiAAAADPkc94coPiaMQo1EI1uuJWlVQGxiffff96PyOTGiQhUjkr733UkT8rfdXxuYcb9SMykg8Tlo7LNBB187eI+ymw2AQIKDsPIbQMCAwIKDsPIbQMEBQCpaAHob6hz4kNkbvs5KMIycf+6CW6hGEZ5Qjrj3uftoJtuFwAbM0yWoWMN5aT+uK8qrHCkGgxpOEKbDu0Tui2Fbyh0iAy3GwAAAAAAAAAAAAAAAAAAQAIKDsPIbQMGBwCpaAHob6hz4kNkbvs5KMIycf+6CW6hGEZ5Qjrj3uftoJtuFwAbM0yWoWMN5aT+uK8qrHCkGgxpOEKbDu0Tui2Fbyh0iAx6EgAAAAAAAAAAAAAAAAAAQAAAAKloAehvqHPiQ2Ru+zkowjJx/7oJbqEYRnlCOuPe5+2gm24XABszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSIDD0JAAAAAAAAAAAAAAAAAABA",
			publicKey: mustPubkeyFromHex("406b63856ff6913fe2170a5c128113c6bd8256438a43340ea3bf6e0bbc56f9ca"),
			invalidPublicKeys: []ed25519.PublicKey{
				mustPubkeyFromHex("406b63856ff6913fe2170a5c128113c6bd8256438a43340ea3bf6e0bbc56f9bb"),
				mustPubkeyFromHex("406b63856ff6913fe2170a5c128113c6bd8256438a43340ea3bf6e0bbc560000"),
				mustPubkeyFromHex("cfa50eeb1c3293c92bd33d5aa672c1717accd8a21b96033debb6d30b5bb230df"),
			},
		},
		{
			name:      "wallet v5",
			boc:       "te6ccgECCAEAAZ4AAfGIAVjXuMKpIWGwKJenbsOOEh1AEZo6J5Zu0R8EDI37LVyKA5tLO3f////oAAAAAAADMY9YuAAAAAFs/6Zj178nNgWPsbSM2UaEwrcyYPF0kSqZ4d+fhPMfynWRWKBCiVh2PtDewtHZ5FW1luvfXHDqGX0DtYSHfVwGAQIKDsPIbQMCAwIKDsPIbQMEBQCpaAFY17jCqSFhsCiXp27DjhIdQBGaOieWbtEfBAyN+y1ciwAbM0yWoWMN5aT+uK8qrHCkGgxpOEKbDu0Tui2Fbyh0iAy3GwAAAAAAAAAAAAAAAAAAQAIKDsPIbQMGBwCpaAFY17jCqSFhsCiXp27DjhIdQBGaOieWbtEfBAyN+y1ciwAbM0yWoWMN5aT+uK8qrHCkGgxpOEKbDu0Tui2Fbyh0iAx6EgAAAAAAAAAAAAAAAAAAQAAAAKloAVjXuMKpIWGwKJenbsOOEh1AEZo6J5Zu0R8EDI37LVyLABszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSIDD0JAAAAAAAAAAAAAAAAAABA",
			publicKey: mustPubkeyFromHex("cfa50eeb1c3293c92bd33d5aa672c1717accd8a21b96033debb6d30b5bb230df"),
			invalidPublicKeys: []ed25519.PublicKey{
				mustPubkeyFromHex("406b63856ff6913fe2170a5c128113c6bd8256438a43340ea3bf6e0bbc56f9bb"),
				mustPubkeyFromHex("406b63856ff6913fe2170a5c128113c6bd8256438a43340ea3bf6e0bbc56f9ca"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := mustFromHex(tt.boc)
			var m tlb.Message
			if err := tlb.Unmarshal(cell, &m); err != nil {
				t.Fatalf("Unmarshal() failed: %v", err)
			}
			msgBody := boc.Cell(m.Body.Value)
			err := MessageV5VerifySignature(msgBody, tt.publicKey)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("MessageV5VerifySignature() had to fail but it didn't")
				}
				if err.Error() != ErrBadSignature.Error() {
					t.Fatalf("MessageV5VerifySignature() failed: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("MessageV5VerifySignature() failed: %v", err)
			}

			for _, publicKey := range tt.invalidPublicKeys {
				if err = MessageV5VerifySignature(msgBody, publicKey); err == nil {
					t.Fatalf("MessageV5VerifySignature() had to fail but it didn't")
				}
			}
		})
	}
}
