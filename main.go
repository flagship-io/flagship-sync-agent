package main

import (
	"flagship-io/flagship-sync-agent/lib"
	"fmt"
	"net/http"
	"time"
)

func main() {

	var flagshipConfig lib.FlagshipConfig
	var BucketingPolling lib.BucketingPolling
	var HttpClient http.Client

	flagshipConfig.New()

	BucketingPolling.New(&flagshipConfig, &HttpClient)

	for {
		fmt.Println("Polling start")

		err := BucketingPolling.Polling()

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Polling success")
		}

		fmt.Println("Polling end")

		if flagshipConfig.PollingInterval == 0 {
			break
		}
		duration := time.Duration(flagshipConfig.PollingInterval) * time.Millisecond
		time.Sleep(duration)
	}
}
