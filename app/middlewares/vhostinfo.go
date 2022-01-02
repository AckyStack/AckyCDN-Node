package middlewares

import (
	"ackycdn-node/app/vhost"
	"ackycdn-node/app/view"
	"github.com/gofiber/fiber/v2"
)

func InjectHostInfoMiddleware(ctx *fiber.Ctx) error {
	vhost := vhost.GetConfigMem(ctx.Hostname())
	if vhost == nil {
		return view.SendDefaultPage(ctx)
	}
	ctx.Locals("vhostinfo", vhost)

	return ctx.Next()
}
