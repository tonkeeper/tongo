package utils

func MethodIdFromName(methodName string) int {
	return int(Crc16String(methodName)&0xffff) | 0x10000
}
