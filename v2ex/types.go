package v2ex

import (
	"liewell.fun/v2ex/models"
	"time"
)

var (
	rateLimitQuota  = "x-rate-limit-limit"     // 请求速率配额,每IP每小时600
	rateLimitRemain = "x-rate-limit-remaining" // 配额剩余量
	rateLimitReset  = "x-rate-limit-reset"     // 重置时间戳,单位秒
)

type ResponseError struct {
	err            error
	statusCode     int
	rateLimit      bool
	rateLimitReset int
}

type ResponseState struct {
	TopicMax  int `json:"topic_max"`
	MemberMax int `json:"member_max"`
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

type ResponseMember struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Url          string `json:"url"`
	Website      string `json:"website"`
	Twitter      string `json:"twitter"`
	Psn          string `json:"psn"`
	Github       string `json:"github"`
	Btc          string `json:"btc"`
	Location     string `json:"location"`
	Tagline      string `json:"tagline"`
	Bio          string `json:"bio"`
	AvatarMini   string `json:"avatar_mini"`
	AvatarNormal string `json:"avatar_normal"`
	AvatarLarge  string `json:"avatar_large"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"last_modified"`
	Status       string `json:"status"`
}

func (rm *ResponseMember) toModel() *models.Member {
	return &models.Member{
		Number:     rm.Id,
		Name:       rm.Username,
		Website:    rm.Website,
		Twitter:    rm.Twitter,
		Github:     rm.Github,
		Location:   rm.Location,
		Tagline:    rm.Tagline,
		Avatar:     rm.AvatarNormal,
		Status:     rm.Status,
		CreateTime: time.Unix(rm.Created, 0),
	}
}
