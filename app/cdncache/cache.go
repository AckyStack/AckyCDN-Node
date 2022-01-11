package cdncache

import (
	"ackycdn-node/app/types"
	"github.com/anxuanzi/goutils/pkg/ftaconv"
	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
	"github.com/gookit/slog"
	"time"
)

type CdnCache struct {
	db         *badger.DB
	gcInterval time.Duration
	done       chan struct{}
}

func InitCdnCache() *CdnCache {
	dbOptions := badger.DefaultOptions("./data/cache.db")
	dbOptions.SyncWrites = true
	dbOptions.Logger = nil
	dbOptions.Compression = options.ZSTD
	dbOptions.ZSTDCompressionLevel = 2
	db, err := badger.Open(dbOptions)
	if err != nil {
		slog.Panic(err)
	}
	store := &CdnCache{
		db:         db,
		gcInterval: 10 * time.Second,
		done:       make(chan struct{}),
	}
	// Start garbage collector
	go store.gc()
	return store
}

// GetCacheItem
// @Description: get a cache item from a cache database
// @receiver c
// @param cacheKey cache key, normally domain-md5val
// @param cacheItem item pointer, use pool to get an item pointer
// @return ok bool weather result is valid, if true there is an item found, if false means error or not found
func (c *CdnCache) GetCacheItem(cacheKey []byte, cacheItem *types.CdnCacheItem) (ok bool) {
	if len(cacheKey) <= 0 {
		return false
	}
	var data []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(cacheKey)
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		return err
	})
	// If no value was found return false
	if err == badger.ErrKeyNotFound {
		return false
	}
	if err != nil {
		slog.Error(err)
		return false
	}
	cacheItem.UnmarshalMsg(data)
	return true
}

func (c *CdnCache) SaveCacheItem(cacheItem *types.CdnCacheItem, exp time.Duration) (ok bool) {
	if len(cacheItem.CacheKey) <= 0 || cacheItem == nil {
		return false
	}
	data, _ := cacheItem.MarshalMsg(nil)
	entry := badger.NewEntry(cacheItem.CacheKey, data)
	if exp != 0 {
		entry.WithTTL(exp)
	}
	err := c.db.Update(func(tx *badger.Txn) error {
		return tx.SetEntry(entry)
	})
	if err != nil {
		slog.Error(err)
		return false
	}
	return true
}

func (c *CdnCache) DeleteSiteCacheItems(domainName string) (ok bool) {
	err := c.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := append(ftaconv.S2B(domainName), []byte("-")...)
		for it.Rewind(); it.ValidForPrefix(prefix); it.Next() {
			txn.Delete(it.Item().Key())
		}
		return nil
	})
	if err != nil {
		slog.Error(err)
		return false
	}
	return true
}

func (c *CdnCache) DeleteSingleCacheItem(cacheKey []byte) (ok bool) {
	err := c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(cacheKey)
	})
	if err != nil {
		return false
	}
	return true
}

func (c *CdnCache) Reset() error {
	return c.db.DropAll()
}

func (c *CdnCache) Close() error {
	c.done <- struct{}{}
	return c.db.Close()
}

func (c *CdnCache) gc() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			_ = c.db.RunValueLogGC(0.7)
		}
	}
}
