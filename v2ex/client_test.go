package v2ex

import (
	"fmt"
	"testing"
)

func TestQueryMemberById(t *testing.T) {
	rm, err := QueryMemberById(87)
	if err != nil {
		t.Error(err.StatusCode())
		t.Error(err.Error())
		t.Error(err.RateLimit())
		t.Error(err.RateLimitReset())
		return
	}
	fmt.Printf("%+v", rm)
}

func TestQuerySiteState(t *testing.T) {
	state, err := QuerySiteState()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v", state)
}
