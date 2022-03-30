package devmanage

import (
	"context"
	"etri-sfpoc-controller/notifier"

	"github.com/gofrs/uuid"
)

func NewManager() (func(), func()) {
	ctx, cancel := context.WithCancel(context.Background())
	run := func() {
		go run(ctx)
	}

	return run, cancel
}

func run(ctx context.Context) {
	_uuid, _ := uuid.NewV4()
	_discoverCh := make(chan notifier.IEvent)
	notifier.Box.AddSubscriber(
		notifier.NewChanSubscriber(
			_uuid.String(),
			notifier.SubtokenDiscoveryDevice,
			notifier.SubtypeCont,
			_discoverCh,
		),
	)

	notifier.Box.AddSubscriber(
		notifier.NewCallbackSubscriber(
			_uuid.String(),
			notifier.SubtokenRcvCtrlMsg,
			notifier.SubtypeCont,
			HandleCtrlMsg,
		),
	)

	_disconntedCh := make(chan notifier.IEvent)
	notifier.Box.AddSubscriber(
		notifier.NewChanSubscriber(
			_uuid.String(),
			notifier.SubtokenDisconnted,
			notifier.SubtypeOnce,
			_disconntedCh,
		),
	)

	notifier.Box.AddSubscriber(
		notifier.NewCallbackSubscriber(
			_uuid.String(),
			notifier.SubtokenStatusChanged,
			notifier.SubtypeOnce,
			StatusReport,
		),
	)
	// notifier.Box.AddSubscriber(
	// 	notifier.NewChanSubscriber(
	// 		_uuid.String(),
	// 		notifier.SubtokenStatusChanged,
	// 		notifier.SubtypeCont,
	// 		_statusChangedCh,
	// 	),
	// )
}
