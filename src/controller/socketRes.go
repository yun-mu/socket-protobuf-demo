package controller

import (
	"model"
	"net"
	"net/http"

	proto "github.com/golang/protobuf/proto"
)

/****************************************** res handler ****************************************/

func PingRes(conn *net.Conn, err error) {
	resPB := model.PingRes{
		Header: &model.ResHeader{
			ResMethod: MethodPing,
		},
	}
	out, _ := proto.Marshal(&resPB)
	(*conn).Write(out)
}

func SocketErrorRes(conn *net.Conn, err error) {
	resPB := model.ErrorRes{
		Header: &model.ResHeader{
			ResMethod: MethodError,
		},
		Status: http.StatusBadRequest,
		ErrMsg: err.Error(),
	}
	out, _ := proto.Marshal(&resPB)
	(*conn).Write(out)
}
