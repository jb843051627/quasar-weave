package notify

import (
	"context"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type Message struct {
	Alert   model.Alert `json:"alert"`
	Channel string      `json:"channel"`
}

type Sink interface {
	Send(context.Context, Message) error
}

type MultiSink []Sink

func (m MultiSink) Send(ctx context.Context, message Message) error {
	for _, sink := range m {
		if err := sink.Send(ctx, message); err != nil {
			return err
		}
	}
	return nil
}
