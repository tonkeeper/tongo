package abi

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGenerateGetMethodsGolang(t *testing.T) {
	b, err := ioutil.ReadFile("known.xml")
	if err != nil {
		t.Fatal(err)
	}
	interfaces, err := ParseInterface(b)
	if err != nil {
		t.Fatal(err)
	}
	g := NewGenerator(nil, "MethodsScanner")
	for _, i := range interfaces {
		s, err := g.GetMethods(i.Methods)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(s)
	}
	fmt.Println(g.CollectedTypes())

}
