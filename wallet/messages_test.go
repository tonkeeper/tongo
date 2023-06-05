package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
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
		want    PayloadV1toV4
		wantErr bool
	}{
		{
			boc: "te6ccgECAwEAAQUAAeGIANmaZLULGG8tJ/XFeVVjhSDQY0nCFNh3aJ3RbCt5Q6RABSMjS4x6Gq0Zqdbt/8u9KDhBmpjeDE1mJwmaGkKpoKmNpuFpsf2j6g/KVbw9kWLcEdc/rCcX6euh2ksWAyZx6AFNTRi7I89J2AAAASAAHAEBaGIAS1ZNypaCh7zgPRcvBcpDlS3gxPwxnEFWGfVBevyzhRwhMS0AAAAAAAAAAAAAAAAAAAECALAPin6lAAAAAAAAAAAxtgM4AKZ+YbyuRCr3COPqoHc/iwAZGwcvzy6H7y1iPME1tc0/ABszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSIAgIAAAAA",
			ver: V4R1,
			want: []RawMessage{
				{
					Message: mustFromHex("te6ccgEBAgEAkQABaGIAS1ZNypaCh7zgPRcvBcpDlS3gxPwxnEFWGfVBevyzhRwhMS0AAAAAAAAAAAAAAAAAAAEBALAPin6lAAAAAAAAAAAxtgM4AKZ+YbyuRCr3COPqoHc/iwAZGwcvzy6H7y1iPME1tc0/ABszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSIAgIAAAAA"),
					Mode:    3,
				},
			},
		},
		{
			boc: "te6ccgEBAgEAqgAB4YgA2ZpktQsYby0n9cV5VWOFINBjScIU2HdondFsK3lDpEAAQ+B903cV6YIMdtd4QtdyekehadSk+QjIgoIiRgjZD9v81PVGEXBKHPgPUknVvxvr/LGcKkLNhY+I1Wuwi/7ACU1NGLsi5dhQAAAA8AAcAQBoQgApn5hvK5EKvcI4+qgdz+LABkbBy/PLofvLWI8wTW1zT6WWgvAAAAAAAAAAAAAAAAAAAA==",
			ver: V4R1,
			want: []RawMessage{
				{
					Message: mustFromHex("te6ccgEBAQEANgAAaEIAKZ+YbyuRCr3COPqoHc/iwAZGwcvzy6H7y1iPME1tc0+lloLwAAAAAAAAAAAAAAAAAAA="),
					Mode:    3,
				},
			},
		},
		{
			boc: "te6ccgECAwEAAQAAAeGIANmaZLULGG8tJ/XFeVVjhSDQY0nCFNh3aJ3RbCt5Q6RAAR/y7WiDk/zi6/QObgK7qDZRawFY0k5TaspQuK98GHfLWcVcMgc/kdpXj+nNrmpWHO2mJ6nyxhuxwzzphZVmuBlNTRi7I88W+AAAARAAHAEBaGIAYeITnAruocV3ZaCBjfbcIK27S8GFMv5jOh6XPwNuAUkgFykzCAAAAAAAAAAAAAAAAAECAKVfzD0UAAAAAAAAAACACmfmG8rkQq9wjj6qB3P4sAGRsHL88uh+8tYjzBNbXNPwAbM0yWoWMN5aT+uK8qrHCkGgxpOEKbDu0Tui2Fbyh0iAcxLQCA==",
			ver: V4R1,
			want: []RawMessage{
				{
					Message: mustFromHex("te6ccgEBAgEAjAABaGIAYeITnAruocV3ZaCBjfbcIK27S8GFMv5jOh6XPwNuAUkgFykzCAAAAAAAAAAAAAAAAAEBAKVfzD0UAAAAAAAAAACACmfmG8rkQq9wjj6qB3P4sAGRsHL88uh+8tYjzBNbXNPwAbM0yWoWMN5aT+uK8qrHCkGgxpOEKbDu0Tui2Fbyh0iAcxLQCA=="),
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
				for _, msg := range rawMessages {
					bocBase64, err := msg.Message.ToBocBase64()
					if err != nil {
						t.Fatal(err)
					}
					fmt.Printf("got message: %v\n", bocBase64)
					fmt.Printf("got mode: %v\n", msg.Mode)
				}
				t.Fatalf("wrong raw messages")
			}
		})
	}
}
