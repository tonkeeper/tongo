package utils

func ConcatPrefixAndSuffixIfExists(prefix, suffix []byte) []byte {
	if len(suffix) == 0 {
		return prefix
	}
	prefix = prefix[:len(prefix)-1] // remove '}'
	suffix[0] = ','                 // replace '{' with ','
	result := make([]byte, 0, len(prefix)+len(suffix))
	result = append(result, prefix...)
	result = append(result, suffix...)
	return result
}
