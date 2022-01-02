package proxy

import (
	"github.com/valyala/fasthttp"
	"sync"
)

var clientPool = sync.Pool{
	New: func() interface{} {
		return new(fasthttp.Client)
	},
}

func acquireClient() *fasthttp.Client {
	return clientPool.Get().(*fasthttp.Client)
}

func releaseClient(c *fasthttp.Client) {
	c.Name = ""
	c.NoDefaultUserAgentHeader = false
	c.Dial = nil
	c.DialDualStack = false
	c.TLSConfig = nil
	c.MaxConnsPerHost = 0
	c.MaxIdleConnDuration = 0
	c.MaxConnDuration = 0
	c.MaxIdemponentCallAttempts = 0
	c.ReadBufferSize = 0
	c.WriteBufferSize = 0
	c.ReadTimeout = 0
	c.WriteTimeout = 0
	c.MaxResponseBodySize = 0
	c.DisableHeaderNamesNormalizing = false
	c.DisablePathNormalizing = false
	c.MaxConnWaitTimeout = 0
	c.RetryIf = nil
	clientPool.Put(c)
}
