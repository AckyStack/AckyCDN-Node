package initialization

import (
	"ackycdn-node/app/logging"
	"ackycdn-node/app/waf"
	"github.com/gookit/slog"
	"github.com/jptosso/coraza-waf/v2"
	"github.com/jptosso/coraza-waf/v2/seclang"
)

func initWaf() {
	//session store for waf sshield
	//sessionStoreShield := session.New(session.Config{
	//	Expiration:     5 * time.Minute,
	//	KeyLookup:      "cookie:ackycdn_sshsid",
	//	CookiePath:     "/",
	//	CookieSecure:   true,
	//	CookieHTTPOnly: true,
	//	KeyGenerator:   ftananoid.GenerateNanoUUID,
	//	CookieSameSite: "None",
	//})
	//sessionStoreShield.RegisterType(waf.ClearanceTokenPayload{})
	filterEngine := coraza.NewWaf()
	parser, _ := seclang.NewParser(filterEngine)
	files := []string{
		"./conf/coraza.conf",
		"./conf/crs-setup.conf",
		"./conf/rules/*.conf",
	}
	for _, f := range files {
		if err := parser.FromFile(f); err != nil {
			slog.Panic(err)
		}
	}
	filterEngine.SetErrorLogCb(logging.LogWafFinalize)
	waf.WAF = &waf.Waf{
		WafFilterEngine: filterEngine,
	}
}
