module etri-sfpoc-controller

go 1.17

replace etrisfpocctnmgmt => ../ETRI-SFPOC-CTNMGMT

replace etrisfpocdatamodel => ../ETRI-SFPOC-DATAMODEL

require (
	etrisfpocdatamodel v0.0.0-00010101000000-000000000000
	github.com/centrifugal/centrifuge-go v0.9.3
	github.com/eclipse/paho.mqtt.golang v1.4.2
	github.com/gin-gonic/gin v1.7.7
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/golang/glog v1.0.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/hashicorp/consul/api v1.15.3
	github.com/jacobsa/go-serial v0.0.0-20180131005756-15cf729a72d4
	github.com/magiconair/properties v1.8.6
	github.com/rjeczalik/notify v0.9.2
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f
)

require (
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/centrifugal/protocol v0.8.11 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.9.7 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/segmentio/encoding v0.3.5 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220829200755-d48e67d00261 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	gorm.io/gorm v1.23.3 // indirect
)
