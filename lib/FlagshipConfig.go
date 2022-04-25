package lib

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	FS_ENV_ID                 = "FS_ENV_ID"
	FS_POLLING_INTERVAL       = "FS_POLLING_INTERVAL"
	FS_BUCKETING_DIRECTORY    = "FS_BUCKETING_DIRECTORY"
	FS_PORT                   = "FS_PORT"
	FS_ADDRESS                = "FS_ADDRESS"
	FlagshipEnvIdErrorMessage = "argument envId is required"
	BucketingDirectoryError   = "environment variable bucketingPath is empty or not set, default value will be used"
	DEFAULT_PORT              = 8080
	DEFAULT_ADDRESS           = "0.0.0.0"
)

type FlagshipConfig struct {
	EnvId           string `json:"envId"`
	PollingInterval int    `json:"pollingInterval"`
	BucketingPath   string `json:"bucketingPath"`
	Port            int    `json:"port"`
	Address         string `json:"address"`
}

func (flagshipConfig *FlagshipConfig) New() (*FlagshipConfig, error) {

	envId := flag.String("envId", "", "environment Id")

	port := flag.Int("port", DEFAULT_PORT, "Http endpoint port. Default: 8080")
	address := flag.String("address", "0.0.0.0", "Http endpoint address. Default: 0.0.0.0")

	pollingInterval := flag.Int("pollingInterval", -1, "pollingInterval")

	bucketingDirectory := flag.String("bucketingPath", "", "bucketing Directory path")

	flag.Parse()

	if *envId != "" {
		_ = os.Setenv(FS_ENV_ID, *envId)
	}

	if *port != DEFAULT_PORT {
		_ = os.Setenv(FS_PORT, strconv.Itoa(*port))
	}

	if *address != DEFAULT_ADDRESS {
		_ = os.Setenv(FS_ADDRESS, *address)
	}

	if *pollingInterval > -1 {
		_ = os.Setenv(FS_POLLING_INTERVAL, strconv.Itoa(*pollingInterval))
	}

	if *bucketingDirectory != "" {
		_ = os.Setenv(FS_BUCKETING_DIRECTORY, *bucketingDirectory)
	}

	return flagshipConfig.GetConfig()
}

func (flagshipConfig *FlagshipConfig) GetConfig() (*FlagshipConfig, error) {

	envId := os.Getenv(FS_ENV_ID)
	if envId == "" {
		return nil, fmt.Errorf(FlagshipEnvIdErrorMessage)
	}

	flagshipConfig.EnvId = envId

	bucketingDirectory := os.Getenv(FS_BUCKETING_DIRECTORY)

	if bucketingDirectory == "" {
		fmt.Println(BucketingDirectoryError)
	}

	flagshipConfig.BucketingPath = bucketingDirectory

	envPollingInterval := os.Getenv(FS_POLLING_INTERVAL)

	if envPollingInterval == "" {
		flagshipConfig.PollingInterval = 2000
		fmt.Println("argument pollingInterval is empty or not set, default value will be used 2000ms")
	} else {
		pollingInterval, err := strconv.Atoi(envPollingInterval)
		if err != nil {
			pollingInterval = 2000
			fmt.Printf("argument pollingInterval is not an int, default value will be used 2000ms")
		}
		flagshipConfig.PollingInterval = pollingInterval
	}

	envPort := os.Getenv(FS_PORT)
	if envPort == "" {
		flagshipConfig.Port = DEFAULT_PORT
		fmt.Println("argument port is empty or not set, default value will be used 8080")
	} else {
		port, err := strconv.Atoi(envPort)
		if err != nil {
			port = DEFAULT_PORT
			fmt.Printf("argument port is not an int, default value will be used 8080")
		}
		flagshipConfig.Port = port
	}

	address := os.Getenv(FS_ADDRESS)
	if address == "" {
		address = DEFAULT_ADDRESS
		fmt.Println("argument address is empty or not set, default value will be used 0.0.0.0")
	}
	flagshipConfig.Address = address

	return flagshipConfig, nil
}
