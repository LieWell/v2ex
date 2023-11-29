package v2ex

import (
	"net/http"
	"net/url"
)

func proxyTransport() *http.Transport {
	proxyURL, _ := url.Parse("http://127.0.0.1:1082")
	return &http.Transport{Proxy: http.ProxyURL(proxyURL)}
}
