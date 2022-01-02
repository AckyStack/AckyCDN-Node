package types

//go:generate msgp
type WafLog struct {

	//  WafType
	//  @Description: type of wafutils module, supports: filter, ratelimit
	WafType string

	//  Message
	//  @Description: macro expanded message
	Message string

	//  Data
	//  @Description: macro expanded logdata
	Data string

	//  AdditionalData
	//  @Description:
	AdditionalData []byte

	//  URI
	//  @Description: full request uri unparsed
	URI string

	//  TransactionID
	//  @Description: transaction id, it's request id, 19 characters
	TransactionID string

	//  Disruptive
	//  @Description: is disruptive
	Disruptive bool

	//  ClientIPAddress
	//  @Description: client IP address
	ClientIPAddress string

	//  CreateTime
	//  @Description: action trigger time
	CreateTime int64
}

type RequestLog struct {
	NodeId          string
	ClientId        string
	ClientIp        string
	ReqId           string
	ReqUA           string
	ReqReferer      string
	ReqMethod       string
	ReqProtocol     string
	ReqHost         string
	ReqUriScheme    string
	ReqUriPath      string
	ReqUriQss       string
	ReqFullUrl      string
	ReqTime         int64
	ResTime         int64
	UpstreamFullUrl string
	UpstreamReqTime int64
	UpstreamResTime int64
	CacheHit        bool
	ByteSend        int
}

func (z *RequestLog) Reset() {
	z.NodeId = ""
	z.ClientId = ""
	z.ClientIp = ""
	z.ReqId = ""
	z.ReqUA = ""
	z.ReqReferer = ""
	z.ReqMethod = ""
	z.ReqProtocol = ""
	z.ReqHost = ""
	z.ReqUriScheme = ""
	z.ReqUriPath = ""
	z.ReqUriQss = ""
	z.ReqFullUrl = ""
	z.ReqTime = 0
	z.ResTime = 0
	z.UpstreamFullUrl = ""
	z.UpstreamReqTime = 0
	z.UpstreamResTime = 0
	z.CacheHit = false
	z.ByteSend = 0
}

func (z *WafLog) Reset() {
	z.WafType = ""
	z.Message = ""
	z.Data = ""
	z.AdditionalData = nil
	z.URI = ""
	z.TransactionID = ""
	z.Disruptive = false
	z.ClientIPAddress = ""
	z.CreateTime = 0
}
