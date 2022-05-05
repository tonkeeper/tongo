package adnl

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	pubkey, err := base64.StdEncoding.DecodeString("Z3X5IRueR4Lbdc0I+1SZwyWmnuDNHdUf14JwIPsGgRw=")
	if err != nil {
		panic(err)
	}
	c, err := NewConnection(context.Background(), pubkey, "127.0.0.1:7742")
	if err != nil {
		panic(err)
	}

	client := NewClient(c)
	b, _ := hex.DecodeString("df068c7978250e896bffffffff00000000000000801dcf3501db6f7082c2ea79e22a0ad7305022bd38300a3731f3dace87c2ba16b3582b89e5bb66984ce75d5f868bec2f0fb6bf3c6c3e0a506e3f15559a94ea58583251571e00000000195a6d69afdc08a11c38b729bfad789b957a50650f9bb84b1947fcc0acc9a41b000000")
	resp, err := client.Request(context.Background(), b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", resp)
}
