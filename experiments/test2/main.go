package main

import (
	"fmt"
	"github.com/startfellows/tongo/tl/parser"
)

var SOURCE = `

`

func main() {
	t, err := parser.Parse(SOURCE)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(t.Declarations))

}
