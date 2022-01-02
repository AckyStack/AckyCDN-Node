package types

//go:generate msgp
type CdnCacheItem struct {
	CacheKey    []byte
	Encoding    []byte
	Body        []byte
	StatusCode  int
	ContentType []byte
}

func (z *CdnCacheItem) Reset() {
	z.CacheKey = z.CacheKey[:0]
	z.Encoding = z.Encoding[:0]
	z.Body = z.Body[:0]
	z.StatusCode = 0
	z.ContentType = z.ContentType[:0]
}
