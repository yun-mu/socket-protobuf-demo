package model

import (
	"constant"
	"errors"
	"fmt"
	"strconv"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User the struct of the table user
type User struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"` // indexes
	// user status:
	Status int `bson:"status" json:"status"`

	// weixin info
	Openid     string   `json:"openid" bson:"openid"` // unique indexes
	Unionid    string   `json:"unionid" bson:"unionid"`
	Nickname   string   `json:"nickname" bson:"nickname"`
	Sex        int      `json:"sex" bson:"sex"`
	Province   string   `json:"province" bson:"province"`
	City       string   `json:"city" bson:"city"`
	Country    string   `json:"country" bson:"country"`
	HeadImgURL string   `json:"headimgurl" bson:"headimgurl"`
	Privilege  []string `json:"privilege" bson:"privilege"`

	// private information
	Name     string `bson:"name" json:"name"`           // real name
	PhoneNum string `bson:"phone_num" json:"phone_num"` // unique indexes
	Email    string `bson:"email" json:"email"`

	EmergCT []NamePhone `bson:"emerg_ct" json:"emerg_ct"` // emergency contact
}

type NamePhone struct {
	Name     string `bson:"name" json:"name"`
	PhoneNum string `bson:"phone_num" json:"phone_num"`
}

func CreateUser(user User) (string, error) {
	if user.Openid == "" {
		return "", errors.New("lack openid")
	}
	query := bson.M{
		"openid": user.Openid,
	}
	oldUser, err := findUser(query, bson.M{
		"headimgurl": 1,
		"nickname":   1,
	})
	if err != nil {
		user.ID = bson.NewObjectId()
		return user.ID.Hex(), insertUser(user)
	}
	if user.HeadImgURL != oldUser.HeadImgURL || user.Nickname != oldUser.Nickname {
		update := bson.M{
			"headimgurl": user.HeadImgURL,
			"nickname":   user.Nickname,
			"sex":        user.Sex,
			"province":   user.Province,
			"city":       user.City,
			"country":    user.Country,
			"privilege":  user.Privilege,
		}
		return oldUser.ID.Hex(), updateUser(query, update)
	}
	return oldUser.ID.Hex(), nil
}

// 一个名字只能有一个联系人
func AddEmergCT(userID string, emergCT NamePhone) error {
	if !bson.IsObjectIdHex(userID) {
		return constant.ErrorIDFormatWrong
	}

	query := bson.M{
		"_id":           bson.ObjectIdHex(userID),
		"emerg_ct.name": emergCT.Name,
	}

	_, err := findUser(query, bson.M{
		"openid": 1,
	})
	if err != nil {
		update := bson.M{
			"$addToSet": emergCT,
		}
		return updateUser(query, update)
	}

	return constant.ErrorHasExist
}

func AddShareCT(userID, shareUserID string) error {
	return AddListenCT(shareUserID, userID)
}

// 一个名字只能有一个联系人
func DelEmergCT(userID, name string) error {
	if !bson.IsObjectIdHex(userID) {
		return constant.ErrorIDFormatWrong
	}

	query := bson.M{
		"_id":           bson.ObjectIdHex(userID),
		"emerg_ct.name": name,
	}

	_, err := findUser(query, bson.M{
		"openid": 1,
	})
	if err == nil {
		update := bson.M{
			"$pull": bson.M{
				"emerg_ct.name": name,
			},
		}
		return updateUser(query, update)
	}

	return constant.ErrorNotExist
}

/****************************************** user db action ****************************************/

func wrapUserDB(f func(*mgo.Collection) (interface{}, error)) (interface{}, error) {
	sess := globalSess.Clone()
	defer sess.Close()
	table := sess.DB(DBNAME).C(USERTABLE)

	return f(table)
}

func findUser(query, selector interface{}) (User, error) {
	sess := globalSess.Copy()
	defer sess.Close()
	userTable := sess.DB(DBNAME).C(USERTABLE)

	user := User{}
	err := userTable.Find(query).Select(selector).One(&user)
	return user, err
}

func findUsers(query, selector interface{}) ([]User, error) {
	sess := globalSess.Copy()
	defer sess.Close()
	userTable := sess.DB(DBNAME).C(USERTABLE)

	users := []User{}
	err := userTable.Find(query).Select(selector).All(&users)
	return users, err
}

func updateUser(query, update interface{}) error {
	sess := globalSess.Clone()
	defer sess.Close()
	userTable := sess.DB(DBNAME).C(USERTABLE)

	return userTable.Update(query, update)
}

func insertUser(docs ...interface{}) error {
	sess := globalSess.Clone()
	defer sess.Close()
	userTable := sess.DB(DBNAME).C(USERTABLE)

	return userTable.Insert(docs...)
}

/****************************************** user redis action ****************************************/

func SetRedisUserLocEQ(userID string, longitude, latitude float64, eq int32) (err error) {
	cntlr := NewRedisDBCntlr()
	defer cntlr.Close()

	key := fmt.Sprintf(constant.RedisUserLoc, userID)
	value := fmt.Sprintf(constant.RedisUserLocValue, longitude, latitude)
	_, err = cntlr.SET(key, value)
	if err != nil {
		return
	}

	key = fmt.Sprintf(constant.RedisUserEQ, userID)
	_, err = cntlr.SET(key, eq)
	return
}

func GetRedisUserLocEQ(userID string) (location Location, eq int, err error) {
	cntlr := NewRedisDBCntlr()
	defer cntlr.Close()

	location = Location{}
	eq = -1

	key := fmt.Sprintf(constant.RedisUserLoc, userID)
	value, err := cntlr.GET(key)
	if err != nil {
		return
	}
	loc := strings.Split(value, ",")
	if len(loc) != 2 {
		err = constant.ErrorNotExist
		return
	}
	location.Longitude, _ = strconv.ParseFloat(loc[0], 64)
	location.Latitude, _ = strconv.ParseFloat(loc[1], 64)

	key = fmt.Sprintf(constant.RedisUserEQ, userID)
	eq, err = cntlr.GETINT(key)
	return
}

func SetRedisUserStatus(userID string, status int32) error {
	cntlr := NewRedisDBCntlr()
	defer cntlr.Close()

	key := fmt.Sprintf(constant.RedisUserStatus, userID)
	_, err := cntlr.SET(key, status)
	return err
}

func GetRedisUserStatus(userID string) (int, error) {
	cntlr := NewRedisDBCntlr()
	defer cntlr.Close()

	key := fmt.Sprintf(constant.RedisUserStatus, userID)
	return cntlr.GETINT(key)
}
