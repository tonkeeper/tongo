package main

import (
	"context"
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/liteclient"
)

func main() {
	// options, err := config.ParseConfigFile("path/to/config.json")
	tongoClient, err := liteclient.NewClientWithDefaultTestnet()
	if err != nil {
		fmt.Printf("Unable to create tongo client: %v", err)
	}
	accountId, _ := tongo.AccountIDFromRaw("0:E2D41ED396A9F1BA03839D63C5650FAFC6FCFB574FD03F2E67D6555B61A3ACD9")
	state, err := tongoClient.GetAccountState(context.Background(), *accountId)
	if err != nil {
		fmt.Printf("Get account state error: %v", err)
	}
	fmt.Printf("Account status: %v\nBalance: %v\n", state.Status, state.Balance)
}
