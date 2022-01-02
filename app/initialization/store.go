package initialization

import (
	"ackycdn-node/app"
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
	cachedb := badger.New(badger.Config{
		Database:  "./data/cache.db",
		Reset:     false,
		Logger:    nil,
		UseLogger: false,
	})

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

	app.G.CacheStore = cachedb
	app.G.VhostConfigsMem = vcf
	app.G.PersistenceVhostDB = stormdb
	app.G.SessionStorage = sessiondb
}
