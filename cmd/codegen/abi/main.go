package main

import (
	"flag"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tonkeeper/tongo/abi/parser"
)

const HEADER = `// Code generated - DO NOT EDIT.

package abi

import (
%v
)

`

func mergeMethods(methods []parser.GetMethod) ([]parser.GetMethod, error) {
	methodsMap := map[string]parser.GetMethod{}
	var golangNamedMethods []parser.GetMethod
	for _, method := range methods {
		current, ok := methodsMap[method.Name]
		if !ok {
			methodsMap[method.Name] = method
			continue
		}
		if len(current.Input.StackValues) > 0 || len(method.Input.StackValues) > 0 {
			return nil, fmt.Errorf("method '%s' has a version with input params, it has to be defined with golang_name to avoid collision", method.Name)
		}
		current.Output = append(current.Output, method.Output...)
		methodsMap[method.Name] = current
	}
	var results []parser.GetMethod
	for _, method := range methodsMap {
		results = append(results, method)
	}
	results = append(results, golangNamedMethods...)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})
	return results, nil
}

func main() {
	schemasDir := flag.String("schemas", "abi/schemas", "directory containing XML TLB schemas")
	outputDir := flag.String("output", "abi", "directory to write generated files")
	flag.Parse()

	var abi parser.ABI
	var methods []parser.GetMethod
	filepath.Walk(*schemasDir, func(path string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(info.Name(), ".xml") {
			return nil
		}
		scheme, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		a, err := parser.ParseABI(scheme)
		if err != nil {
			panic(err)
		}
		methods = append(methods, a.Methods...)
		abi.ExtOut = append(abi.ExtOut, a.ExtOut...)
		abi.ExtIn = append(abi.ExtIn, a.ExtIn...)
		abi.Internals = append(abi.Internals, a.Internals...)
		abi.JettonPayloads = append(abi.JettonPayloads, a.JettonPayloads...)
		abi.NFTPayloads = append(abi.NFTPayloads, a.NFTPayloads...)
		abi.Types = append(abi.Types, a.Types...)
		abi.Interfaces = append(abi.Interfaces, a.Interfaces...)
		return nil
	})

	methods, err := mergeMethods(methods)
	if err != nil {
		panic(err)
	}
	abi.Methods = methods

	gen, err := parser.NewGenerator(nil, abi)
	if err != nil {
		panic(err)
	}

	types := gen.CollectedTypes()
	msgDecoder := gen.GenerateMsgDecoder()

	getMethods, simpleMethods, err := gen.GetMethods()
	if err != nil {
		panic(err)
	}
	invocationOrder, err := gen.RenderInvocationOrderList(simpleMethods)
	if err != nil {
		panic(err)
	}
	messagesMD, err := gen.RenderMessagesMD()
	if err != nil {
		panic(err)
	}

	jettons, err := gen.RenderJetton()
	if err != nil {
		panic(err)
	}

	nfts, err := gen.RenderNFT()
	if err != nil {
		panic(err)
	}

	contractErrors, err := gen.RenderContractErrors()
	if err != nil {
		panic(err)
	}

	for _, f := range [][]string{
		{types, "types.go", `"github.com/tonkeeper/tongo/tlb"`, `"fmt"`, `"encoding/json"`},
		{msgDecoder, "messages_generated.go", `"github.com/tonkeeper/tongo/tlb"`},
		{getMethods, "get_methods.go", `"context"`, `"fmt"`, `"github.com/tonkeeper/tongo/ton"`, `"github.com/tonkeeper/tongo/boc"`, `"github.com/tonkeeper/tongo/tlb"`},
		{invocationOrder, "interfaces.go", `"github.com/tonkeeper/tongo/ton"`},
		{jettons, "jetton_msg_types.go", `"github.com/tonkeeper/tongo/boc"`, `"github.com/tonkeeper/tongo/tlb"`},
		{nfts, "nfts_msg_types.go", `"github.com/tonkeeper/tongo/boc"`, `"github.com/tonkeeper/tongo/tlb"`},
		{contractErrors, "contracts_errors.go"},
	} {
		file, err := os.Create(filepath.Join(*outputDir, f[1]))
		if err != nil {
			panic(err)
		}
		code := []byte(fmt.Sprintf(HEADER, strings.Join(f[2:], "\n")) + f[0])
		formatedCode, err := format.Source(code)
		if err != nil {
			formatedCode = code
			//panic(err)
		}
		_, err = file.Write(formatedCode)
		if err != nil {
			panic(err)
		}
		err = file.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println(filepath.Join(*outputDir, f[1]))
	}
	mdPath := filepath.Join(*outputDir, "messages.md")
	if err := os.WriteFile(mdPath, []byte(messagesMD), 0644); err != nil {
		panic(err)
	}
	fmt.Println(mdPath)
}
