package tolkgen

import "strings"

func safeGoIdent(name string) string {
	return name
}
func safePublicField(name string) string {
	// capitalize the first letter
	return strings.ToUpper(name[:1]) + name[1:]
}

func enumItemIdent(typeIdent, name string) string {
	return typeIdent + name
}
