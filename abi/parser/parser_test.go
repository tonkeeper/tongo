package parser

import (
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
	b, err := os.ReadFile("../known.xml")
	if err != nil {
		t.Fatal(err)
	}
	i, err := ParseInterface(b)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(i)
}
