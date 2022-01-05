package types

import (
	"ackycdn-node/app/cdncache"
	"github.com/asdine/storm/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/badger"
	"github.com/nats-io/nats.go"
	"sync"
)

// NodeConfig
// @Description: configuration of cdn node
type NodeConfig struct {
	CfgKey     string `storm:"id"`
	NodeId     string `storm:"unique"`
	MainIP     string
	CreateTime int64
}

type GlobalResource struct {
	FiberServer *fiber.App

	CdnCache *cdncache.CdnCache

	PersistenceVhostDB *storm.DB

	SessionStorage *badger.Storage

	VhostConfigsMem *sync.Map

	MqConnection *nats.Conn

	Mq nats.JetStreamContext

	NodeInfo *NodeConfig
}
