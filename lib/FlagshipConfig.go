package lib

import (
	"flag"
	"fmt"
	"strconv"
)

const (
	FS_ENV_ID                 = "FS_ENV_ID"
	FS_POLLING_INTERVAL       = "FS_POLLING_INTERVAL"
	FS_PORT                   = "FS_PORT"
	FS_ADDRESS                = "FS_ADDRESS"
	FlagshipEnvIdErrorMessage = "Argument `--envId` is required"
	DEFAULT_PORT              = 8080
	DEFAULT_ADDRESS           = "0.0.0.0"
	DEFAULT_POLLING_INTERVAL  = 10000
)

type FlagshipConfig struct {
	EnvId           string `json:"envId"`
	PollingInterval int    `json:"pollingInterval"`
	Port            int    `json:"port"`
	Address         string `json:"address"`
}

func (flagshipConfig *FlagshipConfig) New() (*FlagshipConfig, error) {

	envId := flag.String("envId", "", "Environment Id")
	port := flag.Int("port", 0, "Http endpoint port. Default: "+strconv.Itoa(DEFAULT_PORT))
	address := flag.String("address", "", "Http endpoint address. Default: "+DEFAULT_ADDRESS)
	pollingInterval := flag.Int("pollingInterval", 0, "pollingInterval")

	flag.Parse()

	if *envId == "" {
		return nil, fmt.Errorf(FlagshipEnvIdErrorMessage)
	}
	flagshipConfig.EnvId = *envId

	if *port == 0 {
		fmt.Println("Argument `--port` is not set, default value will be used " + strconv.Itoa(DEFAULT_PORT))
		flagshipConfig.Port = DEFAULT_PORT
	} else {
		flagshipConfig.Port = *port
	}

	if *address == "" {
		fmt.Println("Argument `--address` is not set, default value will be used " + DEFAULT_ADDRESS)
		flagshipConfig.Address = DEFAULT_ADDRESS
	} else {
		flagshipConfig.Address = *address
	}

	if *pollingInterval == 0 {
		fmt.Println("Argument `--pollingInterval` is not set, default value will be used " + strconv.Itoa(DEFAULT_POLLING_INTERVAL) + "ms")
		flagshipConfig.PollingInterval = DEFAULT_POLLING_INTERVAL
	} else {
		flagshipConfig.PollingInterval = *pollingInterval
	}

	return flagshipConfig, nil
}
