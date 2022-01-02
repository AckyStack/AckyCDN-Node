package main

import (
	"ackycdn-node/app/initialization"
	"ackycdn-node/app/types"
	"ackycdn-node/app/vhost"
	"github.com/gookit/slog"
	"testing"
)

func TestInitData(t *testing.T) {
	slog.Info("starting ackycdn...")
	//begin initialization
	initialization.InitializeApplication()
	slog.Info("init data...")

	vhostinfo := &types.VHostConfig{
		DomainName: "test.ackycdn.com",
		TlsConfig: &types.TlsConfig{
			SSLEnabled:           false,
			RedirectHttpsEnabled: false,
			HSTSEnabled:          false,
			Certificate:          nil,
			Key:                  nil,
		},
		DestinationProtocol:      "https",
		DestinationHeaderRewrite: nil,
		Upstreams: []*types.Upstream{
			{
				Host:   "www.google.com",
				Port:   443,
				Weight: 0,
			},
			{
				Host:   "www.youtube.com",
				Port:   443,
				Weight: 0,
			},
			{
				Host:   "www.bing.com",
				Port:   443,
				Weight: 0,
			},
		},
		LoadBalanceMethod: "hash",
		SecurityControl: &types.SecurityConfig{
			SShieldEnabled:       true,
			OwaspCRSEnabled:      false,
			AlwaysCaptchaEnabled: false,
			RateLimitEnabled:     false,
			RateLimitRate:        0,
		},
		CacheControl: &types.CacheConfig{
			CacheEnabled:         false,
			CacheFileUrlSuffixes: nil,
			CacheExpiration:      0,
		},
		CompressionEnabled:     true,
		SeoOptimizationEnabled: true,
	}

	vhost.PutConfigDB(vhostinfo)

	ShutdownAll()

	//goleak.VerifyNone(t)
}
