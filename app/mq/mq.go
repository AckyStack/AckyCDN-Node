package mq

import (
	"ackycdn-node/app"
	"github.com/gookit/slog"
	"github.com/nats-io/nats.go"
)

type MQSubscribers struct{}

func NewSubscribers() *MQSubscribers {
	return &MQSubscribers{}
}

// StartSubscribers
// @Description: start listening on events from mq broker
// @receiver sub
func (mqs *MQSubscribers) StartSubscribers() {
	err := initSubscribers()
	if err != nil {
		slog.Panic(err)
	}
}

// InitSubscribers
// @Description: defines all the handlers based on the subject
// @return error
func initSubscribers() error {
	mqjs := app.G.Mq
	if _, err := mqjs.Subscribe("CONFIG.vhostAdd", addSite, nats.MaxDeliver(3)); err != nil {
		return err
	}

	if _, err := mqjs.Subscribe("CONFIG.vhostUpdate", updateSite, nats.MaxDeliver(3)); err != nil {
		return err
	}

	if _, err := mqjs.Subscribe("CONFIG.vhostDelete", deleteSite, nats.MaxDeliver(3)); err != nil {
		return err
	}

	if _, err := mqjs.Subscribe("CONFIG.vhostSync", syncSite, nats.MaxDeliver(3)); err != nil {
		return err
	}

	if _, err := mqjs.Subscribe("CONFIG.vhostSyncBatch", syncSites, nats.MaxDeliver(3)); err != nil {
		return err
	}

	//if _, err := JetStream.Subscribe("CMD.purgeCacheVhost", syncSites, nats.MaxDeliver(3)); err != nil {
	//	return err
	//}
	//
	//if _, err := JetStream.Subscribe("CMD.purgeCacheVhost", syncSites, nats.MaxDeliver(3)); err != nil {
	//	return err
	//}

	return nil
}
