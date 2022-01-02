package waf

import (
	"ackycdn-node/app"
	"ackycdn-node/app/types"
	"ackycdn-node/app/view"
	"ackycdn-node/pkg/ackyutils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-module/carbon/v2"
	"time"
)

func SshieldMiddleware(ctx *fiber.Ctx) error {
	//check and get host info
	if ctx.Locals("vhostinfo") == nil {
		return view.SendDefaultPage(ctx)
	}

	vhost := ctx.Locals("vhostinfo").(*types.VHostConfig)
	if !vhost.SecurityControl.SShieldEnabled {
		return ctx.Next()
	}

	//if seo optimization is enabled, skip sshield
	if vhost.SeoOptimizationEnabled && ackyutils.WafUtils().IsCrawler(ctx.Request().Header.UserAgent()) && ackyutils.WafUtils().IsSearchEngine(ctx.Request().Header.UserAgent()) {
		return ctx.Next()
	}

	//cookie exists? if not send shield
	if ctx.Cookies(CookieAckycdnClearance) == "" {
		return processSendSshield(ctx)
	}

	//get data from server
	codedToken, _ := app.G.SessionStorage.Get(ctx.Cookies(CookieAckycdnClearance))

	//token doesn't exist in the server, send shield
	if codedToken == nil {
		return processSendSshield(ctx)
	}

	//decode data
	token := acquireClearanceToken()
	token.UnmarshalMsg(codedToken)

	//check clearance
	if token.Cleared {
		releaseClearanceToken(token)
		return ctx.Next()
	}

	//not cleared directly, based on the condition check if it can be cleared
	//check time, current time more or equals than predicted pass time, then process next
	if carbon.Now().Gte(carbon.CreateFromTimestamp(token.CreateTime).AddSeconds(5)) {
		token.Cleared = true
		token.ClearanceTime = carbon.Now().TimestampWithMillisecond()
		codedToken, _ = token.MarshalMsg(nil)
		app.G.SessionStorage.Set(ctx.Cookies(CookieAckycdnClearance), codedToken, 10*time.Minute)
		releaseClearanceToken(token)
		return ctx.Next()
	}

	releaseClearanceToken(token)
	return processSendSshield(ctx)
}
