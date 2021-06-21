package lib

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func resetForEachTest(oldArgs []string) {
	os.Args = oldArgs
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	_ = os.Unsetenv(FlagshipBucketingDirectory)
	_ = os.Unsetenv(FlagshipConfigFile)
	_ = os.Unsetenv(FlagshipEnvId)
	_ = os.Unsetenv(FlagshipPollingInterval)
}

func messageError(t *testing.T, name string, expected interface{}, result interface{}) {
	t.Errorf(name+" failed, expected %v but got %v", expected, result)
}
func TestNewConfigFile(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//Test FLAGSHIP_CONFIG_FILE
	configFile := "flagship.json"
	os.Args = []string{"cmd", "-config=" + configFile}
	flagshipConfig = FlagshipConfig{}
	flagshipConfig.New()
	if os.Getenv(FlagshipConfigFile) != configFile {
		t.Fail()
	}
}

func TestNewConfigFileShort(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//Test FLAGSHIP_CONFIG_FILE
	configFile := "flagship.json"
	os.Args = []string{"cmd", "-c=" + configFile}
	flagshipConfig = FlagshipConfig{}
	flagshipConfig.New()
	if os.Getenv(FlagshipConfigFile) != configFile {
		t.Fail()
	}
}

func TestNewConfigFileShortAndLong(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//Test FLAGSHIP_CONFIG_FILE
	configFileLong := "flagship.json"
	configFileShort := "shortConfig.json"
	os.Args = []string{"cmd", "-c=" + configFileShort, "-config=" + configFileLong}
	flagshipConfig = FlagshipConfig{}
	flagshipConfig.New()
	if os.Getenv(FlagshipConfigFile) != configFileLong {
		t.Fail()
	}

}

func TestNewEnvId(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//test FLAGSHIP_ENV_ID
	envId := "envId_id"
	os.Args = []string{"cmd", "-envId=" + envId}
	flagshipConfig.New()
	if os.Getenv(FlagshipEnvId) != envId {
		t.Fail()
	}
}

func TestNewPollingLong(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//test FLAGSHIP_ENV_ID
	pollingInterval := 5200
	os.Args = []string{"cmd", "-pollingInterval=" + strconv.Itoa(pollingInterval)}
	flagshipConfig.New()
	if os.Getenv(FlagshipPollingInterval) != strconv.Itoa(pollingInterval) {
		t.Fail()
	}
}

func TestNewPollingShort(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//test FLAGSHIP_ENV_ID
	pollingInterval := 5200
	os.Args = []string{"cmd", "-p=" + strconv.Itoa(pollingInterval)}
	flagshipConfig.New()
	if os.Getenv(FlagshipPollingInterval) != strconv.Itoa(pollingInterval) {
		t.Fail()
	}
}

func TestNewPollingShortAndLong(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//Test FLAGSHIP_POLLING_INTERVAL
	pollingIntervalLong := 5000
	pollingIntervalShort := 2500
	os.Args = []string{"cmd", "-p=" + strconv.Itoa(pollingIntervalShort), "-pollingInterval=" + strconv.Itoa(pollingIntervalLong)}
	flagshipConfig = FlagshipConfig{}
	flagshipConfig.New()
	if os.Getenv(FlagshipPollingInterval) != strconv.Itoa(pollingIntervalLong) {
		t.Fail()
	}
}

func TestNewBucketingDirectoryLong(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//Test FLAGSHIP_BUCKETING_DIRECTORY
	bucketingDirectory := "flagshipDirectory"
	os.Args = []string{"cmd", "-bucketingPath=" + bucketingDirectory}
	flagshipConfig = FlagshipConfig{}
	flagshipConfig.New()
	if os.Getenv(FlagshipBucketingDirectory) != bucketingDirectory {
		t.Fail()
	}
}

func TestNewBucketingDirectoryShort(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//Test FLAGSHIP_BUCKETING_DIRECTORY
	bucketingDirectory := "flagshipDirectory"
	os.Args = []string{"cmd", "-b=" + bucketingDirectory}
	flagshipConfig = FlagshipConfig{}
	flagshipConfig.New()
	if os.Getenv(FlagshipBucketingDirectory) != bucketingDirectory {
		t.Fail()
	}
}

func TestNewBucketingDirectoryShortAndLong(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//Test FLAGSHIP_BUCKETING_DIRECTORY
	bucketingDirectoryLong := "flagshipDirectoryLong"
	bucketingDirectoryShort := "flagshipDirectoryShort"
	os.Args = []string{"cmd", "-b=" + bucketingDirectoryShort, "-bucketingPath=" + bucketingDirectoryLong}
	flagshipConfig = FlagshipConfig{}
	flagshipConfig.New()
	if os.Getenv(FlagshipBucketingDirectory) != bucketingDirectoryLong {
		t.Fail()
	}
}

func TestGetConfigFromEnv(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)
	bucketingDirectory := "flagshipDirectory"
	pollingInterval := 5200
	envId := "env_id"

	_ = os.Setenv(FlagshipEnvId, envId)
	_ = os.Setenv(FlagshipPollingInterval, strconv.Itoa(pollingInterval))
	_ = os.Setenv(FlagshipBucketingDirectory, bucketingDirectory)

	_, err := flagshipConfig.GetConfig()

	if err != nil {
		t.Fail()
	}

	if flagshipConfig.EnvId != envId {
		messageError(t, "flagshipConfig.EnvId", envId, flagshipConfig.EnvId)
	}

	if flagshipConfig.PollingInterval != pollingInterval {
		messageError(t, "flagshipConfig.PollingInterval", pollingInterval, flagshipConfig.PollingInterval)
	}

	if flagshipConfig.BucketingDirectory != bucketingDirectory {
		messageError(t, "flagshipConfig.BucketingDirectory", bucketingDirectory, flagshipConfig.BucketingDirectory)
	}
}

//Test without envId
func TestGetConfigFromEnv1(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	_, err := flagshipConfig.GetConfig()

	if err == nil {
		t.FailNow()
	}

	if err.Error() != FlagshipEnvIdErrorMessage {
		t.Fail()
	}
}

func TestGetFlagshipConfigFile(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)
	filePath := "flagship.json"
	_ = os.Setenv(FlagshipConfigFile, filePath)
	envId := "env_id"
	pollingInterval := "5300"
	bucketingDirectory := "flagship"

	file, _ := os.Create(filePath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
		_ = os.Remove(filePath)
	}(file)

	_, err := file.WriteString("{\n  \"envId\" : \"" + envId + "\",\n  \"pollingInterval\": " + pollingInterval + ",\n  \"bucketingDirectory\": \"" + bucketingDirectory + "\"\n}")

	if err != nil {
		t.Error(err)
	}
	_, err = flagshipConfig.GetConfig()
	if err != nil {
		t.Error(err)
	}

	if flagshipConfig.EnvId != envId {
		messageError(t, "flagshipConfig.EnvId", envId, flagshipConfig.EnvId)
	}
	if intPolling, _ := strconv.Atoi(pollingInterval); flagshipConfig.PollingInterval != intPolling {
		messageError(t, "flagshipConfig.EnvId", pollingInterval, flagshipConfig.PollingInterval)
	}
	if flagshipConfig.EnvId != envId {
		messageError(t, "flagshipConfig.EnvId", envId, flagshipConfig.EnvId)
	}
}

func TestGetFlagshipConfigFileError(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)
	filePath := "flagship.json"
	_ = os.Setenv(FlagshipConfigFile, filePath)

	_, err := flagshipConfig.GetConfig()
	if err == nil {
		t.FailNow()
	}
}

func TestGetFlagshipConfigFileErrorEnvId(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)
	filePath := "flagship.json"
	_ = os.Setenv(FlagshipConfigFile, filePath)
	pollingInterval := "5300"
	bucketingDirectory := "flagship"

	file, _ := os.Create(filePath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
		_ = os.Remove(filePath)
	}(file)

	_, err := file.WriteString("{\n  \"pollingInterval\": " + pollingInterval + ",\n  \"bucketingDirectory\": \"" + bucketingDirectory + "\"\n}")

	if err != nil {
		t.Error(err)
	}
	_, err = flagshipConfig.GetConfig()
	if err == nil {
		t.FailNow()
	}

	if err.Error() != FlagshipConfigEnvIdError {
		t.Error(err)
	}
}
