package lib

import (
	"flag"
	"fmt"
)

const (
	FS_ENV_ID                 = "FS_ENV_ID"
	FS_POLLING_INTERVAL       = "FS_POLLING_INTERVAL"
	FS_PORT                   = "FS_PORT"
	FS_ADDRESS                = "FS_ADDRESS"
	FlagshipEnvIdErrorMessage = "argument envId is required"
	DEFAULT_PORT              = 8080
	DEFAULT_ADDRESS           = "0.0.0.0"
	DEFAULT_POLLING_INTERVAL  = 2000
)

type FlagshipConfig struct {
	EnvId           string `json:"envId"`
	PollingInterval int    `json:"pollingInterval"`
	Port            int    `json:"port"`
	Address         string `json:"address"`
}

func (flagshipConfig *FlagshipConfig) New() (*FlagshipConfig, error) {

	envId := flag.String("envId", "", "environment Id")
	port := flag.Int("port", DEFAULT_PORT, "Http endpoint port. Default: 8080")
	address := flag.String("address", DEFAULT_ADDRESS, "Http endpoint address. Default: 0.0.0.0")
	pollingInterval := flag.Int("pollingInterval", DEFAULT_POLLING_INTERVAL, "pollingInterval")

	flag.Parse()

	if *envId == "" {
		return nil, fmt.Errorf(FlagshipEnvIdErrorMessage)
	}
	flagshipConfig.EnvId = *envId

	if *port == DEFAULT_PORT {
		fmt.Println("argument port is empty or not set, default value will be used 8080")
	}
	flagshipConfig.Port = *port

	if *address == DEFAULT_ADDRESS {
		fmt.Println("argument address is empty or not set, default value will be used 0.0.0.0")
	}
	flagshipConfig.Address = *address

	if *pollingInterval == DEFAULT_POLLING_INTERVAL {
		fmt.Println("argument pollingInterval is empty or not set, default value will be used 2000ms")
	}
	flagshipConfig.PollingInterval = *pollingInterval

	return flagshipConfig, nil
}
