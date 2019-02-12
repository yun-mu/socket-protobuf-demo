package context

import (
	"constant"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/imroc/req"
	"github.com/labstack/echo"
)

// ErrorRes ErrorResponse
type ErrorRes struct {
	Status int    `json:"status"`
	ErrMsg string `json:"err_msg"`
}

// DataRes DataResponse
type DataRes struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// RetError response error, wrong response
func RetError(c echo.Context, code, status int, errMsg string) error {
	return c.JSON(code, ErrorRes{
		Status: status,
		ErrMsg: errMsg,
	})
}

// RetData response data, correct response
func RetData(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, DataRes{
		Status: 200,
		Data:   data,
	})
}

// BindGetJSONData bind the json data of method GET
// body must be a point
func BindGetJSONData(url string, param req.Param, body interface{}) error {
	r, err := req.Get(url, param)
	if err != nil {
		return err
	}
	err = r.ToJSON(body)
	if err != nil {
		return err
	}
	return nil
}

// jwt
func GetJWTUserID(c echo.Context) string {
	return c.Get(constant.JWTContextKey).(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(string)
}
