package main

import (
	"flagship-io/flagship-sync-agent/controller"
	"flagship-io/flagship-sync-agent/lib"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	var flagshipConfig lib.FlagshipConfig
	var BucketingPolling lib.BucketingPolling
	var HttpClient http.Client

	_, err := flagshipConfig.New()

	if err != nil {
		fmt.Println(err)
		return
	}

	BucketingPolling.New(&flagshipConfig, &HttpClient)

	go BucketingPolling.StartPolling()

	bucketingController := controller.BucketingController{}

	gin.SetMode(gin.DebugMode)
	server := gin.Default()

	server.GET("/bucketing", bucketingController.GetBucketing)

	server.Run(flagshipConfig.Address + ":" + strconv.Itoa(flagshipConfig.Port))

}
