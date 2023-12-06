package v2ex

import (
	"liewell.fun/v2ex/models"
	"time"
)

type Member struct {
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

func (rm *Member) toModel() *models.Member {
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
