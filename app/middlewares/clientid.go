package middlewares

import (
	realip "github.com/Ferluci/fast-realip"
	"github.com/anxuanzi/goutils/pkg/ftaconv"
	"github.com/anxuanzi/goutils/pkg/ftamd5"
	"github.com/gofiber/fiber/v2"
)

func ClientIdMiddleware(ctx *fiber.Ctx) error {
	cid := ctx.Get("X-Ackycdn-Client-Id", genId(ctx))
	// Set new id to response header
	ctx.Set("X-Ackycdn-Client-Id", cid)

	// Add the request ID to locals
	ctx.Locals("clientid", cid)

	return ctx.Next()
}

func genId(ctx *fiber.Ctx) string {
	clientInfo := append(ftaconv.S2B(realip.FromRequest(ctx.Context())), ctx.Request().Header.UserAgent()...)
	return ftamd5.Md5Hash(clientInfo)
}
