package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ParseLiteServersEnvVar parses the given string and returns a list of lite servers.
// The string is a comma-separated list of servers in the following format: "ip:port:public-key".
// An example of such an env variable:
// LITE_SERVERS="127.0.0.1:22095:6PGkPQSbyFp12esf1NqmDOaLoFA8i9+Mp5+cAx5wtTU=,192.168.0.17:14095:6PGkPQSbyFp12esf1NqmDOaLoFA8i9+Mp5+cAx5wtTU="
func ParseLiteServersEnvVar(str string) ([]LiteServer, error) {
	if len(str) == 0 {
		return []LiteServer{}, nil
	}
	var servers []LiteServer
	for _, s := range strings.Split(str, ",") {
		params := strings.Split(s, ":")
		if len(params) != 3 {
			return nil, fmt.Errorf("invalid liteserver string: %v", s)
		}
		ip := net.ParseIP(params[0])
		if ip == nil {
			return nil, fmt.Errorf("invalid lite server ip")
		}
		if ip.To4() == nil {
			return nil, fmt.Errorf("IPv6 not supported")
		}
		_, err := strconv.ParseInt(params[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid lite server port: %v", params[1])
		}
		servers = append(servers, LiteServer{
			Host: fmt.Sprintf("%v:%v", params[0], params[1]),
			Key:  params[2],
		})
	}
	return servers, nil
}
