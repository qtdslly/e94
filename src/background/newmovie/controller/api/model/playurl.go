package apimodel

import (
	"background/newmovie/model"
)

type PlayUrl struct {
	Id             uint32  `json:"id"`
	Provider       uint32  `json:"provider"`
	Url            string  `json:"url"`
	IsPlay         bool    `json:"is_play"`
}

func PlayUrlFromDb(src model.PlayUrl) *Video {
	dst := PlayUrl{}
	dst.Id = src.Id
	dst.Provider = src.Provider
	dst.Url = src.Url
	if src.OnLine{
		dst.IsPlay = true
	}else{
		dst.IsPlay = false
	}

	return &dst
}
