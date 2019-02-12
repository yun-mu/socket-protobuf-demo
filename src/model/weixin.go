package model

import (
	"config"
	"constant"
	"util/context"
	"util/log"

	"github.com/imroc/req"
)

type WeixinTokenRes struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}

var weixinLogger = log.GetLogger()

// GetWeixinAccessToken get the access token of WeChat Official Account
func GetWeixinAccessToken(code string) (WeixinTokenRes, error) {
	conf := config.Conf
	param := req.Param{
		"appid":      conf.Wechat.AppID,
		"secret":     conf.Wechat.AppSecret,
		"code":       code,
		"grant_type": "authorization_code",
	}

	weixinTokenRes := WeixinTokenRes{}
	err := context.BindGetJSONData(constant.URLToken, param, &weixinTokenRes)
	if err != nil {
		return weixinTokenRes, err
	}
	return weixinTokenRes, nil
}

// GetUserInfo get the userInfo by weixin access_token and openid
func GetUserInfo(accessToken, openid string) (User, error) {
	param := req.Param{
		"access_token": accessToken,
		"openid":       openid,
		"lang":         "zh_CN",
	}
	user := User{}
	err := context.BindGetJSONData(constant.URLUserInfo, param, &user)
	if err != nil {
		return user, err
	}

	if user.HeadImgURL == "" {
		user.HeadImgURL = constant.URLWeixinDefaultURL
	}
	return user, nil
}
