package middlewares

import (
	"github.com/anxuanzi/goutils/pkg/ftananoid"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func RequestIdMiddleware() fiber.Handler {
	return requestid.New(requestid.Config{
		Header: "X-Ackycdn-Req-Id",
		Generator: func() string {
			return ftananoid.GenerateNanoId("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 19)
		},
	})
}
