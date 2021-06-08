package main

import (
	"flagship-io/flagship-sync-agent/lib"
	"fmt"
	"time"
)

func main() {

	var flagshipConfig lib.FlagshipConfig
	var BucketingPolling lib.BuctingPolling

	flagshipConfig.New()

	BucketingPolling.New(&flagshipConfig)

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
