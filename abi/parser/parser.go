package parser

import (
	"encoding/xml"

	"github.com/tonkeeper/tongo/utils"
)

type ABI struct {
	Methods   []GetMethod `xml:"get_method"`
	Internals []Message   `xml:"internal"`
	Externals []Message   `xml:"external"`
	Types     []string    `xml:"types"`
}

type Message struct {
	Name       string   `xml:"name,attr"`
	Input      string   `xml:",chardata"`
	Interfaces []string `xml:"interface,attr"`
}

type GetMethod struct {
	Tag   xml.Name
	Input struct {
		StackValues []StackRecord `xml:",any"`
	} `xml:"input"`
	Name       string            `xml:"name,attr"`
	Interfaces []string          `xml:"interface,attr"`
	ID         int               `xml:"id,attr"`
	Output     []GetMethodOutput `xml:"output"`
	// GolangName defines a name of a golang function generated to execute this get method.
	GolangName string `xml:"golang_name,attr"`
}

type GetMethodOutput struct {
	Version     string        `xml:"version,attr"`
	FixedLength bool          `xml:"fixed_length,attr"`
	Stack       []StackRecord `xml:",any"`
	Interface   string        `xml:"interface,attr"`
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
	if len(m.GolangName) > 0 {
		return m.GolangName
	}
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
