package proxy

import (
	"ackycdn-node/app"
	"ackycdn-node/app/cdncache"
	"ackycdn-node/app/logging"
	"ackycdn-node/app/types"
	"ackycdn-node/app/view"
	"ackycdn-node/pkg/ackyutils"
	"github.com/anxuanzi/goutils/pkg/ftaconv"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/gookit/slog"
	"github.com/valyala/fasthttp"
	"regexp"
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
	//cache
	if vhost.CacheControl.CacheEnabled {
		ruleMatched, _ := regexp.Match("\\.?(eot|otf|ttf|woff|woff2|html|htm|css|js|jsx|less|scss|ppt|odp|doc|docx|ebook|log|md|msg|odt|org|pages|pdf|rtf|rst|tex|txt|wpd|wps|mobi|epub|azw1|azw3|azw4|azw6|azw|cbr|cbz|aac|aiff|ape|au|flac|gsm|it|m3u|m4a|mid|mod|mp3|mpa|pls|ra|s3m|sid|wav|wma|xm|3g2|3gp|aaf|asf|avchd|avi|drc|flv|m2v|m4p|m4v|mkv|mng|mov|mp2|mp4|mpe|mpeg|mpg|mpv|mxf|nsv|ogg|ogv|ogm|qt|rm|rmvb|roq|srt|svi|vob|webm|wmv|yuv|3dm|3ds|max|bmp|dds|gif|jpg|jpeg|png|psd|xcf|tga|thm|tif|tiff|ai|eps|ps|svg|dwg|dxf|gpx|kml|kmz|webp|ods|xls|xlsx|csv|ics|vcf)$", ctx.Request().URI().FullURI())
		if ruleMatched {
			cacheItem := cdncache.AcquireCacheItem()
			cacheItem.CacheKey = ackyutils.CacheUtils().BuildCacheKeyByte(ctx)
			cacheItem.StatusCode = proxyResponse.StatusCode()
			cacheItem.ContentType = ftaconv.CopyBytes(proxyResponse.Header.ContentType())
			cacheItem.Encoding = ftaconv.CopyBytes(proxyResponse.Header.Peek(fiber.HeaderContentEncoding))
			cacheItem.Body = ftaconv.CopyBytes(proxyResponse.Body())
			exp := time.Minute * 30
			if vhost.CacheControl.CacheExpiration > 0 {
				d := ftaconv.ToString(vhost.CacheControl.CacheExpiration) + "s"
				exp, _ = time.ParseDuration(d)
			}
			app.G.CdnCache.SaveCacheItem(cacheItem, exp)
			cdncache.ReleaseCacheItem(cacheItem)
		}
	}
	ctx.Response().Reset()
	proxyResponse.CopyTo(ctx.Response())
	logging.LogReqFinalize(ctx, false)
	releaseClient(client)
	fasthttp.ReleaseResponse(proxyResponse)
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
