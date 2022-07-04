package main

import (
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
	"github.com/startfellows/tongo/tvm"
	"math/big"
)

func main() {
	// Execute nft collection contract to get nft item address by index
	// https://github.com/ton-blockchain/TIPs/issues/62
	// get_nft_address_by_index(int index) returns slice address

	codeCell, _ := boc.DeserializeBocBase64("te6ccsECFAEAAh8AAAAADQASABcAkQDDARgBMAFQAVUBWgFzAYMBpAGpAa4B0gHXAfECCgEU/wD0pBP0vPLICwECAWICDQICzQMIBOfRBjgEit8ADoaYGAuNhIrfB9IBgA6Y/pn/aiaH0gaZ/qamoYQQg0npyoKUBdRxQbr4KA6EEIVGWAVrhACGRlgqgC54sSfQEKZbUJ5Y/ln4Dni2TAIH2AcBsoueOC+XDIkuAA8YES4AFxgRLgAfGBGBmB4AJAQFBgcAYDUC0z9TE7vy4ZJTE7oB+gDUMCgQNFnwBo4SAaRDQ8hQBc8WE8s/zMzMye1Ukl8F4gCmNXAD1DCON4BA9JZvpSCOKQakIIEA+r6T8sGP3oEBkyGgUyW78vQC+gDUMCJUSzDwBiO6kwKkAt4Ekmwh4rPmMDJQREMTyFAFzxYTyz/MzMzJ7VQALDI0AfpAMEFEyFAFzxYTyz/MzMzJ7VQAPI4V1NQwEDRBMMhQBc8WE8s/zMzMye1U4F8EhA/y8AIBIAkMAgEgCgsALQByMs/+CjPFslwIMjLARP0APQAywDJgABs+QB0yMsCEsoHy//J0IAA9Ra8ARwIfAFd4AYyMsFWM8WUAT6AhPLaxLMzMlx+wCAIBIA4TAgEgDxAAQ7i10x7UTQ+kDTP9TU1DAQJF8E0NQx1DDQccjLBwHPFszJgCASAREgAvtdr9qJofSBpn+pqahg2IOhph+mH/SAYQAC209H2omh9IGmf6mpqGAovgngCOAD4AsAAlvILfaiaH0gaZ/qamoYLehqGCxOlkGk8=")
	dataCell, _ := boc.DeserializeBocBase64("te6ccsECEgEAAmcAAAAALwAzAFcAbwB8AIEAhgCLAPsBeQG8AfcCAgIHAicCOAI/A1OAH+KPIWfXRAHhzc8BIGKAZ7CGFDhMB09Wc+npbBemPgcgAAAAAAAAaBABBBECAAIDAEQBaHR0cHM6Ly9sb3Rvbi5mdW4vY29sbGVjdGlvbi5qc29uACxodHRwczovL2xvdG9uLmZ1bi9uZnQvART/APSkE/S88sgLBQIBYgYQAgLOBw0CASAIDALXDIhxwCSXwPg0NMDAXGwkl8D4PpA+kAx+gAxcdch+gAx+gAw8AIEs44UMGwiNFIyxwXy4ZUB+kDUMBAj8APgBtMf0z+CEF/MPRRSMLqOhzIQN14yQBPgMDQ0NTWCEC/LJqISuuMCXwSED/LwgCQsB9lE1xwXy4ZH6QCHwAfpA0gAx+gCCCvrwgBuhIZRTFaCh3iLXCwHDACCSBqGRNuIgwv/y4ZIhjj6CEAUTjZHIUAnPFlALzxZxJEkUVEagcIAQyMsFUAfPFlAF+gIVy2oSyx/LPyJus5RYzxcBkTLiAckB+wAQR5QQKjdb4goAggKONSbwAYIQ1TJ22xA3RABtcXCAEMjLBVAHzxZQBfoCFctqEssfyz8ibrOUWM8XAZEy4gHJAfsAkzAyNOJVAvADAHJwghCLdxc1BcjL/1AEzxYQJIBAcIAQyMsFUAfPFlAF+gIVy2oSyx/LPyJus5RYzxcBkTLiAckB+wAAET6RDBwuvLhTYAIBIA4PADs7UTQ0z/6QCDXScIAmn8B+kDUMBAkECPgMHBZbW2AAHQDyMs/WM8WAc8WzMntVIAAJoR+f4AUASwBkA+iAH+KPIWfXRAHhzc8BIGKAZ7CGFDhMB09Wc+npbBemPgcw9Pr6lQ==")
	index := big.NewInt(100)

	args := []tvm.StackEntry{ // prepare input parameters
		tvm.NewBigIntStackEntry(*index),
	}

	// It is necessary to provide the address of the contract. In some cases it is used by tvm.
	account, _ := tongo.AccountIDFromRaw("0:4ccba08d80193c3eb4f92cd8cf10bc425ff2d705a552aad6f3453a141e51b7b7")

	result, err := tvm.RunTvm(codeCell[0], dataCell[0], "get_nft_address_by_index", args, account)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Exit code: %v\n", result.ExitCode)
	fmt.Printf("Gas consumed: %v\n", result.GasConsumed)

	if result.ExitCode != 0 && result.ExitCode != 1 { // 1 - alternative success code
		panic("TVM execution failed")
	}
	if len(result.Stack) != 1 || !result.Stack[0].IsCellSlice() {
		panic("invalid stack data")
	}

	var msgAddress tongo.MsgAddress
	err = tlb.Unmarshal(result.Stack[0].CellSlice(), &msgAddress)
	if err != nil {
		panic(err)
	}
	addr, _ := msgAddress.AccountId()
	fmt.Printf("Nft index: %v\nNft item address: %v\n", index, addr.ToRaw())
}
