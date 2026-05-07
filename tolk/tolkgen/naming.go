package tolkgen

import "strings"

func safeGoIdent(name string) string {
	return name
}
func safePublicField(name string) string {
	parts := strings.Split(name, "_")

	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}

	return strings.Join(parts, "")
}
func enumItemIdent(typeIdent, name string) string {
	return typeIdent + name
}
