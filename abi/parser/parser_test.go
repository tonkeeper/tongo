package parser

import (
	"fmt"
	"os"
	"testing"
)

func TestParseMethod(t *testing.T) {
	i, err := ParseMethod([]byte(METHOD))
	fmt.Println(i, err)
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
