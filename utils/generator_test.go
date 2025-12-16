package utils

import "testing"

func TestToSnakeCaseString(t *testing.T) {
	tests := []struct {
		data string
		want string
	}{
		{
			data: "CocoonTest",
			want: "cocoon_test",
		},
		{
			data: "ton.cocoon.test",
			want: "ton_cocoon_test",
		},
		{
			data: "cocoonTest",
			want: "cocoon_test",
		},
		{
			data: "cocoon123Test",
			want: "cocoon123_test",
		},
		{
			data: "Cocoon123test",
			want: "cocoon123test",
		},
		{
			data: "COcoon123Test",
			want: "c_ocoon123_test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			if got := ToSnakeCase(tt.data); got != tt.want {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
