package router

import (
	"etri-sfpoc-controller/statmgmt"
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

	v2 := apiEngine.Group("api/v2")
	{
		v2.POST("/init", func(c *gin.Context) {
			// c.String(http.StatusOK, "Hello world")
			c.Status(http.StatusOK)
		})
	}

	r := gin.New()
	r.Use(func(ctx *gin.Context) {
		if statmgmt.Status() == statmgmt.STATUS_INIT {
			ctx.JSON(http.StatusTemporaryRedirect, gin.H{
				"path": "/init",
			})
			return
		}
		ctx.Next()
	})

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
