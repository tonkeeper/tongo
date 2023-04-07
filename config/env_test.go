package config

import (
	"reflect"
	"testing"
)

func TestParseLiteServersEnvVar(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    []LiteServer
		wantErr string
	}{
		{
			name: "empty string",
			str:  "",
			want: []LiteServer{},
		},
		{
			name: "two servers",
			str:  "127.0.0.1:22095:6PGkPQSbyFp12esf1+Mp5+cAx5wtTU=,192.168.0.17:14095:NqmDOaLoFA8i9+Mp5+cAx5wtTU=",
			want: []LiteServer{
				{Host: "127.0.0.1:22095", Key: "6PGkPQSbyFp12esf1+Mp5+cAx5wtTU="},
				{Host: "192.168.0.17:14095", Key: "NqmDOaLoFA8i9+Mp5+cAx5wtTU="},
			},
		},
		{
			name:    "error - invalid port",
			str:     "127.0.0.1:xxx:6PGkPQSbyFp12esf1+Mp5+cAx5wtTU=",
			wantErr: `invalid lite server port: xxx`,
		},
		{
			name:    "error - no key",
			str:     "127.0.0.1:999",
			wantErr: `invalid liteserver string: 127.0.0.1:999`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLiteServersEnvVar(tt.str)
			if len(tt.wantErr) > 0 {
				if err == nil {
					t.Fatalf("ParseLiteServersEnvVar() error is nil")
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("ParseLiteServersEnvVar() error = %v, wantErr %v", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseLiteServersEnvVar() error: %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParseLiteServersEnvVar() got = %v, want %v", got, tt.want)
			}
		})
	}
}
