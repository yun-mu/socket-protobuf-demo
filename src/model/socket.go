package model

import (
	"constant"
	"fmt"
	"net"

	proto "github.com/golang/protobuf/proto"
)

var (
	userToConn map[string]*net.Conn
	connToUser map[*net.Conn]string
)

func init() {
	userToConn = map[string]*net.Conn{}
	connToUser = map[*net.Conn]string{}
}

func AddSocketConn(userID string, conn *net.Conn) {
	userToConn[userID] = conn
	connToUser[conn] = userID
}

func DelSocketConnWithUserID(userID string) {
	delete(connToUser, userToConn[userID])
	delete(userToConn, userID)
}

func DelSocketConnWithConn(conn *net.Conn) {
	delete(userToConn, connToUser[conn])
	delete(connToUser, conn)
}

func GetSocketConn(userID string) (conn *net.Conn, ok bool) {
	conn, ok = userToConn[userID]
	return
}

func AddListenCT(userID, shareCT string) error {
	resPB := &AddShareCTRes{
		Header: &ResHeader{
			ResMethod: "AddListenCT",
		},
		UserID: shareCT,
	}
	out, err := proto.Marshal(resPB)
	if err != nil {
		return err
	}

	if conn, ok := userToConn[userID]; ok {
		(*conn).Write(out)
	} else {
		AddRedisMsgQuene(userID, out)
	}
	return nil
}

func Ping(conn *net.Conn, userID string) error {
	if _, ok := userToConn[userID]; !ok {
		userToConn[userID] = conn
		connToUser[conn] = userID
	}
	status, err := GetRedisUserStatus(userID)
	if err != nil {
		return err
	}

	if status == constant.UserStatusActive {
		return nil
	}
	return SetRedisUserStatus(userID, constant.UserStatusActive)
}

/****************************************** socket redis action ****************************************/

func AddRedisMsgQuene(userID string, data []byte) (err error) {
	cntlr := NewRedisDBCntlr()
	defer cntlr.Close()

	key := fmt.Sprintf(constant.RedisUserMsgQuene, userID)
	_, err = cntlr.RPUSH(key, data)
	return
}

func GetRedisMsgQuene(userID string) (data [][]byte, err error) {
	cntlr := NewRedisDBCntlr()
	defer cntlr.Close()

	key := fmt.Sprintf(constant.RedisUserMsgQuene, userID)
	data, err = cntlr.LRANGEGetBytes(key, 0, -1)
	if err != nil {
		return
	}
	_, err = cntlr.DEL(key)
	return
}
