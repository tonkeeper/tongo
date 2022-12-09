package main

import (
	"fmt"
	"github.com/startfellows/tongo/tlb/parser"
)

var SOURCE = `
                text#_ {n:#} data:(SnakeData ~n) = Text;
      
`

func main() {
	t, err := parser.Parse(SOURCE)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(t.Declarations))

}
