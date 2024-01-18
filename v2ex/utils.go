package v2ex

import (
	"fmt"
	"io"
	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/httpclient"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func GetImageAndSave(imageURL string, filename, directory string) error {

	// 获取图片
	response, err := httpclient.Get(imageURL).DoRequest()
	if err != nil {
		return err
	}
	defer func() {
		// TODO 需要关闭 body 并处理关闭时的错误吗 ???
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("status code[%d] not 200", response.StatusCode)
	}

	// 创建图片文件
	filePath := filepath.Join(directory, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// 写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
