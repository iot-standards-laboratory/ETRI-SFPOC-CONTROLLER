package centrifuge_api

import (
	"errors"
	"log"

	"github.com/centrifugal/centrifuge-go"
)

var client *centrifuge.Client = nil

func NewClient(
	s_url, token string,
	onConnected centrifuge.ConnectedHandler,
	onDisconnected centrifuge.DisconnectHandler,
) error {
	client = centrifuge.NewJsonClient(
		s_url,
		centrifuge.Config{
			// Uncomment to make it work with Centrifugo JWT auth.
			Token: token,
		},
	)

	client.OnConnected(onConnected)
	client.OnDisconnected(onDisconnected)

	client.OnConnecting(func(e centrifuge.ConnectingEvent) {
		log.Printf("Connecting - %d (%s)", e.Code, e.Reason)
	})

	client.OnError(func(e centrifuge.ErrorEvent) {
		log.Printf("Error: %s", e.Error.Error())
	})

	client.OnMessage(func(e centrifuge.MessageEvent) {
		log.Printf("Message from server: %s", string(e.Data))
	})

	client.OnSubscribed(func(e centrifuge.ServerSubscribedEvent) {
		log.Printf("Subscribed to server-side channel %s: (was recovering: %v, recovered: %v)", e.Channel, e.WasRecovering, e.Recovered)
	})
	client.OnSubscribing(func(e centrifuge.ServerSubscribingEvent) {
		log.Printf("Subscribing to server-side channel %s", e.Channel)
	})
	client.OnUnsubscribed(func(e centrifuge.ServerUnsubscribedEvent) {
		log.Printf("Unsubscribed from server-side channel %s", e.Channel)
	})

	client.OnPublication(func(e centrifuge.ServerPublicationEvent) {
		log.Printf("Publication from server-side channel %s: %s (offset %d)", e.Channel, e.Data, e.Offset)
	})
	client.OnJoin(func(e centrifuge.ServerJoinEvent) {
		log.Printf("Join to server-side channel %s: %s (%s)", e.Channel, e.User, e.Client)
	})
	client.OnLeave(func(e centrifuge.ServerLeaveEvent) {
		log.Printf("Leave from server-side channel %s: %s (%s)", e.Channel, e.User, e.Client)
	})

	err := client.Connect()
	if err != nil {
		return err
	}
	return nil
}

func ResetClient() {
	client.Close()
}

func AddSubscription(channel string) (*centrifuge.Subscription, error) {
	if client == nil {
		return nil, errors.New("centrifuge client is nil error")
	}
	sub, err := client.NewSubscription(channel, centrifuge.SubscriptionConfig{
		Recoverable: true,
		JoinLeave:   true,
	})

	if err != nil {
		return nil, err
	}

	sub.OnSubscribing(func(e centrifuge.SubscribingEvent) {
		log.Printf("Subscribing on channel %s - %d (%s)", sub.Channel, e.Code, e.Reason)
	})
	sub.OnSubscribed(func(e centrifuge.SubscribedEvent) {
		log.Printf("Subscribed on channel %s, (was recovering: %v, recovered: %v)", sub.Channel, e.WasRecovering, e.Recovered)
	})
	sub.OnUnsubscribed(func(e centrifuge.UnsubscribedEvent) {
		log.Printf("Unsubscribed from channel %s - %d (%s)", sub.Channel, e.Code, e.Reason)
	})
	sub.OnError(func(e centrifuge.SubscriptionErrorEvent) {
		log.Printf("Subscription error %s: %s", sub.Channel, e.Error)
	})

	return sub, nil
}
