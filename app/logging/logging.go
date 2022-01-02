package logging

import (
	"ackycdn-node/app"
	"ackycdn-node/app/types"
	realip "github.com/Ferluci/fast-realip"
	"github.com/anxuanzi/goutils/pkg/ftaconv"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/gookit/slog"
	"github.com/jptosso/coraza-waf/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"sync"
)

var reqLogPool = sync.Pool{
	New: func() interface{} {
		return new(types.RequestLog)
	},
}

var wafLogPool = sync.Pool{
	New: func() interface{} {
		return new(types.WafLog)
	},
}

func acquireReqLog() *types.RequestLog {
	return reqLogPool.Get().(*types.RequestLog)
}

func releaseReqLog(reql *types.RequestLog) {
	reql.Reset()
	reqLogPool.Put(reql)
}

func acquireWafLog() *types.WafLog {
	return reqLogPool.Get().(*types.WafLog)
}

func releaseWafLog(wafl *types.WafLog) {
	wafl.Reset()
	reqLogPool.Put(wafl)
}

func LogReqStart(ctx *fiber.Ctx) {
	reqCopy := fasthttp.AcquireRequest()
	ctx.Request().CopyTo(reqCopy)
	l := acquireReqLog()
	l.NodeId = app.G.NodeInfo.NodeId
	l.ClientId = "none"
	l.ClientIp = realip.FromRequest(ctx.Context())
	l.ReqId = "none"
	l.ReqUA = ftaconv.CopyString(ftaconv.B2S(reqCopy.Header.UserAgent()))
	l.ReqReferer = ftaconv.CopyString(ftaconv.B2S(reqCopy.Header.Referer()))
	l.ReqMethod = ftaconv.CopyString(ftaconv.B2S(reqCopy.Header.Method()))
	l.ReqProtocol = ftaconv.CopyString(ctx.Protocol())
	l.ReqHost = ftaconv.CopyString(ctx.Hostname())
	l.ReqUriScheme = ftaconv.CopyString(ftaconv.B2S(reqCopy.URI().Scheme()))
	l.ReqUriPath = ftaconv.CopyString(ctx.Path())
	l.ReqUriQss = ftaconv.CopyString(ftaconv.B2S(reqCopy.URI().QueryString()))
	l.ReqFullUrl = ftaconv.CopyString(ftaconv.CopyString(ctx.OriginalURL()))
	l.ReqTime = carbon.Now().TimestampWithMillisecond()
	l.ResTime = 0
	l.UpstreamFullUrl = ""
	l.UpstreamReqTime = 0
	l.UpstreamResTime = 0
	l.CacheHit = false
	l.ByteSend = 0
	ctx.Locals("cdnrequestlog", l)
	fasthttp.ReleaseRequest(reqCopy)
}

func LogReqUpstream(ctx *fiber.Ctx, upstreamFullUrl string, upstreamReqTime int64, upstreamResTime int64) {
	l := ctx.Locals("cdnrequestlog").(*types.RequestLog)
	l.UpstreamFullUrl = upstreamFullUrl
	l.UpstreamReqTime = upstreamReqTime
	l.UpstreamResTime = upstreamResTime
}

func LogReqFinalize(ctx *fiber.Ctx, cacheHit ...bool) {
	now := carbon.Now().TimestampWithMillisecond()
	isCacheHit := false
	if len(cacheHit) > 0 && cacheHit[0] {
		isCacheHit = true
	}
	l := ctx.Locals("cdnrequestlog").(*types.RequestLog)
	l.ClientId = ctx.Locals("clientid").(string)
	l.ReqId = ctx.Locals("requestid").(string)
	l.ResTime = now
	l.CacheHit = isCacheHit
	l.ByteSend = len(ctx.Response().Body())

	//do send log then release log
	codedLog, err := l.MarshalMsg(nil)
	if err != nil {
		slog.Error(err)
	}

	_, err = app.G.Mq.PublishAsync("LOG.req", codedLog)
	if err != nil {
		slog.Error(err)
	}
	defer releaseReqLog(l)
}

func LogReqFinalizeEmpty(ctx *fiber.Ctx) {
	l := ctx.Locals("cdnrequestlog").(*types.RequestLog)
	releaseReqLog(l)
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func LogWafFinalize(rule coraza.MatchedRule) {
	logData := acquireWafLog()
	ruleData, err := json.Marshal(rule.Rule)
	if err != nil {
		slog.Error(err)
	}
	logData.WafType = "filter"
	logData.Message = rule.Message
	logData.Data = rule.Data
	logData.URI = rule.URI
	logData.TransactionID = rule.ID
	logData.Disruptive = rule.Disruptive
	logData.ClientIPAddress = rule.ClientIPAddress
	logData.AdditionalData = ruleData
	logData.CreateTime = carbon.Now().TimestampWithMillisecond()

	codedLog, err := logData.MarshalMsg(nil)
	if err != nil {
		slog.Error(err)
	}

	_, err = app.G.Mq.PublishAsync("LOG.waf", codedLog)
	if err != nil {
		slog.Error(err)
	}
	defer releaseWafLog(logData)
}

func LogWafRateLimitFinalize(ctx *fiber.Ctx) {
	logData := acquireWafLog()
	logData.WafType = "ratelimit"
	logData.Message = "Too Many Requests"
	logData.Data = ftaconv.CopyString(ctx.Request().Header.String())
	logData.URI = ftaconv.CopyString(ctx.OriginalURL())
	logData.TransactionID = ctx.Locals("requestid").(string)
	logData.Disruptive = false
	logData.ClientIPAddress = realip.FromRequest(ctx.Context())
	logData.AdditionalData = nil
	logData.CreateTime = carbon.Now().TimestampWithMillisecond()

	codedLog, err := logData.MarshalMsg(nil)
	if err != nil {
		slog.Error(err)
	}

	_, err = app.G.Mq.PublishAsync("LOG.waf", codedLog)
	if err != nil {
		slog.Error(err)
	}
	defer releaseWafLog(logData)
}
