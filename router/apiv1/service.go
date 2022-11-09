package apiv1

// func GetServiceList(c *gin.Context) {
// 	defer utils.HandleError(c)

// 	w := c.Writer
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

// 	svcs := cache.GetSvcList()

// 	c.JSON(http.StatusOK, svcs)

// }

// func GetServiceIds(c *gin.Context) {
// 	defer utils.HandleError(c)

// 	w := c.Writer
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

// 	c.JSON(http.StatusOK, cache.GetSvcIds())
// }
