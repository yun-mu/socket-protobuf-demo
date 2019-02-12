package controller

import (
	"io"
	"model"
	"net"
	"time"

	proto "github.com/golang/protobuf/proto"

	"util/token"

	"github.com/sirupsen/logrus"
)

type ReqBasic struct {
	Method string
	UserID string
}

type SocketFunc func(*net.Conn, []byte, int, ReqBasic) error

var methodToFunc map[string]SocketFunc

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	MethodPing         = "Ping"
	MethodGetUserLocEQ = "GetUserLocEQ"
	MethodSetUserLocEQ = "SetUserLocEQ"
	MethodError        = "Error"
	MethodAddListenCT  = "AddListenCT"
	MethodGetMsgQuene  = "GetMsgQuene"
)

func init() {
	methodToFunc = map[string]SocketFunc{
		MethodPing:         PingReq,
		MethodGetUserLocEQ: GetUserLocEQReq,
		MethodSetUserLocEQ: SetUserLocEQReq,
	}
}

/****************************************** basic ****************************************/

// TODO: 心跳维持, 可参考：https://github.com/gorilla/websocket/blob/master/examples/chat/client.go
func HandleSocketRequest(conn net.Conn) {
	defer conn.Close()

	for {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetWriteDeadline(time.Now().Add(writeWait))

		data, n, reqBasic, err := BasicMid(&conn)
		if err != nil {
			if err == io.EOF {
				// 断线重连
				model.DelSocketConnWithConn(&conn)
				break
			}
			SocketErrorRes(&conn, err)
			writeSocketLog("Un Auth", err)
			continue
		}
		err = methodToFunc[reqBasic.Method](&conn, data, n, reqBasic)
		if err != nil {
			SocketErrorRes(&conn, err)
			writeSocketLog(reqBasic.Method, err)
			continue
		}
	}
}

func BasicMid(conn *net.Conn) (data []byte, n int, reqBasic ReqBasic, err error) {
	data = make([]byte, 0)
	reqBasic = ReqBasic{}
	n, err = (*conn).Read(data)
	if err != nil {
		return
	}

	reqPB := new(model.PingReq)
	err = proto.Unmarshal(data[:n], reqPB)
	if err != nil {
		return
	}

	reqBasic.Method = reqPB.GetHeader().GetReqMethod()

	claims, err := token.GetJWTClaim(reqPB.GetHeader().GetAuthToken())
	if err != nil {
		return
	}

	reqBasic.UserID = claims["user_id"].(string)
	return
}

func writeSocketLog(funcName string, err error) {
	logger.WithFields(logrus.Fields{
		"package":    "controller",
		"file":       "socket.go",
		"req_method": funcName,
		"err":        err,
	}).Warnln("")
}
