package abi

// tolkIntMsgs, tolkExtInMsgs, tolkExtOutMsgs are added to
// IntMsgs / ExtInMsgs / ExtOutMsgs, so that ABI generated with Tolk can be plugged in
var tolkIntMsgs = map[ContractInterface][]msgDecoderFunc{}
var tolkExtInMsgs = map[ContractInterface][]msgDecoderFunc{}
var tolkExtOutMsgs = map[ContractInterface][]msgDecoderFunc{}

// tolkMethods (->methodInvocationOrder) and tolkInterfaceOrder(->contractInterfacesOrder)
// are for contractInspector in NewContractInspector
var tolkMethods []MethodDescription
var tolkInterfaceOrder []InterfaceDescription

// tolkContractInterfaceStrings and tolkContractInterfaceFromString are for the
// ContractInterface.String()  and ContractInterfaceFromString()
var tolkContractInterfaceStrings = map[ContractInterface]string{}
var tolkContractInterfaceFromString = map[string]ContractInterface{}
