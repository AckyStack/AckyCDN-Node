package ackyutils

import (
	"github.com/anxuanzi/goutils/pkg/ftaconv"
	"github.com/x-way/crawlerdetect"
	"regexp"
)

type wafUtils struct{}

func WafUtils() *wafUtils {
	return &wafUtils{}
}

func (w *wafUtils) IsSearchEngine(userAgent []byte) bool {
	match, _ := regexp.Match("Google|Baidu|MicroMessenger|miniprogram|bing|sogou|Yisou|360spider|soso|duckduck|Yandex|Yahoo|AOL|teoma", userAgent)
	if match {
		return true
	}
	return false
}

func (w *wafUtils) IsCrawler(userAgent []byte) bool {
	return crawlerdetect.IsCrawler(ftaconv.B2S(userAgent))
}
