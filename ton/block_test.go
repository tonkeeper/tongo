package ton

import (
	"os"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

func TestBlockInfo_GetParents(t *testing.T) {
	tests := []struct {
		name     string
		want     []BlockIDExt
		filename string
		wantErr  bool
	}{
		{
			filename: "testdata/raw-13516764.bin",
			want: []BlockIDExt{
				{
					BlockID: BlockID{
						Workchain: 0,
						Shard:     0xa000000000000000, // a000000000000000
						Seqno:     13516763},
					RootHash: MustParseHash("617f643f15a42f28018e3e3c89f14b952a0d67fa90968ae5360a51b96c6a1c42"),
					FileHash: MustParseHash("563aa5f3d51585b95c0c89bf6c4e39455f4c121269521c1c5b6dc07f03c5d230"),
				},
				{
					BlockID: BlockID{
						Workchain: 0,
						Shard:     0xe000000000000000, // e000000000000000
						Seqno:     13516699,
					},
					RootHash: MustParseHash("032b1bf3016c9b71816c52f207c4cd79d75541f78eacb11cac2ea7b77d2a603d"),
					FileHash: MustParseHash("8dcab64721513f3db73a081dd61cdf51d7fec79347ab348d43ffb8bc052a8db3"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawData, err := os.ReadFile(tt.filename)
			if err != nil {
				t.Errorf("ReadFile() failed: %v", err)
			}
			cell, err := boc.DeserializeBoc(rawData)
			if err != nil {
				t.Errorf("DeserializeBoc() failed: %v", err)
			}
			var data tlb.Block
			if err = tlb.Unmarshal(cell[0], &data); err != nil {
				t.Errorf("Unmarshal() failed: %v", err)
			}
			got, err := GetParents(data.Info)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParents() got = %#v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBlockID(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    BlockID
		wantErr bool
	}{
		{
			name: "all good",
			s:    "(-1,8000000000000000,29537038)",
			want: BlockID{
				Workchain: -1,
				Shard:     0x8000000000000000,
				Seqno:     29537038,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBlockID(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBlockID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseBlockID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
