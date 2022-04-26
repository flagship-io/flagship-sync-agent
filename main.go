package main

import (
	"flagship-io/flagship-sync-agent/controller"
	"flagship-io/flagship-sync-agent/lib"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setupRouter(bucketingPolling *lib.BucketingPolling) *gin.Engine {
	bucketingController := controller.BucketingController{
		BucketingPolling: bucketingPolling,
	}
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	router.GET("/bucketing", bucketingController.GetBucketing)

	return router
}

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

	server := setupRouter(&BucketingPolling)

	server.Run(flagshipConfig.Address + ":" + strconv.Itoa(flagshipConfig.Port))

}
