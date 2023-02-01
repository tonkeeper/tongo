package parser

import (
	"bytes"
	"encoding/xml"

	"github.com/tonkeeper/tongo/utils"
)

type Interface struct {
	Name      string      `xml:"name,attr"`
	Group     string      `xml:"group,attr"`
	Methods   []GetMethod `xml:"get_method"`
	Internals []Internal  `xml:"internal"`
	Externals []Internal  `xml:"external"`
	Types     string      `xml:"types"`
}

type External struct {
	Name    string   `xml:"name,attr"`
	Input   string   `xml:"input"`
	Outputs []string `xml:"output"`
}

type Internal struct {
	Name    string   `xml:"name,attr"`
	Input   string   `xml:"input"`
	Outputs []string `xml:"output"`
}

type GetMethod struct {
	Tag   xml.Name
	Input struct {
		StackValues []StackRecord `xml:",any"`
	} `xml:"input"`
	Name        string            `xml:"name,attr"`
	Callback    bool              `xml:"callback,attr"`
	FixedLength bool              `xml:"fixed_length,attr"`
	ID          int               `xml:"id,attr"`
	Output      []GetMethodOutput `xml:"output"`
	// GolangName defines a name of a golang function generated to execute this get method.
	GolangName string `xml:"golang_name,attr"`
}

type GetMethodOutput struct {
	Version     string        `xml:"version,attr"`
	FixedLength bool          `xml:"fixed_length,attr"`
	Stack       []StackRecord `xml:",any"`
}

type StackRecord struct {
	XMLName  xml.Name
	Name     string `xml:"name,attr"`
	Nullable bool   `xml:"nullable,attr"`
	Type     string `xml:",chardata"`
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

func ParseInterface(s []byte) ([]Interface, error) {
	var i struct {
		List []Interface `xml:"interface"`
	}

	if !bytes.HasPrefix(bytes.TrimSpace(s), []byte("<interfaces>")) {
		s = append(append([]byte("<interfaces>"), s...), []byte("</interfaces>")...)
	}
	err := xml.Unmarshal(s, &i)
	return i.List, err
}

func ParseMethod(s []byte) (GetMethod, error) {
	var m GetMethod
	err := xml.Unmarshal(s, &m)
	return m, err
}
