package abi

import (
	"strings"
)

func (res GetFullDomainResult) EscapedDomain() string {
	return escapedDomain(res.Domain)
}

func (res GetDomainResult) EscapedDomain() string {
	return escapedDomain(res.Domain)
}

func escapedDomain(domain string) string {
	sl := strings.Split(domain, "\u0000")
	for i := 0; i < len(sl)/2; i++ {
		sl[i], sl[len(sl)-i-1] = sl[len(sl)-i-1], sl[i]
	}
	return strings.Trim(strings.Join(sl, "."), ".")
}
