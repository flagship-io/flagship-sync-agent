package lib

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BucketingApiUrl   = "https://cdn.flagship.io/%s/bucketing.json"
	IF_MODIFIED_SINCE = "if-modified-since"
	LAST_MODIFIED     = "last-modified"
)

var BucktingFile []byte

type BucketingPolling struct {
	FlagshipConfig *FlagshipConfig
	HttpClient     *http.Client
	BaseUrl        string
	lastModified   []string
}

func (bucketingPolling *BucketingPolling) New(flagshipConfig *FlagshipConfig, httpClient *http.Client) *BucketingPolling {
	bucketingPolling.FlagshipConfig = flagshipConfig
	bucketingPolling.HttpClient = httpClient
	return bucketingPolling
}

func (bucketingPolling *BucketingPolling) Polling() error {

	bucketingApiUrl := fmt.Sprintf(BucketingApiUrl, bucketingPolling.FlagshipConfig.EnvId)

	req, err := http.NewRequest("GET", bucketingApiUrl, nil)

	if err != nil {
		return err
	}

	if len(bucketingPolling.lastModified) > 0 {
		req.Header = http.Header{
			IF_MODIFIED_SINCE: bucketingPolling.lastModified,
		}
	}

	response, err := bucketingPolling.HttpClient.Do(req)

	if err != nil {
		return err
	}

	lastModified := response.Header[LAST_MODIFIED]
	if len(lastModified) > 0 {
		bucketingPolling.lastModified = lastModified
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	if response.StatusCode != 200 && response.StatusCode != 304 {
		return fmt.Errorf("%s", response.Body)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	if len(body) > 0 {
		BucktingFile = body
	}

	return nil
}

func (bucketingPolling *BucketingPolling) StartPolling() {
	for {
		fmt.Println("Polling start")

		err := bucketingPolling.Polling()

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Polling success")
		}

		fmt.Println("Polling end")

		if bucketingPolling.FlagshipConfig.PollingInterval == 0 {
			break
		}
		duration := time.Duration(bucketingPolling.FlagshipConfig.PollingInterval) * time.Millisecond
		time.Sleep(duration)
	}
}
