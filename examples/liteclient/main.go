package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/tonkeeper/tongo/liteclient"
)

func main() {
	myPayload, _ := hex.DecodeString("7af98bb435263e6c95d6fecb497dfd0aa5f031e7d412986b5ce720496db512052e8f2d100cdf068c7904345aad16000000000000")
	serverPubkey, err := base64.StdEncoding.DecodeString("wQE0MVhXNWUXpWiW5Bk8cAirIh5NNG3cZM1/fSVKIts=")
	if err != nil {
		panic(err)
	}
	c, err := liteclient.NewConnection(context.Background(), serverPubkey, "135.181.140.221:46995")
	if err != nil {
		panic(err)
	}
	packet, err := liteclient.NewPacket(myPayload)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Send: %x\n", myPayload)
	err = c.Send(packet)
	if err != nil {
		panic(err)
	}
	for p := range c.Responses() {
		fmt.Printf("Received: %x\n", p.Payload)
		return
	}
}
