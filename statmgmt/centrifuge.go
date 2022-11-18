package statmgmt

import (
	"errors"
	"etri-sfpoc-controller/centrifuge_api"
	"etri-sfpoc-controller/config"
	"fmt"
	"log"

	"github.com/centrifugal/centrifuge-go"
)

func connectCentrifuge(wsAddr string) error {
	cid, ok := config.Params["cid"]
	if !ok {
		return errors.New("cid is invalid error")
	}

	cname, ok := config.Params["cname"]
	if !ok {
		return errors.New("cname is invalid error")
	}

	token, _ := centrifuge_api.IssueJWT(cid.(string), cname.(string), "ctrl", "/", nil)

	connectedHandler := func(e centrifuge.ConnectedEvent) {
	}

	disconnectedHandler := func(e centrifuge.DisconnectedEvent) {
	}

	err := centrifuge_api.NewClient(
		wsAddr,
		token,
		connectedHandler,
		disconnectedHandler,
	)
	if err != nil {
		centrifuge_api.ResetClient()
		return err
	}

	sub, err := centrifuge_api.AddSubscription(fmt.Sprintf("public:%s", cid))
	if err != nil {
		return err
	}

	sub.OnPublication(func(e centrifuge.PublicationEvent) {
		// control device
		log.Printf("Someone says via channel %s: %s (offset %d)", sub.Channel, e.Data, e.Offset)
	})

	err = sub.Subscribe()
	if err != nil {
		return err
	}

	return nil
}
