package v2ex

import (
	"liewell.fun/v2ex/core"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func ReadProxyTransportFromConfig() *http.Transport {

	proxy := core.GlobalConfig.Http.Proxy
	if !strings.HasPrefix(proxy, core.HttpProtocolPrefix) && !strings.HasPrefix(proxy, core.HttpsProtocolPrefix) {
		core.Logger.Warnf("config Http.Proxy[%s] invalid schema,ignore", proxy)
		return nil
	}

	proxyURL, err := url.Parse(core.GlobalConfig.Http.Proxy)
	if err != nil {
		core.Logger.Warnf("config Http.Proxy[%s] invalid URL,ignore", proxy)
		return nil
	}

	core.Logger.Debugf("config Http.Proxy[%s] check pass ", proxy)
	return &http.Transport{
		Proxy:                 http.ProxyURL(proxyURL),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}
}
