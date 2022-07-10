package utils

import "testing"

func TestCrc16String(t *testing.T) {
	if Crc16String("get_seq") != 38947 {
		t.Fail()
	}
}
