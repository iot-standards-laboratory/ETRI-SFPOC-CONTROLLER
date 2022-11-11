package apiv2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func handleError(c *gin.Context) {
	if r := recover(); r != nil {
		glog.Info(r)
		c.String(http.StatusBadRequest, r.(error).Error())
	}
}
