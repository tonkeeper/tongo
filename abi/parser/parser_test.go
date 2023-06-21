package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"
)

var METHOD = `



        <get_method name="list_nominators">
            <output>
                <tuple list="true">
                    <int name="address">bits256</int>
                    <tinyint name="amount">uint64</tinyint>
                    <tinyint name="pending_deposit_amount">uint64</tinyint>
                    <tinyint name="withdraw_requested">bool</tinyint>
                </tuple>
            </output>
        </get_method>
`

func TestParseMethod(t *testing.T) {
	i, err := ParseMethod([]byte(METHOD))
	fmt.Println(err)
	fmt.Printf("%#v", i)
}

func TestParseInterface(t *testing.T) {
	b, err := os.ReadFile("../schemas/known.xml")
	if err != nil {
		t.Fatal(err)
	}
	i, err := ParseABI(b)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(i)
}

func TestParseArray(t *testing.T) {
	a := `
<ololo interface="A" interface="B">
text
</ololo>
`
	var A struct {
		Interface []string `xml:"interface,attr"`
		Text      string   `xml:",cdata"`
	}
	err := xml.Unmarshal([]byte(a), &A)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", A.Interface[1])

}
