package tlb

import (
	"reflect"
	"testing"
)

func Test_decodeHashmapTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		want    int
		wantErr bool
	}{
		{
			name:    "invalid tag",
			tag:     "32bytes",
			wantErr: true,
		},
		{
			name: "all good",
			tag:  "32bits",
			want: 32,
		},
		{
			name:    "failed to parse 160x",
			tag:     "160xbits",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeHashmapTag(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeHashmapTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("decodeHashmapTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		want    tag
		wantErr bool
	}{
		{
			name: "parse bits",
			tag:  "24bits",
			want: tag{IsRef: false, Len: 24},
		},
		{
			name: "parse bytes",
			tag:  "56bytes",
			want: tag{IsRef: false, Len: 56},
		},
		{
			name:    "invalid tag",
			tag:     "56hello",
			want:    tag{IsRef: false, Len: 56},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTag(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeVarUIntegerTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		want    int
		wantErr bool
	}{
		{
			name: "all good",
			tag:  "45bytes",
			want: 45,
		},
		{
			name:    "invalid tag",
			tag:     "45hello",
			wantErr: true,
		},
		{
			name:    "failed to parse 45x",
			tag:     "45xbytes",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeVarUIntegerTag(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeVarUIntegerTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("decodeVarUIntegerTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}
