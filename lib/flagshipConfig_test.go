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

	_ = os.Unsetenv(FS_ENV_ID)
	_ = os.Unsetenv(FS_POLLING_INTERVAL)
	_ = os.Unsetenv(FS_PORT)
	_ = os.Unsetenv(FS_ADDRESS)
}

func messageError(t *testing.T, name string, expected interface{}, result interface{}) {
	t.Errorf(name+" failed, expected %v but got %v", expected, result)
}

func TestEnvID(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)
	_, err := flagshipConfig.New()
	if err == nil {
		t.FailNow()
	}
}

func TestNewDefault(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//test FLAGSHIP_ENV_ID
	envId := "envId_id"
	os.Args = []string{"cmd", "-envId=" + envId}
	_, _ = flagshipConfig.New()

	if flagshipConfig.EnvId != envId {
		t.Fail()
	}

	if flagshipConfig.PollingInterval != DEFAULT_POLLING_INTERVAL {
		t.Fail()
	}

	if flagshipConfig.Port != DEFAULT_PORT {
		t.Fail()
	}

	if flagshipConfig.Address != DEFAULT_ADDRESS {
		t.Fail()
	}
}
func TestNew(t *testing.T) {
	var flagshipConfig FlagshipConfig
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	//test FLAGSHIP_ENV_ID
	envId := "envId_id"
	pollingInterval := 5200
	port := 3001
	address := "127.0.0.1"
	os.Args = []string{"cmd", "-envId=" + envId, "-pollingInterval=" + strconv.Itoa(pollingInterval), "-port=" + strconv.Itoa(port), "-address=" + address}
	_, _ = flagshipConfig.New()

	if flagshipConfig.EnvId != envId {
		t.Fail()
	}

	if flagshipConfig.PollingInterval != pollingInterval {
		t.Fail()
	}

	if flagshipConfig.Port != port {
		t.Fail()
	}

	if flagshipConfig.Address != address {
		t.Fail()
	}
}
