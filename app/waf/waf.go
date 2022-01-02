package waf

import (
	"ackycdn-node/app"
	"ackycdn-node/app/types"
	"ackycdn-node/app/view"
	"github.com/anxuanzi/goutils/pkg/ftaconv"
	"github.com/anxuanzi/goutils/pkg/ftananoid"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/gookit/slog"
	"github.com/jptosso/coraza-waf/v2"
	corazatypes "github.com/jptosso/coraza-waf/v2/types"
	"github.com/valyala/fasthttp"
	"sync"
	"time"
)

var WAF *Waf

type Waf struct {
	WafFilterEngine *coraza.Waf
}

const (
	CookieAckycdnClearance = "ackycdn_clearance"
)

var clearanceTokenPool = sync.Pool{
	New: func() interface{} {
		return new(types.SshieldClearanceToken)
	},
}

func acquireClearanceToken() *types.SshieldClearanceToken {
	return clearanceTokenPool.Get().(*types.SshieldClearanceToken)
}

func releaseClearanceToken(t *types.SshieldClearanceToken) {
	t.Reset()
	clearanceTokenPool.Put(t)
}

//func SshieldRedirectHandler(ctx *fiber.Ctx) error {
//	if ctx.Params("key", "none") == "none" {
//		slog.Info("param problem")
//		return view.Send10xxErrorPage(app.ErrRestricted, ctx)
//	}
//
//	if !cache.Has(ctx.Params("key")) {
//		slog.Info("cache problem")
//		return view.Send10xxErrorPage(app.ErrRestricted, ctx)
//	}
//
//	//check if client have a valid sshield token
//	sessionStore, err := app.G.ClearanceSessionStore.Get(ctx.Context())
//	if err != nil {
//		slog.Error(err)
//		return view.Send10xxErrorPage(app.ErrSystemInternal, ctx)
//	}
//
//	sToken := sessionStore.Get("sshield")
//	if sToken == nil {
//		slog.Error("client doesn't have sshield clearance record")
//		return DoSshieldBlock(ctx)
//	}
//
//	token := sToken.(*types2.SshieldClearanceToken)
//	if token == nil {
//		slog.Error("ctp not existence")
//		return DoSshieldBlock(ctx)
//	}
//
//	if !carbon.CreateFromTimestamp(token.CreateTime).AddSeconds(5).Lte(carbon.Now()) {
//		slog.Info("time problem")
//		return view.Send10xxErrorPage(app.ErrRestricted, ctx)
//	}
//	token.Cleared = true
//	token.ClearanceTime = carbon.Now().TimestampWithMillisecond()
//	sessionStore.Set("sshield", token)
//	sessionId := ftaconv.CopyBytes(sessionStore.GetSessionID())
//	err = app.G.ClearanceSessionStore.Save(ctx.Context(), sessionStore)
//	if err != nil {
//		slog.Error(err)
//		return view.Send10xxErrorPage(app.ErrSystemInternal, ctx)
//	}
//	ctx.Response().Header.DelCookie("ackycdn_clearance")
//	ck := fasthttp.AcquireCookie()
//	ck.SetKey("ackycdn_clearance")
//	ck.SetPath("/")
//	ck.SetHTTPOnly(true)
//	ck.SetDomain(ctx.Hostname())
//	ck.SetValueBytes(sessionId)
//	if ctx.Secure() {
//		ck.SetSameSite(fasthttp.CookieSameSiteNoneMode)
//		ck.SetSecure(true)
//	} else {
//		ck.SetSameSite(fasthttp.CookieSameSiteLaxMode)
//		ck.SetSecure(false)
//	}
//	ctx.Request().Header.SetCookieBytesKV(ck.Key(), ck.Value())
//	ctx.Response().Header.SetCookie(ck)
//	fasthttp.ReleaseCookie(ck)
//	return ctx.Redirect(cache.Get(ctx.Params("key")).(string), fiber.StatusTemporaryRedirect)
//}

func processFilterInterruption(it *corazatypes.Interruption, ctx *fiber.Ctx) error {
	return view.SendWafBlockedPage(it.RuleID, ctx)
}

func processSendSshield(ctx *fiber.Ctx) error {
	token := acquireClearanceToken()
	token.CreateTime = carbon.Now().TimestampWithMillisecond()
	token.Cleared = false
	token.ClearanceTime = 0
	sid := ftananoid.GenerateNanoUUID() + ftaconv.ToString(carbon.Now().TimestampWithNanosecond())
	tokenCoded, _ := token.MarshalMsg(nil)
	err := app.G.SessionStorage.Set(sid, tokenCoded, 10*time.Minute)
	if err != nil {
		slog.Error(err)
		return view.Send10xxErrorPage(app.ErrSystemInternal, ctx)
	}
	releaseClearanceToken(token)
	ctx.Response().Header.DelCookie(CookieAckycdnClearance)
	ctx.Request().Header.DelCookie(CookieAckycdnClearance)
	ck := fasthttp.AcquireCookie()
	ck.SetKey(CookieAckycdnClearance)
	ck.SetPath("/")
	ck.SetHTTPOnly(true)
	ck.SetDomain(ctx.Hostname())
	ck.SetValue(sid)
	if ctx.Secure() {
		ck.SetSameSite(fasthttp.CookieSameSiteNoneMode)
		ck.SetSecure(true)
	} else {
		ck.SetSameSite(fasthttp.CookieSameSiteLaxMode)
		ck.SetSecure(false)
	}
	ctx.Request().Header.SetCookieBytesKV(ck.Key(), ck.Value())
	ctx.Response().Header.SetCookie(ck)
	fasthttp.ReleaseCookie(ck)
	return view.Send5sShieldPage(ctx)
}
