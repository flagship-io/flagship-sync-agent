package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	BucketingApiUrl = "https://cdn.flagship.io/%s/bucketing.json"
)

var osCreate = os.Create

type BucketingPolling struct {
	FlagshipConfig *FlagshipConfig
	HttpClient     *http.Client
	BaseUrl        string
	lastModified   string
}

func (bucketingPolling *BucketingPolling) New(flagshipConfig *FlagshipConfig, httpClient *http.Client) *BucketingPolling {
	bucketingPolling.FlagshipConfig = flagshipConfig
	bucketingPolling.HttpClient = httpClient
	return bucketingPolling
}

func (bucketingPolling *BucketingPolling) writeBucketingFile(buffer []byte) error {

	bucketingDirectory := "flagship"
	if bucketingPolling.FlagshipConfig.BucketingPath != "" {
		bucketingDirectory = bucketingPolling.FlagshipConfig.BucketingPath
	}

	filePath := bucketingDirectory + "/bucketing.json"

	if len(buffer) == 0 {
		return fmt.Errorf("response content null")
	}

	if _, err := os.Stat(bucketingDirectory); os.IsNotExist(err) {
		err := os.Mkdir(bucketingDirectory, os.ModeDir)
		if err != nil {
			return fmt.Errorf("mkdir directory %s error", bucketingPolling.FlagshipConfig.BucketingPath)
		}
	}

	bucketingFile, err := osCreate(filePath)
	if err != nil {
		return err
	}

	defer func() {
		err := bucketingFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	_, err = bucketingFile.Write(buffer)

	if err != nil {
		return err
	}

	return nil
}

func (bucketingPolling *BucketingPolling) Polling() error {

	_, err := bucketingPolling.FlagshipConfig.GetConfig()

	if err != nil {
		return err
	}

	bucketingApiUrl := fmt.Sprintf(BucketingApiUrl, bucketingPolling.FlagshipConfig.EnvId)

	req, err := http.NewRequest("GET", bucketingApiUrl, nil)

	if err != nil {
		return err
	}

	if bucketingPolling.lastModified != "" {
		req.Header = http.Header{
			"if-modified-since": []string{bucketingPolling.lastModified},
		}
	}

	response, err := bucketingPolling.HttpClient.Do(req)

	if err != nil {
		return err
	}

	lastModified := response.Header.Get("last-modified")
	if lastModified != "" {
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
		err = bucketingPolling.writeBucketingFile(body)
	}

	if err != nil {
		return err
	}
	return nil
}
