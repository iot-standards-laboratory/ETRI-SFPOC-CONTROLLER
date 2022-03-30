package model

import (
	"errors"
	"etri-sfpoc-controller/config"
	"fmt"
	"io/ioutil"
	"net/http"
)

// func (s *_DBHandler) StatusCheck(did string, new map[string]interface{}) bool {
// 	origin, ok := s.states[did]
// 	if !ok {
// 		fmt.Println(did)
// 		fmt.Println("insert origin, before", s.states[did])
// 		s.states[did] = new
// 		// origin = map[string]interface{}{}
// 		// s.states[did] = origin
// 		// for k, v := range new {
// 		// 	origin[k] = v
// 		// }
// 		fmt.Println("insert origin, after", s.states[did])
// 		fmt.Println("insert origin, new", new)
// 		return true
// 	}

// 	changed := false
// 	for k, v := range new {
// 		if v.(float64) != origin[k].(float64) {
// 			fmt.Println("origin, new", v, origin[k])
// 			origin[k] = v
// 			changed = true
// 		}
// 	}

// 	return changed
// }

func (s *_DBHandler) GetSID(sname string) (string, error) {
	sid, ok := s.sidCache[sname]
	if !ok {
		req, err := http.NewRequest("GET",
			fmt.Sprintf("http://%s/%s", config.Params["serverAddr"], "services"),
			nil,
		)

		req.Header.Set("sname", sname)

		if err != nil {
			return "", err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		} else if resp.ContentLength == 0 {
			return "", errors.New("not exist service")
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		sid = string(b)
	}

	return sid, nil
}
