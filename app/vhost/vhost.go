package vhost

import (
	"ackycdn-node/app"
	"ackycdn-node/app/types"
	"github.com/asdine/storm/v3"
	"github.com/gookit/slog"
)

func GetConfigMem(domainName string) *types.VHostConfig {
	val, ok := app.G.VhostConfigsMem.Load(domainName)
	if !ok {
		return nil
	}
	if val == nil {
		slog.Error("failed to get vhost config from memory.")
		return nil
	}
	return val.(*types.VHostConfig)
}

func PutConfigMem(vhostConf *types.VHostConfig) {
	app.G.VhostConfigsMem.Delete(vhostConf.DomainName)
	app.G.VhostConfigsMem.Store(vhostConf.DomainName, vhostConf)
}

func GetConfigDB(domainName string) *types.VHostConfig {
	vhost := &types.VHostConfig{}
	err := app.G.PersistenceVhostDB.One("DomainName", domainName, vhost)
	if err != nil && err == storm.ErrNotFound {
		slog.Error("failed to get vhostutils config from DB, not found.")
		return nil
	}
	if err != nil {
		slog.Error(err)
		return nil
	}
	return vhost
}

func PutConfigDB(vhostConf *types.VHostConfig) {
	err := app.G.PersistenceVhostDB.Save(vhostConf)
	if err != nil {
		slog.Error(err)
	}
}
