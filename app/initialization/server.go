package initialization

import (
	"ackycdn-node/app"
	"ackycdn-node/app/logging"
	"ackycdn-node/app/middlewares"
	"ackycdn-node/app/proxy"
	"ackycdn-node/app/view"
	"ackycdn-node/app/waf"
	"ackycdn-node/pkg/ratelimiter"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gookit/slog"
	"time"
)

// initFiberServer
// @Description: initialize a fiber app with all configuration that ackycdn needed
// @return *fiber.App
func initFiberServer() {
	fapp := fiber.New(fiber.Config{
		Prefork:      false, //can't be enabled!!!
		ServerHeader: "AckyCDN/v1",
		Concurrency:  1024 * 1024,
		IdleTimeout:  30 * time.Second,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			slog.Error(err)
			return view.Send10xxErrorPage(app.ErrSystemInternal, ctx)
		},
		DisableDefaultDate:           false,
		DisableDefaultContentType:    false,
		DisableHeaderNormalizing:     false,
		DisableStartupMessage:        false,
		AppName:                      "AckyCDN Node",
		StreamRequestBody:            false,
		DisablePreParseMultipartForm: true,
		ReduceMemoryUsage:            true,
		Network:                      fiber.NetworkTCP,
		EnableTrustedProxyCheck:      true,
	})

	//  logging middleware for debug
	fapp.Use(logger.New())
	// monitor middleware for debug
	fapp.Get("/.dev", monitor.New())
	fapp.Use(middlewares.LogInfoInjectMiddleware)
	fapp.Use(middlewares.RequestIdMiddleware())
	fapp.Use(middlewares.ClientIdMiddleware)
	fapp.Use(middlewares.InjectHostInfoMiddleware)
	fapp.Use(ratelimiter.New(ratelimiter.Config{LimitReached: func(ctx *fiber.Ctx) error {
		logging.LogWafRateLimitFinalize(ctx)
		return view.SendRateLimitPage(ctx)
	}}))
	fapp.Use(middlewares.HstsMiddleware)
	fapp.Use(middlewares.RedirectHttpsMiddleware)
	fapp.Use(middlewares.CompressMiddleware())
	fapp.Use(waf.SshieldMiddleware)
	fapp.Use(waf.NewRuleFilterMiddleware)
	fapp.Use(middlewares.CacheMiddleware)
	fapp.All("/*", proxy.ReverseProxyHandler)
	app.G.FiberServer = fapp
}
