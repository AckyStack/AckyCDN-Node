package middlewares

import (
	"ackycdn-node/app"
	"ackycdn-node/app/cdncache"
	"ackycdn-node/app/logging"
	"ackycdn-node/app/types"
	"ackycdn-node/app/view"
	"ackycdn-node/pkg/ackyutils"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
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
	if !vhost.CacheControl.CacheEnabled {
		cacheKey := ackyutils.CacheUtils().BuildCacheKey(ctx)
		cacheData, _ := app.G.CacheStore.Get(cacheKey)

		//cache doesn't exist
		if cacheData == nil {
			return ctx.Next()
		}

		//cache exists
		cacheItem := cdncache.AcquireCacheItem()
		_, err := cacheItem.UnmarshalMsg(cacheData)
		if err != nil {
			slog.Error(err)
			return view.Send10xxErrorPage(app.ErrSystemInternal, ctx)
		}

		//put cache to the response
		ctx.Response().SetBodyRaw(cacheItem.Body)
		ctx.Response().SetStatusCode(cacheItem.StatusCode)
		ctx.Response().Header.SetContentTypeBytes(cacheItem.ContentType)
		if len(cacheItem.Encoding) > 0 {
			ctx.Response().Header.SetBytesV(fiber.HeaderContentEncoding, cacheItem.Encoding)
		}

		//finish logging
		logging.LogReqFinalize(ctx, true)

		//release cache item
		defer cdncache.ReleaseCacheItem(cacheItem)

		//return response
		return nil
	}
	//no cache exists or cache feature not enabled
	return ctx.Next()
}
