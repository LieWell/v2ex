package v2ex

import (
	"testing"
)

func TestStartAvatarSpider(t *testing.T) {
	prepareMysql()
	StartAvatarSpider()
}

func TestCheckMissingAvatarSpider(t *testing.T) {
	prepareMysql()
	CheckMissingAvatarSpider()
}
