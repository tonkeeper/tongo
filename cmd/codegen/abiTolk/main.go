package main

import (
	"log"
	"os"

	"github.com/tonkeeper/tongo/tolk/tolkgen"
)

func main() {
	cfg := tolkgen.DefaultCodegenPipelineConfig()

	if len(os.Args) > 1 {
		cfg.SchemasDir = os.Args[1]
	}
	if len(os.Args) > 2 {
		cfg.OutputDir = os.Args[2]
	}
	if len(os.Args) > 3 {
		cfg.ABIOutputDir = os.Args[3]
	}

	if err := tolkgen.GenerateFromSchemas(cfg); err != nil {
		log.Fatal(err)
	}
}
