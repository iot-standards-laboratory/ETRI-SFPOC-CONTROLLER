package devmanage

// func RegisterHandler(e manager.Event) {
// 	param := e.Params()
// 	payload := map[string]interface{}{}
// 	payload["dname"] = param["uuid"]
// 	payload["type"] = "device"
// 	payload["sname"] = param["sname"]

// 	// b, err := json.Marshal(payload)
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// 	return
// 	// }
// 	// resp, err := http.Post("http://localhost:4000/devices", "application/json", bytes.NewReader(b))
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// 	return
// 	// }

// 	respCh := make(chan []byte)
// 	ctx := context.WithValue(context.Background(), managerKey(parameterKey), payload)
// 	ctx = context.WithValue(ctx, managerKey(waitResponseKey), respCh)
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()
// 	taskQueue <- &task{Event: DISCOVERY, Ctx: ctx}

// 	fmt.Println(string(<-respCh))
// }
