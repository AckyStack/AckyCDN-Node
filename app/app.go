package app

import (
	"ackycdn-node/app/cdncache"
	"ackycdn-node/app/types"
	"github.com/asdine/storm/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/badger"
	"github.com/nats-io/nats.go"
	"sync"
)

var G *GlobalResource

type GlobalResource struct {
	FiberServer *fiber.App

	CdnCache *cdncache.CdnCache

	PersistenceVhostDB *storm.DB

	SessionStorage *badger.Storage

	VhostConfigsMem *sync.Map

	MqConnection *nats.Conn

	Mq nats.JetStreamContext

	NodeInfo *types.NodeConfig
}
