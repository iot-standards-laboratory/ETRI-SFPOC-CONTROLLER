package router

import (
	"etri-sfpoc-controller/router/apiv1"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// var db model.DBHandlerI

func init() {
	var err error
	// db, err = model.NewSqliteHandler("dump.db")
	if err != nil {
		panic(err)
	}
}

func NewRouter() *gin.Engine {
	apiEngine := gin.New()

	v1 := apiEngine.Group("api/v1")
	{
		v1.PUT("/devs", apiv1.PutDevice)
		// // for debug api
		// apiv1.GET("/svcs", GetServiceList)
		// apiv1.GET("/svcids", GetServiceIds)
		// apiv1.POST("/devs/discover", PostDevice)
	}

	v2 := apiEngine.Group("api/v2")
	{
		v2.GET("/init", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello")
		})
	}

	// pushEngine := gin.New()
	// pushEngine.Any("/*any", func(c *gin.Context) {
	// 	// GetPublish(c)
	// 	Test(c)
	// })
	r := gin.New()

	assetEngine := gin.New()
	assetEngine.Static("/", "./static")
	r.Any("/*any", func(c *gin.Context) {
		path := c.Param("any")
		if strings.HasPrefix(path, "/api/v1") || strings.HasPrefix(path, "/api/v2") {
			apiEngine.HandleContext(c)
		} else {
			assetEngine.HandleContext(c)
		}
	})

	return r
}
