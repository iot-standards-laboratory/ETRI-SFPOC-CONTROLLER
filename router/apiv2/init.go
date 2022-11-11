package apiv2

import (
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/statmgmt"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func POST_init(c *gin.Context) {
	defer handleError(c)
	if statmgmt.Status() != statmgmt.STATUS_INIT {
		panic(errors.New("invalid url"))
	}
	accessTkn := c.GetHeader("access_token")
	if len(accessTkn) <= 0 {
		panic(errors.New("access token is invalid error"))
	}
	glog.Infof("accessTkn is", accessTkn)

	var body = map[string]interface{}{}
	dec := json.NewDecoder(c.Request.Body)
	dec.Decode(&body)

	cname, ok := body["cname"]
	if !ok {
		panic(errors.New("container name is invalid error"))
	}

	config.Set("cname", cname.(string))

	err := statmgmt.Register(accessTkn)
	if err != nil {
		panic(err)
	}

	go statmgmt.Connect()
}
