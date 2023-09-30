package parser

import (
	"encoding/xml"

	"github.com/tonkeeper/tongo/utils"
)

type ABI struct {
	Methods        []GetMethod `xml:"get_method"`
	Internals      []Message   `xml:"internal"`
	Externals      []Message   `xml:"external"`
	JettonPayloads []Message   `xml:"jetton_payload"`
	NFTPayloads    []Message   `xml:"nft_payload"`
	Interfaces     []Interface `xml:"interface"`
	Types          []string    `xml:"types"`
}

type Interface struct {
	Name    string `xml:"name,attr"`
	Methods []struct {
		Name    string `xml:"name,attr"`
		Version string `xml:"version,attr"`
	} `xml:"get_method"`
	Input struct {
		Internals []InterfaceMessage `xml:"internal"`
		Externals []InterfaceMessage `xml:"ext_in"`
	} `xml:"msg_in"`
	Output struct {
		Internals []InterfaceMessage `xml:"internal"`
		Externals []InterfaceMessage `xml:"ext_out"`
	} `xml:"msg_out"`
	CodeHashes []string `xml:"code_hash"`
	Inherits   string   `xml:"inherits,attr"`
}

type InterfaceMessage struct {
	Name string `xml:"name,attr"`
}

type Message struct {
	Name  string `xml:"name,attr"`
	Input string `xml:",chardata"`
	// FixedLength means that a destination type must have the same size in bits as the number of bits in a cell.
	FixedLength bool `xml:"fixed_length,attr"`
}

type GetMethod struct {
	Tag   xml.Name
	Input struct {
		StackValues []StackRecord `xml:",any"`
	} `xml:"input"`
	Name   string            `xml:"name,attr"`
	ID     int               `xml:"id,attr"`
	Output []GetMethodOutput `xml:"output"`
}

type GetMethodOutput struct {
	Version     string        `xml:"version,attr"`
	FixedLength bool          `xml:"fixed_length,attr"`
	Stack       []StackRecord `xml:",any"`
}

func (o GetMethodOutput) FullResultName(methodName string) string {
	version := ""
	if len(o.Version) > 0 {
		version = "_" + utils.ToCamelCase(o.Version)
	}
	return methodName + version + "Result"
}

type StackRecord struct {
	XMLName       xml.Name
	Name          string        `xml:"name,attr"`
	Nullable      bool          `xml:"nullable,attr"`
	List          bool          `xml:"list,attr"`
	Type          string        `xml:",chardata"`
	RequiredValue string        `xml:"required_value,attr"`
	SubStack      []StackRecord `xml:",any"`
}

func (m GetMethod) UsedByIntrospection() bool {
	return len(m.Input.StackValues) == 0
}

func (m GetMethod) GolangFunctionName() string {
	return utils.ToCamelCase(m.Name)
}

func ParseABI(s []byte) (ABI, error) {
	var abi ABI
	err := xml.Unmarshal(s, &abi)
	return abi, err
}

func ParseMethod(s []byte) (GetMethod, error) {
	var m GetMethod
	err := xml.Unmarshal(s, &m)
	return m, err
}
