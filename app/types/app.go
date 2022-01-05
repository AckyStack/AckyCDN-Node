package types

// NodeConfig
// @Description: configuration of cdn node
type NodeConfig struct {
	CfgKey     string `storm:"id"`
	NodeId     string `storm:"unique"`
	MainIP     string
	CreateTime int64
}
