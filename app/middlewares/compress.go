package middlewares

import (
	"ackycdn-node/app/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func CompressMiddleware() fiber.Handler {
	return compress.New(compress.Config{
		Next: func(c *fiber.Ctx) bool {
			if c.Locals("vhostinfo") == nil {
				return true
			}
			vhost := c.Locals("vhostinfo").(*types.VHostConfig)
			if vhost.CompressionEnabled {
				return false
			}
			return true
		},
		Level: compress.LevelDefault,
	})
}
