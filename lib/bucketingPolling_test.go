package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestBucketingPolling_New(t *testing.T) {
	var bucketingPolling BucketingPolling
	var flagshipConfig FlagshipConfig
	var HttpClient http.Client

	bucketingPollingOut := bucketingPolling.New(&flagshipConfig, &HttpClient)

	if &bucketingPolling != bucketingPollingOut {
		t.FailNow()
	}

	if &flagshipConfig != bucketingPolling.FlagshipConfig {
		t.FailNow()
	}

	if bucketingPolling.HttpClient != &HttpClient {
		t.FailNow()
	}
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestBucketingPolling_Polling(t *testing.T) {

	oldArgs := os.Args
	defer resetForEachTest(oldArgs)

	envId := "env_id"
	pollingInterval := 0
	_ = os.Setenv(FS_ENV_ID, envId)
	_ = os.Setenv(FS_POLLING_INTERVAL, strconv.Itoa(pollingInterval))

	osCreate = func(name string) (*os.File, error) {
		//Test default bucketing File
		if name != "flagship/bucketing.json" {
			t.FailNow()
		}
		f, _ := os.CreateTemp("", "example")
		defer func(name string) {
			_ = os.Remove(name)
		}(f.Name())
		return f, nil
	}

	defer func() {
		osCreate = os.Create
	}()

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		bucketingApiUrl := fmt.Sprintf(BucketingApiUrl, envId)
		fmt.Println(req.URL.String())
		if req.URL.String() != bucketingApiUrl {
			messageError(t, "bucketingApiUrl", bucketingApiUrl, req.URL.String())
		}
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	var flagshipConfig FlagshipConfig

	var bucketingPolling BucketingPolling

	flagshipConfig.New()

	bucketingPolling.New(&flagshipConfig, client)
	err := bucketingPolling.Polling()
	if err != nil {
		t.Error(err)
	}
}

func TestBucketingPolling_Polling2(t *testing.T) {
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)
	envId := "env_id"
	pollingInterval := 0
	BucketingDirectory := "myDirector"
	_ = os.Setenv(FS_ENV_ID, envId)
	_ = os.Setenv(FS_POLLING_INTERVAL, strconv.Itoa(pollingInterval))

	_ = os.Setenv(FS_BUCKETING_DIRECTORY, BucketingDirectory)
	osCreate = func(name string) (*os.File, error) {
		//Test default bucketing File
		if name != BucketingDirectory+"/bucketing.json" {
			t.FailNow()
		}
		f, _ := os.CreateTemp("", "example")
		defer func(name string) {
			_ = os.Remove(name)
		}(f.Name())
		return f, nil
	}

	defer func() {
		osCreate = os.Create
	}()

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		bucketingApiUrl := fmt.Sprintf(BucketingApiUrl, envId)
		fmt.Println(req.URL.String())
		if req.URL.String() != bucketingApiUrl {
			messageError(t, "bucketingApiUrl", bucketingApiUrl, req.URL.String())
		}
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	var flagshipConfig FlagshipConfig

	var bucketingPolling BucketingPolling

	bucketingPolling.New(&flagshipConfig, client)
	err := bucketingPolling.Polling()
	if err != nil {
		t.Error(err)
	}
}

func TestBucketingPolling_Error(t *testing.T) {

	oldArgs := os.Args
	defer resetForEachTest(oldArgs)
	var flagshipConfig FlagshipConfig

	var bucketingPolling BucketingPolling
	var client http.Client
	bucketingPolling.New(&flagshipConfig, &client)
	err := bucketingPolling.Polling()
	if err == nil {
		t.Error(err)
	}
}

func TestBucketingPolling_Polling_http_400(t *testing.T) {
	oldArgs := os.Args
	defer resetForEachTest(oldArgs)
	envId := "env_id"
	pollingInterval := 0
	BucketingDirectory := "myDirector"
	_ = os.Setenv(FS_ENV_ID, envId)
	_ = os.Setenv(FS_POLLING_INTERVAL, strconv.Itoa(pollingInterval))
	_ = os.Setenv(FS_BUCKETING_DIRECTORY, BucketingDirectory)

	httpBody := `{"message": "Forbidden"}`
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		bucketingApiUrl := fmt.Sprintf(BucketingApiUrl, envId)
		fmt.Println(req.URL.String())
		if req.URL.String() != bucketingApiUrl {
			messageError(t, "bucketingApiUrl", bucketingApiUrl, req.URL.String())
		}
		return &http.Response{
			StatusCode: 403,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(httpBody)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	var flagshipConfig FlagshipConfig

	var bucketingPolling BucketingPolling

	bucketingPolling.New(&flagshipConfig, client)
	err := bucketingPolling.Polling()
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "{"+httpBody+"}" {
		t.Error(err)
	}
}
