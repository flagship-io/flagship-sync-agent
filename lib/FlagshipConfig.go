package lib

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	FlagshipConfigFile         = "FLAGSHIP_CONFIG_FILE"
	FlagshipEnvId              = "FLAGSHIP_ENV_ID"
	FlagshipPollingInterval    = "FLAGSHIP_POLLING_INTERVAL"
	FlagshipBucketingDirectory = "FLAGSHIP_BUCKETING_DIRECTORY"
	FlagshipEnvIdErrorMessage  = "argument envId is required"
	BucketingDirectoryError    = "environment variable bucketingPath is empty or not set, default value will be used"
	FlagshipConfigEnvIdError   = "flagshipConfig file envId field is required"
)

type FlagshipConfig struct {
	EnvId           string `json:"envId"`
	PollingInterval int    `json:"pollingInterval"`
	BucketingPath   string `json:"bucketingPath"`
}

func (flagshipConfig *FlagshipConfig) New() *FlagshipConfig {

	flagshipConfigFile := flag.String("config", "", "flagship ConfigFile")
	flagshipConfigFileShort := flag.String("c", "", "flagshipConfigFile short argument")
	envId := flag.String("envId", "", "environment Id")

	pollingInterval := flag.Int("pollingInterval", -1, "pollingInterval")
	pollingIntervalShort := flag.Int("p", -1, "pollingInterval short argument")

	bucketingDirectory := flag.String("bucketingPath", "", "bucketing Directory path")
	bucketingDirectoryShort := flag.String("b", "", "bucketing Directory path short argument")

	flag.Parse()

	if *flagshipConfigFile != "" {
		_ = os.Setenv(FlagshipConfigFile, *flagshipConfigFile)
	} else if *flagshipConfigFileShort != "" {
		_ = os.Setenv(FlagshipConfigFile, *flagshipConfigFileShort)
	}

	if *envId != "" {
		_ = os.Setenv(FlagshipEnvId, *envId)
	}

	if *pollingInterval > -1 {
		_ = os.Setenv(FlagshipPollingInterval, strconv.Itoa(*pollingInterval))
	} else if *pollingIntervalShort > -1 {
		_ = os.Setenv(FlagshipPollingInterval, strconv.Itoa(*pollingIntervalShort))
	}

	if *bucketingDirectory != "" {
		_ = os.Setenv(FlagshipBucketingDirectory, *bucketingDirectory)
	} else if *bucketingDirectoryShort != "" {
		_ = os.Setenv(FlagshipBucketingDirectory, *bucketingDirectoryShort)
	}
	return flagshipConfig
}

func (flagshipConfig *FlagshipConfig) getFlagshipConfigFile(flagshipConfigPath string) (*FlagshipConfig, error) {

	file, err := os.Open(flagshipConfigPath)

	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

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
		return nil, fmt.Errorf(FlagshipConfigEnvIdError)
	}

	return flagshipConfig, nil
}

func (flagshipConfig *FlagshipConfig) getConfigFromEnv() (*FlagshipConfig, error) {

	envId := os.Getenv(FlagshipEnvId)
	if envId == "" {
		return nil, fmt.Errorf(FlagshipEnvIdErrorMessage)
	}

	flagshipConfig.EnvId = envId

	bucketingDirectory := os.Getenv(FlagshipBucketingDirectory)

	if bucketingDirectory == "" {
		fmt.Println(BucketingDirectoryError)
	}

	flagshipConfig.BucketingPath = bucketingDirectory

	envPollingInterval := os.Getenv(FlagshipPollingInterval)

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

	return flagshipConfig, nil
}

func (flagshipConfig *FlagshipConfig) GetConfig() (*FlagshipConfig, error) {

	flagshipConfigPath := os.Getenv(FlagshipConfigFile)
	if flagshipConfigPath != "" {
		flagshipConfig, err := flagshipConfig.getFlagshipConfigFile(flagshipConfigPath)
		if err != nil {
			return flagshipConfig, err
		}
		return flagshipConfig, nil
	}
	flagshipConfig, err := flagshipConfig.getConfigFromEnv()
	if err != nil {
		return flagshipConfig, err
	}
	return flagshipConfig, nil

}
