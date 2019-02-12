/*
Package main package is the entry file
*/
package main

import (
	"config"
	"controller"
	"flag"
	"log"
	"net"

	mid "middleware"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	validator "gopkg.in/go-playground/validator.v9"
)

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "debug the api interface")
	flag.Parse()
}

func main() {
	go socketMain()
	httpMain()
}

func socketMain() {
	log.Println("begin socket listening.")
	l, err := net.Listen("tcp", config.Conf.AppInfo.SocketAddr)

	// l, err := net.Listen("tcp", config.Conf.AppInfo.SocketAddr)
	if err != nil {
		log.Println("Error listening:", err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		// Handle connections in a new goroutine.
		go controller.HandleSocketRequest(conn)
	}
}

func httpMain() {
	e := echo.New()
	e.Use(middleware.Recover())

	// debug 为 true时不开jwt等验证
	// 运行时使用：go run src/main.go -debug=true 或者 bin/main -debug=true
	if !debug {
		skipperPath := []string{
			"/api/v1/weixin/code",
		}
		e.Use(middleware.JWTWithConfig(mid.CustomJWTConfig(skipperPath, "Bearer")))
	}

	// 参数验证器
	e.Validator = &mid.DefaultValidator{Validator: validator.New()}

	v1 := e.Group("/api/v1")
	v1.GET("/slogan", controller.GetSlogan)
	v1.GET("/uptoken/qiniu", controller.GetQiniuImgUpToken)

	weixin := v1.Group("/weixin")
	weixin.POST("/code", controller.SetUserInfoByCode)

	user := v1.Group("/user")
	user.POST("/action/add_emerg_ct", controller.AddEmergCT)
	user.POST("/action/del_emerg_ct", controller.DelEmergCT)
	user.POST("/action/add_share_ct", controller.AddShareCT)

	feedback := v1.Group("/feedback")
	feedback.POST("", controller.CreateFeedback)

	// 启动
	e.Logger.Fatal(e.Start(config.Conf.AppInfo.HTTPAddr))
}
