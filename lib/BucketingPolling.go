package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type BuctingPolling struct {
	FlagshipConfig *FlagshipConfig
}

func (buctingPolling *BuctingPolling) New(flagshipConfig *FlagshipConfig) *BuctingPolling {
	buctingPolling.FlagshipConfig = flagshipConfig
	return buctingPolling
}

func (buctingPolling *BuctingPolling) writeBucktingFile(buffer []byte) error {

	bucketingDirectory := "flagship"
	if buctingPolling.FlagshipConfig.BucketingDirectory != "" {
		bucketingDirectory = buctingPolling.FlagshipConfig.BucketingDirectory
	}

	filePath := bucketingDirectory + "/Bucketing.json"

	if len(buffer) == 0 {
		return fmt.Errorf("response content null")
	}

	if _, err := os.Stat(bucketingDirectory); os.IsNotExist(err) {
		err := os.Mkdir(bucketingDirectory, os.ModeDir)
		if err != nil {
			return fmt.Errorf("mkdir directory %s error", buctingPolling.FlagshipConfig.BucketingDirectory)
		}
	}

	bucketingFile, err := os.Create(filePath)
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

func (buctingPolling *BuctingPolling) Polling() error {

	_, err := buctingPolling.FlagshipConfig.GetConfig()

	if err != nil {
		return err
	}

	bucketingApiUrl := "https://cdn.flagship.io/%s/bucketing.json"
	bucketingApiUrl = fmt.Sprintf(bucketingApiUrl, buctingPolling.FlagshipConfig.EnvId)
	response, err := http.Get(bucketingApiUrl)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 && response.StatusCode != 304 {
		return fmt.Errorf("%s", response.Body)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	err = buctingPolling.writeBucktingFile(body)

	if err != nil {
		return err
	}
	return nil
}
