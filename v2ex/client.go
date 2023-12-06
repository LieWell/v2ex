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

// ResponseInterface 定义一个接口类型,包含了 API 所有可能返回的类型
type ResponseInterface interface {
	Member | SiteState | Node | []Node | []Topic | []Reply
}

func buildGetRequest(rawUrl string) *httpclient.Request {
	r := httpclient.NewRequest(rawUrl, http.MethodGet)
	transport := ReadProxyTransportFromConfig()
	if transport != nil {
		r.SetTransport(transport)
	}
	return r
}

func doRequest[T ResponseInterface](req *httpclient.Request, t *T) *ResponseError {

	// 直接发起请求,如果报错的话,那么直接返回包装后的错误
	resp, err := req.DoRequest()
	if err != nil {
		return newResponseError(err, http.StatusTeapot, false, 0)
	}

	statusCode := resp.StatusCode
	// 请求正常响应,检查是否达到配额
	// 有时候 Header 中并不会返回配额消息,但是状态码是 403
	if statusCode == http.StatusForbidden || resp.Header.Get(rateLimitRemain) == "0" {
		reset, _ := strconv.Atoi(resp.Header.Get(rateLimitReset))
		return newResponseError(ReachLimitError, statusCode, true, reset)
	}

	// 读取响应
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return newResponseError(err, statusCode, false, 0)
	}

	// 尝试判断请求是否是异常请求
	// 这里认为大于 300 的状态码都算作异常
	if resp.StatusCode >= http.StatusMultipleChoices {
		return newResponseError(errors.New(string(all)), statusCode, false, 0)
	}

	// 反序列化为对象
	err = json.Unmarshal(all, &t)
	if err != nil {
		return newResponseError(err, statusCode, false, 0)
	}
	return nil
}

func QuerySiteState() (*SiteState, *ResponseError) {
	rawUrl := fmt.Sprintf("%s/site/stats.json", apiV1)
	req := buildGetRequest(rawUrl)

	var ss SiteState
	err := doRequest(req, &ss)
	return &ss, err
}

func QueryMember(id int) (*Member, *ResponseError) {
	rawUrl := fmt.Sprintf("%s/members/show.json?id=%d", apiV1, id)
	req := buildGetRequest(rawUrl)

	var m Member
	err := doRequest(req, &m)
	return &m, err
}

func QueryNodes() ([]Node, *ResponseError) {
	rawUrl := fmt.Sprintf("%s/nodes/all.json", apiV1)
	req := buildGetRequest(rawUrl)
	var n []Node
	err := doRequest(req, &n)
	return n, err
}

// QueryNode 可以通过节点的ID或者名称查询详情
// 编号或者名称选其一
func QueryNode(idOr int, nameOr string) (*Node, *ResponseError) {
	var rawUrl string
	if idOr != 0 {
		rawUrl = fmt.Sprintf("%s/nodes/show.json?id=%d", apiV1, idOr)
	} else {
		rawUrl = fmt.Sprintf("%s/nodes/show.json?name=%s", apiV1, nameOr)
	}
	req := buildGetRequest(rawUrl)

	var n Node
	err := doRequest(req, &n)
	return &n, err
}

// QueryNodesTopics 获取特定节点下的所有帖子
// node 编号或者名称选其一
func QueryNodesTopics(idOr int, nameOr string) ([]Topic, *ResponseError) {
	var rawUrl string
	if idOr != 0 {
		rawUrl = fmt.Sprintf("%s/topics/show.json?node_id=%d", apiV1, idOr)
	} else {
		rawUrl = fmt.Sprintf("%s/topics/show.json?node_name=%s", apiV1, nameOr)
	}
	req := buildGetRequest(rawUrl)

	var tArr []Topic
	err := doRequest(req, &tArr)
	return tArr, err
}

// QueryMemberTopics 获取用户发表的所有帖子
func QueryMemberTopics(username string) ([]Topic, *ResponseError) {
	rawUrl := fmt.Sprintf("%s/topics/show.json?username=%s", apiV1, username)
	req := buildGetRequest(rawUrl)

	var tArr []Topic
	err := doRequest(req, &tArr)
	return tArr, err
}

func QueryTopic(id int) (*Topic, *ResponseError) {
	rawUrl := fmt.Sprintf("%s/topics/show.json?id=%d", apiV1, id)
	req := buildGetRequest(rawUrl)

	var n []Topic
	err := doRequest(req, &n)

	// 先判断是否有正确得响应
	if len(n) > 0 {
		return &n[0], nil
	}
	return nil, err
}

// QueryReplies 返回所有的回复,包括 顶/赞之类的无效回复也会返回
// 因此此接口结果长度可能不等于 topic 接口下的 replies 数
func QueryReplies(topicId int) ([]Reply, *ResponseError) {
	rawUrl := fmt.Sprintf("%s/replies/show.json?topic_id=%d", apiV1, topicId)
	req := buildGetRequest(rawUrl)

	var tArr []Reply
	err := doRequest(req, &tArr)
	return tArr, err
}
