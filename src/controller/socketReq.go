package controller

import (
	"model"
	"net"

	proto "github.com/golang/protobuf/proto"
)

/****************************************** req handler ****************************************/

func PingReq(conn *net.Conn, data []byte, n int, reqBasic ReqBasic) error {
	return model.Ping(conn, reqBasic.UserID)
}

func SetUserLocEQReq(conn *net.Conn, data []byte, n int, reqBasic ReqBasic) (err error) {
	reqPB := new(model.SetUserLocEQReq)
	err = proto.Unmarshal(data[:n], reqPB)
	if err != nil {
		return
	}
	userID := reqBasic.UserID
	loc := reqPB.GetLoc()

	err = model.SetRedisUserLocEQ(userID, loc.GetLongitude(), loc.GetLatitude(), reqPB.GetEQ())
	return
}

func GetUserLocEQReq(conn *net.Conn, data []byte, n int, reqBasic ReqBasic) (err error) {
	reqPB := new(model.GetUserLocEQReq)
	err = proto.Unmarshal(data[:n], reqPB)
	if err != nil {
		return
	}
	userID := reqPB.GetUserID()

	location, eq, err := model.GetRedisUserLocEQ(userID)
	resPB := &model.GetUserLocEQRes{
		Header: &model.ResHeader{
			ResMethod: MethodGetUserLocEQ,
		},
		Loc: &location,
		EQ:  int32(eq),
	}
	out, err := proto.Marshal(resPB)
	if err != nil {
		return
	}
	(*conn).Write(out)
	return
}

func GetMsgQueneReq(conn *net.Conn, data []byte, n int, reqBasic ReqBasic) (err error) {
	resData, err := model.GetRedisMsgQuene(reqBasic.UserID)
	resPB := &model.GetMsgQueneRes{
		Header: &model.ResHeader{
			ResMethod: MethodGetMsgQuene,
		},
		Data: resData,
	}
	out, err := proto.Marshal(resPB)
	if err != nil {
		return
	}
	(*conn).Write(out)
	return
}
