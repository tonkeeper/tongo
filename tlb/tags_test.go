package tlb

import (
	"testing"
)

func TestTag(t *testing.T) {
	type test struct {
		Tag    string
		Result uint64
	}
	tags := []test{{"test#abcdabcd", 2882382797}, {"test$10101", 21}, {"#abcd", 43981}, {"$101", 5}}

	for _, tg := range tags {
		t1, err := ParseTag(tg.Tag)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if t1.Val != tg.Result {
			t.Fatalf("invalid tag")
		}
	}
}
