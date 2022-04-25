package controller

import (
	"flagship-io/flagship-sync-agent/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BucketingController struct {
}

func (controller *BucketingController) GetBucketing(c *gin.Context) {
	c.Data(http.StatusOK, gin.MIMEJSON, lib.BucktingFile)
}
