package ratelimiter

import (
	"ackycdn-node/app/types"
	"ackycdn-node/app/view"
	"github.com/anxuanzi/goutils/pkg/ftaconv"
	"github.com/anxuanzi/goutils/pkg/ftamd5"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Max number of recent connections during `Duration` seconds before sending a 429 response
	//
	// Default: 5
	Max int

	// KeyGenerator allows you to generate custom keys, by default c.IP() is used
	//
	// Default: func(c *fiber.Ctx) string {
	//   return c.IP()
	// }
	KeyGenerator func(*fiber.Ctx) string

	// Expiration is the time on how long to keep records of requests in memory
	//
	// Default: 1 * time.Minute
	Expiration time.Duration

	// LimitReached is called when a request hits the limit
	//
	// Default: func(c *fiber.Ctx) error {
	//   return c.SendStatus(fiber.StatusTooManyRequests)
	// }
	LimitReached fiber.Handler

	// When set to true, requests with StatusCode >= 400 won't be counted.
	//
	// Default: false
	SkipFailedRequests bool

	// When set to true, requests with StatusCode < 400 won't be counted.
	//
	// Default: false
	SkipSuccessfulRequests bool

	// LimiterMiddleware is the struct that implements a limiter middleware.
	//
	// Default: a new Fixed Window Rate Limiter
	LimiterMiddleware LimiterHandler
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Max:        5,
	Expiration: time.Second,
	Next: func(c *fiber.Ctx) bool {
		if c.Locals("vhostinfo") == nil {
			return true
		}
		vhost := c.Locals("vhostinfo").(*types.VHostConfig)
		if !vhost.SecurityControl.RateLimitEnabled {
			return true
		}
		return false
	},
	KeyGenerator: func(c *fiber.Ctx) string {
		cid := c.Locals("clientid")
		if cid == nil {
			cid := c.Get("X-Ackycdn-Client-Id", genId(c))
			// Set new id to response header
			c.Set("X-Ackycdn-Client-Id", cid)
			// Add the request ID to locals
			c.Locals("clientid", cid)
		}
		return cid.(string)
	},
	LimitReached: func(c *fiber.Ctx) error {
		return view.SendRateLimitPage(c)
	},
	SkipFailedRequests:     false,
	SkipSuccessfulRequests: false,
	LimiterMiddleware:      SlidingWindow{},
}

func genId(ctx *fiber.Ctx) string {
	clientInfo := append(ftaconv.S2B(ctx.IP()), ctx.Request().Header.UserAgent()...)
	return ftamd5.Md5Hash(clientInfo)
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}
	if cfg.Max <= 0 {
		cfg.Max = ConfigDefault.Max
	}
	if int(cfg.Expiration.Seconds()) <= 0 {
		cfg.Expiration = ConfigDefault.Expiration
	}
	if cfg.KeyGenerator == nil {
		cfg.KeyGenerator = ConfigDefault.KeyGenerator
	}
	if cfg.LimitReached == nil {
		cfg.LimitReached = ConfigDefault.LimitReached
	}
	if cfg.LimiterMiddleware == nil {
		cfg.LimiterMiddleware = ConfigDefault.LimiterMiddleware
	}
	return cfg
}
