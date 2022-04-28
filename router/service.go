package router

import (
	"errors"
	"etri-sfpoc-controller/model/cache"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetServiceList(c *gin.Context) {
	defer handleError(c)

	w := c.Writer
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	list, ok := cache.GetSvcList("devicemanagerb")

	if !ok {
		panic(errors.New("not exist device"))
	}

	c.JSON(http.StatusOK, list)

}
