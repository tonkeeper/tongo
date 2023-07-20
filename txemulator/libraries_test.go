package txemulator

import (
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
)

func TestFindLibraries(t *testing.T) {
	tests := []struct {
		name string
		boc  string
		want []tongo.Bits256
	}{
		{
			name: "with library",
			boc:  "te6ccgEBAQEAIwAIQgJYfMeJ7/HIT0bsN5fkX8gJoU/1riTx4MemqZzJ3JBh/w==",
			want: []tongo.Bits256{
				tongo.MustParseHash("587CC789EFF1C84F46EC3797E45FC809A14FF5AE24F1E0C7A6A99CC9DC9061FF"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell, err := boc.DeserializeSinglRootBase64(tt.boc)
			if err != nil {
				t.Fatalf("DeserializeSinglRootBase64() failed: %v", err)
			}
			got, err := FindLibraries(cell)
			if err != nil {
				t.Fatalf("FindLibraries() failed: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindLibraries() got = %v, want %v", got, tt.want)
			}
		})
	}
}
