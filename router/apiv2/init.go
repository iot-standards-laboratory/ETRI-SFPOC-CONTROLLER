package apiv2

import (
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/statmgmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func POST_init(c *gin.Context) {
	defer handleError(c)

	w := c.Writer
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if statmgmt.Status() != statmgmt.STATUS_INIT {
		panic(errors.New("invalid url"))
	}

	accessTkn := c.GetHeader("access_token")
	if len(accessTkn) <= 0 || strings.Compare(accessTkn, "etrismartfarm") != 0 {
		panic(errors.New("access token is invalid error"))
	}
	// glog.Infof("accessTkn is", accessTkn)

	var body = map[string]interface{}{}
	dec := json.NewDecoder(c.Request.Body)
	dec.Decode(&body)

	name, ok := body["name"]
	if !ok {
		panic(errors.New("agent name is invalid error"))
	}
	edgeAddress, ok := body["edgeAddress"]
	if !ok {
		panic(errors.New("edge address is invalid error"))
	}

	config.Set("name", name.(string))
	config.Set("edgeAddress", edgeAddress.(string))

	err := statmgmt.Register(accessTkn)
	if err != nil {
		panic(err)
	}

	// go statmgmt.Connect()
}

func DELETE_init(c *gin.Context) {
	defer handleError(c)

	w := c.Writer
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	accessTkn := c.GetHeader("access_token")
	if len(accessTkn) <= 0 || strings.Compare(accessTkn, "etrismartfarm") != 0 {
		panic(errors.New("access token is invalid error"))
	}
	// glog.Infof("accessTkn is", accessTkn)

	os.Remove("./config.properties")

	// db 초기화
	// Edge 초기화
}
