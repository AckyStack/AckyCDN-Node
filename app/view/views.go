package view

import (
	"ackycdn-node/app"
	"ackycdn-node/app/logging"
	"ackycdn-node/pkg/views"
	realip "github.com/Ferluci/fast-realip"
	"github.com/gofiber/fiber/v2"
)

func Send5xxErrorPage() error {
	return nil
}
func Send10xxErrorPage(errorCode int, ctx *fiber.Ctx) error {
	ctx.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	logging.LogReqFinalizeEmpty(ctx)
	return ctx.SendString(views.Build10xxErrorPageHtml(errorCode, realip.FromRequest(ctx.Context()), app.G.NodeInfo.NodeId))
}

func SendRateLimitPage(ctx *fiber.Ctx) error {
	ctx.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	logging.LogReqFinalizeEmpty(ctx)
	return ctx.SendString(views.BuildRateLimitedPageHtml(ctx.Locals("clientid").(string), realip.FromRequest(ctx.Context()), app.G.NodeInfo.NodeId))
}

func SendWafBlockedPage(policyId int, ctx *fiber.Ctx) error {
	ctx.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	logging.LogReqFinalizeEmpty(ctx)
	return ctx.SendString(views.BuildWafBlockedPageHtml(policyId, realip.FromRequest(ctx.Context()), app.G.NodeInfo.NodeId))
}

func Send5sShieldPage(ctx *fiber.Ctx) error {
	ctx.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	logging.LogReqFinalizeEmpty(ctx)
	return ctx.SendString(views.BuildSShieldHtml(realip.FromRequest(ctx.Context()), app.G.NodeInfo.NodeId))
}

func SendDefaultPage(ctx *fiber.Ctx) error {
	ctx.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	logging.LogReqFinalizeEmpty(ctx)
	return ctx.SendString(views.BuildDefaultPageHtml(realip.FromRequest(ctx.Context()), app.G.NodeInfo.NodeId))
}
