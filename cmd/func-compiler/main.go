package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/tonkeeper/tongo/code"
)

func main() {
	pathToCode := "./"
	if len(os.Args) > 1 {
		pathToCode = os.Args[1]
	}

	entries, err := os.ReadDir(pathToCode)
	if err != nil {
		log.Fatal(err)
	}
	files := make(map[string]string, len(entries))
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".fc") {
			continue
		}
		if e.Type()&os.ModeType != 0 {
			continue
		}
		b, err := os.ReadFile(path.Join(pathToCode, e.Name()))
		if err != nil {
			log.Fatal(err)
		}
		files[e.Name()] = string(b)
	}
	if len(files) == 0 {
		log.Fatal("no .fc in this directory")
	}
	compiler := code.NewFunCCompiler()
	fift, boc, err := compiler.Compile(files)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("FIFT:\n%v\n\nBOC:\n%x\n", fift, boc)

}
