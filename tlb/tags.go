package tlb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/startfellows/tongo/boc"
)

type sumTag struct {
	Name string
	Len  int
	Val  uint64
}

type tag struct {
	Len        int
	IsRef      bool
	IsOptional bool
}

func parseTag(s string) (tag, error) {
	var t tag
	if len(s) == 0 {
		return t, nil
	}
	if s[0] == '^' {
		t.IsRef = true
		s = strings.TrimSpace(s[1:])
	}
	if len(s) == 0 {
		return t, nil
	}
	if strings.Contains(s, "#") || strings.Contains(s, "$") {
		return t, nil
	}
	_, err := fmt.Sscanf(s, "%dbits", &t.Len)
	if err != nil {
		_, err = fmt.Sscanf(s, "%dbytes", &t.Len)
		return t, err
	}
	return t, err
}

func encodeSumTag(c *boc.Cell, tag string) error {
	t, err := parseSumTag(tag)
	if err != nil {
		return err
	}
	err = c.WriteUint(t.Val, t.Len)
	if err != nil {
		return err
	}
	return nil
}

func parseSumTag(s string) (sumTag, error) {
	a := strings.Split(s, "$")
	if len(a) == 2 {
		x, err := strconv.ParseUint(a[1], 2, 32)
		if err != nil {
			return sumTag{}, err
		}
		return sumTag{a[0], len(a[1]), x}, nil
	}
	a = strings.Split(s, "#")
	if len(a) == 2 {
		x, err := strconv.ParseUint(a[1], 16, 32)
		if err != nil {
			return sumTag{}, err
		}
		return sumTag{a[0], len(a[1]) * 4, x}, nil
	}
	return sumTag{}, fmt.Errorf("invalid tag")
}

func decodeHashmapTag(tag string) (int, error) {
	var ln int
	if tag == "" {
		return 0, fmt.Errorf("empty hashmap tag")
	}
	_, err := fmt.Sscanf(tag, "%dbits", &ln)
	if err != nil {
		return 0, err
	}
	return ln, nil
}

func decodeVarUIntegerTag(tag string) (int, error) {
	var n int
	if tag == "" {
		return 0, fmt.Errorf("empty varuint tag")
	}
	_, err := fmt.Sscanf(tag, "%dbytes", &n)
	if err != nil {
		return 0, err
	}
	return n, nil
}
