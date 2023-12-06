package v2ex

import "fmt"

var (
	rateLimitQuota  = "x-rate-limit-limit"     // 请求速率配额,每IP每小时600
	rateLimitRemain = "x-rate-limit-remaining" // 配额剩余量
	rateLimitReset  = "x-rate-limit-reset"     // 重置时间戳,单位秒
)

var (
	ReachLimitError = fmt.Errorf("rate limit remaining 0")
)

type ResponseError struct {
	err            error
	statusCode     int
	rateLimit      bool
	rateLimitReset int
}

func (re *ResponseError) Error() string {
	return re.err.Error()
}

func (re *ResponseError) StatusCode() int {
	return re.statusCode
}

func (re *ResponseError) RateLimit() bool {
	return re.rateLimit
}

func (re *ResponseError) RateLimitReset() int {
	return re.rateLimitReset
}

func newResponseError(err error, code int, limit bool, reset int) *ResponseError {
	return &ResponseError{
		err:            err,
		statusCode:     code,
		rateLimit:      limit,
		rateLimitReset: reset,
	}
}
