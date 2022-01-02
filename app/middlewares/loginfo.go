package middlewares

import (
	"ackycdn-node/app/logging"
	"github.com/gofiber/fiber/v2"
)

func LogInfoInjectMiddleware(ctx *fiber.Ctx) error {
	logging.LogReqStart(ctx)
	return ctx.Next()
}
