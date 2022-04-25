package lib

import (
	"flag"
	"os"
	"strconv"
	"testing"
)

func resetForEachTest(oldArgs []string) {
	os.Args = oldArgs
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	_ = os.Unsetenv(FS_BUCKETING_DIRECTORY)
	_ = os.Unsetenv(FS_ENV_ID)
	_ = os.Unsetenv(FS_POLLING_INTERVAL)
}

func messageError(t *testing.T, name string, expected interface{}, result interface{}) {
	t.Errorf(name+" failed, expected %v but got %v", expected, result)
}

func TestNewEnvId(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//test FLAGSHIP_ENV_ID
	envId := "envId_id"
	os.Args = []string{"cmd", "-envId=" + envId}
	flagshipConfig.New()
	if os.Getenv(FS_ENV_ID) != envId {
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
	if os.Getenv(FS_POLLING_INTERVAL) != strconv.Itoa(pollingInterval) {
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
	if os.Getenv(FS_BUCKETING_DIRECTORY) != bucketingDirectory {
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

	_ = os.Setenv(FS_ENV_ID, envId)
	_ = os.Setenv(FS_POLLING_INTERVAL, strconv.Itoa(pollingInterval))
	_ = os.Setenv(FS_BUCKETING_DIRECTORY, bucketingDirectory)

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

	if flagshipConfig.BucketingPath != bucketingDirectory {
		messageError(t, "flagshipConfig.BucketingPath", bucketingDirectory, flagshipConfig.BucketingPath)
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
