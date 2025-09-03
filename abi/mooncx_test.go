package abi

import (
	"encoding/json"
	"testing"
)

func TestMooncxSwapParams(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{
			name: "null next_fulfill cell",
			json: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":null,"NextFulfill":null,"NextReject":null}`,
			want: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":"","NextFulfill":null,"NextReject":null}`,
		},
		{
			name: "have next_fulfill cell with null payload inside",
			json: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":null,"NextFulfill":"b5ee9c7201010101002400004380059e76b719123ca2f935bb7b438070d4139959761c7f5ec616a39f902a464b5968","NextReject":null}`,
			want: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":"","NextFulfill":{"Recipient":"0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb","Payload":null},"NextReject":null}`,
		},
		{
			name: "have next_fulfill cell with payload some inside",
			json: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":null,"NextFulfill":"b5ee9c7201010201002d00014380059e76b719123ca2f935bb7b438070d4139959761c7f5ec616a39f902a464b597801000b0000177beb90","NextReject":null}`,
			want: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":"","NextFulfill":{"Recipient":"0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb","Payload":"b5ee9c7201010101000800000b0000177beb90"},"NextReject":null}`,
		},
		{
			name: "null next_fulfill data",
			json: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":null,"NextFulfill":null,"NextReject":null}`,
			want: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":"","NextFulfill":null,"NextReject":null}`,
		},
		{
			name: "have next_fulfill data with null payload inside",
			json: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":null,"NextFulfill":{"Recipient":"0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb","Payload":null},"NextReject":null}`,
			want: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":"","NextFulfill":{"Recipient":"0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb","Payload":null},"NextReject":null}`,
		},
		{
			name: "have next_fulfill data with payload some inside",
			json: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":null,"NextFulfill":{"Recipient":"0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb","Payload":"b5ee9c7201010101000800000b0000177beb90"},"NextReject":null}`,
			want: `{"MinOut":"3","Deadline":0,"Excess":"","Referral":"","NextFulfill":{"Recipient":"0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb","Payload":"b5ee9c7201010101000800000b0000177beb90"},"NextReject":null}`,
		},
		{
			name: "normal",
			json: `{"MinOut":"3","Deadline":0,"Excess":"0:05ea635b2a168cadfca174d72b12744a5b57d70378e6912e8a33b6b39bd3ee9d","Referral":"0:05ea635b2a168cadfca174d72b12744a5b57d70378e6912e8a33b6b39bd3ee9d","NextFulfill":{"Recipient":"0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb","Payload":"b5ee9c7201010101000800000b0000177beb90"},"NextReject":null}`,
			want: `{"MinOut":"3","Deadline":0,"Excess":"0:05ea635b2a168cadfca174d72b12744a5b57d70378e6912e8a33b6b39bd3ee9d","Referral":"0:05ea635b2a168cadfca174d72b12744a5b57d70378e6912e8a33b6b39bd3ee9d","NextFulfill":{"Recipient":"0:2cf3b5b8c891e517c9addbda1c0386a09ccacbb0e3faf630b51cfc8152325acb","Payload":"b5ee9c7201010101000800000b0000177beb90"},"NextReject":null}`,
		},
		{
			name: "problem",
			json: `{"MinOut":"1","Deadline":0,"Excess":"0:84c49bf79ed9a5720dbe973e2c654521cb4959415a3d003a16b8bc65c09e8bd6","Referral":"","NextFulfill":{"Payload":null,"Recipient":""},"NextReject":null}`,
			want: `{"MinOut":"1","Deadline":0,"Excess":"0:84c49bf79ed9a5720dbe973e2c654521cb4959415a3d003a16b8bc65c09e8bd6","Referral":"","NextFulfill":{"Recipient":"","Payload":null},"NextReject":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := MoonSwapParams{}
			if err := json.Unmarshal([]byte(tt.json), &actual); err != nil {
				t.Fatalf("Unmarshall failed: %v", err)
			}
			data, err := json.Marshal(actual)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if string(data) != tt.want {
				t.Errorf("Want %v, got %v", tt.want, string(data))
			}
		})
	}
}
