package parser

import (
	"os"
	"testing"
)

func TestGenerateGetMethodsGolang(t *testing.T) {
	b, err := os.ReadFile("../known.xml")
	if err != nil {
		t.Fatal(err)
	}
	interfaces, err := ParseInterface(b)
	if err != nil {
		t.Fatal(err)
	}
	g := NewGenerator(nil, "MethodsScanner")
	for _, i := range interfaces {
		if i.Types != "" {
			err = g.RegisterTypes(i.Types)
			if err != nil {
				t.Fatal(err)
			}
		}
		_, err := g.GetMethods(i.Methods)
		if err != nil {
			t.Fatal(err)
		}
	}

}
