package constant

const (
	URLCode     = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	URLToken    = "https://api.weixin.qq.com/sns/oauth2/access_token"
	URLRefresh  = "https://api.weixin.qq.com/sns/oauth2/refresh_token"
	URLUserInfo = "https://api.weixin.qq.com/sns/userinfo"

	URLWeixinDefaultURL = "<default weixin url>"
)
