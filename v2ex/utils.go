package v2ex

import (
	"liewell.fun/v2ex/core"
	"net/http"
	"net/url"
	"strings"
)

func ReadProxyTransportFromConfig() *http.Transport {

	proxy := core.GlobalConfig.Http.Proxy
	if !strings.HasPrefix(proxy, core.HttpProtocolPrefix) || !strings.HasPrefix(proxy, core.HttpsProtocolPrefix) {
		return nil
	}

	proxyURL, err := url.Parse(core.GlobalConfig.Http.Proxy)
	if err != nil {
		core.Logger.Warnf("config Http.Proxy[%s] invalid,ignore", proxy)
		return nil
	}

	return &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
}
