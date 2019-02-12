package controller

import (
	"config"
	"constant"
	"fmt"
	"model"
	"net/http"
	"time"
	"util"
	"util/context"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

/**
 * @apiDefine CreateFeedback CreateFeedback
 * @apiDescription 添加反馈
 *
 * @apiParam  {String} contact_way 联系方式
 * @apiParam  {String} content 内容
 * @apiParam  {object[]} imgs 图片
 * @apiParam  {String} imgs.url 评论图片URL
 * @apiParam  {String} imgs.micro_url 评论图片缩略图URL
 *
 * @apiParamExample  {json} Request-Example:
 *     {
 *       "contact_way": "contact_way",
 *       "content": "content",
 *       "imgs": [{
 *           "url": "反馈图片链接"
 *           "micro_url": "反馈缩略图链接",
 *         }]
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
 * @api {post} /api/v1/feedback CreateFeedback
 * @apiVersion 1.0.0
 * @apiName CreateFeedback
 * @apiGroup Feedback
 * @apiUse CreateFeedback
 */
func CreateFeedback(c echo.Context) error {
	feedback := model.Feedback{}
	err := c.Bind(&feedback)
	if err != nil {
		writeFeedbackLog("CreateFeedback", constant.ErrorMsgParamWrong, err)
		return context.RetError(c, http.StatusBadRequest, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}

	feedback.UserID = context.GetJWTUserID(c)
	err = model.CreateFeedback(feedback)
	if err != nil {
		writeFeedbackLog("CreateFeedback", "插入反馈时服务器内部错误", err)
		return context.RetError(c, http.StatusBadGateway, http.StatusBadGateway, "插入反馈时服务器内部错误")
	}

	// 提醒管理员有人反馈了
	go func() {
		now := time.Now().Local().String()
		imgHTML := "<img src=\"%s\"  alt=\"反馈图片\" />"

		content := fmt.Sprintf(constant.EmailFeedbackNotice, feedback.Content, now)
		for _, img := range feedback.Imgs {
			content += "<br />" + fmt.Sprintf(imgHTML, img.URL)
		}
		util.SendEmail("安全APP", "安全APP 提醒", content, config.Conf.EmailInfo.To)
	}()

	return context.RetData(c, "")
}

func writeFeedbackLog(funcName, errMsg string, err error) {
	logger.WithFields(logrus.Fields{
		"package":  "controller",
		"file":     "feedback.go",
		"function": funcName,
		"err":      err,
	}).Warn(errMsg)
}
