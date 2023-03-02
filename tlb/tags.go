package tlb

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/boc"
)

var ErrInvalidTag = errors.New("invalid tag")

type Tag struct {
	Name string
	Len  int
	Val  uint64
}

type tag struct {
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
	if strings.Contains(s, "bits") || strings.Contains(s, "bytes") {
		return t, fmt.Errorf("tag format '%v' is deprecated", s)
	}
	return t, nil
}

func encodeSumTag(c *boc.Cell, tag string) error {
	t, err := ParseTag(tag)
	if err != nil {
		return err
	}
	err = c.WriteUint(t.Val, t.Len)
	if err != nil {
		return err
	}
	return nil
}

func ParseTag(tagString string) (Tag, error) {
	var separatorPlace int
	var base int
	var length int
	for i, symbol := range tagString {
		if symbol == '$' {
			base = 2
			separatorPlace = i
			length = len(tagString) - separatorPlace - 1
			break
		}
		if symbol == '#' {
			base = 16
			separatorPlace = i
			length = (len(tagString) - separatorPlace - 1) * 4
			break
		}
	}
	if base == 0 || len(tagString) == separatorPlace+1 {
		return Tag{}, ErrInvalidTag
	}

	val, err := strconv.ParseUint(tagString[separatorPlace+1:], base, 32)
	return Tag{Name: tagString[:separatorPlace], Len: length, Val: val}, err
}
