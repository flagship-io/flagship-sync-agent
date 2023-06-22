package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
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

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestBucketingPolling_Polling(t *testing.T) {

	envId := "env_id"
	pollingInterval := 0

	httpBodyString := `{"campaigns":[]}`
	lastModified := []string{"2022-04-26 11:54:00"}

	count := 1

	client := NewTestClient(func(req *http.Request) (*http.Response, error) {
		// Test request parameters
		bucketingApiUrl := fmt.Sprintf(BucketingApiUrl, envId)
		fmt.Println(req.URL.String())
		if req.URL.String() != bucketingApiUrl {
			messageError(t, "bucketingApiUrl", bucketingApiUrl, req.URL.String())
		}

		if count > 1 && reflect.DeepEqual(req.Header[http.CanonicalHeaderKey(LAST_MODIFIED)], lastModified) {
			messageError(t, "LAST_MODIFIED", lastModified, req.Header[http.CanonicalHeaderKey(LAST_MODIFIED)])
		}

		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(httpBodyString)),
			// Must be set to non-nil value or it panics
			Header: http.Header{
				http.CanonicalHeaderKey(LAST_MODIFIED): lastModified,
			},
		}, nil
	})

	var flagshipConfig FlagshipConfig

	flagshipConfig.EnvId = envId

	flagshipConfig.PollingInterval = pollingInterval

	var bucketingPolling BucketingPolling

	bucketingPolling.New(&flagshipConfig, client)
	err := bucketingPolling.Polling()
	if err != nil {
		t.Error(err)
	}

	bucketingFileString := string(bucketingPolling.BucketingFile)

	if bucketingFileString != httpBodyString {
		messageError(t, "BucktingFile", httpBodyString, bucketingFileString)
	}
	if !reflect.DeepEqual(lastModified, bucketingPolling.lastModified) {
		messageError(t, "BucktingFile", lastModified, bucketingPolling.lastModified)
	}

	count++

	err = bucketingPolling.Polling()
	if err != nil {
		t.Error(err)
	}

	bucketingFileString = string(bucketingPolling.BucketingFile)

	if bucketingFileString != httpBodyString {
		messageError(t, "BucktingFile", httpBodyString, bucketingFileString)
	}
	if !reflect.DeepEqual(lastModified, bucketingPolling.lastModified) {
		messageError(t, "BucktingFile", lastModified, bucketingPolling.lastModified)
	}
}

func TestBucketingPolling_Polling_Error(t *testing.T) {

	envId := "env_id"
	pollingInterval := 0

	netError := fmt.Errorf("net error")

	client := NewTestClient(func(req *http.Request) (*http.Response, error) {

		return nil, netError
	})

	var flagshipConfig FlagshipConfig

	flagshipConfig.EnvId = envId
	flagshipConfig.PollingInterval = pollingInterval

	var bucketingPolling BucketingPolling

	bucketingPolling.New(&flagshipConfig, client)
	err := bucketingPolling.Polling()
	if err == nil {
		t.FailNow()
	}
}

func TestBucketingPolling_Polling_http_400(t *testing.T) {
	envId := "env_id"
	pollingInterval := 0

	httpBody := `{"message": "Forbidden"}`
	client := NewTestClient(func(req *http.Request) (*http.Response, error) {
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
		}, nil
	})

	var flagshipConfig FlagshipConfig

	flagshipConfig.EnvId = envId
	flagshipConfig.PollingInterval = pollingInterval

	var bucketingPolling BucketingPolling

	bucketingPolling.New(&flagshipConfig, client)

	err := bucketingPolling.Polling()
	if err == nil {
		t.FailNow()
	}
}

func TestBucketingPolling_StartPolling(t *testing.T) {

	envId := "env_id"
	pollingInterval := 0

	httpBodyString := `{"campaigns":[]}`
	lastModified := []string{"2022-04-26 11:54:00"}

	count := 1

	client := NewTestClient(func(req *http.Request) (*http.Response, error) {
		// Test request parameters
		bucketingApiUrl := fmt.Sprintf(BucketingApiUrl, envId)
		fmt.Println(req.URL.String())
		if req.URL.String() != bucketingApiUrl {
			messageError(t, "bucketingApiUrl", bucketingApiUrl, req.URL.String())
		}

		if count > 1 && reflect.DeepEqual(req.Header[http.CanonicalHeaderKey(LAST_MODIFIED)], lastModified) {
			messageError(t, "LAST_MODIFIED", lastModified, req.Header[http.CanonicalHeaderKey(LAST_MODIFIED)])
		}

		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(httpBodyString)),
			// Must be set to non-nil value or it panics
			Header: http.Header{
				http.CanonicalHeaderKey(LAST_MODIFIED): lastModified,
			},
		}, nil
	})

	var flagshipConfig FlagshipConfig

	flagshipConfig.EnvId = envId

	flagshipConfig.PollingInterval = pollingInterval

	var bucketingPolling BucketingPolling

	bucketingPolling.New(&flagshipConfig, client)
	bucketingPolling.StartPolling()

	bucketingFileString := string(bucketingPolling.BucketingFile)

	if bucketingFileString != httpBodyString {
		messageError(t, "BucktingFile", httpBodyString, bucketingFileString)
	}
	if !reflect.DeepEqual(lastModified, bucketingPolling.lastModified) {
		messageError(t, "BucktingFile", lastModified, bucketingPolling.lastModified)
	}
}
