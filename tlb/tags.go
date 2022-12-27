package tlb

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/startfellows/tongo/boc"
)

var ErrInvalidTag = errors.New("invalid tag")

type sumTag struct {
	Name string
	Len  int
	Val  uint64
}

type tag struct {
	Len   int
	IsRef bool
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
	if strings.HasSuffix(s, "bits") {
		length, err := strconv.Atoi(s[:len(s)-4])
		if err != nil {
			return t, ErrInvalidTag
		}
		t.Len = length
		return t, nil
	} else if strings.HasSuffix(s, "bytes") {
		length, err := strconv.Atoi(s[:len(s)-5])
		if err != nil {
			return t, ErrInvalidTag
		}
		t.Len = length
		return t, nil
	}
	return t, ErrInvalidTag
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
	return sumTag{}, ErrInvalidTag
}

func decodeHashmapTag(tag string) (int, error) {
	if tag == "" {
		return 0, fmt.Errorf("empty hashmap tag")
	}
	if !strings.HasSuffix(tag, "bits") {
		return 0, ErrInvalidTag
	}
	n, err := strconv.Atoi(tag[:len(tag)-4])
	if err != nil {
		return 0, ErrInvalidTag
	}
	return n, nil
}

func decodeVarUIntegerTag(tag string) (int, error) {
	if tag == "" {
		return 0, fmt.Errorf("empty varuint tag")
	}
	if !strings.HasSuffix(tag, "bytes") {
		return 0, ErrInvalidTag
	}
	n, err := strconv.Atoi(tag[:len(tag)-5])
	if err != nil {
		return 0, ErrInvalidTag
	}
	return n, err
}
