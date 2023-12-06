package v2ex

import (
	"encoding/json"
	"flag"
	"fmt"
	"liewell.fun/v2ex/core"
	"testing"
)

func prepareSetEnv() {
	c := flag.String("c", "../config.yaml", "")
	flag.Parse()
	core.LoadYamlConfig(*c)
	core.InitZap()
}

func TestQueryMemberById(t *testing.T) {
	prepareSetEnv()
	member, err := QueryMember(87)
	if err != nil {
		t.Error(err.StatusCode())
		t.Error(err.Error())
		t.Error(err.RateLimit())
		t.Error(err.RateLimitReset())
		return
	}
	marshal, _ := json.Marshal(member)
	fmt.Println(string(marshal))
}

func TestQuerySiteState(t *testing.T) {
	prepareSetEnv()
	state, err := QuerySiteState()
	if err != nil {
		t.Error(err)
		return
	}
	marshal, _ := json.Marshal(state)
	fmt.Println(string(marshal))
}

func TestQueryNodes(t *testing.T) {
	prepareSetEnv()
	nodeList, err := QueryNodes()
	if err != nil {
		t.Error(err)
		return
	}
	marshal, _ := json.Marshal(nodeList)
	fmt.Println(string(marshal))
}

func TestQueryNodeByIdOrName(t *testing.T) {
	prepareSetEnv()
	node, err := QueryNode(44, "")
	if err != nil {
		t.Error(err)
		return
	}
	marshal, _ := json.Marshal(node)
	fmt.Println(string(marshal))
}

func TestQueryTopic(t *testing.T) {
	prepareSetEnv()
	topic, err := QueryTopic(1)
	if err != nil {
		t.Error(err)
		return
	}
	marshal, _ := json.Marshal(topic)
	fmt.Println(string(marshal))
}

func TestQueryNodesTopics(t *testing.T) {
	prepareSetEnv()
	topic, err := QueryNodesTopics(1, "")
	if err != nil {
		t.Error(err)
		return
	}
	marshal, _ := json.Marshal(topic)
	fmt.Println(string(marshal))
}

func TestQueryMemberTopics(t *testing.T) {
	prepareSetEnv()
	topic, err := QueryMemberTopics("Livid")
	if err != nil {
		t.Error(err)
		return
	}
	marshal, _ := json.Marshal(topic)
	fmt.Println(string(marshal))
}

func TestQueryReplies(t *testing.T) {
	prepareSetEnv()
	replies, err := QueryReplies(1)
	if err != nil {
		t.Error(err)
		return
	}
	marshal, _ := json.Marshal(replies)
	fmt.Println("Replies length:", len(replies))
	fmt.Println(string(marshal))
}
