package main

import (
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	outputDir := flag.String("output", "tlb", "directory to write generated files")
	to := flag.Int("to", 10, "maximum tensor size to generate")
	flag.Parse()

	var b strings.Builder
	fmt.Fprintf(&b, `// Code generated - DO NOT EDIT.

package tlb

const MaxTensorSize = %d

`, *to)

	for size := 0; size <= *to; size++ {
		generateTensor(&b, size)
	}

	formatted, err := format.Source([]byte(b.String()))
	if err != nil {
		panic(err)
	}
	outPath := filepath.Join(*outputDir, "tensor.go")
	if err := os.WriteFile(outPath, formatted, 0644); err != nil {
		panic(err)
	}
	fmt.Println(outPath)
}

func generateTensor(b *strings.Builder, size int) {
	typeParams := make([]string, size)
	fields := make([]string, size)
	values := make([]string, size)
	assignments := make([]string, size)
	for i := 0; i < size; i++ {
		typeParams[i] = fmt.Sprintf("T%d", i)
		fields[i] = fmt.Sprintf("\tV%d T%d", i, i)
		values[i] = fmt.Sprintf("v%d T%d", i, i)
		assignments[i] = fmt.Sprintf("V%d: v%d", i, i)
	}

	data := struct {
		Size        int
		TypeParams  string
		TypeArgs    string
		Fields      []string
		Values      string
		Assignments string
	}{
		Size:        size,
		TypeParams:  fmtTypeParams(typeParams),
		TypeArgs:    fmtTypeArgs(typeParams),
		Fields:      fields,
		Values:      strings.Join(values, ", "),
		Assignments: strings.Join(assignments, ", "),
	}

	if err := tensorTemplate.Execute(b, data); err != nil {
		panic(err)
	}
}

var tensorTemplate = template.Must(template.New("tensor").Parse(
	`type Tensor{{.Size}}{{.TypeParams}}{{if .Fields}} struct {
{{range .Fields}}{{.}}
{{end}}}{{else}} struct{}{{end}}

type ShapedTuple{{.Size}}{{.TypeParams}} = Tensor{{.Size}}{{.TypeArgs}}

func MakeTensor{{.Size}}{{.TypeParams}}({{.Values}}) Tensor{{.Size}}{{.TypeArgs}} {
	return Tensor{{.Size}}{{.TypeArgs}}{ {{- .Assignments -}} }
}

`))

func fmtTypeParams(params []string) string {
	if len(params) == 0 {
		return ""
	}
	return fmt.Sprintf("[%s any]", strings.Join(params, ", "))
}

func fmtTypeArgs(params []string) string {
	if len(params) == 0 {
		return ""
	}
	return fmt.Sprintf("[%s]", strings.Join(params, ", "))
}
