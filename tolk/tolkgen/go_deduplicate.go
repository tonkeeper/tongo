package tolkgen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/tonkeeper/tongo/utils"
)

type generatedGoDef struct {
	key           string
	code          string
	path          string
	order         int
	constGroupKey string
	shared        bool
	removed       bool
}

type generatedGoFile struct {
	path    string
	pkgName string
	defs    []*generatedGoDef
}

func deduplicateGoDefinitions(paths []string, sharedPath string) error {
	files, err := parseGeneratedGoFiles(paths)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return removeGeneratedSharedFile(sharedPath)
	}

	byKey := make(map[string][]*generatedGoDef)
	var allDefs []*generatedGoDef
	for _, file := range files {
		for _, def := range file.defs {
			byKey[def.key] = append(byKey[def.key], def)
			allDefs = append(allDefs, def)
		}
	}

	sharedDefs := make(map[string]*generatedGoDef)
	for key, defs := range byKey {
		if len(defs) < 2 {
			continue
		}
		first := defs[0]
		for _, def := range defs[1:] {
			if def.code != first.code {
				return fmt.Errorf("duplicate Go definition %q differs between %s and %s", key, first.path, def.path)
			}
		}
		sharedDefs[key] = first
		for _, def := range defs {
			def.removed = true
		}
	}

	for _, file := range files {
		if !hasRemovedDefs(file.defs) {
			continue
		}
		if err := writeGeneratedGoFile(file.path, file.pkgName, remainingDefs(file.defs)); err != nil {
			return err
		}
	}

	var orderedShared []*generatedGoDef
	for _, def := range allDefs {
		sharedDef, ok := sharedDefs[def.key]
		if !ok || sharedDef.shared {
			continue
		}
		sharedDef.shared = true
		orderedShared = append(orderedShared, sharedDef)
	}
	if len(orderedShared) == 0 {
		return removeGeneratedSharedFile(sharedPath)
	}
	return writeGeneratedGoFile(sharedPath, files[0].pkgName, orderedShared)
}

func parseGeneratedGoFiles(paths []string) ([]generatedGoFile, error) {
	var files []generatedGoFile
	order := 0
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		fset := token.NewFileSet()
		parsed, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		file := generatedGoFile{path: path, pkgName: parsed.Name.Name}
		for _, decl := range parsed.Decls {
			defs, err := splitTopLevelDecl(fset, path, decl)
			if err != nil {
				return nil, err
			}
			for _, def := range defs {
				def.order = order
				order++
				file.defs = append(file.defs, def)
			}
		}
		files = append(files, file)
	}
	return files, nil
}

func splitTopLevelDecl(fset *token.FileSet, path string, decl ast.Decl) ([]*generatedGoDef, error) {
	switch d := decl.(type) {
	case *ast.GenDecl:
		if d.Tok == token.IMPORT {
			return nil, nil
		}
		defs := make([]*generatedGoDef, 0, len(d.Specs))
		for _, spec := range d.Specs {
			key, err := specDefinitionKey(d.Tok, spec)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
			code, err := renderGeneratedDecl(fset, &ast.GenDecl{Tok: d.Tok, Specs: []ast.Spec{spec}})
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
			defs = append(defs, &generatedGoDef{
				key:           key,
				code:          code,
				path:          path,
				constGroupKey: constGroupKey(fset, d.Tok, spec),
			})
		}
		return defs, nil
	case *ast.FuncDecl:
		key, err := funcDefinitionKey(fset, d)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		code, err := renderGeneratedDecl(fset, d)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		return []*generatedGoDef{{key: key, code: code, path: path}}, nil
	default:
		code, err := renderGeneratedDecl(fset, d)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		return []*generatedGoDef{{key: fmt.Sprintf("%T:%s", d, code), code: code, path: path}}, nil
	}
}

func specDefinitionKey(tok token.Token, spec ast.Spec) (string, error) {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		return "type:" + s.Name.Name, nil
	case *ast.ValueSpec:
		names := make([]string, 0, len(s.Names))
		for _, name := range s.Names {
			names = append(names, name.Name)
		}
		return strings.ToLower(tok.String()) + ":" + strings.Join(names, ","), nil
	default:
		return "", fmt.Errorf("unsupported %s spec %T", tok, spec)
	}
}

func constGroupKey(fset *token.FileSet, tok token.Token, spec ast.Spec) string {
	if tok != token.CONST {
		return ""
	}
	valueSpec, ok := spec.(*ast.ValueSpec)
	if !ok || valueSpec.Type == nil || len(valueSpec.Names) != 1 || len(valueSpec.Values) != 1 {
		return ""
	}
	typ, err := exprString(fset, valueSpec.Type)
	if err != nil {
		return ""
	}
	return "const:" + typ
}

func funcDefinitionKey(fset *token.FileSet, decl *ast.FuncDecl) (string, error) {
	if decl.Recv == nil || len(decl.Recv.List) == 0 {
		return "func:" + decl.Name.Name, nil
	}
	recv, err := receiverTypeKey(fset, decl.Recv.List[0].Type)
	if err != nil {
		return "", err
	}
	return "method:" + recv + "." + decl.Name.Name, nil
}

func receiverTypeKey(fset *token.FileSet, expr ast.Expr) (string, error) {
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	return exprString(fset, expr)
}

func exprString(fset *token.FileSet, node ast.Node) (string, error) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderGeneratedDecl(fset *token.FileSet, node any) (string, error) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		return "", err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(formatted)) + "\n", nil
}

func remainingDefs(defs []*generatedGoDef) []*generatedGoDef {
	remaining := make([]*generatedGoDef, 0, len(defs))
	for _, def := range defs {
		if !def.removed {
			remaining = append(remaining, def)
		}
	}
	return remaining
}

func hasRemovedDefs(defs []*generatedGoDef) bool {
	for _, def := range defs {
		if def.removed {
			return true
		}
	}
	return false
}

func writeGeneratedGoFile(path, pkgName string, defs []*generatedGoDef) error {
	var body strings.Builder
	for i := 0; i < len(defs); {
		if defs[i].constGroupKey == "" {
			body.WriteString(defs[i].code)
			body.WriteString("\n")
			i++
			continue
		}

		j := i + 1
		for j < len(defs) && defs[j].constGroupKey == defs[i].constGroupKey {
			j++
		}
		if j-i == 1 {
			body.WriteString(defs[i].code)
			body.WriteString("\n")
			i = j
			continue
		}

		body.WriteString(groupConstDefs(defs[i:j]))
		body.WriteString("\n")
		i = j
	}

	bodyCode := strings.TrimSpace(body.String())
	var content string
	if bodyCode == "" {
		content = fmt.Sprintf("// Code generated - DO NOT EDIT.\n\npackage %s\n", pkgName)
	} else if imports := deriveImportBlock(bodyCode, nil); imports != "" {
		content = fmt.Sprintf("// Code generated - DO NOT EDIT.\n\npackage %s\n\nimport (\n%s\n)\n\n%s\n", pkgName, imports, bodyCode)
	} else {
		content = fmt.Sprintf("// Code generated - DO NOT EDIT.\n\npackage %s\n\n%s\n", pkgName, bodyCode)
	}
	return utils.WriteFormattedGoCode(path, content)
}

func groupConstDefs(defs []*generatedGoDef) string {
	var b strings.Builder
	b.WriteString("const (\n")
	for _, def := range defs {
		line := strings.TrimSpace(def.code)
		line = strings.TrimPrefix(line, "const ")
		b.WriteString("\t")
		b.WriteString(line)
		b.WriteString("\n")
	}
	b.WriteString(")\n")
	return b.String()
}

func removeGeneratedSharedFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !strings.HasPrefix(string(data), "// Code generated - DO NOT EDIT.") {
		return fmt.Errorf("refuse to remove non-generated shared file %s", path)
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	fmt.Printf("%s\n", filepath.Clean(path))
	return nil
}
