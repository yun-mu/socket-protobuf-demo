package controller

import (
	"constant"
	"model"
	"net/http"
	"util/context"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

/**
 * @apiDefine AddEmergCT AddEmergCT
 * @apiDescription 添加紧急联系人
 *
 * @apiParam {String} name 名字
 * @apiParam {String} phone_num 手机号
 *
 * @apiParamExample  {json} Request-Example:
 *     {
 *       "name": "名字",
 *       "phone_num": "手机号",
 *     }
 *
 * @apiSuccess {Number} status=200 状态码
 * @apiSuccess {Object} data 正确返回数据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": 200,
 *       "data": ""
 *     }
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
 * @api {post} /api/v1/user/action/add_emerg_ct AddEmergCT
 * @apiVersion 1.0.0
 * @apiName AddEmergCT
 * @apiGroup User
 * @apiUse AddEmergCT
 */
func AddEmergCT(c echo.Context) error {
	param := model.NamePhone{}
	err := c.Bind(&param)
	if err != nil {
		writeUserLog("AddEmergCT", constant.ErrorMsgParamWrong, err)
		return context.RetError(c, http.StatusBadRequest, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}

	userID := context.GetJWTUserID(c)
	if err != nil || userID == "" {
		writeUserLog("AddEmergCT", "unAuth", err)
		return context.RetError(c, http.StatusUnauthorized, http.StatusUnauthorized, "unAuth")
	}
	err = model.AddEmergCT(userID, param)
	if err != nil {
		writeUserLog("AddEmergCT", "添加紧急联系人时服务器错误", err)
		return context.RetError(c, http.StatusBadGateway, http.StatusBadGateway, "添加紧急联系人时服务器错误")
	}
	return context.RetData(c, "")
}

/**
 * @apiDefine AddShareCT AddShareCT
 * @apiDescription 添加共享位置联系人
 *
 * @apiParam {String} user_id 共享联系人的user_id
 *
 * @apiParamExample  {json} Request-Example:
 *     {
 *       "user_id": "user_id",
 *     }
 *
 * @apiSuccess {Number} status=200 状态码
 * @apiSuccess {Object} data 正确返回数据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": 200,
 *       "data": ""
 *     }
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
 * @api {post} /api/v1/user/action/add_share_ct AddShareCT
 * @apiVersion 1.0.0
 * @apiName AddShareCT
 * @apiGroup User
 * @apiUse AddShareCT
 */
func AddShareCT(c echo.Context) error {
	param := UserIDParam{}
	err := c.Bind(&param)
	if err != nil {
		writeUserLog("AddShareCT", constant.ErrorMsgParamWrong, err)
		return context.RetError(c, http.StatusBadRequest, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}

	userID := context.GetJWTUserID(c)
	if err != nil || userID == "" {
		writeUserLog("AddShareCT", "unAuth", err)
		return context.RetError(c, http.StatusUnauthorized, http.StatusUnauthorized, "unAuth")
	}
	err = model.AddShareCT(userID, param.UserID)
	if err != nil {
		writeUserLog("AddShareCT", "添加共享联系人错误", err)
		return context.RetError(c, http.StatusBadGateway, http.StatusBadGateway, "添加共享联系人错误")
	}
	return context.RetData(c, "")
}

/**
 * @apiDefine DelEmergCT DelEmergCT
 * @apiDescription 删除紧急联系人
 *
 * @apiParam {String} name 名字
 *
 * @apiParamExample  {json} Request-Example:
 *     {
 *       "name": "名字",
 *     }
 *
 * @apiSuccess {Number} status=200 状态码
 * @apiSuccess {Object} data 正确返回数据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": 200,
 *       "data": ""
 *     }
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
 * @api {delete} /api/v1/user/action/del_emerg_ct DelEmergCT
 * @apiVersion 1.0.0
 * @apiName DelEmergCT
 * @apiGroup User
 * @apiUse DelEmergCT
 */
func DelEmergCT(c echo.Context) error {
	param := model.NamePhone{}
	err := c.Bind(&param)
	if err != nil {
		writeUserLog("DelEmergCT", constant.ErrorMsgParamWrong, err)
		return context.RetError(c, http.StatusBadRequest, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}

	userID := context.GetJWTUserID(c)
	if err != nil || userID == "" {
		writeUserLog("DelEmergCT", "unAuth", err)
		return context.RetError(c, http.StatusUnauthorized, http.StatusUnauthorized, "unAuth")
	}
	err = model.DelEmergCT(userID, param.Name)
	if err != nil || userID == "" {
		writeUserLog("DelEmergCT", "删除紧急联系人时服务器错误", err)
		return context.RetError(c, http.StatusBadGateway, http.StatusBadGateway, "删除紧急联系人时服务器错误")
	}
	return context.RetData(c, "")
}

func writeUserLog(funcName, errMsg string, err error) {
	logger.WithFields(logrus.Fields{
		"package":  "controller",
		"file":     "user.go",
		"function": funcName,
		"err":      err,
	}).Warn(errMsg)
}
