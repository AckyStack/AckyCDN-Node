package ackyutils

import (
	"ackycdn-node/app"
)

type vhostUtils struct{}

func VhostUtils() *vhostUtils {
	return &vhostUtils{}
}

func (v *vhostUtils) BuildVhostCfgKey(domainName string) string {
	return app.PREFIX_VHOST + domainName
}
