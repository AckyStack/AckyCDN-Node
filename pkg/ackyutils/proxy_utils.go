package ackyutils

import (
	"ackycdn-node/app/types"
	"github.com/anxuanzi/goutils/pkg/ftabalancer"
	"github.com/anxuanzi/goutils/pkg/ftaconv"
)

type proxyUtils struct{}

func ProxyUtils() *proxyUtils {
	return &proxyUtils{}
}

func (p *proxyUtils) SelectHost(vhost *types.VHostConfig) string {
	upHostList := make([]string, len(vhost.Upstreams))
	upHostMap := make(map[string]int, len(vhost.Upstreams))
	for i, upstream := range vhost.Upstreams {
		upHostList[i] = upstream.Host + ":" + ftaconv.ToString(upstream.Port)
		upHostMap[upstream.Host+":"+ftaconv.ToString(upstream.Port)] = upstream.Weight
	}
	var balancer ftabalancer.Balancer
	//select destination
	switch vhost.LoadBalanceMethod {
	case "hash":
		balancer = ftabalancer.New(ftabalancer.ConsistentHash, upHostMap, upHostList)
	case "random":
		balancer = ftabalancer.New(ftabalancer.Random, upHostMap, upHostList)
	case "rr":
		balancer = ftabalancer.New(ftabalancer.RoundRobin, upHostMap, upHostList)
	case "swrr":
		balancer = ftabalancer.New(ftabalancer.SmoothWeightedRoundRobin, upHostMap, upHostList)
	case "wr":
		balancer = ftabalancer.New(ftabalancer.WeightedRand, upHostMap, upHostList)
	case "wrr":
		balancer = ftabalancer.New(ftabalancer.WeightedRoundRobin, upHostMap, upHostList)
	}
	return balancer.Select()
}

//func (p *proxyUtils)()  {
//
//}
