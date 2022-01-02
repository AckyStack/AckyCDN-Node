package middlewares

import (
	"ackycdn-node/app/types"
	"ackycdn-node/app/view"
	"ackycdn-node/pkg/ackyutils"
	"github.com/gofiber/fiber/v2"
)

func HstsMiddleware(ctx *fiber.Ctx) error {

	if ctx.Locals("vhostinfo") == nil {
		return view.SendDefaultPage(ctx)
	}
	vhost := ctx.Locals("vhostinfo").(*types.VHostConfig)

	if vhost.SeoOptimizationEnabled && ackyutils.WafUtils().IsCrawler(ctx.Request().Header.UserAgent()) && ackyutils.WafUtils().IsSearchEngine(ctx.Request().Header.UserAgent()) {
		return ctx.Next()
	}

	if vhost.TlsConfig.HSTSEnabled {
		ctx.Response().Header.Set("Strict-Transport-Security", "max-age=2592000")
		if ctx.Protocol() == "http" {
			ctx.Request().URI().SetScheme("https")
			return ctx.Redirect(ctx.Request().URI().String(), fiber.StatusTemporaryRedirect)
		}
	}

	return ctx.Next()
}
