package router

import (
	"etri-sfpoc-controller/model"
	"strings"

	"github.com/gin-gonic/gin"
)

var db model.DBHandlerI

func init() {
	var err error
	db, err = model.NewSqliteHandler("dump.db")
	if err != nil {
		panic(err)
	}
}

func NewRouter() *gin.Engine {
	apiEngine := gin.New()

	apiv1 := apiEngine.Group("api/v1")
	{
		apiv1.PUT("/devs", PutDevice)
		apiv1.POST("/devs/discover", PostDevice)
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
		if strings.HasPrefix(path, "/api/v1") {
			apiEngine.HandleContext(c)
		} else {
			assetEngine.HandleContext(c)
		}
	})

	return r
}
