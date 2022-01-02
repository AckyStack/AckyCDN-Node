package types

//go:generate msgp

// VHostConfig
// @Description: main config of a virtual host
type VHostConfig struct {

	//  DomainName
	//  @Description: the domain of the host
	DomainName string `storm:"id"`

	//  TlsConfig
	//  @Description: ssl configurations
	TlsConfig *TlsConfig

	//  DestinationProtocol
	//  @Description: the protocol uses to request backend, support: http, https, follow
	DestinationProtocol string

	//  DestinationHeaderRewrite
	//  @Description: Rewrite headers to destination
	DestinationHeaderRewrite map[string]string

	//  Upstreams
	//  @Description: proxy backend, could be multiple and load balanced
	Upstreams []*Upstream

	//  LoadBalanceMethod
	//  @Description: method of load balance, support: hash, random, rr, swrr, wr, wrr
	LoadBalanceMethod string

	//  SecurityControl
	//  @Description: Web Application Firewall configurations
	SecurityControl *SecurityConfig

	//  CacheControl
	//  @Description: Advanced cache config
	CacheControl *CacheConfig

	//  CompressionEnabled
	//  @Description: compress the response, supports gzip and br
	CompressionEnabled bool

	//  SeoOptimizationEnabled
	//  @Description: enable seo optimization with, bypass all security middleware send bot request to origin directly
	SeoOptimizationEnabled bool
}

// Upstream
// @Description: proxy backend
type Upstream struct {
	Host   string
	Port   int
	Weight int
}

type TlsConfig struct {
	//  SSLEnabled
	//  @Description: enable ssl for the site
	SSLEnabled bool

	//  RedirectHttpsEnabled
	//  @Description: force redirect https
	RedirectHttpsEnabled bool

	//  HSTSEnabled
	//  @Description: hsts enable
	HSTSEnabled bool

	//TODO do I need this?
	////  OCSPEnabled
	////  @Description: oscp header enable
	//OCSPEnabled bool `json:"ocsp_enabled"`

	//  Certificate
	//  @Description: certificate content
	Certificate []byte

	//  Key
	//  @Description: the key of certificate
	Key []byte
}

// SecurityConfig
// @Description: configuration of security settings
type SecurityConfig struct {

	//  SShieldEnabled
	//  @Description: enable 5 second shield to limit the request rate
	SShieldEnabled bool

	//  OwaspCRSEnabled
	//  @Description: the CRS provides protection against many common attack categories
	OwaspCRSEnabled bool

	//  AlwaysCaptchaEnabled
	//  @Description: when enabled, to visit the page, always require captcha
	AlwaysCaptchaEnabled bool

	//  RateLimitEnabled
	//  @Description: Enable rate limiter
	RateLimitEnabled bool

	//  RateLimitRate
	//  @Description: rate limiter value, works only when rate limiter enabled
	RateLimitRate int
}

// CacheConfig
// @Description: configuration of caches
type CacheConfig struct {

	//  CacheEnabled
	//  @Description: enable cache
	CacheEnabled bool

	//  CacheFileUrlSuffixes
	//  @Description: the request file suffix that want to cache (using regex to match url)
	CacheFileUrlSuffixes []string

	//  CacheExpiration
	//  @Description: cache expiration time, default unit second
	CacheExpiration int
}
