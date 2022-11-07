package apiv1

import (
	"etri-sfpoc-controller/notifier"
	"etri-sfpoc-controller/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PutDevice(c *gin.Context) {
	defer utils.HandleError(c)

	w := c.Writer
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	notifier.Box.Publish(notifier.NewEvent("test", map[string]string{"Hello": "Wrold"}, notifier.SubtokenStatusChanged))
	c.Status(http.StatusOK)
}

// func PostDevice(c *gin.Context) {
// 	defer handleError(c)

// 	w := c.Writer
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

// 	param := map[string]interface{}{}
// 	c.BindJSON(&param)

// 	respCh := make(chan bool)
// 	go devmanage.RegisterDevice(param, respCh)

// 	select {
// 	case b := <-respCh:
// 		if !b {
// 			panic(errors.New("something went wrong"))
// 		}
// 		c.String(http.StatusCreated, "OK")
// 	case <-c.Request.Context().Done():
// 		panic(errors.New("request is canceled"))
// 	}
// }
