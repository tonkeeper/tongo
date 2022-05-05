## ADNL pure-golang implementation.

### Installation

```shell

go get github.com/startfellows/tongo

```

### Usage 

LiteServer query

```go
package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/lightclient"
)

func main() {
	serverPubkey := "Z3X5IRueR4Lbdc0I+1SZwyWmnuDNHdUf14JwIPsGgRw="
	host := "127.0.0.1:7742"
	client, err := liteclient.NewClient(serverPubkey, host)
	state, err := client.GetLastRawAccountState(*tongo.MustParseAccountId("0:0c307d4bf558ca82f33dda0db140bfa2a8a511c61993582e69d4b834e6495e3c"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", state.Balance)
}

```