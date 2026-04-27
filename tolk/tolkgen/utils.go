package tolkgen

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func quoteList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	var b strings.Builder
	for i, item := range items {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strconv.Quote(item))
	}
	return b.String()
}

// deriveImportBlock infers the import block by scanning the generated code body
// for package-qualified references (e.g. "boc.", "tlb.", "ton.").
func deriveImportBlock(code string, extra map[string]string) string {
	type entry struct {
		marker string
		line   string
	}
	stdlib := []entry{
		{`context.`, `"context"`},
		{`fmt.`, `"fmt"`},
	}
	tongo := []entry{
		{`boc.`, `"github.com/tonkeeper/tongo/boc"`},
		{`liteclient.`, `"github.com/tonkeeper/tongo/liteclient"`},
		{`tlb.`, `"github.com/tonkeeper/tongo/tlb"`},
		{`ton.`, `"github.com/tonkeeper/tongo/ton"`},
	}
	var lines []string
	for _, c := range stdlib {
		if strings.Contains(code, c.marker) {
			lines = append(lines, "\t"+c.line)
		}
	}
	extraAliases := make([]string, 0, len(extra))
	for alias := range extra {
		extraAliases = append(extraAliases, alias)
	}
	slices.Sort(extraAliases)
	for _, alias := range extraAliases {
		if strings.Contains(code, alias+".") {
			lines = append(lines, fmt.Sprintf("\t%s %q", alias, extra[alias]))
		}
	}
	for _, c := range tongo {
		if strings.Contains(code, c.marker) {
			lines = append(lines, "\t"+c.line)
		}
	}
	return strings.Join(lines, "\n")
}
