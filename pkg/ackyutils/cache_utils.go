package ackyutils

import (
	"github.com/anxuanzi/goutils/pkg/ftaconv"
	"github.com/anxuanzi/goutils/pkg/ftamd5"
	"github.com/gofiber/fiber/v2"
)

type cacheUtils struct{}

func CacheUtils() *cacheUtils {
	return &cacheUtils{}
}

func (cu *cacheUtils) BuildCacheKey(ctx *fiber.Ctx) string {
	return ftamd5.Md5Hash(append(append(ctx.Request().Host(), []byte("-")...), append(append(ctx.Request().URI().Scheme(), ctx.Request().URI().Host()...), ctx.Request().URI().RequestURI()...)...))
}

func (cu *cacheUtils) BuildCacheKeyByte(ctx *fiber.Ctx) []byte {
	return ftaconv.S2B(cu.BuildCacheKey(ctx))
}
