package router

import (
	"etri-sfpoc-controller/notifier"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PutDevice(c *gin.Context) {
	defer handleError(c)

	w := c.Writer
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	notifier.Box.Publish(notifier.NewEvent("test", map[string]string{"Hello": "Wrold"}, notifier.SubtokenStatusChanged))
	c.Status(http.StatusOK)
}
