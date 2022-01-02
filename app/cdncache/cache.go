package cdncache

import (
	"ackycdn-node/app/types"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		return new(types.CdnCacheItem)
	},
}

func AcquireCacheItem() *types.CdnCacheItem {
	return pool.Get().(*types.CdnCacheItem)
}

func ReleaseCacheItem(cci *types.CdnCacheItem) {
	cci.Reset()
	pool.Put(cci)
}
