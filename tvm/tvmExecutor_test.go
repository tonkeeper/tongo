package tvm

import (
	"fmt"
	"github.com/startfellows/tongo/boc"
	"testing"
)

func TestExec(t *testing.T) {
	//  () main() {
	//		;; noop
	//	}
	//
	//	(int) sum(int a, int b) method_id {
	//		return (a + b);
	//	}
	code, _ := boc.DeserializeBocBase64("te6cckEBBAEAGwABFP8A9KQT9LzyyAsBAgFiAwIAB6GX/0EAAtCnICBl")
	// Empty data
	data, _ := boc.DeserializeBocBase64("te6cckEBAQEAAgAAAEysuc0=")

	code, _ := boc.DeserializeBocBase64("te6cckECEwEAAf4AART/APSkE/S88sgLAQIBYgIDAgLNBAUCASANDgIBIAYHAgFICwwD7UIMcAkVvgAdDTAwFxsJFb4PpAMO1E0PpA0z/U1NQwBtMf0z+CEGk9OVBSMLqOKRZfBgLQEoIQqMsArXCAEMjLBVAFzxYk+gIUy2oTyx/LPwHPFsmAQPsA4DFRZccF8uGRIMAB4wIgwALjAjQDwAPjAl8FhA/y8ICAkKAC1QHIyz/4KM8WyXAgyMsBE/QA9ADLAMmABiMATTP1MTu/LhklMTugH6ANQwJxA0WfAFjhMBpEQzAshQBc8WE8s/zMzMye1Ukl8F4gCiMHAF1DCON4BA9JZvpSCOKQikIIEA+r6T8sGP3oEBkyGgUye78vQC+gDUMCJURjDwBSW6kwSkBN4Gkmwh4rPmMDRANMhQBc8WE8s/zMzMye1UACgD+kAwQzTIUAXPFhPLP8zMzMntVAAbPkAdMjLAhLKB8v/ydCAAPRa8ANwIfAEd4AYyMsFWM8WUAT6AhPLaxLMzMlx+wCACASAPEAAlvILfaiaH0gaZ/qamoYLehqGCxABDuLXTHtRND6QNM/1NTUMBAkXwTQ1DHUMNBxyMsHAc8WzMmAIBIBESAC+12v2omh9IGmf6mpqGDYg6GmH6Yf9IBhAALbT0faiaH0gaZ/qamoYCi+CeAG4APgCQFlTuZA==")
	invalidContractData, _ := boc.DeserializeBocBase64("te6cckEBAQEAAgAAAEysuc0=")
	args := make([]StackEntry, 0)

	runTvm, err := RunTvm(code[0], invalidContractData[0], "get_collection_data", args)
	_ = runTvm

	fmt.Println(runTvm)
	fmt.Println(err)
}
