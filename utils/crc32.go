package utils

import (
	"hash/crc32"
)

// Crc32String returns a crc32 checksum of a string.
// Crc32String(value) returns the same uint32 as FunC's "<value>"c string literal.
func Crc32String(data string) uint32 {
	return crc32.Checksum([]byte(data), crc32.IEEETable)
}
