package tlb

import "github.com/tonkeeper/tongo/boc"

func cloneCell(c *boc.Cell) *boc.Cell {
	if c == nil {
		return nil
	}
	data, err := c.ToBoc()
	if err != nil {
		panic(err)
	}
	cloned, err := boc.DeserializeBoc(data)
	if err != nil || len(cloned) == 0 {
		panic(err)
	}
	return cloned[0]
}
