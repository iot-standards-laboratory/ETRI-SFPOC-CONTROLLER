package cachestorage

import (
	"sync"
)

// services to registered device
var svcMap = map[string]string{} // svcIds[sname] = meta of device list
var svcMutex sync.Mutex

// when AddService is called : Device registration is done
// always called with AddDeviceController
func AddSvc(sname, sid string) error {
	// start mutex
	svcMutex.Lock()
	defer svcMutex.Unlock()

	// get list for sname
	svcMap[sname] = sid

	return nil
}

func QuerySvcs() error {
	svcMutex.Lock()
	defer svcMutex.Unlock()

	// sname, sid, endpoint
	return nil
}

// func GetSvcUrls(sname, path string) (string, error) {

// 	if path[0] != '/' {
// 		return "", errors.New("path should start '/'")
// 	}

// 	var sid string
// 	var ok bool

// 	if strings.Compare(config.Params["mode"].(string), string(config.STANDALONE)) == 0 {
// 		sid, ok = config.Params["sid"].(string)
// 		if !ok {
// 			return "", errors.New("sid is blank")
// 		} else if strings.Compare(sid, "blank") == 0 {
// 			return "", errors.New("sid is blank")
// 		}
// 	} else {
// 		sid, ok = GetSvcId(sname)
// 		if !ok {
// 			return "", errors.New("sid is blank")
// 		}
// 	}

// 	return fmt.Sprintf("http://%s/svc/%s%s", config.Params["serverAddr"], sid, path), nil
// }

// func (s *_DBHandler) GetSID(sname string) (string, error) {
// 	sid, ok := s.sidCache[sname]
// 	if !ok {
// 		req, err := http.NewRequest("GET",
// 			fmt.Sprintf("http://%s/%s", config.Params["serverAddr"], "api/v1/svcs"),
// 			nil,
// 		)

// 		req.Header.Set("sname", sname)

// 		if err != nil {
// 			return "", err
// 		}

// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			return "", err
// 		} else if resp.ContentLength == 0 {
// 			return "", errors.New("not exist service")
// 		}

// 		b, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			return "", err
// 		}

// 		sid = string(b)
// 	}

// 	return sid, nil
// }
