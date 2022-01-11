package middlewares

import (
	"ackycdn-node/app"
	"ackycdn-node/app/cdncache"
	"ackycdn-node/app/logging"
	"ackycdn-node/app/types"
	"ackycdn-node/app/view"
	"ackycdn-node/pkg/ackyutils"
	"github.com/anxuanzi/goutils/pkg/ftaconv"
	"github.com/gofiber/fiber/v2"
)

func CacheMiddleware(ctx *fiber.Ctx) error {
	//cache only works on the get request
	if !ctx.Request().Header.IsGet() {
		return ctx.Next()
	}

	//try to get vhost information from context
	if ctx.Locals("vhostinfo") == nil {
		return view.SendDefaultPage(ctx)
	}
	vhost := ctx.Locals("vhostinfo").(*types.VHostConfig)

	//if cache is enabled
	if vhost.CacheControl.CacheEnabled {
		cacheKey := ackyutils.CacheUtils().BuildCacheKeyByte(ctx)
		cacheItem := cdncache.AcquireCacheItem()
		ok := app.G.CdnCache.GetCacheItem(cacheKey, cacheItem)
		//maybe cache doesn't exist
		if !ok {
			return ctx.Next()
		}
		//put cache to the response
		ctx.Response().SetBodyRaw(ftaconv.CopyBytes(cacheItem.Body))
		ctx.Response().SetStatusCode(cacheItem.StatusCode)
		ctx.Response().Header.SetContentTypeBytes(ftaconv.CopyBytes(cacheItem.ContentType))
		if len(cacheItem.Encoding) > 0 {
			ctx.Response().Header.SetBytesV(fiber.HeaderContentEncoding, ftaconv.CopyBytes(cacheItem.Encoding))
		}
		//finish logging
		logging.LogReqFinalize(ctx, true)
		//release cache item
		cdncache.ReleaseCacheItem(cacheItem)
		//return response
		return nil
	}

	//no cache exists or cache feature not enabled
	return ctx.Next()
}
