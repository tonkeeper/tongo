package tolk

import (
	"errors"
	"fmt"
	"strconv"
)

func binHexToUint64(s string) (uint64, error) {
	if len(s) <= 2 {
		return 0, errors.New("number length must be greater than 2")
	}

	if s[1] == 'b' {
		val, err := strconv.ParseUint(s[2:], 2, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid bin number: %v", err)
		}
		return val, nil
	} else if s[1] == 'x' {
		val, err := strconv.ParseUint(s[2:], 16, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid hex number: %v", err)
		}
		return val, nil
	} else {
		return 0, fmt.Errorf("invalid number base, must be either bin or hex, got")
	}
}

func PrefixToUint(prefix string) (uint64, error) {
	if prefix == "" {
		return 0, errors.New("invalid prefix")
	}

	if len(prefix) == 1 {
		intPrefix, err := strconv.ParseUint(prefix, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid dec prefix: %v", err)
		}

		return intPrefix, nil
	}

	if len(prefix) == 2 {
		return 0, fmt.Errorf("prefix tag len must be either 1 or >2")
	}

	var intPrefix uint64
	var err error
	if prefix[1] == 'b' {
		intPrefix, err = strconv.ParseUint(prefix[2:], 2, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid bin prefix: %v", err)
		}
	} else if prefix[1] == 'x' {
		intPrefix, err = strconv.ParseUint(prefix[2:], 16, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid hex prefix: %v", err)
		}
	} else {
		return 0, fmt.Errorf("prefix tag must be either binary or hex format")
	}

	return intPrefix, nil
}
