package abitolk

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func TestDecodeAndEncodeStonfiSwapPayload(t *testing.T) {
	data := "b5ee9c720101030100fd0001b1178d4519ca4d2d8911062e4253087a4ad408017ff5ba4e14e70a63b4f627751f880704d84b0c2e7f3dbd3862df56a778b3590d0018cafe48d8e51fc3a324d17dfd74f55cc5ad294a1420bf80e39620cb7212f09c905cbead030101e16664de2a800c597b2db76f2c8617ac5f9802ab6d934b9d61ef13b9ab84ba05cae018fdd4a550023cb6103b0174463e61166815ec7f2ad4d715ee4dc488af85599ffba4cb3abd4e004796c207602e88c7cc22cd02bd8fe55a9ae2bdc9b89115f0ab33ff74996757a9800000007fffffffc0020055533ba4e6ad58011e5b081d80ba231f308b340af63f956a6b8af726e24457c2accffdd2659d5ea600000010"
	boc1, _ := boc.DeserializeBocHex(data)

	var x InMsgBody
	if err := tlb.Unmarshal(boc1[0], &x); err != nil {
		t.Fatalf("Unable to unmarshal: %v", err)
	}

	boc2 := boc.NewCell()
	if err := tlb.Marshal(boc2, x); err != nil {
		t.Fatalf("Unable to marshal: %v", err)
	}

	b, _ := boc2.ToBoc()
	res := fmt.Sprintf("%x", b)
	if res != data {
		t.Fatalf("got different result")
	}
}

func TestDecodeAndEncodeDedustSwapPayloadWithContractIfaces(t *testing.T) {
	data := "b5ee9c72010206010001d40001ae0f8a7ea5a9d57b6fda293e0932191c080031551c5ddaa2e8fb5c06780f3727107b28395b1eca8b3e5dd39ae8e96d71dabb0029af18b3e5aaa14e6fd77d47290de6f1a8bd6078d4d6eb3eca9f0dc1f8dd96b0c82160ec01010155e3a0d4828007cbff951bbf9e6d86d93fe8de62ac5556a37322908b5ad84d97bcc93ab14ab1044a0ec2a44002024b00000000801384bfb142ca3fb8918eeee76e56e38a22c9b94832508c035fcfe938a41dac3ba7030401d11a3c2fc38014d78c59f2d550a737ebbea39486f378d45eb03c6a6b759f654f86e0fc6ecb5870029af18b3e5aaa14e6fd77d47290de6f1a8bd6078d4d6eb3eca9f0dc1f8dd96b0e00535e3167cb55429cdfaefa8e521bcde3517ac0f1a9add67d953e1b83f1bb2d61c005008d3c37d7e98014d78c59f2d550a737ebbea39486f378d45eb03c6a6b759f654f86e0fc6ecb5870029af18b3e5aaa14e6fd77d47290de6f1a8bd6078d4d6eb3eca9f0dc1f8dd96b0e00d58012bd3391b8d4870a4dbe92c326dfea2413e60f45186f8d8f06219c28156f156ca001500000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000040"
	boc1, _ := boc.DeserializeBocHex(data)

	decoder := tlb.NewDecoder()
	decoder.WithContractInterfaces([]tlb.ContractInterface{tlb.ContractInterface(TonTep74JettonWallet), tlb.ContractInterface(TonDedustPool)})

	var x InMsgBody
	if err := decoder.Unmarshal(boc1[0], &x); err != nil {
		t.Fatalf("Unable to unmarshal: %v", err)
	}

	boc2 := boc.NewCell()
	if err := tlb.Marshal(boc2, x); err != nil {
		t.Fatalf("Unable to marshal: %v", err)
	}

	b, _ := boc2.ToBoc()
	res := fmt.Sprintf("%x", b)
	if res != data {
		t.Fatalf("got different result")
	}
}

func TestDecodeAndEncodeStonfiSwapPayloadWithContractIfaces(t *testing.T) {
	data := "b5ee9c720101030100fd0001b1178d4519ca4d2d8911062e4253087a4ad408017ff5ba4e14e70a63b4f627751f880704d84b0c2e7f3dbd3862df56a778b3590d0018cafe48d8e51fc3a324d17dfd74f55cc5ad294a1420bf80e39620cb7212f09c905cbead030101e16664de2a800c597b2db76f2c8617ac5f9802ab6d934b9d61ef13b9ab84ba05cae018fdd4a550023cb6103b0174463e61166815ec7f2ad4d715ee4dc488af85599ffba4cb3abd4e004796c207602e88c7cc22cd02bd8fe55a9ae2bdc9b89115f0ab33ff74996757a9800000007fffffffc0020055533ba4e6ad58011e5b081d80ba231f308b340af63f956a6b8af726e24457c2accffdd2659d5ea600000010"
	boc1, _ := boc.DeserializeBocHex(data)

	decoder := tlb.NewDecoder()
	decoder.WithContractInterfaces([]tlb.ContractInterface{tlb.ContractInterface(TonTep74JettonWallet), tlb.ContractInterface(TonStonfiV2PoolWeightedStableSwap)})

	var x InMsgBody
	if err := decoder.Unmarshal(boc1[0], &x); err != nil {
		t.Fatalf("Unable to unmarshal: %v", err)
	}

	boc2 := boc.NewCell()
	if err := tlb.Marshal(boc2, x); err != nil {
		t.Fatalf("Unable to marshal: %v", err)
	}

	b, _ := boc2.ToBoc()
	res := fmt.Sprintf("%x", b)
	if res != data {
		t.Fatalf("got different result")
	}
}

func TestDecodeAndEncodeJsonDedustSwapPayloadWithContractIfaces(t *testing.T) {
	data := "b5ee9c72010206010001d40001ae0f8a7ea5a9d57b6fda293e0932191c080031551c5ddaa2e8fb5c06780f3727107b28395b1eca8b3e5dd39ae8e96d71dabb0029af18b3e5aaa14e6fd77d47290de6f1a8bd6078d4d6eb3eca9f0dc1f8dd96b0c82160ec01010155e3a0d4828007cbff951bbf9e6d86d93fe8de62ac5556a37322908b5ad84d97bcc93ab14ab1044a0ec2a44002024b00000000801384bfb142ca3fb8918eeee76e56e38a22c9b94832508c035fcfe938a41dac3ba7030401d11a3c2fc38014d78c59f2d550a737ebbea39486f378d45eb03c6a6b759f654f86e0fc6ecb5870029af18b3e5aaa14e6fd77d47290de6f1a8bd6078d4d6eb3eca9f0dc1f8dd96b0e00535e3167cb55429cdfaefa8e521bcde3517ac0f1a9add67d953e1b83f1bb2d61c005008d3c37d7e98014d78c59f2d550a737ebbea39486f378d45eb03c6a6b759f654f86e0fc6ecb5870029af18b3e5aaa14e6fd77d47290de6f1a8bd6078d4d6eb3eca9f0dc1f8dd96b0e00d58012bd3391b8d4870a4dbe92c326dfea2413e60f45186f8d8f06219c28156f156ca001500000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000040"
	boc1, _ := boc.DeserializeBocHex(data)

	decoder := tlb.NewDecoder()
	decoder.WithContractInterfaces([]tlb.ContractInterface{tlb.ContractInterface(TonTep74JettonWallet), tlb.ContractInterface(TonDedustPool)})

	var x InMsgBody
	if err := decoder.Unmarshal(boc1[0], &x); err != nil {
		t.Fatalf("Unable to unmarshal tlb: %v", err)
	}

	val, err := json.Marshal(x)
	if err != nil {
		t.Fatalf("Unable to marshal json: %v", err)
	}

	var x2 InMsgBody
	if err = json.Unmarshal(val, &x2); err != nil {
		t.Fatalf("Unable to unmarshal json: %v", err)
	}

	boc2 := boc.NewCell()
	if err := tlb.Marshal(boc2, x2); err != nil {
		t.Fatalf("Unable to marshal tlb: %v", err)
	}

	b, _ := boc2.ToBoc()
	res := fmt.Sprintf("%x", b)
	if res != data {
		t.Fatalf("got different result")
	}
}

func TestDecodeAndEncodeJsonStonfiSwapPayloadWithContractIfaces(t *testing.T) {
	data := "b5ee9c720101030100fd0001b1178d4519ca4d2d8911062e4253087a4ad408017ff5ba4e14e70a63b4f627751f880704d84b0c2e7f3dbd3862df56a778b3590d0018cafe48d8e51fc3a324d17dfd74f55cc5ad294a1420bf80e39620cb7212f09c905cbead030101e16664de2a800c597b2db76f2c8617ac5f9802ab6d934b9d61ef13b9ab84ba05cae018fdd4a550023cb6103b0174463e61166815ec7f2ad4d715ee4dc488af85599ffba4cb3abd4e004796c207602e88c7cc22cd02bd8fe55a9ae2bdc9b89115f0ab33ff74996757a9800000007fffffffc0020055533ba4e6ad58011e5b081d80ba231f308b340af63f956a6b8af726e24457c2accffdd2659d5ea600000010"
	boc1, _ := boc.DeserializeBocHex(data)

	decoder := tlb.NewDecoder()
	decoder.WithContractInterfaces([]tlb.ContractInterface{tlb.ContractInterface(TonTep74JettonWallet), tlb.ContractInterface(TonStonfiV2PoolWeightedStableSwap)})

	var x InMsgBody
	if err := decoder.Unmarshal(boc1[0], &x); err != nil {
		t.Fatalf("Unable to unmarshal: %v", err)
	}

	val, err := json.Marshal(x)
	if err != nil {
		t.Fatalf("Unable to marshal: %v", err)
	}

	var x2 InMsgBody
	if err = json.Unmarshal(val, &x2); err != nil {
		t.Fatalf("Unable to unmarshal json: %v", err)
	}

	boc2 := boc.NewCell()
	if err := tlb.Marshal(boc2, x2); err != nil {
		t.Fatalf("Unable to marshal tlb: %v", err)
	}

	b, _ := boc2.ToBoc()
	res := fmt.Sprintf("%x", b)
	if res != data {
		t.Fatalf("got different result")
	}
}
