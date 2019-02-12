package model

import (
	"config"
	"constant"
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// mgo controller
type MgoDBCntlr struct {
	sess *mgo.Session
	db   *mgo.Database
	c    *mgo.Collection
}

// DBNAME 数据库名字
var (
	DBNAME          = config.Conf.DB.DBName
	globalSess      *mgo.Session
	mongoURL        string
	DefaultSelector = bson.M{}
)

const (
	// USERTABLE user 表名
	USERTABLE     = "user"
	FEEDBACKTABLE = "feedback"

	MongoCopyType  = "1"
	MongoCloneType = "2"
)

func init() {
	dbConf := config.Conf.DB
	if dbConf.User != "" && dbConf.PW != "" {
		mongoURL = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", dbConf.User, dbConf.PW, dbConf.Host, dbConf.Port, dbConf.AdminDBName)
	} else {
		mongoURL = fmt.Sprintf("mongodb://%s:%s", dbConf.Host, dbConf.Port)
	}

	var err error
	globalSess, err = GetDBSession()
	if err != nil {
		panic(err)
	}
}

/****************************************** db session manage ****************************************/

// GetSession get the db session
func GetDBSession() (*mgo.Session, error) {
	globalMgoSession, err := mgo.DialWithTimeout(mongoURL, 10*time.Second)
	if err != nil {
		return nil, err
	}
	globalMgoSession.SetMode(mgo.Monotonic, true)
	//default is 4096
	globalMgoSession.SetPoolLimit(1000)
	return globalMgoSession, nil
}

func GetCloneSess() *mgo.Session {
	return globalSess.Clone()
}

func GetCopySess() *mgo.Session {
	return globalSess.Copy()
}

/********************************************* MgoDBCntlr *******************************************/

/*
	args 说明：
	一个参数：tableName 表名字，这时采用默认db 使用copy方式
	两个参数：tableName, sessType "1" 表示copy方式 "2" 表示clone方式
	三个参数：tableName, sessType, dbName 数据库名字
*/
func NewMgoDBCntlr(args ...string) *MgoDBCntlr {
	mgoSess := &MgoDBCntlr{}

	if len(args) <= 3 {
		if len(args) >= 2 && args[1] == "2" {
			mgoSess.sess = globalSess.Clone()
		} else {
			mgoSess.sess = globalSess.Copy()
		}

		if len(args) == 3 {
			mgoSess.db = mgoSess.sess.DB(args[2])

		} else {
			mgoSess.db = mgoSess.sess.DB(DBNAME)
		}
		mgoSess.c = mgoSess.db.C(args[0])
	}

	return mgoSess
}

func (this *MgoDBCntlr) Close() {
	this.sess.Close()
}

func (this *MgoDBCntlr) SetTableName(tableName string) {
	this.c = this.db.C(tableName)
}

/*
	args: query(interface{}), result(interface{}), select(interface{})
	说明：参数必须要求query、result，其他的参数依次递增，即若想使用后面的参数必须有前面的所有参数
*/
func (this *MgoDBCntlr) Find(args ...interface{}) error {
	var mgoQuery *mgo.Query
	if len(args) < 2 || len(args) > 3 {
		return constant.ErrorOutOfRange
	}
	for i := 0; i < len(args); i++ {
		switch i {
		case 0:
			mgoQuery = this.c.Find(args[i])
		case 2:
			mgoQuery = mgoQuery.Select(args[i])
		}
	}
	return mgoQuery.One(args[1])
}

/*
	args: query(interface{}), result(interface{}), select(interface{}), limit(int), skip(int), sort([]string)
	说明：参数必须要求query、result，其他的参数依次递增，即若想使用后面的参数必须有前面的所有参数
*/
func (this *MgoDBCntlr) FindAll(args ...interface{}) error {
	var mgoQuery *mgo.Query
	if len(args) < 2 || len(args) > 6 {
		return constant.ErrorOutOfRange
	}
	for i := 0; i < len(args); i++ {
		switch i {
		case 0:
			mgoQuery = this.c.Find(args[i])
		case 2:
			mgoQuery = mgoQuery.Select(args[i])
		case 3:
			if limit, ok := args[i].(int); ok {
				mgoQuery = mgoQuery.Limit(limit)
			}
		case 4:
			if skip, ok := args[4].(int); ok {
				mgoQuery = mgoQuery.Skip(skip)
			}
		case 5:
			if sort, ok := args[i].([]string); ok {
				mgoQuery = mgoQuery.Sort(sort...)
			}
		}
	}
	return mgoQuery.All(args[1])
}

func (this *MgoDBCntlr) FindCount(query interface{}) (int, error) {
	return this.c.Find(query).Count()
}

func (this *MgoDBCntlr) Update(query, update interface{}) error {
	return this.c.Update(query, update)
}

func (this *MgoDBCntlr) UpdateAll(query, update interface{}) (*mgo.ChangeInfo, error) {
	return this.c.UpdateAll(query, update)
}
