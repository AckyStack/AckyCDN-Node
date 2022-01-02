package initialization

import (
	"ackycdn-node/app"
	"ackycdn-node/app/mq"
	"github.com/gookit/slog"
	"github.com/nats-io/nats.go"
)

func initMq() {
	//initialization message queue
	mqconn, err := nats.Connect("104.194.243.158:4222", nats.Token("anxuanzi"))
	if err != nil {
		slog.Panic(err)
	}

	mqjs, err := mqconn.JetStream() //nats.PublishAsyncMaxPending(256)
	if err != nil {
		slog.Panic(err)
	}

	app.G.MqConnection = mqconn
	app.G.Mq = mqjs

	mq.NewSubscribers().StartSubscribers()
}
