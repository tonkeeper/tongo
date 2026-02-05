package tolk

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

func binHexToUint64(s string) (uint64, error) {
	if len(s) <= 2 {
		return 0, errors.New("number length must be greater than 2")
	}

	if s[1] == 'b' {
		val, err := strconv.ParseUint(s[2:], 2, 64)
		if err != nil {
			return 0, err
		}
		return val, nil
	} else if s[1] == 'x' {
		val, err := strconv.ParseUint(s[2:], 16, 64)
		if err != nil {
			return 0, err
		}
		return val, nil
	} else {
		return 0, fmt.Errorf("invalid number base, must be either bin or hex, got")
	}
}

func binDecHexToUint(num string) (*big.Int, error) {
	if len(num) == 0 {
		return nil, fmt.Errorf("number string is empty")
	}

	if len(num) == 1 {
		val, ok := new(big.Int).SetString(num, 10)
		if !ok {
			return nil, fmt.Errorf("canont convert %s to int", num)
		}
		return val, nil
	}

	if num[1] == 'b' {
		val, ok := new(big.Int).SetString(num, 2)
		if !ok {
			return nil, fmt.Errorf("canont convert %s to int", num)
		}
		return val, nil
	} else if num[1] == 'x' {
		val, ok := new(big.Int).SetString(num, 16)
		if !ok {
			return nil, fmt.Errorf("canont convert %s to int", num)
		}
		return val, nil
	} else {
		val, ok := new(big.Int).SetString(num, 10)
		if !ok {
			return nil, fmt.Errorf("canont convert %s to int", num)
		}
		return val, nil
	}
}
