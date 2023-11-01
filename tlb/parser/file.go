package parser

import (
	"bytes"
	_ "embed"
	"go/format"
	"os"
	"text/template"
)

type File struct {
	Name    string
	Package string
	Imports []string
	Code    string
}

//go:embed file.tmpl
var fileTemplateContent string

var (
	fileTmpl = template.Must(template.New("file").Parse(fileTemplateContent))
)

func (f *File) Save() error {
	file, err := os.Create(f.Name)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := fileTmpl.Execute(&buf, f); err != nil {
		return err
	}
	code, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	if _, err := file.Write(code); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil

}
