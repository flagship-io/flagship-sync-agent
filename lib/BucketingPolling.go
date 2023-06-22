package lib

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	BucketingApiUrl    = "https://cdn.flagship.io/%s/bucketing.json"
	IF_MODIFIED_SINCE  = "if-modified-since"
	LAST_MODIFIED      = "last-modified"
	HTTP_ERROR_MESSAGE = "Error status code: %d with message: %s"
)

type BucketingPolling struct {
	FlagshipConfig *FlagshipConfig
	HttpClient     *http.Client
	BaseUrl        string
	lastModified   []string
	BucketingFile  []byte
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

	lastModified := response.Header[http.CanonicalHeaderKey(LAST_MODIFIED)]
	if len(lastModified) > 0 {
		bucketingPolling.lastModified = lastModified
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 && response.StatusCode != 304 {
		return fmt.Errorf(HTTP_ERROR_MESSAGE, response.StatusCode, body)
	}

	if len(body) > 0 {
		bucketingPolling.BucketingFile = body
		fmt.Printf("Polling event with code status 200 : %s", body)
	} else {
		fmt.Println("Polling event with code status " + strconv.Itoa(response.StatusCode))
	}

	return nil
}

func (bucketingPolling *BucketingPolling) StartPolling() {
	fmt.Println("Polling start")
	for {
		err := bucketingPolling.Polling()
		if err != nil {
			fmt.Println(err.Error())
		}
		if bucketingPolling.FlagshipConfig.PollingInterval == 0 {
			break
		}
		duration := time.Duration(bucketingPolling.FlagshipConfig.PollingInterval) * time.Millisecond
		time.Sleep(duration)
	}
}
