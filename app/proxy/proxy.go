package proxy

import (
	"ackycdn-node/app"
	"ackycdn-node/app/logging"
	"ackycdn-node/app/types"
	"ackycdn-node/app/view"
	"ackycdn-node/pkg/ackyutils"
	"github.com/anxuanzi/goutils/pkg/ftaconv"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/gookit/slog"
	"github.com/valyala/fasthttp"
	"time"
)

// ReverseProxyHandler
// @Description: reverse proxy for CDN
// @param ctx fiber context
// @return error fiber error
func ReverseProxyHandler(ctx *fiber.Ctx) error {
	//nowTime := carbon.Now()
	if ctx.Locals("vhostinfo") == nil {
		return view.SendDefaultPage(ctx)
	}
	vhost := ctx.Locals("vhostinfo").(*types.VHostConfig)

	//create a new request client
	client := acquireClient()
	client.Name = "AckyCDN-Proxy"
	client.NoDefaultUserAgentHeader = true
	client.DialDualStack = false
	client.MaxConnsPerHost = 256
	client.MaxIdleConnDuration = 10 * time.Second
	client.MaxConnDuration = 30 * time.Second
	client.MaxIdemponentCallAttempts = 5
	client.ReadTimeout = 10 * time.Second
	client.WriteTimeout = 10 * time.Second
	client.DisableHeaderNamesNormalizing = false
	client.DisablePathNormalizing = true
	client.MaxConnWaitTimeout = 5 * time.Second

	//modify current fiber request to match the requirements

	//build request url
	ctx.Request().URI().SetHost(ackyutils.ProxyUtils().SelectHost(vhost))
	switch vhost.DestinationProtocol {
	case "https":
		ctx.Request().URI().SetScheme("https")
	case "http":
		ctx.Request().URI().SetScheme("http")
	default:
		ctx.Request().URI().SetSchemeBytes(ctx.Request().URI().Scheme())
	}

	//remove hop headers
	for _, key := range hopHeaders {
		ctx.Request().Header.Del(key)
	}

	//add proxy header if not exist
	if ctx.Request().Header.Peek("X-Forwarded-For") == nil {
		ctx.Request().Header.Add("X-Forwarded-For", ctx.IP())
	}

	if ctx.Request().Header.Peek("X-Real-IP") == nil {
		ctx.Request().Header.Add("X-Real-IP", ctx.IP())
	}

	if ctx.Request().Header.Peek("X-Forwarded-Proto") == nil {
		ctx.Request().Header.AddBytesV("X-Forwarded-Proto", ctx.Request().URI().Scheme())
	}

	if ctx.Request().Header.Peek("Host") == nil {
		ctx.Request().Header.Add("Host", ctx.Hostname())
	}

	proxyReqStartTime := carbon.Now().TimestampWithMillisecond()
	proxyResponse := fasthttp.AcquireResponse()
	err := client.DoTimeout(ctx.Request(), proxyResponse, 20*time.Second)
	if err != nil {
		slog.Error(err)
		return view.Send10xxErrorPage(app.ErrSystemInternal, ctx)
	}
	proxyReqEndTime := carbon.Now().TimestampWithMillisecond()
	logging.LogReqUpstream(ctx, ftaconv.B2S(ctx.Request().URI().FullURI()), proxyReqStartTime, proxyReqEndTime)

	ctx.Response().SetStatusCode(proxyResponse.StatusCode())
	ctx.Response().Header.SetContentTypeBytes(proxyResponse.Header.ContentType())
	if len(proxyResponse.Header.Peek(fiber.HeaderContentEncoding)) > 0 {
		ctx.Response().Header.SetBytesV(fiber.HeaderContentEncoding, proxyResponse.Header.Peek(fiber.HeaderContentEncoding))
	}
	ctx.Response().SetBodyRaw(proxyResponse.Body())
	logging.LogReqFinalize(ctx, false)

	defer fasthttp.ReleaseResponse(proxyResponse)
	defer releaseClient(client)
	return nil
}

// Hop-by-hop headers. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop headers are required to appear in the
// Connection header field. These are the headers defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.
var hopHeaders = []string{
	"Connection",          // Connection
	"Proxy-Connection",    // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",          // Keep-Alive
	"Proxy-Authenticate",  // Proxy-Authenticate
	"Proxy-Authorization", // Proxy-Authorization
	"Te",                  // canonicalized version of "TE"
	"Trailer",             // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",   // Transfer-Encoding
	"Upgrade",             // Upgrade
}
