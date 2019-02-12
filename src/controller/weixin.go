package controller

import (
	"errors"
	"model"
	"net/http"
	"util/context"
	"util/token"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// SetUserInfoByCode get wechat access
/**
 * @apiDefine SetUserInfoByCode SetUserInfoByCode
 * @apiDescription 通过微信 code 获取用户信息并设置cookie之后进行302跳转
 *
 * @apiParam  {String} code 微信code
 * @apiParam  {String} state index url
 *
 * @apiParamExample   {json} Request-Example:
 *    {
 *      "code": "code",
 *      "state": "state"
 *    }
 *
 * @apiSuccess {Number} status=200 状态码
 * @apiSuccess {Object} data 正确返回数据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": 200,
 *       "data": {
 *	         "jwt_token": "jwt_token" // 有效时间为七天，发过来的时候需要在前面加上"Bearer "
 *         }
 *     }
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 302 StatusFound
 *
 * @apiError {Number} status 状态码
 * @apiError {String} err_msg 错误信息
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 401 Unauthorized
 *     {
 *       "status": 401,
 *       "err_msg": "Unauthorized"
 *     }
 */
/**
 * @api {post} /api/v1/weixin/code SetUserInfoByCode
 * @apiVersion 1.0.0
 * @apiName SetUserInfoByCode
 * @apiGroup Weixin
 * @apiUse SetUserInfoByCode
 */
func SetUserInfoByCode(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		writeWeixinLog("SetUserInfoByCode", "code is missing", errors.New("code is missing"))
		return context.RetError(c, http.StatusBadRequest, http.StatusBadRequest, "code is missing")
	}

	weixinTokenRes, err := model.GetWeixinAccessToken(code)
	if err != nil {
		writeWeixinLog("SetUserInfoByCode", "get weixin accessTokenRes faild", err)
		return context.RetError(c, http.StatusBadGateway, http.StatusBadGateway, "get weixin accessTokenRes faild")
	}

	userInfo, err := model.GetUserInfo(weixinTokenRes.AccessToken, weixinTokenRes.Openid)
	if err != nil {
		writeWeixinLog("SetUserInfoByCode", "get userInfo faild", err)
		return context.RetError(c, http.StatusBadGateway, http.StatusBadGateway, "get userInfo faild")
	}

	userID, err := model.CreateUser(userInfo)
	if err != nil {
		writeWeixinLog("SetUserInfoByCode", "get userInfo faild", err)
		return context.RetError(c, http.StatusBadGateway, http.StatusBadGateway, "get userInfo faild")

	}

	jwtAuth := map[string]interface{}{}
	jwtAuth["user_id"] = userID

	jwtToken := token.GetJWTToken(jwtAuth)
	return context.RetData(c, map[string]string{
		"jwt_token": jwtToken,
	})
}

func writeWeixinLog(funcName, errMsg string, err error) {
	logger.WithFields(logrus.Fields{
		"package":  "controller",
		"file":     "weixin.go",
		"function": funcName,
		"err":      err,
	}).Warnln(errMsg)
}
