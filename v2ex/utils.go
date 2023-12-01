package v2ex

import (
	"liewell.fun/v2ex/core"
	"net/http"
	"net/url"
)

func proxyTransport() *http.Transport {
	proxyURL, _ := url.Parse(core.GlobalConfig.Http.Proxy)
	return &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
}
