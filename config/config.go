package config

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type liteServerConfig struct {
	Ip   int64        `json:"ip"`
	Port int64        `json:"port"`
	ID   liteServerId `json:"id"`
}

type liteServerId struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type configGlobal struct {
	LiteServers []liteServerConfig `json:"liteservers"`
	//Validator   ValidatorConfig  `json:"validator"`
}

// GlobalConfigurationFile contains global configuration of the TON Blockchain.
// It is shared by all nodes and includes information about network, init block, hardforks, etc.
type GlobalConfigurationFile struct {
	LiteServers []LiteServer
}

// LiteServer TODO: clarify struct
type LiteServer struct {
	Host string
	Key  string
}

func ParseConfigFile(path string) (*GlobalConfigurationFile, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	return ParseConfig(jsonFile)
}

func convertToLiteServerOptions(server liteServerConfig) (LiteServer, error) {
	if server.ID.Type != "pub.ed25519" {
		return LiteServer{}, fmt.Errorf("not pub.ed25519 liteserver ID. Other types not supported")
	}
	if server.Ip > 0xFFFF_FFFF {
		return LiteServer{}, fmt.Errorf("only IPv4 supported")
	}
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, uint32(server.Ip))
	return LiteServer{
		Host: fmt.Sprintf("%v.%v.%v.%v:%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3], server.Port),
		Key:  server.ID.Key,
	}, nil
}

func ParseConfig(data io.Reader) (*GlobalConfigurationFile, error) {
	var conf configGlobal
	err := json.NewDecoder(data).Decode(&conf)
	if err != nil {
		return nil, err
	}
	var options GlobalConfigurationFile
	for _, server := range conf.LiteServers {
		ls, err := convertToLiteServerOptions(server)
		if err != nil {
			continue
		}
		options.LiteServers = append(options.LiteServers, ls)
	}
	if len(options.LiteServers) == 0 {
		return nil, fmt.Errorf("no one supported liteservers")
	}
	return &options, nil
}
