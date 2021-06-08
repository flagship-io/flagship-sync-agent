package lib

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

type FlagshipConfig struct {
	EnvId              string `json:"envId"`
	PollingInterval    int    `json:"pollingInterval"`
	BucketingDirectory string `json:"bucketingDirectory"`
}

func (flagshipConfig *FlagshipConfig) New() *FlagshipConfig {

	flagshipConfigFile := flag.String("flagshipConfigFile", "", "flagshipConfigFile")
	flagshipConfigFileShort := flag.String("f", "", "flagshipConfigFile short argument")
	envId := flag.String("envId", "", "environment Id")

	pollingInterval := flag.Int("pollingInterval", -1, "pollingInterval")
	pollingIntervalShort := flag.Int("p", -1, "pollingInterval short argument")

	bucketingDirectory := flag.String("bucketingDirectory", "", "bucketingDirectory")
	bucketingDirectoryShort := flag.String("b", "", "bucketingDirectory short argument")

	flag.Parse()

	if *flagshipConfigFile != "" {
		os.Setenv("FLAGSHIP_CONFIG_FILE", *flagshipConfigFile)
	} else if *flagshipConfigFileShort != "" {
		os.Setenv("FLAGSHIP_CONFIG_FILE", *flagshipConfigFileShort)
	}

	if *envId != "" {
		os.Setenv("FLAGSHIP_ENV_ID", *envId)
	}

	if *pollingInterval > -1 {
		os.Setenv("FLAGSHIP_POLLING_INTERVAL", strconv.Itoa(*pollingInterval))
	} else if *pollingIntervalShort > -1 {
		os.Setenv("FLAGSHIP_POLLING_INTERVAL", strconv.Itoa(*pollingIntervalShort))
	}

	if *bucketingDirectory != "" {
		os.Setenv("FLAGSHIP_BUCKETING_DIRECTORY", *bucketingDirectory)
	} else if *bucketingDirectoryShort != "" {
		os.Setenv("FLAGSHIP_BUCKETING_DIRECTORY", *bucketingDirectoryShort)
	}
	return flagshipConfig
}

func (flagshipConfig *FlagshipConfig) getFlagshipConfigFile(flagshipConfigPath string) (*FlagshipConfig, error) {

	if flagshipConfigPath == "" {
		return nil, fmt.Errorf("FlagshipConfig path is empty")
	}

	file, err := os.Open(flagshipConfigPath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	flagshipConfig.PollingInterval = 2000

	err = json.Unmarshal(fileBytes, &flagshipConfig)

	if err != nil {
		return nil, err
	}

	if flagshipConfig.EnvId == "" {
		return nil, fmt.Errorf("flagshipConfig file envId field is empty")
	}

	return flagshipConfig, nil
}

func (flagshipConfig *FlagshipConfig) getConfigFromEnv() (*FlagshipConfig, error) {

	envId := os.Getenv("FLAGSHIP_ENV_ID")
	if envId == "" {
		return nil, fmt.Errorf("environment variable \"FLAGSHIP_ENV_ID\" is empty")
	}

	flagshipConfig.EnvId = envId

	bucketingDirectory := os.Getenv("FLAGSHIP_BUCKETING_DIRECTORY")

	if bucketingDirectory == "" {
		fmt.Println("environement variable \"FLAGSHIP_BUCKETING_DIRECTORY\" is empty or not set, default value will be used")
	}

	flagshipConfig.BucketingDirectory = bucketingDirectory

	envPollingInterval := os.Getenv("FLAGSHIP_POLLING_INTERVAL")

	if envPollingInterval == "" {
		flagshipConfig.PollingInterval = 2000
		fmt.Println("environement variable \"FLAGSHIP_POLLING_INTERVAL\" is empty or not set, default value will be usedm 2000ms")
	} else {
		pollingInterval, err := strconv.Atoi(envPollingInterval)
		if err != nil {
			pollingInterval = 2000
			fmt.Printf("environement variable \"FLAGSHIP_POLLING_INTERVAL\" is not an int, default value will be usedm 2000ms")
		}
		flagshipConfig.PollingInterval = pollingInterval
	}

	return flagshipConfig, nil
}

func (flagshipConfig *FlagshipConfig) GetConfig() (*FlagshipConfig, error) {

	flagshipConfigPath := os.Getenv("FLAGSHIP_CONFIG_FILE")
	if flagshipConfigPath != "" {
		flagshipConfig, err := flagshipConfig.getFlagshipConfigFile(flagshipConfigPath)
		if err != nil {
			return flagshipConfig, err
		}
		return flagshipConfig, nil
	} else {
		flagshipConfig, err := flagshipConfig.getConfigFromEnv()
		if err != nil {
			return flagshipConfig, err
		}
		return flagshipConfig, nil
	}
}
