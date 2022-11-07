package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context) {
	if r := recover(); r != nil {
		log.Println(r)
		c.String(http.StatusBadRequest, r.(error).Error())
	}
}
