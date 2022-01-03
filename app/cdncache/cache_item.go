package cdncache

import (
	"ackycdn-node/app/types"
	"sync"
)

var cacheItemPool = sync.Pool{
	New: func() interface{} {
		return new(types.CdnCacheItem)
	},
}

func AcquireCacheItem() *types.CdnCacheItem {
	return cacheItemPool.Get().(*types.CdnCacheItem)
}

func ReleaseCacheItem(cci *types.CdnCacheItem) {
	cci.Reset()
	cacheItemPool.Put(cci)
}
