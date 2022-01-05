package initialization

import (
	"ackycdn-node/app"
	"ackycdn-node/app/cdncache"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/codec/msgpack"
	"github.com/gofiber/storage/badger"
	"github.com/gookit/slog"
	"sync"
)

func initStores() {
	//initialization databases
	stormdb, err := storm.Open("./data/ackycdn.db", storm.Codec(msgpack.Codec), storm.Batch())
	if err != nil {
		slog.Panic(err)
	}

	//initialize db for web cache
	cachedb := cdncache.InitCdnCache()

	//initialize db for session storage
	sessiondb := badger.New(badger.Config{
		Database:  "./data/session.db",
		Reset:     true,
		Logger:    nil,
		UseLogger: false,
	})

	//memory storage for site configs
	vcf := new(sync.Map)
	if vcf == nil {
		slog.Panic("failed to initialize vhost config map")
	}

	app.G.CdnCache = cachedb
	app.G.VhostConfigsMem = vcf
	app.G.PersistenceVhostDB = stormdb
	app.G.SessionStorage = sessiondb
}
