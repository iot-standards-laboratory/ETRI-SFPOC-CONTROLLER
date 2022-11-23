package router

import (
	"etri-sfpoc-controller/router/apiv2"
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

	v2 := apiEngine.Group("api/v2")
	{
		v2.POST("/init", apiv2.POST_init)
		v2.DELETE("/init", apiv2.DELETE_init)
	}

	r := gin.New()
	// r.Use(func(ctx *gin.Context) {
	// 	if statmgmt.Status() == statmgmt.STATUS_INIT {
	// 		ctx.JSON(http.StatusTemporaryRedirect, gin.H{
	// 			"path": "/init",
	// 		})
	// 		return
	// 	}
	// 	ctx.Next()
	// })

	assetEngine := gin.New()
	assetEngine.Static("/", "./front/build/web")
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
