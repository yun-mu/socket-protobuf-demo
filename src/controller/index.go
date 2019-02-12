package controller

import (
	"config"
	"constant"
	"model"
	"net/http"
	"strings"
	"util/context"
	"util/log"
	"util/token"

	"github.com/labstack/echo"
	"github.com/qiniu/api.v7/storage"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var logger = log.GetLogger()

/**
 * @apiDefine GetSlogan GetSlogan
 * @apiDescription 获取slogan
 *
 * @apiSuccess {Number} status=200 状态码
 * @apiSuccess {Object} data 正确返回数据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": 200,
 *       "data": "slogan"
 *     }
 *
 * @apiError {Number} status 状态码
 * @apiError {String} err_msg 错误信息
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 502 Bad Gateway
 *     {
 *       "status": 502,
 *       "err_msg": "Bad Gateway"
 *     }
 */
/**
 * @api {get} /api/v1/slogan GetSlogan
 * @apiVersion 1.0.0
 * @apiName GetSlogan
 * @apiGroup Index
 * @apiUse GetSlogan
 */
func GetSlogan(c echo.Context) error {
	return context.RetData(c, config.Conf.AppInfo.Slogan)
}

// GetQiniuImgUpToken 获取上传图片的七牛云upload-token
/**
 * @apiDefine GetQiniuImgUpToken GetQiniuImgUpToken
 * @apiDescription 获取上传图片的七牛云upload-token，地区：华南 链接：https://developer.qiniu.com/kodo/sdk/1236/android
 *
 * @apiParam {Number} type 类型：1->反馈图片
 * @apiParam {String} suffix 后缀，如：.jpg
 *
 * @apiParamExample  {query} Request-Example:
 *     {
 *       "type": 1,
 *       "suffix": ".jpg",
 *     }
 *
 * @apiSuccess {Number} status=200 状态码
 * @apiSuccess {Object} data 正确返回数据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": 200,
 *       "data": {
 *           "upload_token": "upload_token",
 *           "key": "key",
 *           "img": { // 上传到七牛云之后的url和自动持续化的缩略图：160 * 160
 *        	     "url": "url",
 *               "micro_url": "micro_url",
 *             }
 *         }
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
 * @api {get} /api/v1/uptoken/qiniu GetQiniuImgUpToken
 * @apiVersion 1.0.0
 * @apiName GetQiniuImgUpToken
 * @apiGroup Index
 * @apiUse GetQiniuImgUpToken
 */
func GetQiniuImgUpToken(c echo.Context) error {
	data := TypeSuffix{}
	err := c.Bind(&data)
	if err != nil {
		writeIndexLog("GetQiniuImgUpToken", constant.ErrorMsgParamWrong, err)
		return context.RetError(c, http.StatusBadRequest, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}

	imgID := uuid.Must(uuid.NewV4()).String()

	suffix := c.QueryParam("suffix")
	if suffix == "" {
		suffix = constant.ImgSuffix
	}

	imgPrefix, ok := constant.ImgPrefix[data.Type]
	if !ok {
		writeIndexLog("GetQiniuImgUpToken", constant.ErrorMsgParamWrong, err)
		return context.RetError(c, http.StatusBadRequest, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}
	microImgPrefix := constant.ImgPrefixMicro[data.Type]
	keyToOverwrite := imgPrefix + imgID + suffix
	saveAsKey := microImgPrefix + imgID + suffix

	fop := constant.ImgOps + "|saveas/" + storage.EncodedEntry(config.Conf.Qiniu.Bucket, saveAsKey)
	persistentOps := strings.Join([]string{fop}, ";")
	upToken := token.GetCustomUpToken(keyToOverwrite, persistentOps, constant.TokenQiniuExpire)

	img := model.Img{
		URL:      constant.ImgURIPrefix + keyToOverwrite,
		MicroURL: constant.ImgURIPrefix + saveAsKey,
	}
	resData := map[string]interface{}{
		"upload_token": upToken,
		"key":          keyToOverwrite,
		"img":          img,
	}
	return context.RetData(c, resData)
}

func writeIndexLog(funcName, errMsg string, err error) {
	logger.WithFields(logrus.Fields{
		"package":  "controller",
		"file":     "index.go",
		"function": funcName,
		"err":      err,
	}).Warn(errMsg)
}
