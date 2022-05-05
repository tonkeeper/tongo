package config

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type configFile struct {
	Config configFileConfig `json:"config"`
	//Keystore KeyStoreType           `json:"keystore_type"`
}

type configFileConfig struct {
	Config         configServer `json:"config"`
	BlockchainName string       `json:"blockchain_name"`
	//UseCallbacksForNetwork bool         `json:"use_callbacks_for_network"`
	//IgnoreCache            bool         `json:"ignore_cache"`
}

type liteServerConfig struct {
	Type string       `json:"@type"`
	Ip   int64        `json:"ip"`
	Port string       `json:"port"`
	ID   liteServerId `json:"id"`
}

type liteServerId struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type configServer struct {
	LiteServers []liteServerConfig `json:"liteservers"`
	//Validator   ValidatorConfig            `json:"validator"`
}

type Options struct {
	LiteServers []LiteServer
}

// LiteServer TODO: clarify struct
type LiteServer struct {
	Host string
	Key  string
}

func ParseConfigFile(path string) (*Options, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	var conf configFile
	err = json.NewDecoder(jsonFile).Decode(&conf)
	if err != nil {
		return nil, err
	}
	var options Options
	for _, server := range conf.Config.Config.LiteServers {
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

func convertToLiteServerOptions(server liteServerConfig) (LiteServer, error) {
	if server.ID.Type != "pub.ed25519" {
		return LiteServer{}, fmt.Errorf("not pub.ed25519 liteserver ID. Other types not supported")
	}
	if server.Ip > 0xFFFF_FFFF {
		return LiteServer{}, fmt.Errorf("only IPv4 supported")
	}
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, uint32(server.Ip))
	port, err := strconv.Atoi(server.Port)
	if err != nil {
		return LiteServer{}, err
	}
	return LiteServer{
		Host: fmt.Sprintf("%v.%v.%v.%v:%v", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3], port),
		Key:  server.ID.Key,
	}, nil
}
