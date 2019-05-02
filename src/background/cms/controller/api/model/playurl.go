package apimodel

import (
	"background/newmovie/model"
	"background/common/aes1"
	"background/common/logger"
)

type PlayUrl struct {
	Id             uint32  `json:"id"`
	Provider       uint32  `json:"provider"`
	Url            string  `json:"url"`
	Quality        uint8   `json:"quality"`
	IsPlay         bool    `json:"is_play"`
}

func PlayUrlFromDb(src model.PlayUrl) *PlayUrl {
	dst := PlayUrl{}
	dst.Id = src.Id
	dst.Provider = src.Provider
	dst.Url = src.Url
	dst.Quality = src.Quality
	if src.OnLine{
		dst.IsPlay = true
	}else{
		dst.IsPlay = false
	}

	//mac := hmac.New(sha1.New, []byte(secret))
	//mac.Write([]byte(dst.Url))
	//
	//dst.Url = strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))


	data,err := aes1.Encrypt([]byte(dst.Url))
	if err != nil{
		logger.Error(err)
		return nil
	}

	dst.Url = data
	return &dst
}
