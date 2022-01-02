package waf

import (
	"ackycdn-node/app/types"
	"ackycdn-node/app/view"
	"ackycdn-node/pkg/ackyutils"
	realip "github.com/Ferluci/fast-realip"
	"github.com/anxuanzi/goutils/pkg/ftaconv"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

func NewRuleFilterMiddleware(ctx *fiber.Ctx) error {
	if ctx.Locals("vhostinfo") == nil {
		return view.SendDefaultPage(ctx)
	}
	vhost := ctx.Locals("vhostinfo").(*types.VHostConfig)

	if vhost.SeoOptimizationEnabled && ackyutils.WafUtils().IsCrawler(ctx.Request().Header.UserAgent()) && ackyutils.WafUtils().IsSearchEngine(ctx.Request().Header.UserAgent()) {
		return ctx.Next()
	}

	if !vhost.SecurityControl.OwaspCRSEnabled {
		return ctx.Next()
	}

	tx := WAF.WafFilterEngine.NewTransaction()
	tx.ID = ctx.Locals("requestid").(string)
	defer tx.ProcessLogging()

	//step one ProcessConnection
	port, err := ftaconv.ToInt(ctx.Port())
	if err != nil {
		slog.Error(err)
		return view.SendDefaultPage(ctx)
	}

	tx.ProcessConnection(realip.FromRequest(ctx.Context()), port, "", 0)
	tx.ProcessURI(ctx.Request().URI().String(), ctx.Method(), ctx.Protocol())
	ctx.Request().Header.VisitAll(func(key, value []byte) {
		tx.AddRequestHeader(ftaconv.BytesToString(key), ftaconv.ToString(value))
	})

	// We process phase 1 (Request)
	if it := tx.ProcessRequestHeaders(); it != nil {
		return processFilterInterruption(it, ctx)
	}

	tx.RequestBodyBuffer.Write(ctx.Body())
	if it, _ := tx.ProcessRequestBody(); it != nil {
		return processFilterInterruption(it, ctx)
	}

	return ctx.Next()
}
