package v2ex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"liewell.fun/v2ex/httpclient"
	"net/http"
	"strconv"
)

var (
	host  = "https://www.v2ex.com"
	apiV1 = host + "/api"
	apiV2 = host + "/api/v2"
)

func QuerySiteState() (*ResponseState, error) {
	var resp ResponseState
	u := fmt.Sprintf("%s/site/stats.json", apiV1)
	err := httpclient.Get(u).SetTransport(proxyTransport()).ToJSON(&resp)
	return &resp, err
}

func QueryMemberById(id int) (*ResponseMember, *ResponseError) {

	url := fmt.Sprintf("%s/members/show.json?id=%d", apiV1, id)
	resp, err := httpclient.Get(url).SetTransport(proxyTransport()).DoRequest()
	if err != nil {
		return nil, &ResponseError{
			err:            err,
			statusCode:     http.StatusTeapot,
			rateLimit:      false,
			rateLimitReset: 0,
		}
	}

	// 达到配额
	if resp.Header.Get(rateLimitRemain) == "0" {
		reset, _ := strconv.Atoi(resp.Header.Get(rateLimitReset))
		return nil, &ResponseError{
			err:            errors.New("rate limit remaining 0"),
			statusCode:     resp.StatusCode,
			rateLimit:      true,
			rateLimitReset: reset,
		}
	}

	// 请求异常
	// 这里认为大于 400 的状态码都算作异常
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &ResponseError{
			err:            errors.New(http.StatusText(resp.StatusCode)),
			statusCode:     resp.StatusCode,
			rateLimit:      false,
			rateLimitReset: 0,
		}
	}

	// 读取响应
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ResponseError{
			err:            err,
			statusCode:     resp.StatusCode,
			rateLimit:      false,
			rateLimitReset: 0,
		}
	}

	// 返回响应详情
	var rm ResponseMember
	err = json.Unmarshal(all, &rm)
	if err != nil {
		return nil, &ResponseError{
			err:            err,
			statusCode:     resp.StatusCode,
			rateLimit:      false,
			rateLimitReset: 0,
		}
	}
	return &rm, nil
}
