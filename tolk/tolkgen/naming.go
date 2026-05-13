package tolkgen

import (
	"strings"

	"github.com/tonkeeper/tongo/utils"
)

func safeGoIdent(name string) string {
	ident := utils.ToCamelCase(name)
	if ident == "" {
		return "_"
	}
	return ident
}
func safePublicField(name string) string {
	return safeGoIdent(name)
}
func enumItemIdent(typeIdent, name string) string {
	return typeIdent + safeGoIdent(name)
}

func safeErrorIdent(name string) string {
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9')
	})
	if len(parts) == 0 {
		return "_"
	}
	var b strings.Builder
	for _, part := range parts {
		if isAllUpperASCII(part) {
			part = strings.ToLower(part)
		}
		b.WriteString(utils.ToCamelCase(part))
	}
	ident := b.String()
	if ident == "" {
		return "_"
	}
	if ident[0] >= '0' && ident[0] <= '9' {
		return "_" + ident
	}
	return ident
}

func isAllUpperASCII(s string) bool {
	hasLetter := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			return false
		}
		if c >= 'A' && c <= 'Z' {
			hasLetter = true
		}
	}
	return hasLetter
}
